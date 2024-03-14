package http2

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/http2/hpack"
)

type (
	// Conn implements HTTP/2 connection
	Conn struct {
		conn net.Conn

		writeLock sync.Mutex

		recvBuf []byte

		hpackDecoder *hpack.Decoder

		hpackEncoder    *hpack.Encoder
		hpackEncoderBuf bytes.Buffer
	}
)

var (
	endianess = binary.BigEndian
)

const (
	firstClientStreamId = 1
)

func newServerConn(conn net.Conn) *Conn {
	c := &Conn{
		conn: conn,

		recvBuf: make([]byte, RecvBufferSize),
	}

	return c
}

func (c *Conn) Close() error {
	return c.conn.Close()
}

func (c *Conn) ServerHandshake(ctx context.Context) (*Request, error) {
	c.conn.SetReadDeadline(time.Now().Add(1 * time.Second))
	n, err := c.conn.Read(c.recvBuf[:cap(c.recvBuf)])
	if err != nil {
		return nil, err
	}

	reqBuf := c.recvBuf[:n]

	// Is it prior-knowledge (magic PRISM)?
	if bytes.HasPrefix(reqBuf, connectionPreface) {
		n := copy(reqBuf, reqBuf[len(connectionPreface):])
		c.conn = NewBufReader(c.conn, reqBuf[:n])
		return nil, c.serverAfterClientPreface(ctx)
	}

	request := NewHTTPRequest(firstClientStreamId)

	// Read all HTTP/1 headers
	for {
		sepPos := bytes.Index(reqBuf, httpHeaderSep)
		if sepPos == -1 {
			return nil, NewErrWrongHTTPRequest(ErrTooLongRequest)
		}

		hdr := reqBuf[:sepPos]
		reqBuf = reqBuf[sepPos+len(httpHeaderSep):]

		if len(hdr) == 0 {
			break
		}
		if err := request.parseHTTPHeader(hdr); err != nil {
			return nil, NewErrWrongHTTPRequest(err)
		}
	}

	if len(reqBuf) > 0 {
		// This tiny implementation do NOT support reading request from multiple tcp frames
		return nil, NewErrWrongHTTPRequest(fmt.Errorf("there are data after the request headers: %q", reqBuf))
	}

	// ToDo: process request.Headers["connection"]

	if request.Headers["upgrade"] != "h2c" {
		return nil, NewErrWrongHTTPRequest(ErrWrongUpgradeHeader)
	}

	if hdr := request.Headers["http2-settings"]; len(hdr) > 0 {
		frame, err := c.parseHTTP2SettingsHeader(hdr)
		if err != nil {
			return nil, NewErrWrongHTTPRequest(fmt.Errorf("wrong HTTP2-Settings header: %w", err))
		}
		log.Printf("client HTTP2-Settings: %+v", &frame)
		// ToDo: apply
	}

	// ToDo: https://http.dev/426 Upgrade Required

	// Send "101 Switching Protocols" response
	_ = c.conn.SetWriteDeadline(time.Now().Add(1 * time.Second))
	_, err = c.conn.Write(httpSwitchingProtoResponse)
	if err != nil {
		return nil, fmt.Errorf("send 101: %w", err)
	}

	// Read magic PRISM
	_ = c.conn.SetReadDeadline(time.Now().Add(1 * time.Second))
	n, err = c.conn.Read(c.recvBuf[:len(connectionPreface)])
	if err != nil {
		return nil, err
	}
	if !bytes.Equal(c.recvBuf[:n], connectionPreface) {
		return nil, NewErrWrongHTTPRequest(fmt.Errorf("expect %q, got %q", connectionPreface, c.recvBuf[:n]))
	}

	return request, c.serverAfterClientPreface(ctx)
}

func (c *Conn) serverAfterClientPreface(_ context.Context) error {
	// https://datatracker.ietf.org/doc/html/rfc9113#name-starting-http-2-with-prior-

	settings := GetDefaultSettings()
	err := c.sendFrame(FrameTypeSettings, 0, 0, settings.ToSettingsFrame().WritePayload(nil))
	if err != nil {
		return err
	}

	c.hpackEncoder = hpack.NewEncoder(&c.hpackEncoderBuf)
	//c.hpackEncoder.SetMaxDynamicTableSize(1048896)
	c.hpackEncoder.SetMaxDynamicTableSizeLimit(initialHeaderTableSize)

	c.hpackDecoder = hpack.NewDecoder(initialMaxHeaderListSize, nil)

	return nil
}

