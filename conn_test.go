package http2

import (
	"bytes"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	nh2 "golang.org/x/net/http2"
	"golang.org/x/net/http2/hpack"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	"github.com/atercattus/tinygrpc/testdata"
	fastGrpcTestPb "github.com/atercattus/tinygrpc/testdata/proto"
)

func httpHandler(_ context.Context, request *Request, response *Response) error {
	response.Headers = map[string]string{
		"content-type": "text/plain; charset=utf-8",
	}

	if request.Method == http.MethodGet {
		response.Body = []byte(fmt.Sprintf("Hello, %v", request.URI))
	} else {
		response.Body = bytes.ToUpper(request.Body)
	}

	response.Status = 200

	return nil
}

func grpcHandler(_ context.Context, request *Request, response *Response) error {
	switch request.URI {
	case `/fastgrpctest.v1.FastGrpcTestService/Ping`:
		var pingRequest fastGrpcTestPb.PingRequest
		err := proto.Unmarshal(request.Body, &pingRequest)
		if err != nil {
			return fmt.Errorf("unmarshal grpc request: %w", err)
		}

		// log.Printf("grpc request: %d", pingRequest.Val)

		var pingResponse fastGrpcTestPb.PingResponse
		pingResponse.Val = pingRequest.Val * 2

		response.Body, err = proto.Marshal(&pingResponse)
		if err != nil {
			return fmt.Errorf("marshal grpc response: %w", err)
		}

	case `/fastgrpctest.v1.FastGrpcTestService/Sleep`:
		var sleepRequest fastGrpcTestPb.SleepRequest
		err := proto.Unmarshal(request.Body, &sleepRequest)
		if err != nil {
			return fmt.Errorf("unmarshal grpc request: %w", err)
		}

		// log.Printf("grpc request: %d", sleepRequest.Duration)

		time.Sleep(time.Duration(sleepRequest.Duration))

		var sleepResponse fastGrpcTestPb.SleepResponse
		sleepResponse.Duration = sleepRequest.Duration
		sleepResponse.Now = time.Now().UnixNano()

		response.Body, err = proto.Marshal(&sleepResponse)
		if err != nil {
			return fmt.Errorf("marshal grpc response: %w", err)
		}

	default:
		return fmt.Errorf("unknown method %q", request.URI)
	}

	return nil
}

func Test_H2C_HTTP2PriorKnowledge(t *testing.T) {
	l, err := net.Listen("tcp", ":0")
	require.NoError(t, err)
	defer l.Close()

	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	h2s := NewServer()
	h2s.SetHTTPHandler(httpHandler)
	h2s.Listen(l)
	go func() {
		require.NoError(t, h2s.Serve(ctx))
	}()

	var client http.Client
	client.Transport = &nh2.Transport{
		AllowHTTP: true,
		DialTLS: func(netw, addr string, cfg *tls.Config) (net.Conn, error) { // Skip TLS Dial
			return net.Dial(netw, addr)
		},
	}

	uriPath := fmt.Sprintf("/hello?world=%d", rand.Int())
	uri := "http://" + l.Addr().String() + uriPath

	t.Run("GET", func(t *testing.T) {
		resp, err := client.Get(uri)
		require.NoError(t, err)

		defer resp.Body.Close()

		b, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		require.Equal(t, 200, resp.StatusCode)
		require.Equal(t, `Hello, `+uriPath, string(b))
	})

	t.Run("POST", func(t *testing.T) {
		body := []byte(uri)
		bodyRdr := bytes.NewReader(body)

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri, bodyRdr)
		require.NoError(t, err)

		resp, err := client.Do(req)
		require.NoError(t, err)

		defer resp.Body.Close()

		b, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		require.Equal(t, 200, resp.StatusCode)
		require.Equal(t, string(bytes.ToUpper(body)), string(b))
	})

	t.Run("POST large body", func(t *testing.T) {
		body := bytes.Repeat([]byte("Hello, traveler!"), 1000)
		bodyRdr := bytes.NewReader(body)

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri, bodyRdr)
		require.NoError(t, err)

		resp, err := client.Do(req)
		require.NoError(t, err)

		defer resp.Body.Close()

		b, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		require.Equal(t, 200, resp.StatusCode)
		require.Equal(t, string(bytes.ToUpper(body)), string(b))
	})
}

func Test_H2_HTTP2PriorKnowledge(t *testing.T) {
	l, err := net.Listen("tcp", ":0")
	require.NoError(t, err)
	defer l.Close()

	cert, err := tls.X509KeyPair(testdata.CertPEMBlock, testdata.KeyPEMBlock)
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	h2s := NewServer()
	h2s.SetHTTPHandler(httpHandler)
	h2s.ListenTLS(l, cert)
	go func() {
		require.NoError(t, h2s.Serve(ctx))
	}()

	var client http.Client
	client.Transport = &nh2.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	var uri = fmt.Sprintf("/hello?world=%d", rand.Int())

	resp, err := client.Get("https://" + l.Addr().String() + uri)
	require.NoError(t, err)

	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	require.Equal(t, 200, resp.StatusCode)
	require.Equal(t, `Hello, `+uri, string(b))
}

