// Code generated by protoc-gen-go. DO NOT EDIT.
// source: registry/runtime.proto

package registry

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	common "github.com/oasislabs/ekiden/go/grpc/common"
	grpc "google.golang.org/grpc"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type Runtime struct {
	Id                       []byte                          `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Code                     []byte                          `protobuf:"bytes,2,opt,name=code,proto3" json:"code,omitempty"`
	TeeHardware              common.CapabilitiesTEE_Hardware `protobuf:"varint,3,opt,name=tee_hardware,json=teeHardware,proto3,enum=common.CapabilitiesTEE_Hardware" json:"tee_hardware,omitempty"`
	ReplicaGroupSize         uint64                          `protobuf:"varint,4,opt,name=replica_group_size,json=replicaGroupSize,proto3" json:"replica_group_size,omitempty"`
	StorageGroupSize         uint64                          `protobuf:"varint,5,opt,name=storage_group_size,json=storageGroupSize,proto3" json:"storage_group_size,omitempty"`
	ReplicaGroupBackupSize   uint64                          `protobuf:"varint,6,opt,name=replica_group_backup_size,json=replicaGroupBackupSize,proto3" json:"replica_group_backup_size,omitempty"`
	ReplicaAllowedStragglers uint64                          `protobuf:"varint,7,opt,name=replica_allowed_stragglers,json=replicaAllowedStragglers,proto3" json:"replica_allowed_stragglers,omitempty"`
	RegistrationTime         uint64                          `protobuf:"varint,8,opt,name=registration_time,json=registrationTime,proto3" json:"registration_time,omitempty"`
	XXX_NoUnkeyedLiteral     struct{}                        `json:"-"`
	XXX_unrecognized         []byte                          `json:"-"`
	XXX_sizecache            int32                           `json:"-"`
}

func (m *Runtime) Reset()         { *m = Runtime{} }
func (m *Runtime) String() string { return proto.CompactTextString(m) }
func (*Runtime) ProtoMessage()    {}
func (*Runtime) Descriptor() ([]byte, []int) {
	return fileDescriptor_abec0a28ef4e4ff8, []int{0}
}

func (m *Runtime) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Runtime.Unmarshal(m, b)
}
func (m *Runtime) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Runtime.Marshal(b, m, deterministic)
}
func (m *Runtime) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Runtime.Merge(m, src)
}
func (m *Runtime) XXX_Size() int {
	return xxx_messageInfo_Runtime.Size(m)
}
func (m *Runtime) XXX_DiscardUnknown() {
	xxx_messageInfo_Runtime.DiscardUnknown(m)
}

var xxx_messageInfo_Runtime proto.InternalMessageInfo

func (m *Runtime) GetId() []byte {
	if m != nil {
		return m.Id
	}
	return nil
}

func (m *Runtime) GetCode() []byte {
	if m != nil {
		return m.Code
	}
	return nil
}

func (m *Runtime) GetTeeHardware() common.CapabilitiesTEE_Hardware {
	if m != nil {
		return m.TeeHardware
	}
	return common.CapabilitiesTEE_Invalid
}

func (m *Runtime) GetReplicaGroupSize() uint64 {
	if m != nil {
		return m.ReplicaGroupSize
	}
	return 0
}

func (m *Runtime) GetStorageGroupSize() uint64 {
	if m != nil {
		return m.StorageGroupSize
	}
	return 0
}

func (m *Runtime) GetReplicaGroupBackupSize() uint64 {
	if m != nil {
		return m.ReplicaGroupBackupSize
	}
	return 0
}

func (m *Runtime) GetReplicaAllowedStragglers() uint64 {
	if m != nil {
		return m.ReplicaAllowedStragglers
	}
	return 0
}

func (m *Runtime) GetRegistrationTime() uint64 {
	if m != nil {
		return m.RegistrationTime
	}
	return 0
}

type RegisterRuntimeRequest struct {
	// Signed blob should be a CBOR-serialized runtime.
	Runtime              *common.Signed `protobuf:"bytes,1,opt,name=runtime,proto3" json:"runtime,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *RegisterRuntimeRequest) Reset()         { *m = RegisterRuntimeRequest{} }
func (m *RegisterRuntimeRequest) String() string { return proto.CompactTextString(m) }
func (*RegisterRuntimeRequest) ProtoMessage()    {}
func (*RegisterRuntimeRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_abec0a28ef4e4ff8, []int{1}
}

func (m *RegisterRuntimeRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RegisterRuntimeRequest.Unmarshal(m, b)
}
func (m *RegisterRuntimeRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RegisterRuntimeRequest.Marshal(b, m, deterministic)
}
func (m *RegisterRuntimeRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RegisterRuntimeRequest.Merge(m, src)
}
func (m *RegisterRuntimeRequest) XXX_Size() int {
	return xxx_messageInfo_RegisterRuntimeRequest.Size(m)
}
func (m *RegisterRuntimeRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_RegisterRuntimeRequest.DiscardUnknown(m)
}

var xxx_messageInfo_RegisterRuntimeRequest proto.InternalMessageInfo

func (m *RegisterRuntimeRequest) GetRuntime() *common.Signed {
	if m != nil {
		return m.Runtime
	}
	return nil
}

type RegisterRuntimeResponse struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RegisterRuntimeResponse) Reset()         { *m = RegisterRuntimeResponse{} }
func (m *RegisterRuntimeResponse) String() string { return proto.CompactTextString(m) }
func (*RegisterRuntimeResponse) ProtoMessage()    {}
func (*RegisterRuntimeResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_abec0a28ef4e4ff8, []int{2}
}