func (c *Conn) RecvFrame(ctx context.Context) (Frame, error) {
	for {
		f, err := c.recvFrame(ctx)
		if err != nil {
			if strings.Contains(err.Error(), "i/o timeout") {
				after := time.NewTimer(100 * time.Millisecond)
				select {
				case <-after.C:
				case <-ctx.Done():
					if !after.Stop() {
						select {
						case <-after.C:
						default:
						}
					}
					return nil, ctx.Err()
				}
				continue
			}
			return nil, fmt.Errorf("recv stream frame: %w", err)
		}

		return f, nil
	}
}

// recvFrame receives the next frame, independently processing service ones, and returns only DATA or HEADERS frames
func (c *Conn) recvFrame(ctx context.Context) (Frame, error) {
	for {
		f, err := c.recvFrameInternal(ctx)
		if err != nil {
			return nil, err
		}

		switch ft := f.(type) {
		case *SettingsFrame:
			if ft.Flags&FlagAck == 0 {
				// ToDo: apply
				if err := c.sendFrame(FrameTypeSettings, FlagAck, ft.StreamId, nil); err != nil {
					return nil, fmt.Errorf("send SETTINGS[%d] ACK: %w", ft.StreamId, err)
				}
			}

		case *WindowUpdateFrame:
			// ToDo: apply

		case *PingFrame:
			if ft.Flags&FlagAck == 0 {
				if err := c.sendFrame(FrameTypePing, FlagAck, ft.StreamId, ft.WritePayload(nil)); err != nil {
					return nil, fmt.Errorf("send PING ACK: %w", err)
				}
			}

		case *GoawayFrame:
			log.Printf("recv GOAWAY: streamId:%d errorCode:%d errorDescr:%s", ft.StreamId, ft.ErrorCode, string(ft.DebugData))
			return f, io.EOF

		case *RstStreamFrame:
			log.Printf("recv RST_STREAM: streamId:%d errorCode:%d", ft.StreamId, ft.ErrorCode)
			return f, io.EOF

		default:
			return f, nil
		}
	}
}

func (c *Conn) recvFrameInternal(_ context.Context) (Frame, error) {
	var (
		buf      [9]byte
		buf2     [4]byte
		frameHdr FrameHdr
		payload  []byte
	)

	_ = c.conn.SetReadDeadline(time.Now().Add(1 * time.Second))
	if n, err := io.ReadAtLeast(c.conn, buf[:], len(buf)); err != nil {
		return nil, err
	} else if n < len(buf) {
		return nil, fmt.Errorf("read frame header: expect %d bytes, got %d: %q", len(buf), n, buf[:n])
	}
	copy(buf2[1:], buf[0:3])
	frameHdr.Length = endianess.Uint32(buf2[:])

	frameHdr.Type = FrameType(buf[3])
	frameHdr.Flags = FrameFlags(buf[4])
	frameHdr.StreamId = endianess.Uint32(buf[5:9]) & 0x7FFFFFFF

	if frameHdr.Length > 0 {
		payload = make([]byte, frameHdr.Length)
		_ = c.conn.SetReadDeadline(time.Now().Add(1 * time.Second))
		if _, err := io.ReadAtLeast(c.conn, payload, int(frameHdr.Length)); err != nil {
			return nil, err
		}
	}

	var frame Frame

	switch frameHdr.Type {
	case FrameTypeSettings:
		frame = &SettingsFrame{}

	case FrameTypeWindowUpdate:
		frame = &WindowUpdateFrame{}

	case FrameTypeGoaway:
		frame = &GoawayFrame{}

	case FrameTypeHeaders:
		frame = &HeadersFrame{
			hpackDecoder: c.hpackDecoder,
		}

	case FrameTypeData:
		frame = &DataFrame{}

	case FrameTypePing:
		frame = &PingFrame{}

	case FrameTypeRstStream:
		frame = &RstStreamFrame{}

	default:
		err := &ErrProtocolError{Descr: `Unsupported frame type ` + strconv.Itoa(int(frameHdr.Type))}
		return nil, err
	}

	*(frame.Hdr()) = frameHdr

	if err := frame.ReadPayload(&frameHdr, payload); err != nil {
		err := &ErrProtocolError{Descr: err.Error()}
		return nil, err
	}

	// naive implementation
	recvBytes := len(buf) + len(payload)
	winUpdate := WindowUpdateFrame{WindowSizeIncrement: uint32(recvBytes)}
	if err := c.sendFrame(FrameTypeWindowUpdate, 0, frameHdr.StreamId, winUpdate.WritePayload(nil)); err != nil {
		log.Println("Cannot send WINDOW_UPDATE frame:", err.Error()) // or is it fatal?
	}

	return frame, nil
}