func Test_DefaultHTTP2handler(t *testing.T) {
	l, err := net.Listen("tcp", ":0")
	require.NoError(t, err)
	defer l.Close()

	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	h2s := NewServer()
	// Do not set HTTP handler
	h2s.Listen(l)
	go func() {
		require.NoError(t, h2s.Serve(ctx))
	}()

	var client http.Client
	client.Transport = &nh2.Transport{
		AllowHTTP: true,
		DialTLS: func(netw, addr string, cfg *tls.Config) (net.Conn, error) { // Skip TLS Dial
			return net.Dial(netw, addr)
		},
	}

	resp, err := client.Get("http://" + l.Addr().String())
	require.NoError(t, err)

	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	require.Equal(t, 404, resp.StatusCode)
	require.Len(t, b, 0)
}

func Test_BatchOfParallelHTTPRequests(t *testing.T) {
	l, err := net.Listen("tcp", ":0")
	require.NoError(t, err)
	defer l.Close()

	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	h2s := NewServer()
	h2s.SetHTTPHandler(httpHandler)
	h2s.Listen(l)
	go func() {
		require.NoError(t, h2s.Serve(ctx))
	}()

	var client http.Client
	client.Transport = &nh2.Transport{
		AllowHTTP: true,
		DialTLS: func(netw, addr string, cfg *tls.Config) (net.Conn, error) { // Skip TLS Dial
			return net.Dial(netw, addr)
		},
	}

	var wg sync.WaitGroup
	for i := 0; i < 500; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			resp, err := client.Get(fmt.Sprintf("http://%s/hello?world=%d", l.Addr().String(), rand.Int()))
			require.NoError(t, err)

			defer resp.Body.Close()

			_, err = io.Copy(io.Discard, resp.Body)
			require.NoError(t, err)

			require.Equal(t, 200, resp.StatusCode)
		}()
	}
	wg.Wait()
}

func Test_DefaultGRPCHandler(t *testing.T) {
	l, err := net.Listen("tcp", ":0")
	require.NoError(t, err)
	defer l.Close()

	cert, err := tls.X509KeyPair(testdata.CertPEMBlock, testdata.KeyPEMBlock)
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	h2s := NewServer()
	// Do not set gRPC handler
	h2s.ListenTLS(l, cert)
	go func() {
		require.NoError(t, h2s.Serve(ctx))
	}()

	conn, err := grpc.Dial(l.Addr().String(), []grpc.DialOption{
		grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{InsecureSkipVerify: true})),
	}...)
	require.NoError(t, err)
	defer conn.Close()

	client := fastGrpcTestPb.NewFastGrpcTestServiceClient(conn)
	_, err = client.Ping(ctx, &fastGrpcTestPb.PingRequest{Val: 42})

	st, ok := status.FromError(err)
	require.True(t, ok)
	require.Equal(t, codes.Unimplemented, st.Code())
}

func Test_UnaryGRPCCall(t *testing.T) {
	l, err := net.Listen("tcp", ":0")
	require.NoError(t, err)
	defer l.Close()

	cert, err := tls.X509KeyPair(testdata.CertPEMBlock, testdata.KeyPEMBlock)
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	h2s := NewServer()
	h2s.SetGRPPCHandler(fastGrpcTestPb.FastGrpcTestService_ServiceDesc, grpcHandler)
	h2s.ListenTLS(l, cert)
	go func() {
		require.NoError(t, h2s.Serve(ctx))
	}()

	conn, err := grpc.Dial(l.Addr().String(), []grpc.DialOption{
		grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{InsecureSkipVerify: true})),
	}...)
	require.NoError(t, err)
	defer conn.Close()

	const reqVal = 42
	client := fastGrpcTestPb.NewFastGrpcTestServiceClient(conn)
	resp, err := client.Ping(ctx, &fastGrpcTestPb.PingRequest{Val: reqVal})
	require.NoError(t, err)

	// log.Printf("grpc response: %+v", resp.Val)

	require.EqualValues(t, 2*reqVal, resp.Val)
}

func Test_UnaryGRPCCallInParallel(t *testing.T) {
	l, err := net.Listen("tcp", ":0")
	require.NoError(t, err)
	defer l.Close()

	cert, err := tls.X509KeyPair(testdata.CertPEMBlock, testdata.KeyPEMBlock)
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	h2s := NewServer()
	h2s.SetGRPPCHandler(fastGrpcTestPb.FastGrpcTestService_ServiceDesc, grpcHandler)
	h2s.ListenTLS(l, cert)
	go func() {
		require.NoError(t, h2s.Serve(ctx))
	}()

	conn, err := grpc.Dial(l.Addr().String(), []grpc.DialOption{
		grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{InsecureSkipVerify: true})),
	}...)
	require.NoError(t, err)
	defer conn.Close()

	client := fastGrpcTestPb.NewFastGrpcTestServiceClient(conn)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		resp, err := client.Sleep(ctx, &fastGrpcTestPb.SleepRequest{Duration: (time.Second).Nanoseconds()})
		require.NoError(t, err)

		log.Printf("grpc response: %d %s", resp.Duration, time.Unix(0, resp.Now))
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		time.Sleep(100 * time.Millisecond)

		resp, err := client.Sleep(ctx, &fastGrpcTestPb.SleepRequest{Duration: (55 * time.Millisecond).Nanoseconds()})
		require.NoError(t, err)

		log.Printf("grpc response: %d %s", resp.Duration, time.Unix(0, resp.Now))
	}()

	wg.Wait()
}

