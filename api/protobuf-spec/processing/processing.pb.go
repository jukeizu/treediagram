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

type MessageReplyReceived struct {
	Id                   string   `protobuf:"bytes,1,opt,name=Id,proto3" json:"Id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *MessageReplyReceived) Reset()         { *m = MessageReplyReceived{} }
func (m *MessageReplyReceived) String() string { return proto.CompactTextString(m) }
func (*MessageReplyReceived) ProtoMessage()    {}
func (*MessageReplyReceived) Descriptor() ([]byte, []int) {
	return fileDescriptor_processing_696f0d4cd7e00ec9, []int{0}
}
func (m *MessageReplyReceived) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MessageReplyReceived.Unmarshal(m, b)
}
func (m *MessageReplyReceived) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MessageReplyReceived.Marshal(b, m, deterministic)
}
func (dst *MessageReplyReceived) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MessageReplyReceived.Merge(dst, src)
}
func (m *MessageReplyReceived) XXX_Size() int {
	return xxx_messageInfo_MessageReplyReceived.Size(m)
}
func (m *MessageReplyReceived) XXX_DiscardUnknown() {
	xxx_messageInfo_MessageReplyReceived.DiscardUnknown(m)
}

var xxx_messageInfo_MessageReplyReceived proto.InternalMessageInfo

func (m *MessageReplyReceived) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

type MessageRequest struct {
	Id                   string    `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Source               string    `protobuf:"bytes,2,opt,name=source,proto3" json:"source,omitempty"`
	Bot                  *User     `protobuf:"bytes,3,opt,name=bot,proto3" json:"bot,omitempty"`
	Author               *User     `protobuf:"bytes,4,opt,name=author,proto3" json:"author,omitempty"`
	ChannelId            string    `protobuf:"bytes,5,opt,name=channelId,proto3" json:"channelId,omitempty"`
	ServerId             string    `protobuf:"bytes,6,opt,name=serverId,proto3" json:"serverId,omitempty"`
	Servers              []*Server `protobuf:"bytes,7,rep,name=servers,proto3" json:"servers,omitempty"`
	Mentions             []*User   `protobuf:"bytes,8,rep,name=mentions,proto3" json:"mentions,omitempty"`
	Content              string    `protobuf:"bytes,9,opt,name=content,proto3" json:"content,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *MessageRequest) Reset()         { *m = MessageRequest{} }
func (m *MessageRequest) String() string { return proto.CompactTextString(m) }
func (*MessageRequest) ProtoMessage()    {}
func (*MessageRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_processing_696f0d4cd7e00ec9, []int{1}
}
func (m *MessageRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MessageRequest.Unmarshal(m, b)
}
func (m *MessageRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MessageRequest.Marshal(b, m, deterministic)
}
func (dst *MessageRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MessageRequest.Merge(dst, src)
}
func (m *MessageRequest) XXX_Size() int {
	return xxx_messageInfo_MessageRequest.Size(m)
}
func (m *MessageRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_MessageRequest.DiscardUnknown(m)
}

var xxx_messageInfo_MessageRequest proto.InternalMessageInfo

func (m *MessageRequest) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *MessageRequest) GetSource() string {
	if m != nil {
		return m.Source
	}
	return ""
}

func (m *MessageRequest) GetBot() *User {
	if m != nil {
		return m.Bot
	}
	return nil
}

func (m *MessageRequest) GetAuthor() *User {
	if m != nil {
		return m.Author
	}
	return nil
}

func (m *MessageRequest) GetChannelId() string {
	if m != nil {
		return m.ChannelId
	}
	return ""
}

func (m *MessageRequest) GetServerId() string {
	if m != nil {
		return m.ServerId
	}
	return ""
}

func (m *MessageRequest) GetServers() []*Server {
	if m != nil {
		return m.Servers
	}
	return nil
}

func (m *MessageRequest) GetMentions() []*User {
	if m != nil {
		return m.Mentions
	}
	return nil
}

func (m *MessageRequest) GetContent() string {
	if m != nil {
		return m.Content
	}
	return ""
}

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
	return fileDescriptor_processing_696f0d4cd7e00ec9, []int{2}
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

type Server struct {
	Id                   string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Name                 string   `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Server) Reset()         { *m = Server{} }
