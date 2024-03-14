package http2

import (
	"strconv"
	"strings"
)

type (
	FrameFlags uint8
)

const (
	FlagAck = FrameFlags(0x1)

	FlagEndStream  = FrameFlags(0x1)
	FlagEndHeaders = FrameFlags(0x4)
	FlagPadded     = FrameFlags(0x8)
	FlagPriority   = FrameFlags(0x20)
)

func (f FrameFlags) String() string {
	if f == 0 {
		return "FrameFlags(0)"
	}

	var sb strings.Builder

	for i := 0; i < 8; i++ {
		b := FrameFlags(1 << i)
		if f&b == 0 {
			continue
		}

		var name string

		switch b {
		case FlagEndStream:
			name = `EndStream/Ack`
		case FlagEndHeaders:
			name = `EndHeaders`
		case FlagPadded:
			name = `Padded`
		case FlagPriority:
			name = `Priority`
		default:
			name = strconv.Itoa(int(b))
		}

		sb.WriteString(name)
		sb.WriteByte('|')
	}

	return "FrameFlags(" + sb.String()[:sb.Len()-1] + ")"
}