func Test_GRPCStreams(t *testing.T) {
	handler := func(ctx context.Context, request *Request, conn *Conn) error {
		data, _, err := recvGRPCDataFrame(ctx, conn)
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			return fmt.Errorf("recv stream frame: %w", err)
		}

		var req fastGrpcTestPb.SubRequest
		if err := proto.Unmarshal(data, &req); err != nil {
			return fmt.Errorf("streaming unmarshal: %w", err)
		}

		err = conn.SendHeaders(request.StreamId, 0, []hpack.HeaderField{
			{Name: `:status`, Value: "200"},
			{Name: `content-type`, Value: `application/grpc+proto`},
		})
		if err != nil {
			return fmt.Errorf("send grpc response headers: %w", err)
		}

		ii := 0
		for range time.Tick(time.Duration(req.PubInterval)) {
			var r fastGrpcTestPb.SubResponse
			r.CurrentTime = time.Now().UnixNano()
			b, _ := proto.Marshal(&r)

			err := sendGRPCDataFrame(conn, request.StreamId, 0, b)
			if err != nil {
				log.Println(fmt.Errorf("send grpc sreaming response headers: %w", err))
			}

			ii++
			if ii >= 2 {
				break
			}
		}

		return nil
	}

	l, err := net.Listen("tcp", ":0")
	require.NoError(t, err)
	defer l.Close()

	cert, err := tls.X509KeyPair(testdata.CertPEMBlock, testdata.KeyPEMBlock)
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	h2s := NewServer()
	h2s.SetGRPPCStreamHandler(fastGrpcTestPb.FastGrpcTestService_ServiceDesc, handler)
	h2s.ListenTLS(l, cert)
	go func() {
		require.NoError(t, h2s.Serve(ctx))
	}()

	conn, err := grpc.Dial(l.Addr().String(), []grpc.DialOption{
		grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{InsecureSkipVerify: true})),
	}...)
	require.NoError(t, err)
	defer conn.Close()

	client := fastGrpcTestPb.NewFastGrpcTestServiceClient(conn)

	subCli, err := client.Sub(ctx)
	require.NoError(t, err)

	err = subCli.Send(&fastGrpcTestPb.SubRequest{PubInterval: int64(250 * time.Millisecond)})
	require.NoError(t, err)

	now := time.Now()

	for i := 0; i < 2; i++ {
		msg, err := subCli.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		require.NoError(t, err)

		if err != nil {
			st, err := status.FromError(err)
			log.Println("stream client error", st, err)
			break
		}

		recvTime := time.Unix(0, msg.CurrentTime)
		log.Printf("sub recv new: %s", recvTime)

		require.Greater(t, recvTime, now)
	}
}

func sendGRPCDataFrame(conn *Conn, streamId uint32, flags FrameFlags, body []byte) error {
	// ToDo: it's same as in Server.processGRPCUnary

	var responseRaw []byte
	responseRaw = append(responseRaw, 0) // Compressed-Flag
	responseRaw = endianess.AppendUint32(responseRaw, uint32(len(body)))
	responseRaw = append(responseRaw, body...)

	return conn.SendData(streamId, flags, responseRaw)
}

func recvGRPCDataFrame(ctx context.Context, conn *Conn) ([]byte, *FrameHdr, error) {
	f, err := conn.RecvFrame(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("RecvFrame: %w", err)
	}

	df, ok := f.(*DataFrame)
	if !ok {
		return nil, nil, fmt.Errorf("expect DATA frame, got %T", f)
	}

	data := df.Data

	// ToDo: it's same as in Server.processGRPCAsync

	// https://github.com/grpc/grpc/blob/master/doc/PROTOCOL-HTTP2.md
	// Length-Prefixed-Message → Compressed-Flag Message-Length Message
	// Compressed-Flag → 0 / 1 # encoded as 1 byte unsigned integer
	// Message-Length → {length of Message} # encoded as 4 byte unsigned integer (big endian)
	// Message → *{binary octet}
	if data[0] != 0 {
		return nil, nil, fmt.Errorf("compression is not supported")
	}
	messageLength := endianess.Uint32(data[1:5])
	if expect := uint32(len(data)) - 5; messageLength != expect {
		// The repeated sequence of Length-Prefixed-Message items is delivered in DATA frames
		return nil, nil, fmt.Errorf("expect message len %d, got %d", expect, messageLength)
	}

	return data[5:], &df.FrameHdr, nil
}