func (m *Server) String() string { return proto.CompactTextString(m) }
func (*Server) ProtoMessage()    {}
func (*Server) Descriptor() ([]byte, []int) {
	return fileDescriptor_processing_696f0d4cd7e00ec9, []int{3}
}
func (m *Server) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Server.Unmarshal(m, b)
}
func (m *Server) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Server.Marshal(b, m, deterministic)
}
func (dst *Server) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Server.Merge(dst, src)
}
func (m *Server) XXX_Size() int {
	return xxx_messageInfo_Server.Size(m)
}
func (m *Server) XXX_DiscardUnknown() {
	xxx_messageInfo_Server.DiscardUnknown(m)
}

var xxx_messageInfo_Server proto.InternalMessageInfo

func (m *Server) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Server) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

type SendMessageRequestReply struct {
	Id                   string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SendMessageRequestReply) Reset()         { *m = SendMessageRequestReply{} }
func (m *SendMessageRequestReply) String() string { return proto.CompactTextString(m) }
func (*SendMessageRequestReply) ProtoMessage()    {}
func (*SendMessageRequestReply) Descriptor() ([]byte, []int) {
	return fileDescriptor_processing_696f0d4cd7e00ec9, []int{4}
}
func (m *SendMessageRequestReply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SendMessageRequestReply.Unmarshal(m, b)
}
func (m *SendMessageRequestReply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SendMessageRequestReply.Marshal(b, m, deterministic)
}
func (dst *SendMessageRequestReply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SendMessageRequestReply.Merge(dst, src)
}
func (m *SendMessageRequestReply) XXX_Size() int {
	return xxx_messageInfo_SendMessageRequestReply.Size(m)
}
func (m *SendMessageRequestReply) XXX_DiscardUnknown() {
	xxx_messageInfo_SendMessageRequestReply.DiscardUnknown(m)
}

var xxx_messageInfo_SendMessageRequestReply proto.InternalMessageInfo

func (m *SendMessageRequestReply) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

type MessageReplyRequest struct {
	Id                   string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *MessageReplyRequest) Reset()         { *m = MessageReplyRequest{} }
func (m *MessageReplyRequest) String() string { return proto.CompactTextString(m) }
func (*MessageReplyRequest) ProtoMessage()    {}
func (*MessageReplyRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_processing_696f0d4cd7e00ec9, []int{5}
}
func (m *MessageReplyRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MessageReplyRequest.Unmarshal(m, b)
}
func (m *MessageReplyRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MessageReplyRequest.Marshal(b, m, deterministic)
}
func (dst *MessageReplyRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MessageReplyRequest.Merge(dst, src)
}
func (m *MessageReplyRequest) XXX_Size() int {
	return xxx_messageInfo_MessageReplyRequest.Size(m)
}
func (m *MessageReplyRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_MessageReplyRequest.DiscardUnknown(m)
}

var xxx_messageInfo_MessageReplyRequest proto.InternalMessageInfo

func (m *MessageReplyRequest) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

