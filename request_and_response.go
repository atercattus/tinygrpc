package http2

import (
	"bytes"
	"context"
	"fmt"
	"strings"
)

type (
	Request struct {
		isNonFirstHeader bool
		lastFrameFlags   FrameFlags

		Method      string
		URI         string
		HTTPVersion string

		StreamId uint32

		Headers map[string]string

		Body []byte
	}

	Response struct {
		Status  uint64
		Headers map[string]string
		Body    []byte
	}

	Handler       func(ctx context.Context, request *Request, response *Response) error
	StreamHandler func(ctx context.Context, request *Request, conn *Conn) error
)

func NewHTTPRequest(streamId uint32) *Request {
	var r Request
	r.StreamId = streamId
	r.Headers = make(map[string]string)
	return &r
}

func (r *Request) parseHTTPHeader(hdr []byte) error {
	if !r.isNonFirstHeader {
		r.isNonFirstHeader = true
		return r.parseHTTPRequestLine(hdr)
	}

	key, value := splitHeaderValue(hdr)
	if len(value) == 0 {
		return fmt.Errorf("wrong HTTP header line: %q", hdr)
	}
	r.Headers[strings.ToLower(string(key))] = string(value)

	return nil
}

func (r *Request) parseHTTPRequestLine(line []byte) error {
	p := bytes.IndexByte(line, httpSP)
	if p < 1 {
		return NewErrWrongHTTPRequest(ErrWrongRequestLine)
	}

	r.Method = string(line[:p])
	line = line[p+1:]

	p = bytes.IndexByte(line, httpSP)
	if p < 1 {
		return NewErrWrongHTTPRequest(ErrWrongRequestLine)
	}

	r.URI = string(line[:p])

	r.HTTPVersion = string(line[p+1:])

	return nil
}

func splitHeaderValue(hdr []byte) ([]byte, []byte) {
	p := bytes.IndexByte(hdr, ':')
	if p < 1 {
		return nil, nil
	}
	f := p + 1

	t := len(hdr)

	for (f < t) && (hdr[f] == httpSP) {
		f++
	}

	t--
	for (f < t) && (hdr[t] == httpSP) {
		t--
	}

	return hdr[:p], hdr[f : t+1]
}
