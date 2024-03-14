// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.12
// source: scheme.proto

package fastGrpcTestPb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type PingRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Val int32 `protobuf:"varint,1,opt,name=val,proto3" json:"val,omitempty"`
}

func (x *PingRequest) Reset() {
	*x = PingRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_scheme_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PingRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PingRequest) ProtoMessage() {}

func (x *PingRequest) ProtoReflect() protoreflect.Message {
	mi := &file_scheme_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PingRequest.ProtoReflect.Descriptor instead.
func (*PingRequest) Descriptor() ([]byte, []int) {
	return file_scheme_proto_rawDescGZIP(), []int{0}
}

func (x *PingRequest) GetVal() int32 {
	if x != nil {
		return x.Val
	}
	return 0
}

type PingResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Val int32 `protobuf:"varint,1,opt,name=val,proto3" json:"val,omitempty"`
}

func (x *PingResponse) Reset() {
	*x = PingResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_scheme_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PingResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PingResponse) ProtoMessage() {}

func (x *PingResponse) ProtoReflect() protoreflect.Message {
	mi := &file_scheme_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PingResponse.ProtoReflect.Descriptor instead.
func (*PingResponse) Descriptor() ([]byte, []int) {
	return file_scheme_proto_rawDescGZIP(), []int{1}
}

func (x *PingResponse) GetVal() int32 {
	if x != nil {
		return x.Val
	}
	return 0
}

type SleepRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Duration int64 `protobuf:"varint,1,opt,name=duration,proto3" json:"duration,omitempty"`
}

func (x *SleepRequest) Reset() {
	*x = SleepRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_scheme_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SleepRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SleepRequest) ProtoMessage() {}

func (x *SleepRequest) ProtoReflect() protoreflect.Message {
	mi := &file_scheme_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SleepRequest.ProtoReflect.Descriptor instead.
func (*SleepRequest) Descriptor() ([]byte, []int) {
	return file_scheme_proto_rawDescGZIP(), []int{2}
}

func (x *SleepRequest) GetDuration() int64 {
	if x != nil {
		return x.Duration
	}
	return 0
}

type SleepResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Duration int64 `protobuf:"varint,1,opt,name=duration,proto3" json:"duration,omitempty"`
	Now      int64 `protobuf:"varint,2,opt,name=now,proto3" json:"now,omitempty"`
}

func (x *SleepResponse) Reset() {
	*x = SleepResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_scheme_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SleepResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SleepResponse) ProtoMessage() {}

func (x *SleepResponse) ProtoReflect() protoreflect.Message {
	mi := &file_scheme_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SleepResponse.ProtoReflect.Descriptor instead.
func (*SleepResponse) Descriptor() ([]byte, []int) {
	return file_scheme_proto_rawDescGZIP(), []int{3}
}

func (x *SleepResponse) GetDuration() int64 {
	if x != nil {
		return x.Duration
	}
	return 0
}

func (x *SleepResponse) GetNow() int64 {
	if x != nil {
		return x.Now
	}
	return 0
}

type SubRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PubInterval int64 `protobuf:"varint,1,opt,name=pubInterval,proto3" json:"pubInterval,omitempty"`
}

func (x *SubRequest) Reset() {
	*x = SubRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_scheme_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SubRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SubRequest) ProtoMessage() {}

func (x *SubRequest) ProtoReflect() protoreflect.Message {
	mi := &file_scheme_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SubRequest.ProtoReflect.Descriptor instead.
func (*SubRequest) Descriptor() ([]byte, []int) {
	return file_scheme_proto_rawDescGZIP(), []int{4}
}

func (x *SubRequest) GetPubInterval() int64 {
	if x != nil {
		return x.PubInterval
	}
	return 0
}

type SubResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CurrentTime int64 `protobuf:"varint,1,opt,name=currentTime,proto3" json:"currentTime,omitempty"`
}

func (x *SubResponse) Reset() {
	*x = SubResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_scheme_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SubResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SubResponse) ProtoMessage() {}

