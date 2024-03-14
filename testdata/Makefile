.DEFAULT_GOAL := proto

.PHONY: protoc
protoc:
	@go install google.golang.org/protobuf/cmd/protoc-gen-go
	@go get google.golang.org/grpc/cmd/protoc-gen-go-grpc
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc

.PHONY: proto
proto: protoc scheme.proto
	@protoc \
		--proto_path=. \
		--go_out=proto         --plugin protoc-gen-go="${GOPATH}/bin/protoc-gen-go" \
		--go-grpc_out=proto    --plugin protoc-gen-go-grpc="${GOPATH}/bin/protoc-gen-go-grpc" \
		scheme.proto
	@go vet ./...
