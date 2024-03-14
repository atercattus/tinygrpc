# tinygrpc

A tiny implementation of HTTP2 and gRPC server, written for study purposes, **not for practical use** at this moment.

This implementation supports:
- HTTP/1.1 protocol upgrade or prior-knowledge (h2c handshake)
- minimal implementation for HTTP2 flow including multiplexing, WINDOW_UPDATE and PING frames
- gRPC unary calls
- gRPC streaming implementation with minimal support for google.golang.org/grpc.ServiceDesc

This implementation does NOT support for now:
- ALPN TLS h2 handshake
- client settings accounting (HTTP2-Settings header, SETTINGS and WINDOW_UPDATE frames)
- GOAWAY, RST_STREAM frames
- splitting responses into parts (multiple DATA and HEADER frames)
- assembling multiplexed requests from parts for incoming gRPC calls
- correct return of non-fatal errors as statuses
- initial HTTP/1 request greater than RecvBufferSize (16KB)

The main functionality tests are in the file [conn_test.go](conn_test.go).


Feel free to run `$ go test -count 1 -v ./...` ;-)