func (c *Conn) sendFrame(typ FrameType, flags FrameFlags, streamId uint32, payload []byte) error {
	var buf bytes.Buffer

	// ToDo: split large payload for several frames

	if err := c.buildFrame(&buf, typ, flags, streamId, payload); err != nil {
		return err
	}

	_, err := c.conn.Write(buf.Bytes())
	return err
}

func (c *Conn) sendFrameHeaders(streamId uint32, flags FrameFlags, frame *HeadersFrame) error {
	return c.sendFrame(FrameTypeHeaders, flags, streamId, frame.WritePayload(nil, c.hpackEncoder, &c.hpackEncoderBuf))
}

func (c *Conn) buildFrame(buffer *bytes.Buffer, typ FrameType, flags FrameFlags, streamId uint32, payload []byte) error {
	var buf [4]byte

	endianess.PutUint32(buf[0:4], uint32(len(payload))&0x00FFFFFF) // The size of the frame header is not included when describing frame sizes.
	buffer.Write(buf[1:4])

	buffer.Write([]byte{byte(typ), byte(flags)})
	endianess.PutUint32(buf[0:4], streamId&0x7FFFFFFF)
	buffer.Write(buf[0:4])

	if len(payload) > 0 {
		buffer.Write(payload)
	}

	return nil
}

// parseHTTP2SettingsHeader parses HTTP2-Settings header value which is just base64(SETTINGS frame payload)
func (c *Conn) parseHTTP2SettingsHeader(hdr string) (*SettingsFrame, error) {
	settingsFrame, err := base64.StdEncoding.DecodeString(hdr)
	if err != nil {
		return nil, fmt.Errorf("bad base64 value")
	}

	var frame SettingsFrame
	err = frame.ReadPayload(&FrameHdr{Length: uint32(len(settingsFrame))}, settingsFrame)
	if err != nil {
		return nil, fmt.Errorf("cannot parse value")
	}

	return &frame, nil
}

// RecvRequestHeaders receives only HEADER frames until gets FlagEndHeaders
func (c *Conn) RecvRequestHeaders(ctx context.Context) (*Request, error) {
	var request *Request

	for {
		f, err := c.RecvFrame(ctx)
		if err != nil {
			return nil, fmt.Errorf("recvFrame: %w", err)
		}

		ft, ok := f.(*HeadersFrame)
		if !ok {
			return nil, fmt.Errorf("expect HEADERS frame, got %T", f)
		}

		if request == nil {
			request = NewHTTPRequest(f.Hdr().StreamId)
		}

		for _, hf := range ft.HeaderBlockFragment {
			switch hf.Name {
			case ":method":
				request.Method = hf.Value
			case ":path":
				request.URI = hf.Value
			default:
				request.Headers[hf.Name] = hf.Value
			}
		}

		request.lastFrameFlags = f.Hdr().Flags
		if f.Hdr().Flags&FlagEndHeaders != 0 {
			break
		}
	}

	return request, nil
}

// RecvRequestBody receives only DATA frames until gets FlagEndStream
func (c *Conn) RecvRequestBody(ctx context.Context, request *Request) error {
	if request.lastFrameFlags&FlagEndStream != 0 {
		return nil
	}

	for {
		f, err := c.RecvFrame(ctx)
		if err != nil {
			return fmt.Errorf("recvFrame: %w", err)
		}

		ft, ok := f.(*DataFrame)
		if !ok {
			return fmt.Errorf("expect DATA frame, got %T", f)
		}

		request.Body = append(request.Body, ft.Data...)

		if request == nil {
			request = NewHTTPRequest(f.Hdr().StreamId)
		}

		request.lastFrameFlags = f.Hdr().Flags
		if f.Hdr().Flags&FlagEndStream != 0 {
			break
		}
	}

	return nil
}

func (c *Conn) SendTrailerFrame(streamId uint32, headers []hpack.HeaderField) error {
	return c.SendHeaders(streamId, FlagEndStream, headers)
}

func (c *Conn) SendHeaders(streamId uint32, flags FrameFlags, headers []hpack.HeaderField) error {
	c.writeLock.Lock()
	defer c.writeLock.Unlock()

	err := c.sendFrameHeaders(streamId, flags|FlagEndHeaders, &HeadersFrame{
		HeaderBlockFragment: headers,
	})

	if err != nil {
		return fmt.Errorf("send headers: %w", err)
	}

	return nil
}

func (c *Conn) SendData(streamId uint32, flags FrameFlags, body []byte) error {
	c.writeLock.Lock()
	defer c.writeLock.Unlock()

	return c.sendFrame(FrameTypeData, flags, streamId, body)
}