func (m *RegisterRuntimeResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RegisterRuntimeResponse.Unmarshal(m, b)
}
func (m *RegisterRuntimeResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RegisterRuntimeResponse.Marshal(b, m, deterministic)
}
func (m *RegisterRuntimeResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RegisterRuntimeResponse.Merge(m, src)
}
func (m *RegisterRuntimeResponse) XXX_Size() int {
	return xxx_messageInfo_RegisterRuntimeResponse.Size(m)
}
func (m *RegisterRuntimeResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_RegisterRuntimeResponse.DiscardUnknown(m)
}

var xxx_messageInfo_RegisterRuntimeResponse proto.InternalMessageInfo

type RuntimeRequest struct {
	Id                   []byte   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RuntimeRequest) Reset()         { *m = RuntimeRequest{} }
func (m *RuntimeRequest) String() string { return proto.CompactTextString(m) }
func (*RuntimeRequest) ProtoMessage()    {}
func (*RuntimeRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_abec0a28ef4e4ff8, []int{3}
}

func (m *RuntimeRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RuntimeRequest.Unmarshal(m, b)
}
func (m *RuntimeRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RuntimeRequest.Marshal(b, m, deterministic)
}
func (m *RuntimeRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RuntimeRequest.Merge(m, src)
}
func (m *RuntimeRequest) XXX_Size() int {
	return xxx_messageInfo_RuntimeRequest.Size(m)
}
func (m *RuntimeRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_RuntimeRequest.DiscardUnknown(m)
}

var xxx_messageInfo_RuntimeRequest proto.InternalMessageInfo

func (m *RuntimeRequest) GetId() []byte {
	if m != nil {
		return m.Id
	}
	return nil
}

