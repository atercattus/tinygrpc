package http2

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/http2/hpack"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

type (
	Server struct {
		listeners []net.Listener

		handlerHTTP       Handler
		handlerGRPC       Handler
		handlerGRPCStream StreamHandler
		grpcServiceDesc   grpc.ServiceDesc

		activeConnsGroup errgroup.Group
	}
)

func NewServer() *Server {
	s := &Server{}

	s.handlerHTTP = s.defaultHTTPHandler
	s.handlerGRPC = s.defaultGRPCHandler
	s.handlerGRPCStream = s.defaultGRPCStreamHandler

	return s
}

// Listen adds listener for future Serve()
//
// You can initialize a listener:
//
//	l, _ := net.Listen("tcp", ":80")
//	NewServer().Listen(l)
func (s *Server) Listen(l net.Listener) {
	s.listeners = append(s.listeners, l)
}

// ListenTLS wraps listener to TLS with certs
//
// You can load a cert from pair:
//
//	l, _ := net.Listen("tcp", ":443")
//	cert, _ := tls.X509KeyPair([]byte(certPEMBlock), []byte(keyPEMBlock))
//	NewServer().ListenTLS(l, cert)
func (s *Server) ListenTLS(l net.Listener, certs ...tls.Certificate) {
	l = tls.NewListener(l, &tls.Config{
		Certificates: certs,
		NextProtos:   []string{"h2"}, // http2.NextProtoTLS
	})
	s.Listen(l)
}

func (s *Server) ListenTCP(addr string) error {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	s.Listen(l)

	return nil
}

func (s *Server) SetGRPPCHandler(desc grpc.ServiceDesc, h Handler) {
	s.handlerGRPC = h
	s.grpcServiceDesc = desc
}

func (s *Server) SetGRPPCStreamHandler(desc grpc.ServiceDesc, h StreamHandler) {
	s.handlerGRPCStream = h
	s.grpcServiceDesc = desc
}

func (s *Server) SetHTTPHandler(h Handler) {
	s.handlerHTTP = h
}

func (s *Server) Serve(ctx context.Context) error {
	var wg errgroup.Group

	go func() {
		<-ctx.Done()

		var eg errgroup.Group
		for _, l := range s.listeners {
			eg.Go(func() error {
				err := l.Close()
				if (err != nil) && strings.Contains(err.Error(), "use of closed network connection") {
					err = nil
				}
				return err
			})
		}

		if err := eg.Wait(); err != nil {
			log.Printf("Close listeners: %s", err)
		}

		if err := s.activeConnsGroup.Wait(); err != nil {
			log.Printf("Finish active connections: %s", err)
		}
	}()

	for _, l := range s.listeners {
		wg.Go(func() error {
			// log.Println("Serve connections on", l.Addr())
			if err := s.serve(ctx, l); err != nil {
				if ctx.Err() != nil {
					return nil
				}
				return fmt.Errorf("serve %s: %s", l.Addr(), err)
			}
			return nil
		})
	}

	return wg.Wait()
}

func (s *Server) serve(ctx context.Context, l net.Listener) error {
	for {
		conn, err := s.acceptConn(ctx, l)
		if err != nil {
			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				return err
			}

			if ne, ok := err.(interface {
				Temporary() bool
			}); ok && ne.Temporary() {
				s.randSleep(ctx, time.Millisecond)
				continue
			}

			return fmt.Errorf("accept new connection: %w", err)
		}

		s.activeConnsGroup.Go(func() error {
			s.serveConn(ctx, conn)
			return nil
		})
	}
}

func (s *Server) acceptConn(ctx context.Context, l net.Listener) (net.Conn, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	if dl, ok := l.(interface {
		SetDeadline(t time.Time) error
	}); ok {
		_ = dl.SetDeadline(time.Now().Add(time.Second))
	}

	return l.Accept()
}