func (x *SubResponse) ProtoReflect() protoreflect.Message {
	mi := &file_scheme_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SubResponse.ProtoReflect.Descriptor instead.
func (*SubResponse) Descriptor() ([]byte, []int) {
	return file_scheme_proto_rawDescGZIP(), []int{5}
}

func (x *SubResponse) GetCurrentTime() int64 {
	if x != nil {
		return x.CurrentTime
	}
	return 0
}

var File_scheme_proto protoreflect.FileDescriptor

var file_scheme_proto_rawDesc = []byte{
	0x0a, 0x0c, 0x73, 0x63, 0x68, 0x65, 0x6d, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0f,
	0x66, 0x61, 0x73, 0x74, 0x67, 0x72, 0x70, 0x63, 0x74, 0x65, 0x73, 0x74, 0x2e, 0x76, 0x31, 0x22,
	0x1f, 0x0a, 0x0b, 0x50, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x10,
	0x0a, 0x03, 0x76, 0x61, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x03, 0x76, 0x61, 0x6c,
	0x22, 0x20, 0x0a, 0x0c, 0x50, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x10, 0x0a, 0x03, 0x76, 0x61, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x03, 0x76,
	0x61, 0x6c, 0x22, 0x2a, 0x0a, 0x0c, 0x53, 0x6c, 0x65, 0x65, 0x70, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x64, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x64, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0x3d,
	0x0a, 0x0d, 0x53, 0x6c, 0x65, 0x65, 0x70, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x1a, 0x0a, 0x08, 0x64, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x08, 0x64, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x10, 0x0a, 0x03, 0x6e,
	0x6f, 0x77, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x03, 0x6e, 0x6f, 0x77, 0x22, 0x2e, 0x0a,
	0x0a, 0x53, 0x75, 0x62, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x20, 0x0a, 0x0b, 0x70,
	0x75, 0x62, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x76, 0x61, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x0b, 0x70, 0x75, 0x62, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x76, 0x61, 0x6c, 0x22, 0x2f, 0x0a,
	0x0b, 0x53, 0x75, 0x62, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x20, 0x0a, 0x0b,
	0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x0b, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x32, 0xe8,
	0x01, 0x0a, 0x13, 0x46, 0x61, 0x73, 0x74, 0x47, 0x72, 0x70, 0x63, 0x54, 0x65, 0x73, 0x74, 0x53,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x43, 0x0a, 0x04, 0x50, 0x69, 0x6e, 0x67, 0x12, 0x1c,
	0x2e, 0x66, 0x61, 0x73, 0x74, 0x67, 0x72, 0x70, 0x63, 0x74, 0x65, 0x73, 0x74, 0x2e, 0x76, 0x31,
	0x2e, 0x50, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1d, 0x2e, 0x66,
	0x61, 0x73, 0x74, 0x67, 0x72, 0x70, 0x63, 0x74, 0x65, 0x73, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x50,
	0x69, 0x6e, 0x67, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x46, 0x0a, 0x05, 0x53,
	0x6c, 0x65, 0x65, 0x70, 0x12, 0x1d, 0x2e, 0x66, 0x61, 0x73, 0x74, 0x67, 0x72, 0x70, 0x63, 0x74,
	0x65, 0x73, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x6c, 0x65, 0x65, 0x70, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x1e, 0x2e, 0x66, 0x61, 0x73, 0x74, 0x67, 0x72, 0x70, 0x63, 0x74, 0x65,
	0x73, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x6c, 0x65, 0x65, 0x70, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x44, 0x0a, 0x03, 0x53, 0x75, 0x62, 0x12, 0x1b, 0x2e, 0x66, 0x61, 0x73,
	0x74, 0x67, 0x72, 0x70, 0x63, 0x74, 0x65, 0x73, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x75, 0x62,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1c, 0x2e, 0x66, 0x61, 0x73, 0x74, 0x67, 0x72,
	0x70, 0x63, 0x74, 0x65, 0x73, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x75, 0x62, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x28, 0x01, 0x30, 0x01, 0x42, 0x13, 0x5a, 0x11, 0x2e, 0x2f, 0x3b,
	0x66, 0x61, 0x73, 0x74, 0x47, 0x72, 0x70, 0x63, 0x54, 0x65, 0x73, 0x74, 0x50, 0x62, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_scheme_proto_rawDescOnce sync.Once
	file_scheme_proto_rawDescData = file_scheme_proto_rawDesc
)

func file_scheme_proto_rawDescGZIP() []byte {
	file_scheme_proto_rawDescOnce.Do(func() {
		file_scheme_proto_rawDescData = protoimpl.X.CompressGZIP(file_scheme_proto_rawDescData)
	})
	return file_scheme_proto_rawDescData
}

var file_scheme_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_scheme_proto_goTypes = []interface{}{
	(*PingRequest)(nil),   // 0: fastgrpctest.v1.PingRequest
	(*PingResponse)(nil),  // 1: fastgrpctest.v1.PingResponse
	(*SleepRequest)(nil),  // 2: fastgrpctest.v1.SleepRequest
	(*SleepResponse)(nil), // 3: fastgrpctest.v1.SleepResponse
	(*SubRequest)(nil),    // 4: fastgrpctest.v1.SubRequest
	(*SubResponse)(nil),   // 5: fastgrpctest.v1.SubResponse
}
var file_scheme_proto_depIdxs = []int32{
	0, // 0: fastgrpctest.v1.FastGrpcTestService.Ping:input_type -> fastgrpctest.v1.PingRequest
	2, // 1: fastgrpctest.v1.FastGrpcTestService.Sleep:input_type -> fastgrpctest.v1.SleepRequest
	4, // 2: fastgrpctest.v1.FastGrpcTestService.Sub:input_type -> fastgrpctest.v1.SubRequest
	1, // 3: fastgrpctest.v1.FastGrpcTestService.Ping:output_type -> fastgrpctest.v1.PingResponse
	3, // 4: fastgrpctest.v1.FastGrpcTestService.Sleep:output_type -> fastgrpctest.v1.SleepResponse
	5, // 5: fastgrpctest.v1.FastGrpcTestService.Sub:output_type -> fastgrpctest.v1.SubResponse
	3, // [3:6] is the sub-list for method output_type
	0, // [0:3] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_scheme_proto_init() }
func file_scheme_proto_init() {
	if File_scheme_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_scheme_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PingRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_scheme_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PingResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_scheme_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SleepRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_scheme_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SleepResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_scheme_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SubRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_scheme_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SubResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_scheme_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_scheme_proto_goTypes,
		DependencyIndexes: file_scheme_proto_depIdxs,
		MessageInfos:      file_scheme_proto_msgTypes,
	}.Build()
	File_scheme_proto = out.File
	file_scheme_proto_rawDesc = nil
	file_scheme_proto_goTypes = nil
	file_scheme_proto_depIdxs = nil
}
