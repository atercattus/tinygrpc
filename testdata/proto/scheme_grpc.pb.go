// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.21.12
// source: scheme.proto

package fastGrpcTestPb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	FastGrpcTestService_Ping_FullMethodName  = "/fastgrpctest.v1.FastGrpcTestService/Ping"
	FastGrpcTestService_Sleep_FullMethodName = "/fastgrpctest.v1.FastGrpcTestService/Sleep"
	FastGrpcTestService_Sub_FullMethodName   = "/fastgrpctest.v1.FastGrpcTestService/Sub"
)

// FastGrpcTestServiceClient is the client API for FastGrpcTestService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type FastGrpcTestServiceClient interface {
	Ping(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*PingResponse, error)
	Sleep(ctx context.Context, in *SleepRequest, opts ...grpc.CallOption) (*SleepResponse, error)
	Sub(ctx context.Context, opts ...grpc.CallOption) (FastGrpcTestService_SubClient, error)
}

type fastGrpcTestServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewFastGrpcTestServiceClient(cc grpc.ClientConnInterface) FastGrpcTestServiceClient {
	return &fastGrpcTestServiceClient{cc}
}

func (c *fastGrpcTestServiceClient) Ping(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*PingResponse, error) {
	out := new(PingResponse)
	err := c.cc.Invoke(ctx, FastGrpcTestService_Ping_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *fastGrpcTestServiceClient) Sleep(ctx context.Context, in *SleepRequest, opts ...grpc.CallOption) (*SleepResponse, error) {
	out := new(SleepResponse)
	err := c.cc.Invoke(ctx, FastGrpcTestService_Sleep_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *fastGrpcTestServiceClient) Sub(ctx context.Context, opts ...grpc.CallOption) (FastGrpcTestService_SubClient, error) {
	stream, err := c.cc.NewStream(ctx, &FastGrpcTestService_ServiceDesc.Streams[0], FastGrpcTestService_Sub_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &fastGrpcTestServiceSubClient{stream}
	return x, nil
}

type FastGrpcTestService_SubClient interface {
	Send(*SubRequest) error
	Recv() (*SubResponse, error)
	grpc.ClientStream
}

type fastGrpcTestServiceSubClient struct {
	grpc.ClientStream
}

func (x *fastGrpcTestServiceSubClient) Send(m *SubRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *fastGrpcTestServiceSubClient) Recv() (*SubResponse, error) {
	m := new(SubResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// FastGrpcTestServiceServer is the server API for FastGrpcTestService service.
// All implementations must embed UnimplementedFastGrpcTestServiceServer
// for forward compatibility
type FastGrpcTestServiceServer interface {
	Ping(context.Context, *PingRequest) (*PingResponse, error)
	Sleep(context.Context, *SleepRequest) (*SleepResponse, error)
	Sub(FastGrpcTestService_SubServer) error
	mustEmbedUnimplementedFastGrpcTestServiceServer()
}

// UnimplementedFastGrpcTestServiceServer must be embedded to have forward compatible implementations.
type UnimplementedFastGrpcTestServiceServer struct {
}

func (UnimplementedFastGrpcTestServiceServer) Ping(context.Context, *PingRequest) (*PingResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ping not implemented")
}
func (UnimplementedFastGrpcTestServiceServer) Sleep(context.Context, *SleepRequest) (*SleepResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Sleep not implemented")
}
func (UnimplementedFastGrpcTestServiceServer) Sub(FastGrpcTestService_SubServer) error {
	return status.Errorf(codes.Unimplemented, "method Sub not implemented")
}
func (UnimplementedFastGrpcTestServiceServer) mustEmbedUnimplementedFastGrpcTestServiceServer() {}

// UnsafeFastGrpcTestServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to FastGrpcTestServiceServer will
// result in compilation errors.
type UnsafeFastGrpcTestServiceServer interface {
	mustEmbedUnimplementedFastGrpcTestServiceServer()
}

func RegisterFastGrpcTestServiceServer(s grpc.ServiceRegistrar, srv FastGrpcTestServiceServer) {
	s.RegisterService(&FastGrpcTestService_ServiceDesc, srv)
}

func _FastGrpcTestService_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FastGrpcTestServiceServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: FastGrpcTestService_Ping_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FastGrpcTestServiceServer).Ping(ctx, req.(*PingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _FastGrpcTestService_Sleep_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SleepRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FastGrpcTestServiceServer).Sleep(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: FastGrpcTestService_Sleep_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FastGrpcTestServiceServer).Sleep(ctx, req.(*SleepRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _FastGrpcTestService_Sub_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(FastGrpcTestServiceServer).Sub(&fastGrpcTestServiceSubServer{stream})
}

type FastGrpcTestService_SubServer interface {
	Send(*SubResponse) error
	Recv() (*SubRequest, error)
	grpc.ServerStream
}

type fastGrpcTestServiceSubServer struct {
	grpc.ServerStream
}

func (x *fastGrpcTestServiceSubServer) Send(m *SubResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *fastGrpcTestServiceSubServer) Recv() (*SubRequest, error) {
	m := new(SubRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// FastGrpcTestService_ServiceDesc is the grpc.ServiceDesc for FastGrpcTestService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var FastGrpcTestService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "fastgrpctest.v1.FastGrpcTestService",
	HandlerType: (*FastGrpcTestServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Ping",
			Handler:    _FastGrpcTestService_Ping_Handler,
		},
		{
			MethodName: "Sleep",
			Handler:    _FastGrpcTestService_Sleep_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Sub",
			Handler:       _FastGrpcTestService_Sub_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "scheme.proto",
}