type Response struct {
	Messages             []*Message `protobuf:"bytes,1,rep,name=messages,proto3" json:"messages,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *Response) Reset()         { *m = Response{} }
func (m *Response) String() string { return proto.CompactTextString(m) }
func (*Response) ProtoMessage()    {}
func (*Response) Descriptor() ([]byte, []int) {
	return fileDescriptor_processing_696f0d4cd7e00ec9, []int{6}
}
func (m *Response) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Response.Unmarshal(m, b)
}
func (m *Response) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Response.Marshal(b, m, deterministic)
}
func (dst *Response) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Response.Merge(dst, src)
}
func (m *Response) XXX_Size() int {
	return xxx_messageInfo_Response.Size(m)
}
func (m *Response) XXX_DiscardUnknown() {
	xxx_messageInfo_Response.DiscardUnknown(m)
}

var xxx_messageInfo_Response proto.InternalMessageInfo

func (m *Response) GetMessages() []*Message {
	if m != nil {
		return m.Messages
	}
	return nil
}

type Message struct {
	Content              string   `protobuf:"bytes,1,opt,name=content,proto3" json:"content,omitempty"`
	IsPrivateMessage     bool     `protobuf:"varint,2,opt,name=isPrivateMessage,proto3" json:"isPrivateMessage,omitempty"`
	IsRedirect           bool     `protobuf:"varint,3,opt,name=isRedirect,proto3" json:"isRedirect,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Message) Reset()         { *m = Message{} }
func (m *Message) String() string { return proto.CompactTextString(m) }
func (*Message) ProtoMessage()    {}
func (*Message) Descriptor() ([]byte, []int) {
	return fileDescriptor_processing_696f0d4cd7e00ec9, []int{7}
}
func (m *Message) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Message.Unmarshal(m, b)
}
func (m *Message) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Message.Marshal(b, m, deterministic)
}
func (dst *Message) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Message.Merge(dst, src)
}
func (m *Message) XXX_Size() int {
	return xxx_messageInfo_Message.Size(m)
}
func (m *Message) XXX_DiscardUnknown() {
	xxx_messageInfo_Message.DiscardUnknown(m)
}

var xxx_messageInfo_Message proto.InternalMessageInfo

func (m *Message) GetContent() string {
	if m != nil {
		return m.Content
	}
	return ""
}

func (m *Message) GetIsPrivateMessage() bool {
	if m != nil {
		return m.IsPrivateMessage
	}
	return false
}

func (m *Message) GetIsRedirect() bool {
	if m != nil {
		return m.IsRedirect
	}
	return false
}