func (s *Server) serveConn(ctx context.Context, netConn net.Conn) {
	conn := newServerConn(netConn)
	defer conn.Close()

	initialRequest, err := conn.ServerHandshake(ctx)
	if err != nil {
		fmt.Printf("Handshake with %s failed: %s\n", netConn.RemoteAddr(), err)
		return
	}

	// log.Println("New connection with", netConn.RemoteAddr())

	if err := s.handler(ctx, initialRequest, conn); err != nil {
		switch {
		case errors.Is(err, context.Canceled):
		default:
			fmt.Printf("Serve connection with %s failed: %s\n", netConn.RemoteAddr(), err)
		}
	}
}

func (s *Server) randSleep(ctx context.Context, d time.Duration) {
	if d == 0 {
		return
	}

	const jitter = 10
	halfRange := int64(d) * jitter / 100 // +/- jitter percents
	d += time.Duration(rand.Int63n(2*halfRange) - halfRange)

	t := time.NewTicker(d)
	defer t.Stop()

	select {
	case <-t.C:
		return
	case <-ctx.Done():
		return
	}
}

func (s *Server) defaultHTTPHandler(_ context.Context, _ *Request, response *Response) error {
	response.Status = 404
	return nil
}

func (s *Server) defaultGRPCHandler(_ context.Context, _ *Request, response *Response) error {
	response.Status = 12 // UNIMPLEMENTED
	return nil
}

func (s *Server) defaultGRPCStreamHandler(_ context.Context, request *Request, conn *Conn) error {
	const grpcStatus = 12 // UNIMPLEMENTED
	return s.sendGRPCTrailers(request.StreamId, grpcStatus, conn)
}

func (s *Server) handler(ctx context.Context, request *Request, conn *Conn) error {
	var err error

	// request will not be nil after HTTP/1.1 Connection Upgrade only

	for {
		if request == nil {
			request, err = conn.RecvRequestHeaders(ctx)
			if err != nil {
				if errors.Is(err, io.EOF) {
					return nil
				}
				if strings.Contains(err.Error(), "i/o timeout") {
					return nil
				}
				return fmt.Errorf("RecvRequestHeaders: %w", err)
			}
		}

		if request.Headers["content-type"] == "application/grpc" {
			if err := s.processGRPCAsync(ctx, request, conn); err != nil {
				return fmt.Errorf("processGRPCAsync (%v): %w", request, err)
			}
			request = nil
			continue
		}

		if err := conn.RecvRequestBody(ctx, request); err != nil {
			return fmt.Errorf("RecvRequestBody: %w", err)
		}

		go func(request *Request) {
			if err := s.processHTTPAsync(ctx, request, conn); err != nil {
				// ToDo: http2 error
				log.Println("Cannot process HTTP REQUEST:", err.Error())
			}
		}(request)

		request = nil
	}
}

func (s *Server) processHTTPAsync(ctx context.Context, request *Request, conn *Conn) error {
	var response Response
	err := s.handlerHTTP(ctx, request, &response)
	if err != nil {
		return fmt.Errorf("http handler: %w", err)
	}

	var flags FrameFlags
	if len(response.Body) == 0 {
		flags |= FlagEndStream
	}

	respHeaders := []hpack.HeaderField{
		{Name: `:status`, Value: strconv.FormatUint(response.Status, 10)},
		{Name: `content-length`, Value: strconv.Itoa(len(response.Body))},
		{Name: `date`, Value: time.Now().Format(time.RFC1123)},
	}

	if response.Headers != nil {
		for k, v := range response.Headers {
			respHeaders = append(respHeaders, hpack.HeaderField{
				Name:  k,
				Value: v,
			})
		}
	}

	if err := conn.SendHeaders(request.StreamId, flags, respHeaders); err != nil {
		return fmt.Errorf("http handler send headers: %w", err)
	}

	if len(response.Body) > 0 {
		if err := conn.SendData(request.StreamId, FlagEndStream, response.Body); err != nil {
			return fmt.Errorf("http handler send body: %w", err)
		}
	}

	return nil
}

