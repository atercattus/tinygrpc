package http2

import "fmt"

type (
	ErrorCode uint32
)

const (
	ErrorCodeNoError            = ErrorCode(0x0)
	ErrorCodeProtocolError      = ErrorCode(0x1)
	ErrorCodeInternalError      = ErrorCode(0x2)
	ErrorCodeFlowControlError   = ErrorCode(0x3)
	ErrorCodeSettingsTimeout    = ErrorCode(0x4)
	ErrorCodeStreamClosed       = ErrorCode(0x5)
	ErrorCodeFrameSizeError     = ErrorCode(0x6)
	ErrorCodeRefusedStream      = ErrorCode(0x7)
	ErrorCodeCancel             = ErrorCode(0x8)
	ErrorCodeCompressionError   = ErrorCode(0x9)
	ErrorCodeConnectError       = ErrorCode(0xa)
	ErrorCodeEnhanceYourCalm    = ErrorCode(0xb)
	ErrorCodeInadequateSecurity = ErrorCode(0xc)
	ErrorCodeHTTP11Required     = ErrorCode(0xd)
)

func (e ErrorCode) String() string {
	switch e {
	case ErrorCodeNoError:
		return `NO_ERROR`
	case ErrorCodeProtocolError:
		return `PROTOCOL_ERROR`
	case ErrorCodeInternalError:
		return `INTERNAL_ERROR`
	case ErrorCodeFlowControlError:
		return `FLOW_CONTROL_ERROR`
	case ErrorCodeSettingsTimeout:
		return `SETTINGS_TIMEOUT`
	case ErrorCodeStreamClosed:
		return `STREAM_CLOSED`
	case ErrorCodeFrameSizeError:
		return `FRAME_SIZE_ERROR`
	case ErrorCodeRefusedStream:
		return `REFUSED_STREAM`
	case ErrorCodeCancel:
		return `CANCEL`
	case ErrorCodeCompressionError:
		return `COMPRESSION_ERROR`
	case ErrorCodeConnectError:
		return `CONNECT_ERROR`
	case ErrorCodeEnhanceYourCalm:
		return `ENHANCE_YOUR_CALM`
	case ErrorCodeInadequateSecurity:
		return `INADEQUATE_SECURITY`
	case ErrorCodeHTTP11Required:
		return `HTTP_1_1_REQUIRED`
	default:
		return fmt.Sprintf(`ErrorCode(%d)`, e)
	}
}