type MessageReply struct {
	Id                   string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	ProcessingRequestId  string   `protobuf:"bytes,2,opt,name=processingRequestId,proto3" json:"processingRequestId,omitempty"`
	ChannelId            string   `protobuf:"bytes,3,opt,name=channelId,proto3" json:"channelId,omitempty"`
	UserId               string   `protobuf:"bytes,4,opt,name=userId,proto3" json:"userId,omitempty"`
	IsPrivateMessage     bool     `protobuf:"varint,5,opt,name=isPrivateMessage,proto3" json:"isPrivateMessage,omitempty"`
	IsRedirect           bool     `protobuf:"varint,6,opt,name=isRedirect,proto3" json:"isRedirect,omitempty"`
	Content              string   `protobuf:"bytes,7,opt,name=content,proto3" json:"content,omitempty"`
	Created              int32    `protobuf:"varint,8,opt,name=created,proto3" json:"created,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *MessageReply) Reset()         { *m = MessageReply{} }
func (m *MessageReply) String() string { return proto.CompactTextString(m) }
func (*MessageReply) ProtoMessage()    {}
func (*MessageReply) Descriptor() ([]byte, []int) {
	return fileDescriptor_processing_696f0d4cd7e00ec9, []int{8}
}
func (m *MessageReply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MessageReply.Unmarshal(m, b)
}
func (m *MessageReply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MessageReply.Marshal(b, m, deterministic)
}
func (dst *MessageReply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MessageReply.Merge(dst, src)
}
func (m *MessageReply) XXX_Size() int {
	return xxx_messageInfo_MessageReply.Size(m)
}
func (m *MessageReply) XXX_DiscardUnknown() {
	xxx_messageInfo_MessageReply.DiscardUnknown(m)
}

var xxx_messageInfo_MessageReply proto.InternalMessageInfo

func (m *MessageReply) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *MessageReply) GetProcessingRequestId() string {
	if m != nil {
		return m.ProcessingRequestId
	}
	return ""
}

func (m *MessageReply) GetChannelId() string {
	if m != nil {
		return m.ChannelId
	}
	return ""
}

func (m *MessageReply) GetUserId() string {
	if m != nil {
		return m.UserId
	}
	return ""
}

func (m *MessageReply) GetIsPrivateMessage() bool {
	if m != nil {
		return m.IsPrivateMessage
	}
	return false
}

func (m *MessageReply) GetIsRedirect() bool {
	if m != nil {
		return m.IsRedirect
	}
	return false
}

func (m *MessageReply) GetContent() string {
	if m != nil {
		return m.Content
	}
	return ""
}

func (m *MessageReply) GetCreated() int32 {
	if m != nil {
		return m.Created
	}
	return 0
}

type ProcessingRequest struct {
	Id                   string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	IntentId             string   `protobuf:"bytes,2,opt,name=intentId,proto3" json:"intentId,omitempty"`
	Source               string   `protobuf:"bytes,3,opt,name=source,proto3" json:"source,omitempty"`
	ChannelId            string   `protobuf:"bytes,4,opt,name=channelId,proto3" json:"channelId,omitempty"`
	ServerId             string   `protobuf:"bytes,5,opt,name=serverId,proto3" json:"serverId,omitempty"`
	BotId                string   `protobuf:"bytes,6,opt,name=botId,proto3" json:"botId,omitempty"`
	UserId               string   `protobuf:"bytes,7,opt,name=userId,proto3" json:"userId,omitempty"`
	Created              int32    `protobuf:"varint,8,opt,name=created,proto3" json:"created,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ProcessingRequest) Reset()         { *m = ProcessingRequest{} }
func (m *ProcessingRequest) String() string { return proto.CompactTextString(m) }
func (*ProcessingRequest) ProtoMessage()    {}
func (*ProcessingRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_processing_696f0d4cd7e00ec9, []int{9}
}
func (m *ProcessingRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ProcessingRequest.Unmarshal(m, b)
}
func (m *ProcessingRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ProcessingRequest.Marshal(b, m, deterministic)
}
func (dst *ProcessingRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ProcessingRequest.Merge(dst, src)
}
func (m *ProcessingRequest) XXX_Size() int {
	return xxx_messageInfo_ProcessingRequest.Size(m)
}
func (m *ProcessingRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ProcessingRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ProcessingRequest proto.InternalMessageInfo

func (m *ProcessingRequest) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *ProcessingRequest) GetIntentId() string {
	if m != nil {
		return m.IntentId
	}
	return ""
}

func (m *ProcessingRequest) GetSource() string {
	if m != nil {
		return m.Source
	}
	return ""
}

func (m *ProcessingRequest) GetChannelId() string {
	if m != nil {
		return m.ChannelId
	}
	return ""
}

func (m *ProcessingRequest) GetServerId() string {
	if m != nil {
		return m.ServerId
	}
	return ""
}

func (m *ProcessingRequest) GetBotId() string {
	if m != nil {
		return m.BotId
	}
	return ""
}

func (m *ProcessingRequest) GetUserId() string {
	if m != nil {
		return m.UserId
	}
	return ""
}

func (m *ProcessingRequest) GetCreated() int32 {
	if m != nil {
		return m.Created
	}
	return 0
}

type ProcessingEvent struct {
	Id                   string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	ProcessingRequestId  string   `protobuf:"bytes,2,opt,name=processingRequestId,proto3" json:"processingRequestId,omitempty"`
	Description          string   `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
	Type                 string   `protobuf:"bytes,4,opt,name=type,proto3" json:"type,omitempty"`
	Created              int32    `protobuf:"varint,5,opt,name=created,proto3" json:"created,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ProcessingEvent) Reset()         { *m = ProcessingEvent{} }
func (m *ProcessingEvent) String() string { return proto.CompactTextString(m) }
func (*ProcessingEvent) ProtoMessage()    {}
func (*ProcessingEvent) Descriptor() ([]byte, []int) {
	return fileDescriptor_processing_696f0d4cd7e00ec9, []int{10}
}
func (m *ProcessingEvent) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ProcessingEvent.Unmarshal(m, b)
}
func (m *ProcessingEvent) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ProcessingEvent.Marshal(b, m, deterministic)
}
func (dst *ProcessingEvent) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ProcessingEvent.Merge(dst, src)
}
func (m *ProcessingEvent) XXX_Size() int {
	return xxx_messageInfo_ProcessingEvent.Size(m)
}
func (m *ProcessingEvent) XXX_DiscardUnknown() {
	xxx_messageInfo_ProcessingEvent.DiscardUnknown(m)
}