func (s *Server) processGRPCAsync(ctx context.Context, request *Request, conn *Conn) error {
	if s.isStreamingGRPC(request.URI) {
		if err := s.handlerGRPCStream(ctx, request, conn); err != nil {
			return fmt.Errorf("handlerGRPCStream: %w", err)
		}
		return nil
	}

	if err := conn.RecvRequestBody(ctx, request); err != nil {
		return fmt.Errorf("RecvRequestBody: %w", err)
	}

	data := request.Body

	// https://github.com/grpc/grpc/blob/master/doc/PROTOCOL-HTTP2.md
	// Length-Prefixed-Message → Compressed-Flag Message-Length Message
	// Compressed-Flag → 0 / 1 # encoded as 1 byte unsigned integer
	// Message-Length → {length of Message} # encoded as 4 byte unsigned integer (big endian)
	// Message → *{binary octet}
	if data[0] != 0 {
		return fmt.Errorf("compression is not supported")
	}
	messageLength := endianess.Uint32(data[1:5])
	if expect := uint32(len(data)) - 5; messageLength != expect {
		// The repeated sequence of Length-Prefixed-Message items is delivered in DATA frames
		return fmt.Errorf("expect message len %d, got %d", expect, messageLength)
	}

	request.Body = data[5:]

	go func() {
		if err := s.processGRPCUnary(ctx, request, conn); err != nil {
			// ToDo: grpc error
			log.Println("Cannot process GRPC unary call:", err.Error())
		}
	}()

	return nil
}

func (s *Server) isStreamingGRPC(requestURI string) bool {
	for _, stream := range s.grpcServiceDesc.Streams {
		if requestURI == "/"+s.grpcServiceDesc.ServiceName+"/"+stream.StreamName {
			return true
		}
	}

	return false
}

func (s *Server) processGRPCUnary(ctx context.Context, request *Request, conn *Conn) error {
	var response Response
	err := s.handlerGRPC(ctx, request, &response)
	if err != nil {
		return fmt.Errorf("process GRPC unary call: %w", err)
	}

	var responseRaw []byte
	responseRaw = append(responseRaw, 0) // Compressed-Flag
	responseRaw = endianess.AppendUint32(responseRaw, uint32(len(response.Body)))
	responseRaw = append(responseRaw, response.Body...)

	respHeaders := []hpack.HeaderField{
		{Name: `:status`, Value: "200"},
		{Name: `content-type`, Value: `application/grpc+proto`},
	}

	if response.Headers != nil {
		for k, v := range response.Headers {
			respHeaders = append(respHeaders, hpack.HeaderField{
				Name:  k,
				Value: v,
			})
		}
	}

	if err := conn.SendHeaders(request.StreamId, 0, respHeaders); err != nil {
		return fmt.Errorf("send grpc response headers: %w", err)
	}

	if err := conn.SendData(request.StreamId, 0, responseRaw); err != nil {
		return fmt.Errorf("send grpc response headers: %w", err)
	}

	if err := s.sendGRPCTrailers(request.StreamId, response.Status, conn); err != nil {
		return fmt.Errorf("send grpc response headers: %w", err)
	}

	return nil
}

func (s *Server) processGRPCStreaming(ctx context.Context, request *Request, conn *Conn) error {
	if err := s.handlerGRPCStream(ctx, request, conn); err != nil {
		return fmt.Errorf("processGRPCStreaming: %w", err)
	}

	const grpcStatus = 0
	return s.sendGRPCTrailers(request.StreamId, grpcStatus, conn)
}

func (s *Server) sendGRPCTrailers(streamId uint32, status uint64, conn *Conn) error {
	return conn.SendTrailerFrame(streamId, []hpack.HeaderField{
		{Name: `grpc-status`, Value: strconv.FormatUint(status, 10)},
	})
}
