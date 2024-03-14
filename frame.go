package http2

import (
	"bytes"
	"fmt"

	"golang.org/x/net/http2/hpack"
)

type (
	Frame interface {
		Hdr() *FrameHdr

		// ReadPayload fills a frame from payload.
		//
		// payload MUST not be used after method call
		ReadPayload(frameHdr *FrameHdr, payload []byte) error
	}

	FrameHdr struct {
		Length   uint32
		Type     FrameType
		Flags    FrameFlags
		StreamId uint32
	}

	DataFrame struct { // page 33
		FrameHdr
		Data []byte
	}

	HeadersFrame struct { // page 35
		FrameHdr
		Exclusive           bool   // A single bit flag indicating that the stream dependency is exclusive. This field is only present if the PRIORITY flag is set.
		StreamIdDepend      uint32 // This field is only present if the PRIORITY flag is set.
		Weight              uint8  // This field is only present if the PRIORITY flag is set.
		HeaderBlockFragment []hpack.HeaderField

		hpackDecoder *hpack.Decoder
	}

	SettingsFrame struct { // page 39
		FrameHdr
		Params []SettingsFrameParam
	}

	SettingsFrameParam struct {
		Id    SettingsType
		Value uint32
	}

	PingFrame struct { // page 44
		FrameHdr
		Data uint64
	}

	GoawayFrame struct { // page 45
		FrameHdr
		LastStreamId uint32
		ErrorCode    ErrorCode
		DebugData    []byte
	}

	WindowUpdateFrame struct { // page 48
		FrameHdr
		WindowSizeIncrement uint32
	}

	RstStreamFrame struct {
		FrameHdr
		ErrorCode int32
	}
)

func (fh *FrameHdr) Hdr() *FrameHdr {
	return fh
}

func (f *SettingsFrame) ReadPayload(frameHdr *FrameHdr, payload []byte) error {
	if (frameHdr.Length % 6) != 0 {
		return &ErrProtocolError{Descr: `Frame settings payload length not multiply 6`}
	}

	for p := uint32(0); p < frameHdr.Length; p += 6 {
		id := SettingsType(endianess.Uint16(payload[p : p+2]))
		val := endianess.Uint32(payload[p+2 : p+6])
		f.Params = append(f.Params, SettingsFrameParam{Id: id, Value: val})
	}

	return nil
}

func (f *SettingsFrame) WritePayload(dst []byte) []byte {
	for _, p := range f.Params {
		dst = endianess.AppendUint16(dst, uint16(p.Id))
		dst = endianess.AppendUint32(dst, p.Value)
	}
	return dst
}

func (f *WindowUpdateFrame) ReadPayload(frameHdr *FrameHdr, payload []byte) error {
	if frameHdr.Length != 4 {
		return &ErrProtocolError{Descr: `Wrong window update frame payload length`}
	}

	f.WindowSizeIncrement = endianess.Uint32(payload[0:4]) & 0x7FFFFFFF

	return nil
}

func (f *WindowUpdateFrame) WritePayload(dst []byte) []byte {
	return endianess.AppendUint32(dst, f.WindowSizeIncrement)
}

func (f *HeadersFrame) ReadPayload(frameHdr *FrameHdr, payload []byte) error {
	if frameHdr.Flags&(FlagPadded|FlagPriority) != 0 {
		return fmt.Errorf("FlagPadded|FlagPriority are not implemented")
	}

	if frameHdr.Length > 0 {
		f.hpackDecoder.SetEmitFunc(func(hf hpack.HeaderField) {
			f.HeaderBlockFragment = append(f.HeaderBlockFragment, hf)
		})
		if _, err := f.hpackDecoder.Write(payload); err != nil {
			return &ErrProtocolError{Descr: `Cannot decode HPACK headers frame: ` + err.Error()}
		}
		if err := f.hpackDecoder.Close(); err != nil {
			return fmt.Errorf("hpack decoder close: %w", err)
		}
	}

	return nil
}

func (f *HeadersFrame) WritePayload(dst []byte, encoder *hpack.Encoder, encoderBuf *bytes.Buffer) []byte {
	encoderBuf.Reset()

	for _, hf := range f.HeaderBlockFragment {
		_ = encoder.WriteField(hpack.HeaderField{
			Name:      hf.Name,
			Value:     hf.Value,
			Sensitive: hf.Sensitive,
		})
	}

	return append(dst, encoderBuf.Bytes()...)
}

func (f *DataFrame) ReadPayload(frameHdr *FrameHdr, payload []byte) error {
	if frameHdr.Length != uint32(len(payload)) {
		return fmt.Errorf("data frame: expect %d bytes, got %d", frameHdr.Length, len(payload))
	}

	f.Data = append(f.Data[:0], payload...)

	return nil
}

func (f *DataFrame) WritePayload(dst []byte) []byte {
	return append(dst, f.Data...)
}

func (f *PingFrame) ReadPayload(frameHdr *FrameHdr, payload []byte) error {
	if frameHdr.Length != 8 {
		return &ErrProtocolError{Descr: `Wrong PING frame payload length`}
	}

	f.Data = endianess.Uint64(payload)

	return nil
}

func (f *PingFrame) WritePayload(dst []byte) []byte {
	return endianess.AppendUint64(dst, f.Data)
}

func (f *GoawayFrame) ReadPayload(frameHdr *FrameHdr, payload []byte) error {
	if frameHdr.Length < 8 {
		return &ErrProtocolError{Descr: `Wrong GOAWAY frame payload length`}
	}

	f.LastStreamId = endianess.Uint32(payload[0:4])
	f.ErrorCode = ErrorCode(endianess.Uint32(payload[4:8]))
	if frameHdr.Length > 8 {
		f.DebugData = append(f.DebugData, payload[8:]...)
	}

	return nil
}

func (f *GoawayFrame) WritePayload(dst []byte) []byte {
	dst = endianess.AppendUint32(dst, f.LastStreamId)
	dst = endianess.AppendUint32(dst, uint32(f.ErrorCode))
	return append(dst, f.DebugData...)
}

func (f *RstStreamFrame) ReadPayload(frameHdr *FrameHdr, payload []byte) error {
	if frameHdr.Length != 4 {
		return &ErrProtocolError{Descr: `Wrong RST_STREAM frame payload length`}
	}

	f.ErrorCode = int32(endianess.Uint32(payload[0:4]) & 0x7FFFFFFF)

	return nil
}

func (f *RstStreamFrame) WritePayload(dst []byte) []byte {
	return endianess.AppendUint32(dst, uint32(f.ErrorCode))
}