var xxx_messageInfo_ProcessingEvent proto.InternalMessageInfo

func (m *ProcessingEvent) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *ProcessingEvent) GetProcessingRequestId() string {
	if m != nil {
		return m.ProcessingRequestId
	}
	return ""
}

func (m *ProcessingEvent) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func (m *ProcessingEvent) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

func (m *ProcessingEvent) GetCreated() int32 {
	if m != nil {
		return m.Created
	}
	return 0
}

func init() {
	proto.RegisterType((*MessageReplyReceived)(nil), "processing.MessageReplyReceived")
	proto.RegisterType((*MessageRequest)(nil), "processing.MessageRequest")
	proto.RegisterType((*User)(nil), "processing.User")
	proto.RegisterType((*Server)(nil), "processing.Server")
	proto.RegisterType((*SendMessageRequestReply)(nil), "processing.SendMessageRequestReply")
	proto.RegisterType((*MessageReplyRequest)(nil), "processing.MessageReplyRequest")
	proto.RegisterType((*Response)(nil), "processing.Response")
	proto.RegisterType((*Message)(nil), "processing.Message")
	proto.RegisterType((*MessageReply)(nil), "processing.MessageReply")
	proto.RegisterType((*ProcessingRequest)(nil), "processing.ProcessingRequest")
	proto.RegisterType((*ProcessingEvent)(nil), "processing.ProcessingEvent")
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
	SendMessageRequest(ctx context.Context, in *MessageRequest, opts ...grpc.CallOption) (*SendMessageRequestReply, error)
	GetMessageReply(ctx context.Context, in *MessageReplyRequest, opts ...grpc.CallOption) (*MessageReply, error)
}

type processingClient struct {
	cc *grpc.ClientConn
}

func NewProcessingClient(cc *grpc.ClientConn) ProcessingClient {
	return &processingClient{cc}
}

