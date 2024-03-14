package http2

import "fmt"

type (
	ErrWrongHTTPRequest struct {
		Cause error
	}

	ErrProtocolError struct {
		Descr string
	}
)

var (
	ErrWrongRequestLine   = fmt.Errorf("wrong request line")
	ErrTooLongRequest     = fmt.Errorf("request is too long")
	ErrWrongUpgradeHeader = fmt.Errorf("wrong Upgrade header value")
)

func (e *ErrWrongHTTPRequest) Error() string {
	return "wrong HTTP request: " + e.Cause.Error()
}

var (
	_ error = &ErrWrongHTTPRequest{}
)

func NewErrWrongHTTPRequest(cause error) *ErrWrongHTTPRequest {
	return &ErrWrongHTTPRequest{Cause: cause}
}

func (e *ErrProtocolError) Http2Error() string {
	return `Protocol error`
}

func (e *ErrProtocolError) Error() string {
	return e.Http2Error() + `: ` + e.Descr
}
