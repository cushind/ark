// Code generated by protoc-gen-go. DO NOT EDIT.
// source: RestoreItemAction.proto

package generated

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type RestoreExecuteRequest struct {
	Item    []byte `protobuf:"bytes,1,opt,name=item,proto3" json:"item,omitempty"`
	Restore []byte `protobuf:"bytes,2,opt,name=restore,proto3" json:"restore,omitempty"`
}

func (m *RestoreExecuteRequest) Reset()                    { *m = RestoreExecuteRequest{} }
func (m *RestoreExecuteRequest) String() string            { return proto.CompactTextString(m) }
func (*RestoreExecuteRequest) ProtoMessage()               {}
func (*RestoreExecuteRequest) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{0} }

func (m *RestoreExecuteRequest) GetItem() []byte {
	if m != nil {
		return m.Item
	}
	return nil
}

func (m *RestoreExecuteRequest) GetRestore() []byte {
	if m != nil {
		return m.Restore
	}
	return nil
}

type RestoreExecuteResponse struct {
	Item    []byte `protobuf:"bytes,1,opt,name=item,proto3" json:"item,omitempty"`
	Warning string `protobuf:"bytes,2,opt,name=warning" json:"warning,omitempty"`
}

func (m *RestoreExecuteResponse) Reset()                    { *m = RestoreExecuteResponse{} }
func (m *RestoreExecuteResponse) String() string            { return proto.CompactTextString(m) }
func (*RestoreExecuteResponse) ProtoMessage()               {}
func (*RestoreExecuteResponse) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{1} }

func (m *RestoreExecuteResponse) GetItem() []byte {
	if m != nil {
		return m.Item
	}
	return nil
}

func (m *RestoreExecuteResponse) GetWarning() string {
	if m != nil {
		return m.Warning
	}
	return ""
}

func init() {
	proto.RegisterType((*RestoreExecuteRequest)(nil), "generated.RestoreExecuteRequest")
	proto.RegisterType((*RestoreExecuteResponse)(nil), "generated.RestoreExecuteResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for RestoreItemAction service

type RestoreItemActionClient interface {
	AppliesTo(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*AppliesToResponse, error)
	Execute(ctx context.Context, in *RestoreExecuteRequest, opts ...grpc.CallOption) (*RestoreExecuteResponse, error)
}

type restoreItemActionClient struct {
	cc *grpc.ClientConn
}

func NewRestoreItemActionClient(cc *grpc.ClientConn) RestoreItemActionClient {
	return &restoreItemActionClient{cc}
}

func (c *restoreItemActionClient) AppliesTo(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*AppliesToResponse, error) {
	out := new(AppliesToResponse)
	err := grpc.Invoke(ctx, "/generated.RestoreItemAction/AppliesTo", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *restoreItemActionClient) Execute(ctx context.Context, in *RestoreExecuteRequest, opts ...grpc.CallOption) (*RestoreExecuteResponse, error) {
	out := new(RestoreExecuteResponse)
	err := grpc.Invoke(ctx, "/generated.RestoreItemAction/Execute", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for RestoreItemAction service

type RestoreItemActionServer interface {
	AppliesTo(context.Context, *Empty) (*AppliesToResponse, error)
	Execute(context.Context, *RestoreExecuteRequest) (*RestoreExecuteResponse, error)
}

func RegisterRestoreItemActionServer(s *grpc.Server, srv RestoreItemActionServer) {
	s.RegisterService(&_RestoreItemAction_serviceDesc, srv)
}

func _RestoreItemAction_AppliesTo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RestoreItemActionServer).AppliesTo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/generated.RestoreItemAction/AppliesTo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RestoreItemActionServer).AppliesTo(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _RestoreItemAction_Execute_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RestoreExecuteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RestoreItemActionServer).Execute(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/generated.RestoreItemAction/Execute",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RestoreItemActionServer).Execute(ctx, req.(*RestoreExecuteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _RestoreItemAction_serviceDesc = grpc.ServiceDesc{
	ServiceName: "generated.RestoreItemAction",
	HandlerType: (*RestoreItemActionServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AppliesTo",
			Handler:    _RestoreItemAction_AppliesTo_Handler,
		},
		{
			MethodName: "Execute",
			Handler:    _RestoreItemAction_Execute_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "RestoreItemAction.proto",
}

func init() { proto.RegisterFile("RestoreItemAction.proto", fileDescriptor3) }

var fileDescriptor3 = []byte{
	// 210 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x0f, 0x4a, 0x2d, 0x2e,
	0xc9, 0x2f, 0x4a, 0xf5, 0x2c, 0x49, 0xcd, 0x75, 0x4c, 0x2e, 0xc9, 0xcc, 0xcf, 0xd3, 0x2b, 0x28,
	0xca, 0x2f, 0xc9, 0x17, 0xe2, 0x4c, 0x4f, 0xcd, 0x4b, 0x2d, 0x4a, 0x2c, 0x49, 0x4d, 0x91, 0xe2,
	0x09, 0xce, 0x48, 0x2c, 0x4a, 0x4d, 0x81, 0x48, 0x28, 0xb9, 0x72, 0x89, 0x42, 0xf5, 0xb8, 0x56,
	0xa4, 0x26, 0x97, 0x96, 0xa4, 0x06, 0xa5, 0x16, 0x96, 0xa6, 0x16, 0x97, 0x08, 0x09, 0x71, 0xb1,
	0x64, 0x96, 0xa4, 0xe6, 0x4a, 0x30, 0x2a, 0x30, 0x6a, 0xf0, 0x04, 0x81, 0xd9, 0x42, 0x12, 0x5c,
	0xec, 0x45, 0x10, 0xc5, 0x12, 0x4c, 0x60, 0x61, 0x18, 0x57, 0xc9, 0x8d, 0x4b, 0x0c, 0xdd, 0x98,
	0xe2, 0x82, 0xfc, 0xbc, 0xe2, 0x54, 0x5c, 0xe6, 0x94, 0x27, 0x16, 0xe5, 0x65, 0xe6, 0xa5, 0x83,
	0xcd, 0xe1, 0x0c, 0x82, 0x71, 0x8d, 0x16, 0x30, 0x72, 0x09, 0x62, 0xf8, 0x41, 0xc8, 0x9a, 0x8b,
	0xd3, 0xb1, 0xa0, 0x20, 0x27, 0x33, 0xb5, 0x38, 0x24, 0x5f, 0x48, 0x40, 0x0f, 0xee, 0x17, 0x3d,
	0xd7, 0xdc, 0x82, 0x92, 0x4a, 0x29, 0x19, 0x24, 0x11, 0xb8, 0x3a, 0xb8, 0x03, 0xfc, 0xb8, 0xd8,
	0xa1, 0x6e, 0x12, 0x52, 0x40, 0x52, 0x88, 0xd5, 0xd7, 0x52, 0x8a, 0x78, 0x54, 0x40, 0xcc, 0x4b,
	0x62, 0x03, 0x07, 0x9c, 0x31, 0x20, 0x00, 0x00, 0xff, 0xff, 0xb9, 0x08, 0x09, 0x74, 0x6c, 0x01,
	0x00, 0x00,
}