func (c *processingClient) SendMessageRequest(ctx context.Context, in *MessageRequest, opts ...grpc.CallOption) (*SendMessageRequestReply, error) {
	out := new(SendMessageRequestReply)
	err := c.cc.Invoke(ctx, "/processing.Processing/SendMessageRequest", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *processingClient) GetMessageReply(ctx context.Context, in *MessageReplyRequest, opts ...grpc.CallOption) (*MessageReply, error) {
	out := new(MessageReply)
	err := c.cc.Invoke(ctx, "/processing.Processing/GetMessageReply", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ProcessingServer is the server API for Processing service.
type ProcessingServer interface {
	SendMessageRequest(context.Context, *MessageRequest) (*SendMessageRequestReply, error)
	GetMessageReply(context.Context, *MessageReplyRequest) (*MessageReply, error)
}

func RegisterProcessingServer(s *grpc.Server, srv ProcessingServer) {
	s.RegisterService(&_Processing_serviceDesc, srv)
}

func _Processing_SendMessageRequest_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MessageRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProcessingServer).SendMessageRequest(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/processing.Processing/SendMessageRequest",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProcessingServer).SendMessageRequest(ctx, req.(*MessageRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Processing_GetMessageReply_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MessageReplyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProcessingServer).GetMessageReply(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/processing.Processing/GetMessageReply",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProcessingServer).GetMessageReply(ctx, req.(*MessageReplyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Processing_serviceDesc = grpc.ServiceDesc{
	ServiceName: "processing.Processing",
	HandlerType: (*ProcessingServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SendMessageRequest",
			Handler:    _Processing_SendMessageRequest_Handler,
		},
		{
			MethodName: "GetMessageReply",
			Handler:    _Processing_GetMessageReply_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "processing.proto",
}

func init() { proto.RegisterFile("processing.proto", fileDescriptor_processing_696f0d4cd7e00ec9) }

var fileDescriptor_processing_696f0d4cd7e00ec9 = []byte{
	// 571 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x55, 0xdd, 0x6e, 0xd3, 0x30,
	0x14, 0x5e, 0xda, 0xe6, 0xa7, 0xa7, 0x68, 0x2b, 0xee, 0x34, 0xac, 0x0a, 0x41, 0x65, 0x04, 0x2a,
	0x53, 0x35, 0x50, 0xb9, 0xe4, 0x1a, 0xa1, 0x5e, 0x80, 0xa6, 0x4c, 0x88, 0xeb, 0x34, 0x39, 0xda,
	0x22, 0x6d, 0x71, 0xb0, 0xdd, 0x4a, 0x7b, 0x01, 0x5e, 0x84, 0x67, 0x80, 0x77, 0xe1, 0x6d, 0x50,
	0x9c, 0x34, 0xb1, 0x9b, 0x14, 0x4d, 0xda, 0x9d, 0xcf, 0x39, 0x5f, 0xec, 0xf3, 0x7d, 0xe7, 0x7c,
	0x2d, 0x8c, 0x73, 0xc1, 0x63, 0x94, 0x32, 0xcd, 0xae, 0x2f, 0x72, 0xc1, 0x15, 0x27, 0xd0, 0x64,
	0xd8, 0x1b, 0x38, 0xfd, 0x82, 0x52, 0x46, 0xd7, 0x18, 0x62, 0x7e, 0x7b, 0x1f, 0x62, 0x8c, 0xe9,
	0x16, 0x13, 0x72, 0x0c, 0xbd, 0x55, 0x42, 0x9d, 0x99, 0x33, 0x1f, 0x86, 0xbd, 0x55, 0xc2, 0xfe,
	0xf4, 0xe0, 0xb8, 0x06, 0xfe, 0xd8, 0xa0, 0x54, 0x05, 0x24, 0xad, 0x21, 0x69, 0x42, 0xce, 0xc0,
	0x93, 0x7c, 0x23, 0x62, 0xa4, 0x3d, 0x9d, 0xab, 0x22, 0xc2, 0xa0, 0xbf, 0xe6, 0x8a, 0xf6, 0x67,
	0xce, 0x7c, 0xb4, 0x1c, 0x5f, 0x18, 0xed, 0x7c, 0x93, 0x28, 0xc2, 0xa2, 0x48, 0xe6, 0xe0, 0x45,
	0x1b, 0x75, 0xc3, 0x05, 0x1d, 0x1c, 0x80, 0x55, 0x75, 0xf2, 0x1c, 0x86, 0xf1, 0x4d, 0x94, 0x65,
	0x78, 0xbb, 0x4a, 0xa8, 0xab, 0x1f, 0x6a, 0x12, 0x64, 0x0a, 0x81, 0x44, 0xb1, 0x45, 0xb1, 0x4a,
	0xa8, 0xa7, 0x8b, 0x75, 0x4c, 0x16, 0xe0, 0x97, 0x67, 0x49, 0xfd, 0x59, 0x7f, 0x3e, 0x5a, 0x12,
	0xf3, 0x91, 0x2b, 0x5d, 0x0a, 0x77, 0x10, 0xb2, 0x80, 0xe0, 0x0e, 0x33, 0x95, 0xf2, 0x4c, 0xd2,
	0x40, 0xc3, 0xdb, 0x3d, 0xd5, 0x08, 0x42, 0xc1, 0x8f, 0x79, 0xa6, 0x30, 0x53, 0x74, 0xa8, 0x9f,
	0xdd, 0x85, 0xec, 0x1c, 0x06, 0x05, 0xb6, 0xa5, 0x16, 0x81, 0x41, 0x16, 0xdd, 0xed, 0xb4, 0xd2,
	0x67, 0xb6, 0x00, 0xaf, 0x6c, 0xe3, 0x41, 0xe8, 0xb7, 0xf0, 0xec, 0x0a, 0xb3, 0xc4, 0x9e, 0x8a,
	0x9e, 0xe2, 0xfe, 0xe7, 0xec, 0x35, 0x4c, 0xec, 0x29, 0x77, 0x4e, 0x90, 0x7d, 0x84, 0x20, 0x44,
	0x99, 0xf3, 0x4c, 0x22, 0x79, 0x57, 0xf0, 0xd7, 0x9f, 0x48, 0xea, 0x68, 0xfe, 0x13, 0x93, 0xff,
	0xee, 0xba, 0x1a, 0xc4, 0x38, 0xf8, 0x55, 0xd2, 0x54, 0xc3, 0xb1, 0xd4, 0x20, 0xe7, 0x30, 0x4e,
	0xe5, 0xa5, 0x48, 0xb7, 0x91, 0xc2, 0x0a, 0xad, 0x39, 0x05, 0x61, 0x2b, 0x4f, 0x5e, 0x00, 0xa4,
	0x32, 0xc4, 0x24, 0x15, 0x18, 0x97, 0xeb, 0x13, 0x84, 0x46, 0x86, 0xfd, 0xec, 0xc1, 0x13, 0x93,
	0x55, 0x4b, 0xb4, 0xf7, 0x30, 0x69, 0x3a, 0xae, 0x38, 0xaf, 0x92, 0x4a, 0xc3, 0xae, 0x92, 0xbd,
	0x5c, 0xfd, 0xfd, 0xe5, 0x3a, 0x03, 0x6f, 0x23, 0xf5, 0x6a, 0x0d, 0xca, 0x05, 0x2f, 0xa3, 0x4e,
	0x52, 0xee, 0x83, 0x48, 0x79, 0xfb, 0xa4, 0x4c, 0xe9, 0x7c, 0x5b, 0xba, 0xa2, 0x22, 0x30, 0x52,
	0x98, 0xd0, 0x60, 0xe6, 0xcc, 0xdd, 0x70, 0x17, 0xb2, 0xbf, 0x0e, 0x3c, 0xbd, 0xdc, 0x67, 0xd3,
	0x52, 0x63, 0x0a, 0x41, 0xaa, 0x6f, 0xaa, 0x25, 0xa8, 0x63, 0xc3, 0xba, 0x7d, 0xcb, 0xba, 0x96,
	0x1e, 0x83, 0xff, 0x99, 0xcd, 0xdd, 0x33, 0xdb, 0x29, 0xb8, 0x6b, 0xae, 0x6a, 0x17, 0x96, 0x81,
	0xa1, 0xa0, 0x6f, 0x29, 0x78, 0x98, 0xdb, 0x2f, 0x07, 0x4e, 0x1a, 0x6e, 0x9f, 0xb6, 0x85, 0x12,
	0x8f, 0x9f, 0xf3, 0x0c, 0x46, 0x09, 0xca, 0x58, 0xa4, 0x79, 0x61, 0xdf, 0x8a, 0xb4, 0x99, 0x2a,
	0x0c, 0xa7, 0xee, 0x73, 0xac, 0x48, 0xeb, 0xb3, 0xd9, 0xa5, 0x6b, 0x75, 0xb9, 0xfc, 0xed, 0x00,
	0x34, 0x5d, 0x92, 0xef, 0x40, 0xda, 0xce, 0x24, 0xd3, 0x2e, 0xff, 0x94, 0xb5, 0xe9, 0x2b, 0xfb,
	0xa7, 0xa8, 0xd3, 0xd5, 0xec, 0x88, 0x7c, 0x85, 0x93, 0xcf, 0xa8, 0xac, 0xa5, 0x7f, 0xd9, 0x79,
	0x6b, 0x63, 0xf2, 0x29, 0x3d, 0x04, 0x60, 0x47, 0x6b, 0x4f, 0xff, 0x21, 0x7c, 0xf8, 0x17, 0x00,
	0x00, 0xff, 0xff, 0x06, 0xef, 0x54, 0x6a, 0x24, 0x06, 0x00, 0x00,
}
