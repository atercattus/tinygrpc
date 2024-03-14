package http2

import "fmt"

type (
	FrameType uint8
)

const (
	FrameTypeData         = FrameType(0)
	FrameTypeHeaders      = FrameType(1)
	FrameTypePriority     = FrameType(2)
	FrameTypeRstStream    = FrameType(3)
	FrameTypeSettings     = FrameType(4)
	FrameTypePushPromise  = FrameType(5)
	FrameTypePing         = FrameType(6)
	FrameTypeGoaway       = FrameType(7)
	FrameTypeWindowUpdate = FrameType(8)
	FrameTypeContinuation = FrameType(9)
)

func (f FrameType) String() string {
	switch f {
	case FrameTypeData:
		return `DATA`
	case FrameTypeHeaders:
		return `HEADERS`
	case FrameTypePriority:
		return `PRIORITY`
	case FrameTypeRstStream:
		return `RST_STREAM`
	case FrameTypeSettings:
		return `SETTINGS`
	case FrameTypePushPromise:
		return `PUSH_PROMISE`
	case FrameTypePing:
		return `PING`
	case FrameTypeGoaway:
		return `GOAWAY`
	case FrameTypeWindowUpdate:
		return `WINDOW_UPDATE`
	case FrameTypeContinuation:
		return `CONTINUATION`
	default:
		return fmt.Sprintf(`FrameType(%d)`, f)
	}
}