type RuntimeResponse struct {
	Runtime              *Runtime `protobuf:"bytes,1,opt,name=runtime,proto3" json:"runtime,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RuntimeResponse) Reset()         { *m = RuntimeResponse{} }
func (m *RuntimeResponse) String() string { return proto.CompactTextString(m) }
func (*RuntimeResponse) ProtoMessage()    {}
func (*RuntimeResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_abec0a28ef4e4ff8, []int{4}
}

func (m *RuntimeResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RuntimeResponse.Unmarshal(m, b)
}
func (m *RuntimeResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RuntimeResponse.Marshal(b, m, deterministic)
}
func (m *RuntimeResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RuntimeResponse.Merge(m, src)
}
func (m *RuntimeResponse) XXX_Size() int {
	return xxx_messageInfo_RuntimeResponse.Size(m)
}
func (m *RuntimeResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_RuntimeResponse.DiscardUnknown(m)
}

var xxx_messageInfo_RuntimeResponse proto.InternalMessageInfo

func (m *RuntimeResponse) GetRuntime() *Runtime {
	if m != nil {
		return m.Runtime
	}
	return nil
}

type RuntimesRequest struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RuntimesRequest) Reset()         { *m = RuntimesRequest{} }
func (m *RuntimesRequest) String() string { return proto.CompactTextString(m) }
func (*RuntimesRequest) ProtoMessage()    {}
func (*RuntimesRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_abec0a28ef4e4ff8, []int{5}
}

func (m *RuntimesRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RuntimesRequest.Unmarshal(m, b)
}
func (m *RuntimesRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RuntimesRequest.Marshal(b, m, deterministic)
}
func (m *RuntimesRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RuntimesRequest.Merge(m, src)
}
func (m *RuntimesRequest) XXX_Size() int {
	return xxx_messageInfo_RuntimesRequest.Size(m)
}
func (m *RuntimesRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_RuntimesRequest.DiscardUnknown(m)
}

var xxx_messageInfo_RuntimesRequest proto.InternalMessageInfo

type RuntimesResponse struct {
	Runtime              []*Runtime `protobuf:"bytes,1,rep,name=runtime,proto3" json:"runtime,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *RuntimesResponse) Reset()         { *m = RuntimesResponse{} }
func (m *RuntimesResponse) String() string { return proto.CompactTextString(m) }
func (*RuntimesResponse) ProtoMessage()    {}
func (*RuntimesResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_abec0a28ef4e4ff8, []int{6}
}

func (m *RuntimesResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RuntimesResponse.Unmarshal(m, b)
}
func (m *RuntimesResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RuntimesResponse.Marshal(b, m, deterministic)
}
func (m *RuntimesResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RuntimesResponse.Merge(m, src)
}
func (m *RuntimesResponse) XXX_Size() int {
	return xxx_messageInfo_RuntimesResponse.Size(m)
}
func (m *RuntimesResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_RuntimesResponse.DiscardUnknown(m)
}

var xxx_messageInfo_RuntimesResponse proto.InternalMessageInfo

func (m *RuntimesResponse) GetRuntime() []*Runtime {
	if m != nil {
		return m.Runtime
	}
	return nil
}

type WatchRuntimesRequest struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *WatchRuntimesRequest) Reset()         { *m = WatchRuntimesRequest{} }
func (m *WatchRuntimesRequest) String() string { return proto.CompactTextString(m) }
func (*WatchRuntimesRequest) ProtoMessage()    {}
func (*WatchRuntimesRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_abec0a28ef4e4ff8, []int{7}
}

func (m *WatchRuntimesRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_WatchRuntimesRequest.Unmarshal(m, b)
}
func (m *WatchRuntimesRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_WatchRuntimesRequest.Marshal(b, m, deterministic)
}
func (m *WatchRuntimesRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_WatchRuntimesRequest.Merge(m, src)
}
func (m *WatchRuntimesRequest) XXX_Size() int {
	return xxx_messageInfo_WatchRuntimesRequest.Size(m)
}
func (m *WatchRuntimesRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_WatchRuntimesRequest.DiscardUnknown(m)
}

var xxx_messageInfo_WatchRuntimesRequest proto.InternalMessageInfo

