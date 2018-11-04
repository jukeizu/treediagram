// Code generated by protoc-gen-go. DO NOT EDIT.
// source: processing.proto

package processing

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

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type User struct {
	Id                   string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Name                 string   `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *User) Reset()         { *m = User{} }
func (m *User) String() string { return proto.CompactTextString(m) }
func (*User) ProtoMessage()    {}
func (*User) Descriptor() ([]byte, []int) {
	return fileDescriptor_processing_b8799dbbbd904ec6, []int{0}
}
func (m *User) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_User.Unmarshal(m, b)
}
func (m *User) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_User.Marshal(b, m, deterministic)
}
func (dst *User) XXX_Merge(src proto.Message) {
	xxx_messageInfo_User.Merge(dst, src)
}
func (m *User) XXX_Size() int {
	return xxx_messageInfo_User.Size(m)
}
func (m *User) XXX_DiscardUnknown() {
	xxx_messageInfo_User.DiscardUnknown(m)
}

var xxx_messageInfo_User proto.InternalMessageInfo

func (m *User) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *User) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

type TreediagramRequest struct {
	Source               string   `protobuf:"bytes,1,opt,name=source,proto3" json:"source,omitempty"`
	CorrelationId        string   `protobuf:"bytes,2,opt,name=correlationId,proto3" json:"correlationId,omitempty"`
	Bot                  *User    `protobuf:"bytes,3,opt,name=bot,proto3" json:"bot,omitempty"`
	Author               *User    `protobuf:"bytes,4,opt,name=author,proto3" json:"author,omitempty"`
	ChannelId            string   `protobuf:"bytes,5,opt,name=channelId,proto3" json:"channelId,omitempty"`
	ServerId             string   `protobuf:"bytes,6,opt,name=serverId,proto3" json:"serverId,omitempty"`
	Mentions             []*User  `protobuf:"bytes,7,rep,name=mentions,proto3" json:"mentions,omitempty"`
	Content              string   `protobuf:"bytes,8,opt,name=content,proto3" json:"content,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TreediagramRequest) Reset()         { *m = TreediagramRequest{} }
func (m *TreediagramRequest) String() string { return proto.CompactTextString(m) }
func (*TreediagramRequest) ProtoMessage()    {}
func (*TreediagramRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_processing_b8799dbbbd904ec6, []int{1}
}
func (m *TreediagramRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TreediagramRequest.Unmarshal(m, b)
}
func (m *TreediagramRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TreediagramRequest.Marshal(b, m, deterministic)
}
func (dst *TreediagramRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TreediagramRequest.Merge(dst, src)
}
func (m *TreediagramRequest) XXX_Size() int {
	return xxx_messageInfo_TreediagramRequest.Size(m)
}
func (m *TreediagramRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_TreediagramRequest.DiscardUnknown(m)
}

var xxx_messageInfo_TreediagramRequest proto.InternalMessageInfo

func (m *TreediagramRequest) GetSource() string {
	if m != nil {
		return m.Source
	}
	return ""
}

func (m *TreediagramRequest) GetCorrelationId() string {
	if m != nil {
		return m.CorrelationId
	}
	return ""
}

func (m *TreediagramRequest) GetBot() *User {
	if m != nil {
		return m.Bot
	}
	return nil
}

func (m *TreediagramRequest) GetAuthor() *User {
	if m != nil {
		return m.Author
	}
	return nil
}

func (m *TreediagramRequest) GetChannelId() string {
	if m != nil {
		return m.ChannelId
	}
	return ""
}

func (m *TreediagramRequest) GetServerId() string {
	if m != nil {
		return m.ServerId
	}
	return ""
}

func (m *TreediagramRequest) GetMentions() []*User {
	if m != nil {
		return m.Mentions
	}
	return nil
}

func (m *TreediagramRequest) GetContent() string {
	if m != nil {
		return m.Content
	}
	return ""
}

type TreediagramReply struct {
	Id                   string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TreediagramReply) Reset()         { *m = TreediagramReply{} }
func (m *TreediagramReply) String() string { return proto.CompactTextString(m) }
func (*TreediagramReply) ProtoMessage()    {}
func (*TreediagramReply) Descriptor() ([]byte, []int) {
	return fileDescriptor_processing_b8799dbbbd904ec6, []int{2}
}
func (m *TreediagramReply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TreediagramReply.Unmarshal(m, b)
}
func (m *TreediagramReply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TreediagramReply.Marshal(b, m, deterministic)
}
func (dst *TreediagramReply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TreediagramReply.Merge(dst, src)
}
func (m *TreediagramReply) XXX_Size() int {
	return xxx_messageInfo_TreediagramReply.Size(m)
}
func (m *TreediagramReply) XXX_DiscardUnknown() {
	xxx_messageInfo_TreediagramReply.DiscardUnknown(m)
}

var xxx_messageInfo_TreediagramReply proto.InternalMessageInfo

func (m *TreediagramReply) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

type Match struct {
	Id                   string   `protobuf:"bytes,1,opt,name=Id,proto3" json:"Id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Match) Reset()         { *m = Match{} }
func (m *Match) String() string { return proto.CompactTextString(m) }
func (*Match) ProtoMessage()    {}
func (*Match) Descriptor() ([]byte, []int) {
	return fileDescriptor_processing_b8799dbbbd904ec6, []int{3}
}
func (m *Match) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Match.Unmarshal(m, b)
}
func (m *Match) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Match.Marshal(b, m, deterministic)
}
func (dst *Match) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Match.Merge(dst, src)
}
func (m *Match) XXX_Size() int {
	return xxx_messageInfo_Match.Size(m)
}
func (m *Match) XXX_DiscardUnknown() {
	xxx_messageInfo_Match.DiscardUnknown(m)
}

var xxx_messageInfo_Match proto.InternalMessageInfo

func (m *Match) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func init() {
	proto.RegisterType((*User)(nil), "processing.User")
	proto.RegisterType((*TreediagramRequest)(nil), "processing.TreediagramRequest")
	proto.RegisterType((*TreediagramReply)(nil), "processing.TreediagramReply")
	proto.RegisterType((*Match)(nil), "processing.Match")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// ProcessingClient is the client API for Processing service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ProcessingClient interface {
	Request(ctx context.Context, in *TreediagramRequest, opts ...grpc.CallOption) (*TreediagramReply, error)
}

type processingClient struct {
	cc *grpc.ClientConn
}

func NewProcessingClient(cc *grpc.ClientConn) ProcessingClient {
	return &processingClient{cc}
}

func (c *processingClient) Request(ctx context.Context, in *TreediagramRequest, opts ...grpc.CallOption) (*TreediagramReply, error) {
	out := new(TreediagramReply)
	err := c.cc.Invoke(ctx, "/processing.Processing/Request", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ProcessingServer is the server API for Processing service.
type ProcessingServer interface {
	Request(context.Context, *TreediagramRequest) (*TreediagramReply, error)
}

func RegisterProcessingServer(s *grpc.Server, srv ProcessingServer) {
	s.RegisterService(&_Processing_serviceDesc, srv)
}

func _Processing_Request_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TreediagramRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProcessingServer).Request(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/processing.Processing/Request",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProcessingServer).Request(ctx, req.(*TreediagramRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Processing_serviceDesc = grpc.ServiceDesc{
	ServiceName: "processing.Processing",
	HandlerType: (*ProcessingServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Request",
			Handler:    _Processing_Request_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "processing.proto",
}

func init() { proto.RegisterFile("processing.proto", fileDescriptor_processing_b8799dbbbd904ec6) }

var fileDescriptor_processing_b8799dbbbd904ec6 = []byte{
	// 295 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x91, 0xc1, 0x6a, 0x3a, 0x31,
	0x10, 0xc6, 0xff, 0xae, 0xba, 0xea, 0xfc, 0x69, 0x91, 0x39, 0xb4, 0x41, 0xa4, 0x48, 0xe8, 0x41,
	0x4a, 0xf1, 0x60, 0x9f, 0x22, 0x87, 0x42, 0x91, 0x96, 0x9e, 0x63, 0x32, 0xb8, 0x0b, 0xbb, 0xc9,
	0x36, 0xc9, 0x16, 0x7c, 0x9d, 0x3e, 0x69, 0xd9, 0xb0, 0xae, 0xda, 0xd6, 0x5b, 0xe6, 0xfb, 0x4d,
	0x66, 0x3e, 0xbe, 0x81, 0x69, 0xe5, 0xac, 0x22, 0xef, 0x73, 0xb3, 0x5b, 0x55, 0xce, 0x06, 0x8b,
	0x70, 0x54, 0xf8, 0x03, 0x0c, 0xde, 0x3c, 0x39, 0xbc, 0x86, 0x24, 0xd7, 0xac, 0xb7, 0xe8, 0x2d,
	0x27, 0x9b, 0x24, 0xd7, 0x88, 0x30, 0x30, 0xb2, 0x24, 0x96, 0x44, 0x25, 0xbe, 0xf9, 0x57, 0x02,
	0xf8, 0xea, 0x88, 0x74, 0x2e, 0x77, 0x4e, 0x96, 0x1b, 0xfa, 0xa8, 0xc9, 0x07, 0xbc, 0x81, 0xd4,
	0xdb, 0xda, 0x29, 0x6a, 0xbf, 0xb7, 0x15, 0xde, 0xc3, 0x95, 0xb2, 0xce, 0x51, 0x21, 0x43, 0x6e,
	0x8d, 0xd0, 0xed, 0xac, 0x73, 0x11, 0x39, 0xf4, 0xb7, 0x36, 0xb0, 0xfe, 0xa2, 0xb7, 0xfc, 0xbf,
	0x9e, 0xae, 0x4e, 0xcc, 0x36, 0xbe, 0x36, 0x0d, 0xc4, 0x25, 0xa4, 0xb2, 0x0e, 0x99, 0x75, 0x6c,
	0x70, 0xa1, 0xad, 0xe5, 0x38, 0x87, 0x89, 0xca, 0xa4, 0x31, 0x54, 0x08, 0xcd, 0x86, 0x71, 0xdf,
	0x51, 0xc0, 0x19, 0x8c, 0x3d, 0xb9, 0x4f, 0x72, 0x42, 0xb3, 0x34, 0xc2, 0xae, 0xc6, 0x47, 0x18,
	0x97, 0x64, 0x1a, 0x53, 0x9e, 0x8d, 0x16, 0xfd, 0x3f, 0xb7, 0x74, 0x1d, 0xc8, 0x60, 0xa4, 0xac,
	0x09, 0x64, 0x02, 0x1b, 0xc7, 0x41, 0x87, 0x92, 0x73, 0x98, 0x9e, 0x65, 0x54, 0x15, 0xfb, 0x9f,
	0xe1, 0xf2, 0x5b, 0x18, 0x3e, 0xcb, 0xa0, 0xb2, 0x06, 0x88, 0x0e, 0x08, 0xbd, 0x7e, 0x07, 0x78,
	0xe9, 0x76, 0xa2, 0x80, 0xd1, 0x21, 0xe3, 0xbb, 0x53, 0x2f, 0xbf, 0x6f, 0x30, 0x9b, 0x5f, 0xe4,
	0x55, 0xb1, 0xe7, 0xff, 0xb6, 0x69, 0xbc, 0xfc, 0xd3, 0x77, 0x00, 0x00, 0x00, 0xff, 0xff, 0xd1,
	0xa2, 0x38, 0xc1, 0x0d, 0x02, 0x00, 0x00,
}
