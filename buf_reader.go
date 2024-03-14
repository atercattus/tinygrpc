package http2

import (
	"net"
	"time"
)

type (
	// BufReader is a wrapper around net.Conn
	// which for reading first reads from the internal buffer (prepend),
	// and then from the connection (src)
	BufReader struct {
		src net.Conn

		prepend []byte
		pos     int
	}
)

var (
	_ net.Conn = &BufReader{}
)

func NewBufReader(src net.Conn, prepend []byte) *BufReader {
	if len(prepend) > 0 {
		prepend = append([]byte{}, prepend...)
	}

	return &BufReader{
		src:     src,
		prepend: prepend,
	}
}

func (b *BufReader) Read(p []byte) (n int, err error) {
	if len(b.prepend) == 0 {
		return b.src.Read(p)
	}

	return b.readFromBuf(p)
}

func (b *BufReader) readFromBuf(p []byte) (n int, err error) {
	n = copy(p, b.prepend[b.pos:])
	b.pos += n

	if b.pos == len(b.prepend) {
		b.prepend = b.prepend[:0]
		b.pos = 0
	}

	return n, nil
}

func (b *BufReader) Write(p []byte) (n int, err error) {
	return b.src.Write(p)
}

func (b *BufReader) Close() error {
	return b.src.Close()
}

func (b *BufReader) LocalAddr() net.Addr {
	return b.src.LocalAddr()
}

func (b *BufReader) RemoteAddr() net.Addr {
	return b.src.RemoteAddr()
}

func (b *BufReader) SetDeadline(t time.Time) error {
	return b.src.SetDeadline(t)
}

func (b *BufReader) SetReadDeadline(t time.Time) error {
	return b.src.SetReadDeadline(t)
}

func (b *BufReader) SetWriteDeadline(t time.Time) error {
	return b.src.SetWriteDeadline(t)
}