type WatchRuntimesResponse struct {
	Runtime              *Runtime `protobuf:"bytes,1,opt,name=runtime,proto3" json:"runtime,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *WatchRuntimesResponse) Reset()         { *m = WatchRuntimesResponse{} }
func (m *WatchRuntimesResponse) String() string { return proto.CompactTextString(m) }
func (*WatchRuntimesResponse) ProtoMessage()    {}
func (*WatchRuntimesResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_abec0a28ef4e4ff8, []int{8}
}

func (m *WatchRuntimesResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_WatchRuntimesResponse.Unmarshal(m, b)
}
func (m *WatchRuntimesResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_WatchRuntimesResponse.Marshal(b, m, deterministic)
}
func (m *WatchRuntimesResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_WatchRuntimesResponse.Merge(m, src)
}
func (m *WatchRuntimesResponse) XXX_Size() int {
	return xxx_messageInfo_WatchRuntimesResponse.Size(m)
}
func (m *WatchRuntimesResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_WatchRuntimesResponse.DiscardUnknown(m)
}

var xxx_messageInfo_WatchRuntimesResponse proto.InternalMessageInfo

func (m *WatchRuntimesResponse) GetRuntime() *Runtime {
	if m != nil {
		return m.Runtime
	}
	return nil
}

func init() {
	proto.RegisterType((*Runtime)(nil), "registry.Runtime")
	proto.RegisterType((*RegisterRuntimeRequest)(nil), "registry.RegisterRuntimeRequest")
	proto.RegisterType((*RegisterRuntimeResponse)(nil), "registry.RegisterRuntimeResponse")
	proto.RegisterType((*RuntimeRequest)(nil), "registry.RuntimeRequest")
	proto.RegisterType((*RuntimeResponse)(nil), "registry.RuntimeResponse")
	proto.RegisterType((*RuntimesRequest)(nil), "registry.RuntimesRequest")
	proto.RegisterType((*RuntimesResponse)(nil), "registry.RuntimesResponse")
	proto.RegisterType((*WatchRuntimesRequest)(nil), "registry.WatchRuntimesRequest")
	proto.RegisterType((*WatchRuntimesResponse)(nil), "registry.WatchRuntimesResponse")
}

func init() { proto.RegisterFile("registry/runtime.proto", fileDescriptor_abec0a28ef4e4ff8) }

var fileDescriptor_abec0a28ef4e4ff8 = []byte{
	// 501 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x94, 0x51, 0x8f, 0xd2, 0x40,
	0x14, 0x85, 0x85, 0xc5, 0x65, 0x73, 0x59, 0xd9, 0x65, 0x54, 0x2c, 0x7d, 0xd0, 0xda, 0x27, 0x92,
	0xdd, 0xb4, 0x06, 0x9f, 0x4c, 0x8c, 0x46, 0x70, 0x5d, 0x9f, 0x0b, 0x89, 0xc6, 0x17, 0x32, 0xb4,
	0x37, 0x65, 0xb2, 0xa5, 0x53, 0x67, 0x86, 0x6c, 0xdc, 0x1f, 0xe9, 0xaf, 0xf1, 0x07, 0x18, 0xa6,
	0x33, 0xc0, 0x96, 0xaa, 0x89, 0x4f, 0xc0, 0x3d, 0xe7, 0x7c, 0xd3, 0x9e, 0x3b, 0x01, 0xfa, 0x02,
	0x53, 0x26, 0x95, 0xf8, 0x11, 0x8a, 0x75, 0xae, 0xd8, 0x0a, 0x83, 0x42, 0x70, 0xc5, 0xc9, 0x89,
	0x9d, 0xbb, 0x8f, 0x63, 0xbe, 0x5a, 0xf1, 0x3c, 0x2c, 0x3f, 0x4a, 0xd9, 0xff, 0xd5, 0x84, 0x76,
	0x54, 0x06, 0x48, 0x17, 0x9a, 0x2c, 0x71, 0x1a, 0x5e, 0x63, 0x78, 0x1a, 0x35, 0x59, 0x42, 0x08,
	0xb4, 0x62, 0x9e, 0xa0, 0xd3, 0xd4, 0x13, 0xfd, 0x9d, 0x4c, 0xe0, 0x54, 0x21, 0xce, 0x97, 0x54,
	0x24, 0xb7, 0x54, 0xa0, 0x73, 0xe4, 0x35, 0x86, 0xdd, 0x91, 0x17, 0x18, 0xe8, 0x84, 0x16, 0x74,
	0xc1, 0x32, 0xa6, 0x18, 0xca, 0xd9, 0xd5, 0x55, 0xf0, 0xd9, 0xf8, 0xa2, 0x8e, 0x42, 0xb4, 0x3f,
	0xc8, 0x25, 0x10, 0x81, 0x45, 0xc6, 0x62, 0x3a, 0x4f, 0x05, 0x5f, 0x17, 0x73, 0xc9, 0xee, 0xd0,
	0x69, 0x79, 0x8d, 0x61, 0x2b, 0x3a, 0x37, 0xca, 0xf5, 0x46, 0x98, 0xb2, 0x3b, 0xed, 0x96, 0x8a,
	0x0b, 0x9a, 0xe2, 0xbe, 0xfb, 0x61, 0xe9, 0x36, 0xca, 0xce, 0xfd, 0x06, 0x06, 0xf7, 0xd9, 0x0b,
	0x1a, 0xdf, 0xd8, 0xd0, 0xb1, 0x0e, 0xf5, 0xf7, 0x8f, 0x18, 0x6b, 0x59, 0x47, 0xdf, 0x82, 0x6b,
	0xa3, 0x34, 0xcb, 0xf8, 0x2d, 0x26, 0x73, 0xa9, 0x04, 0x4d, 0xd3, 0x0c, 0x85, 0x74, 0xda, 0x3a,
	0xeb, 0x18, 0xc7, 0x87, 0xd2, 0x30, 0xdd, 0xea, 0xe4, 0x02, 0x7a, 0xa6, 0x6a, 0xaa, 0x18, 0xcf,
	0xe7, 0x9b, 0x4a, 0x9d, 0x13, 0xfb, 0x4e, 0x3b, 0x61, 0xc6, 0x56, 0xe8, 0x8f, 0xa1, 0x1f, 0xe9,
	0x19, 0x0a, 0xd3, 0x7e, 0x84, 0xdf, 0xd7, 0x28, 0x15, 0x19, 0x42, 0xdb, 0x2c, 0x50, 0x6f, 0xa2,
	0x33, 0xea, 0xda, 0x6e, 0xa7, 0x2c, 0xcd, 0x31, 0x89, 0xac, 0xec, 0x0f, 0xe0, 0xd9, 0x01, 0x43,
	0x16, 0x3c, 0x97, 0xe8, 0x7b, 0xd0, 0xad, 0x60, 0x2b, 0xbb, 0xf5, 0xdf, 0xc1, 0x59, 0x25, 0x44,
	0x2e, 0xaa, 0x27, 0xf7, 0x02, 0x7b, 0x77, 0x02, 0xeb, 0xdd, 0x1e, 0xde, 0xdb, 0xe6, 0xa5, 0x39,
	0xc2, 0x7f, 0x0f, 0xe7, 0xbb, 0x51, 0x1d, 0xf3, 0xe8, 0x1f, 0xcc, 0x3e, 0x3c, 0xf9, 0x42, 0x55,
	0xbc, 0xac, 0x82, 0x3f, 0xc2, 0xd3, 0xca, 0xfc, 0x3f, 0x9e, 0x78, 0xf4, 0xb3, 0xb9, 0xf7, 0xca,
	0xa5, 0x89, 0x7c, 0x85, 0xb3, 0x4a, 0x85, 0xc4, 0xdb, 0x43, 0xd4, 0x6e, 0xc8, 0x7d, 0xf9, 0x17,
	0x87, 0xe9, 0xff, 0x01, 0x99, 0x00, 0x5c, 0xa3, 0xb2, 0x50, 0xe7, 0xf0, 0xb9, 0x0c, 0x6c, 0x50,
	0xa3, 0x6c, 0x21, 0x9f, 0xa0, 0xb3, 0x83, 0x48, 0x72, 0xe8, 0xb5, 0x15, 0xb9, 0x6e, 0x9d, 0xb4,
	0xe5, 0xcc, 0xe0, 0xd1, 0xbd, 0x02, 0xc9, 0xf3, 0x9d, 0xbd, 0xae, 0x71, 0xf7, 0xc5, 0x1f, 0x75,
	0xcb, 0x7c, 0xd5, 0x18, 0x07, 0xdf, 0x2e, 0x53, 0xa6, 0x96, 0xeb, 0xc5, 0xe6, 0x82, 0x86, 0x9c,
	0x4a, 0x26, 0x33, 0xba, 0x90, 0x21, 0xde, 0xb0, 0x04, 0xf3, 0x30, 0xe5, 0x61, 0x2a, 0x8a, 0x38,
	0xb4, 0xa4, 0xc5, 0xb1, 0xfe, 0xc7, 0x79, 0xfd, 0x3b, 0x00, 0x00, 0xff, 0xff, 0xf3, 0xa7, 0x13,
	0x43, 0xaa, 0x04, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// RuntimeRegistryClient is the client API for RuntimeRegistry service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type RuntimeRegistryClient interface {
	RegisterRuntime(ctx context.Context, in *RegisterRuntimeRequest, opts ...grpc.CallOption) (*RegisterRuntimeResponse, error)
	GetRuntime(ctx context.Context, in *RuntimeRequest, opts ...grpc.CallOption) (*RuntimeResponse, error)
	GetRuntimes(ctx context.Context, in *RuntimesRequest, opts ...grpc.CallOption) (*RuntimesResponse, error)
	WatchRuntimes(ctx context.Context, in *WatchRuntimesRequest, opts ...grpc.CallOption) (RuntimeRegistry_WatchRuntimesClient, error)
}

type runtimeRegistryClient struct {
	cc *grpc.ClientConn
}

func NewRuntimeRegistryClient(cc *grpc.ClientConn) RuntimeRegistryClient {
	return &runtimeRegistryClient{cc}
}

func (c *runtimeRegistryClient) RegisterRuntime(ctx context.Context, in *RegisterRuntimeRequest, opts ...grpc.CallOption) (*RegisterRuntimeResponse, error) {
	out := new(RegisterRuntimeResponse)
	err := c.cc.Invoke(ctx, "/registry.RuntimeRegistry/RegisterRuntime", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *runtimeRegistryClient) GetRuntime(ctx context.Context, in *RuntimeRequest, opts ...grpc.CallOption) (*RuntimeResponse, error) {
	out := new(RuntimeResponse)
	err := c.cc.Invoke(ctx, "/registry.RuntimeRegistry/GetRuntime", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *runtimeRegistryClient) GetRuntimes(ctx context.Context, in *RuntimesRequest, opts ...grpc.CallOption) (*RuntimesResponse, error) {
	out := new(RuntimesResponse)
	err := c.cc.Invoke(ctx, "/registry.RuntimeRegistry/GetRuntimes", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *runtimeRegistryClient) WatchRuntimes(ctx context.Context, in *WatchRuntimesRequest, opts ...grpc.CallOption) (RuntimeRegistry_WatchRuntimesClient, error) {
	stream, err := c.cc.NewStream(ctx, &_RuntimeRegistry_serviceDesc.Streams[0], "/registry.RuntimeRegistry/WatchRuntimes", opts...)
	if err != nil {
		return nil, err
	}
	x := &runtimeRegistryWatchRuntimesClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type RuntimeRegistry_WatchRuntimesClient interface {
	Recv() (*WatchRuntimesResponse, error)
	grpc.ClientStream
}

type runtimeRegistryWatchRuntimesClient struct {
	grpc.ClientStream
}

func (x *runtimeRegistryWatchRuntimesClient) Recv() (*WatchRuntimesResponse, error) {
	m := new(WatchRuntimesResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// RuntimeRegistryServer is the server API for RuntimeRegistry service.
type RuntimeRegistryServer interface {
	RegisterRuntime(context.Context, *RegisterRuntimeRequest) (*RegisterRuntimeResponse, error)
	GetRuntime(context.Context, *RuntimeRequest) (*RuntimeResponse, error)
	GetRuntimes(context.Context, *RuntimesRequest) (*RuntimesResponse, error)
	WatchRuntimes(*WatchRuntimesRequest, RuntimeRegistry_WatchRuntimesServer) error
}

func RegisterRuntimeRegistryServer(s *grpc.Server, srv RuntimeRegistryServer) {
	s.RegisterService(&_RuntimeRegistry_serviceDesc, srv)
}

func _RuntimeRegistry_RegisterRuntime_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterRuntimeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RuntimeRegistryServer).RegisterRuntime(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/registry.RuntimeRegistry/RegisterRuntime",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RuntimeRegistryServer).RegisterRuntime(ctx, req.(*RegisterRuntimeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RuntimeRegistry_GetRuntime_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RuntimeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RuntimeRegistryServer).GetRuntime(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/registry.RuntimeRegistry/GetRuntime",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RuntimeRegistryServer).GetRuntime(ctx, req.(*RuntimeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RuntimeRegistry_GetRuntimes_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RuntimesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RuntimeRegistryServer).GetRuntimes(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/registry.RuntimeRegistry/GetRuntimes",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RuntimeRegistryServer).GetRuntimes(ctx, req.(*RuntimesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RuntimeRegistry_WatchRuntimes_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(WatchRuntimesRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(RuntimeRegistryServer).WatchRuntimes(m, &runtimeRegistryWatchRuntimesServer{stream})
}

type RuntimeRegistry_WatchRuntimesServer interface {
	Send(*WatchRuntimesResponse) error
	grpc.ServerStream
}

type runtimeRegistryWatchRuntimesServer struct {
	grpc.ServerStream
}

func (x *runtimeRegistryWatchRuntimesServer) Send(m *WatchRuntimesResponse) error {
	return x.ServerStream.SendMsg(m)
}

var _RuntimeRegistry_serviceDesc = grpc.ServiceDesc{
	ServiceName: "registry.RuntimeRegistry",
	HandlerType: (*RuntimeRegistryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RegisterRuntime",
			Handler:    _RuntimeRegistry_RegisterRuntime_Handler,
		},
		{
			MethodName: "GetRuntime",
			Handler:    _RuntimeRegistry_GetRuntime_Handler,
		},
		{
			MethodName: "GetRuntimes",
			Handler:    _RuntimeRegistry_GetRuntimes_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "WatchRuntimes",
			Handler:       _RuntimeRegistry_WatchRuntimes_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "registry/runtime.proto",
}