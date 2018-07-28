// Code generated by protoc-gen-go. DO NOT EDIT.
// source: publishing/publishing.proto

package publishing

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

type SendMessageRequest struct {
	Message              *Message `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
	Request              *Request `protobuf:"bytes,2,opt,name=request,proto3" json:"request,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SendMessageRequest) Reset()         { *m = SendMessageRequest{} }
func (m *SendMessageRequest) String() string { return proto.CompactTextString(m) }
func (*SendMessageRequest) ProtoMessage()    {}
func (*SendMessageRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_publishing_805a4042085b2f89, []int{0}
}
func (m *SendMessageRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SendMessageRequest.Unmarshal(m, b)
}
func (m *SendMessageRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SendMessageRequest.Marshal(b, m, deterministic)
}
func (dst *SendMessageRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SendMessageRequest.Merge(dst, src)
}
func (m *SendMessageRequest) XXX_Size() int {
	return xxx_messageInfo_SendMessageRequest.Size(m)
}
func (m *SendMessageRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_SendMessageRequest.DiscardUnknown(m)
}

var xxx_messageInfo_SendMessageRequest proto.InternalMessageInfo

func (m *SendMessageRequest) GetMessage() *Message {
	if m != nil {
		return m.Message
	}
	return nil
}

func (m *SendMessageRequest) GetRequest() *Request {
	if m != nil {
		return m.Request
	}
	return nil
}

type SendMessageReply struct {
	Id                   string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SendMessageReply) Reset()         { *m = SendMessageReply{} }
func (m *SendMessageReply) String() string { return proto.CompactTextString(m) }
func (*SendMessageReply) ProtoMessage()    {}
func (*SendMessageReply) Descriptor() ([]byte, []int) {
	return fileDescriptor_publishing_805a4042085b2f89, []int{1}
}
func (m *SendMessageReply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SendMessageReply.Unmarshal(m, b)
}
func (m *SendMessageReply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SendMessageReply.Marshal(b, m, deterministic)
}
func (dst *SendMessageReply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SendMessageReply.Merge(dst, src)
}
func (m *SendMessageReply) XXX_Size() int {
	return xxx_messageInfo_SendMessageReply.Size(m)
}
func (m *SendMessageReply) XXX_DiscardUnknown() {
	xxx_messageInfo_SendMessageReply.DiscardUnknown(m)
}

var xxx_messageInfo_SendMessageReply proto.InternalMessageInfo

func (m *SendMessageReply) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

type Message struct {
	Content              string   `protobuf:"bytes,1,opt,name=content,proto3" json:"content,omitempty"`
	Embed                *Embed   `protobuf:"bytes,2,opt,name=embed,proto3" json:"embed,omitempty"`
	Tts                  bool     `protobuf:"varint,3,opt,name=tts,proto3" json:"tts,omitempty"`
	Files                []*File  `protobuf:"bytes,4,rep,name=files,proto3" json:"files,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Message) Reset()         { *m = Message{} }
func (m *Message) String() string { return proto.CompactTextString(m) }
func (*Message) ProtoMessage()    {}
func (*Message) Descriptor() ([]byte, []int) {
	return fileDescriptor_publishing_805a4042085b2f89, []int{2}
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

func (m *Message) GetEmbed() *Embed {
	if m != nil {
		return m.Embed
	}
	return nil
}

func (m *Message) GetTts() bool {
	if m != nil {
		return m.Tts
	}
	return false
}

func (m *Message) GetFiles() []*File {
	if m != nil {
		return m.Files
	}
	return nil
}

type Request struct {
	CorrelationId        string   `protobuf:"bytes,1,opt,name=correlationId,proto3" json:"correlationId,omitempty"`
	ChannelId            string   `protobuf:"bytes,2,opt,name=channelId,proto3" json:"channelId,omitempty"`
	User                 *User    `protobuf:"bytes,3,opt,name=user,proto3" json:"user,omitempty"`
	PrivateMessage       bool     `protobuf:"varint,4,opt,name=privateMessage,proto3" json:"privateMessage,omitempty"`
	IsRedirect           bool     `protobuf:"varint,5,opt,name=isRedirect,proto3" json:"isRedirect,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Request) Reset()         { *m = Request{} }
func (m *Request) String() string { return proto.CompactTextString(m) }
func (*Request) ProtoMessage()    {}
func (*Request) Descriptor() ([]byte, []int) {
	return fileDescriptor_publishing_805a4042085b2f89, []int{3}
}
func (m *Request) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Request.Unmarshal(m, b)
}
func (m *Request) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Request.Marshal(b, m, deterministic)
}
func (dst *Request) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Request.Merge(dst, src)
}
func (m *Request) XXX_Size() int {
	return xxx_messageInfo_Request.Size(m)
}
func (m *Request) XXX_DiscardUnknown() {
	xxx_messageInfo_Request.DiscardUnknown(m)
}

var xxx_messageInfo_Request proto.InternalMessageInfo

func (m *Request) GetCorrelationId() string {
	if m != nil {
		return m.CorrelationId
	}
	return ""
}

func (m *Request) GetChannelId() string {
	if m != nil {
		return m.ChannelId
	}
	return ""
}

func (m *Request) GetUser() *User {
	if m != nil {
		return m.User
	}
	return nil
}

func (m *Request) GetPrivateMessage() bool {
	if m != nil {
		return m.PrivateMessage
	}
	return false
}

func (m *Request) GetIsRedirect() bool {
	if m != nil {
		return m.IsRedirect
	}
	return false
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
	return fileDescriptor_publishing_805a4042085b2f89, []int{4}
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

type EmbedFooter struct {
	Text                 string   `protobuf:"bytes,1,opt,name=text,proto3" json:"text,omitempty"`
	IconUrl              string   `protobuf:"bytes,2,opt,name=iconUrl,proto3" json:"iconUrl,omitempty"`
	ProxyIconUrl         string   `protobuf:"bytes,3,opt,name=proxyIconUrl,proto3" json:"proxyIconUrl,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *EmbedFooter) Reset()         { *m = EmbedFooter{} }
func (m *EmbedFooter) String() string { return proto.CompactTextString(m) }
func (*EmbedFooter) ProtoMessage()    {}
func (*EmbedFooter) Descriptor() ([]byte, []int) {
	return fileDescriptor_publishing_805a4042085b2f89, []int{5}
}
func (m *EmbedFooter) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EmbedFooter.Unmarshal(m, b)
}
func (m *EmbedFooter) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EmbedFooter.Marshal(b, m, deterministic)
}
func (dst *EmbedFooter) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EmbedFooter.Merge(dst, src)
}
func (m *EmbedFooter) XXX_Size() int {
	return xxx_messageInfo_EmbedFooter.Size(m)
}
func (m *EmbedFooter) XXX_DiscardUnknown() {
	xxx_messageInfo_EmbedFooter.DiscardUnknown(m)
}

var xxx_messageInfo_EmbedFooter proto.InternalMessageInfo

func (m *EmbedFooter) GetText() string {
	if m != nil {
		return m.Text
	}
	return ""
}

func (m *EmbedFooter) GetIconUrl() string {
	if m != nil {
		return m.IconUrl
	}
	return ""
}

func (m *EmbedFooter) GetProxyIconUrl() string {
	if m != nil {
		return m.ProxyIconUrl
	}
	return ""
}

type EmbedImage struct {
	Url                  string   `protobuf:"bytes,1,opt,name=url,proto3" json:"url,omitempty"`
	ProxyUrl             string   `protobuf:"bytes,2,opt,name=proxyUrl,proto3" json:"proxyUrl,omitempty"`
	Width                int32    `protobuf:"varint,3,opt,name=width,proto3" json:"width,omitempty"`
	Height               int32    `protobuf:"varint,4,opt,name=height,proto3" json:"height,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *EmbedImage) Reset()         { *m = EmbedImage{} }
func (m *EmbedImage) String() string { return proto.CompactTextString(m) }
func (*EmbedImage) ProtoMessage()    {}
func (*EmbedImage) Descriptor() ([]byte, []int) {
	return fileDescriptor_publishing_805a4042085b2f89, []int{6}
}
func (m *EmbedImage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EmbedImage.Unmarshal(m, b)
}
func (m *EmbedImage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EmbedImage.Marshal(b, m, deterministic)
}
func (dst *EmbedImage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EmbedImage.Merge(dst, src)
}
func (m *EmbedImage) XXX_Size() int {
	return xxx_messageInfo_EmbedImage.Size(m)
}
func (m *EmbedImage) XXX_DiscardUnknown() {
	xxx_messageInfo_EmbedImage.DiscardUnknown(m)
}

var xxx_messageInfo_EmbedImage proto.InternalMessageInfo

func (m *EmbedImage) GetUrl() string {
	if m != nil {
		return m.Url
	}
	return ""
}

func (m *EmbedImage) GetProxyUrl() string {
	if m != nil {
		return m.ProxyUrl
	}
	return ""
}

func (m *EmbedImage) GetWidth() int32 {
	if m != nil {
		return m.Width
	}
	return 0
}

func (m *EmbedImage) GetHeight() int32 {
	if m != nil {
		return m.Height
	}
	return 0
}

type EmbedThumbnail struct {
	Url                  string   `protobuf:"bytes,1,opt,name=url,proto3" json:"url,omitempty"`
	ProxyUrl             string   `protobuf:"bytes,2,opt,name=proxyUrl,proto3" json:"proxyUrl,omitempty"`
	Width                int32    `protobuf:"varint,3,opt,name=width,proto3" json:"width,omitempty"`
	Height               int32    `protobuf:"varint,4,opt,name=height,proto3" json:"height,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *EmbedThumbnail) Reset()         { *m = EmbedThumbnail{} }
func (m *EmbedThumbnail) String() string { return proto.CompactTextString(m) }
func (*EmbedThumbnail) ProtoMessage()    {}
func (*EmbedThumbnail) Descriptor() ([]byte, []int) {
	return fileDescriptor_publishing_805a4042085b2f89, []int{7}
}
func (m *EmbedThumbnail) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EmbedThumbnail.Unmarshal(m, b)
}
func (m *EmbedThumbnail) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EmbedThumbnail.Marshal(b, m, deterministic)
}
func (dst *EmbedThumbnail) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EmbedThumbnail.Merge(dst, src)
}
func (m *EmbedThumbnail) XXX_Size() int {
	return xxx_messageInfo_EmbedThumbnail.Size(m)
}
func (m *EmbedThumbnail) XXX_DiscardUnknown() {
	xxx_messageInfo_EmbedThumbnail.DiscardUnknown(m)
}

var xxx_messageInfo_EmbedThumbnail proto.InternalMessageInfo

func (m *EmbedThumbnail) GetUrl() string {
	if m != nil {
		return m.Url
	}
	return ""
}

func (m *EmbedThumbnail) GetProxyUrl() string {
	if m != nil {
		return m.ProxyUrl
	}
	return ""
}

func (m *EmbedThumbnail) GetWidth() int32 {
	if m != nil {
		return m.Width
	}
	return 0
}

func (m *EmbedThumbnail) GetHeight() int32 {
	if m != nil {
		return m.Height
	}
	return 0
}

type EmbedVideo struct {
	Url                  string   `protobuf:"bytes,1,opt,name=url,proto3" json:"url,omitempty"`
	ProxyUrl             string   `protobuf:"bytes,2,opt,name=proxyUrl,proto3" json:"proxyUrl,omitempty"`
	Width                int32    `protobuf:"varint,3,opt,name=width,proto3" json:"width,omitempty"`
	Height               int32    `protobuf:"varint,4,opt,name=height,proto3" json:"height,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *EmbedVideo) Reset()         { *m = EmbedVideo{} }
func (m *EmbedVideo) String() string { return proto.CompactTextString(m) }
func (*EmbedVideo) ProtoMessage()    {}
func (*EmbedVideo) Descriptor() ([]byte, []int) {
	return fileDescriptor_publishing_805a4042085b2f89, []int{8}
}
func (m *EmbedVideo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EmbedVideo.Unmarshal(m, b)
}
func (m *EmbedVideo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EmbedVideo.Marshal(b, m, deterministic)
}
func (dst *EmbedVideo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EmbedVideo.Merge(dst, src)
}
func (m *EmbedVideo) XXX_Size() int {
	return xxx_messageInfo_EmbedVideo.Size(m)
}
func (m *EmbedVideo) XXX_DiscardUnknown() {
	xxx_messageInfo_EmbedVideo.DiscardUnknown(m)
}

var xxx_messageInfo_EmbedVideo proto.InternalMessageInfo

func (m *EmbedVideo) GetUrl() string {
	if m != nil {
		return m.Url
	}
	return ""
}

func (m *EmbedVideo) GetProxyUrl() string {
	if m != nil {
		return m.ProxyUrl
	}
	return ""
}

func (m *EmbedVideo) GetWidth() int32 {
	if m != nil {
		return m.Width
	}
	return 0
}

func (m *EmbedVideo) GetHeight() int32 {
	if m != nil {
		return m.Height
	}
	return 0
}

type EmbedProvider struct {
	Url                  string   `protobuf:"bytes,1,opt,name=url,proto3" json:"url,omitempty"`
	Name                 string   `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *EmbedProvider) Reset()         { *m = EmbedProvider{} }
func (m *EmbedProvider) String() string { return proto.CompactTextString(m) }
func (*EmbedProvider) ProtoMessage()    {}
func (*EmbedProvider) Descriptor() ([]byte, []int) {
	return fileDescriptor_publishing_805a4042085b2f89, []int{9}
}
func (m *EmbedProvider) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EmbedProvider.Unmarshal(m, b)
}
func (m *EmbedProvider) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EmbedProvider.Marshal(b, m, deterministic)
}
func (dst *EmbedProvider) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EmbedProvider.Merge(dst, src)
}
func (m *EmbedProvider) XXX_Size() int {
	return xxx_messageInfo_EmbedProvider.Size(m)
}
func (m *EmbedProvider) XXX_DiscardUnknown() {
	xxx_messageInfo_EmbedProvider.DiscardUnknown(m)
}

var xxx_messageInfo_EmbedProvider proto.InternalMessageInfo

func (m *EmbedProvider) GetUrl() string {
	if m != nil {
		return m.Url
	}
	return ""
}

func (m *EmbedProvider) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

type EmbedAuthor struct {
	Url                  string   `protobuf:"bytes,1,opt,name=url,proto3" json:"url,omitempty"`
	Name                 string   `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	IconUrl              string   `protobuf:"bytes,3,opt,name=iconUrl,proto3" json:"iconUrl,omitempty"`
	ProxyIconUrl         string   `protobuf:"bytes,4,opt,name=proxyIconUrl,proto3" json:"proxyIconUrl,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *EmbedAuthor) Reset()         { *m = EmbedAuthor{} }
func (m *EmbedAuthor) String() string { return proto.CompactTextString(m) }
func (*EmbedAuthor) ProtoMessage()    {}
func (*EmbedAuthor) Descriptor() ([]byte, []int) {
	return fileDescriptor_publishing_805a4042085b2f89, []int{10}
}
func (m *EmbedAuthor) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EmbedAuthor.Unmarshal(m, b)
}
func (m *EmbedAuthor) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EmbedAuthor.Marshal(b, m, deterministic)
}
func (dst *EmbedAuthor) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EmbedAuthor.Merge(dst, src)
}
func (m *EmbedAuthor) XXX_Size() int {
	return xxx_messageInfo_EmbedAuthor.Size(m)
}
func (m *EmbedAuthor) XXX_DiscardUnknown() {
	xxx_messageInfo_EmbedAuthor.DiscardUnknown(m)
}

var xxx_messageInfo_EmbedAuthor proto.InternalMessageInfo

func (m *EmbedAuthor) GetUrl() string {
	if m != nil {
		return m.Url
	}
	return ""
}

func (m *EmbedAuthor) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *EmbedAuthor) GetIconUrl() string {
	if m != nil {
		return m.IconUrl
	}
	return ""
}

func (m *EmbedAuthor) GetProxyIconUrl() string {
	if m != nil {
		return m.ProxyIconUrl
	}
	return ""
}

type EmbedField struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Value                string   `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
	Inline               bool     `protobuf:"varint,3,opt,name=inline,proto3" json:"inline,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *EmbedField) Reset()         { *m = EmbedField{} }
func (m *EmbedField) String() string { return proto.CompactTextString(m) }
func (*EmbedField) ProtoMessage()    {}
func (*EmbedField) Descriptor() ([]byte, []int) {
	return fileDescriptor_publishing_805a4042085b2f89, []int{11}
}
func (m *EmbedField) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EmbedField.Unmarshal(m, b)
}
func (m *EmbedField) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EmbedField.Marshal(b, m, deterministic)
}
func (dst *EmbedField) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EmbedField.Merge(dst, src)
}
func (m *EmbedField) XXX_Size() int {
	return xxx_messageInfo_EmbedField.Size(m)
}
func (m *EmbedField) XXX_DiscardUnknown() {
	xxx_messageInfo_EmbedField.DiscardUnknown(m)
}

var xxx_messageInfo_EmbedField proto.InternalMessageInfo

func (m *EmbedField) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *EmbedField) GetValue() string {
	if m != nil {
		return m.Value
	}
	return ""
}

func (m *EmbedField) GetInline() bool {
	if m != nil {
		return m.Inline
	}
	return false
}

type Embed struct {
	Url                  string          `protobuf:"bytes,1,opt,name=url,proto3" json:"url,omitempty"`
	Type                 string          `protobuf:"bytes,2,opt,name=type,proto3" json:"type,omitempty"`
	Title                string          `protobuf:"bytes,3,opt,name=title,proto3" json:"title,omitempty"`
	Description          string          `protobuf:"bytes,4,opt,name=description,proto3" json:"description,omitempty"`
	Timestamp            string          `protobuf:"bytes,5,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	Color                int32           `protobuf:"varint,6,opt,name=color,proto3" json:"color,omitempty"`
	Footer               *EmbedFooter    `protobuf:"bytes,7,opt,name=footer,proto3" json:"footer,omitempty"`
	Image                *EmbedImage     `protobuf:"bytes,8,opt,name=image,proto3" json:"image,omitempty"`
	Thumbnail            *EmbedThumbnail `protobuf:"bytes,9,opt,name=thumbnail,proto3" json:"thumbnail,omitempty"`
	Video                *EmbedVideo     `protobuf:"bytes,10,opt,name=video,proto3" json:"video,omitempty"`
	Provider             *EmbedProvider  `protobuf:"bytes,11,opt,name=provider,proto3" json:"provider,omitempty"`
	Author               *EmbedAuthor    `protobuf:"bytes,12,opt,name=author,proto3" json:"author,omitempty"`
	Fields               []*EmbedField   `protobuf:"bytes,13,rep,name=fields,proto3" json:"fields,omitempty"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_unrecognized     []byte          `json:"-"`
	XXX_sizecache        int32           `json:"-"`
}

func (m *Embed) Reset()         { *m = Embed{} }
func (m *Embed) String() string { return proto.CompactTextString(m) }
func (*Embed) ProtoMessage()    {}
func (*Embed) Descriptor() ([]byte, []int) {
	return fileDescriptor_publishing_805a4042085b2f89, []int{12}
}
func (m *Embed) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Embed.Unmarshal(m, b)
}
func (m *Embed) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Embed.Marshal(b, m, deterministic)
}
func (dst *Embed) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Embed.Merge(dst, src)
}
func (m *Embed) XXX_Size() int {
	return xxx_messageInfo_Embed.Size(m)
}
func (m *Embed) XXX_DiscardUnknown() {
	xxx_messageInfo_Embed.DiscardUnknown(m)
}

var xxx_messageInfo_Embed proto.InternalMessageInfo

func (m *Embed) GetUrl() string {
	if m != nil {
		return m.Url
	}
	return ""
}

func (m *Embed) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

func (m *Embed) GetTitle() string {
	if m != nil {
		return m.Title
	}
	return ""
}

func (m *Embed) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func (m *Embed) GetTimestamp() string {
	if m != nil {
		return m.Timestamp
	}
	return ""
}

func (m *Embed) GetColor() int32 {
	if m != nil {
		return m.Color
	}
	return 0
}

func (m *Embed) GetFooter() *EmbedFooter {
	if m != nil {
		return m.Footer
	}
	return nil
}

func (m *Embed) GetImage() *EmbedImage {
	if m != nil {
		return m.Image
	}
	return nil
}

func (m *Embed) GetThumbnail() *EmbedThumbnail {
	if m != nil {
		return m.Thumbnail
	}
	return nil
}

func (m *Embed) GetVideo() *EmbedVideo {
	if m != nil {
		return m.Video
	}
	return nil
}

func (m *Embed) GetProvider() *EmbedProvider {
	if m != nil {
		return m.Provider
	}
	return nil
}

func (m *Embed) GetAuthor() *EmbedAuthor {
	if m != nil {
		return m.Author
	}
	return nil
}

func (m *Embed) GetFields() []*EmbedField {
	if m != nil {
		return m.Fields
	}
	return nil
}

type File struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	ContentType          string   `protobuf:"bytes,2,opt,name=contentType,proto3" json:"contentType,omitempty"`
	Bytes                []byte   `protobuf:"bytes,3,opt,name=bytes,proto3" json:"bytes,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *File) Reset()         { *m = File{} }
func (m *File) String() string { return proto.CompactTextString(m) }
func (*File) ProtoMessage()    {}
func (*File) Descriptor() ([]byte, []int) {
	return fileDescriptor_publishing_805a4042085b2f89, []int{13}
}
func (m *File) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_File.Unmarshal(m, b)
}
func (m *File) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_File.Marshal(b, m, deterministic)
}
func (dst *File) XXX_Merge(src proto.Message) {
	xxx_messageInfo_File.Merge(dst, src)
}
func (m *File) XXX_Size() int {
	return xxx_messageInfo_File.Size(m)
}
func (m *File) XXX_DiscardUnknown() {
	xxx_messageInfo_File.DiscardUnknown(m)
}

var xxx_messageInfo_File proto.InternalMessageInfo

func (m *File) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *File) GetContentType() string {
	if m != nil {
		return m.ContentType
	}
	return ""
}

func (m *File) GetBytes() []byte {
	if m != nil {
		return m.Bytes
	}
	return nil
}

func init() {
	proto.RegisterType((*SendMessageRequest)(nil), "publishing.SendMessageRequest")
	proto.RegisterType((*SendMessageReply)(nil), "publishing.SendMessageReply")
	proto.RegisterType((*Message)(nil), "publishing.Message")
	proto.RegisterType((*Request)(nil), "publishing.Request")
	proto.RegisterType((*User)(nil), "publishing.User")
	proto.RegisterType((*EmbedFooter)(nil), "publishing.EmbedFooter")
	proto.RegisterType((*EmbedImage)(nil), "publishing.EmbedImage")
	proto.RegisterType((*EmbedThumbnail)(nil), "publishing.EmbedThumbnail")
	proto.RegisterType((*EmbedVideo)(nil), "publishing.EmbedVideo")
	proto.RegisterType((*EmbedProvider)(nil), "publishing.EmbedProvider")
	proto.RegisterType((*EmbedAuthor)(nil), "publishing.EmbedAuthor")
	proto.RegisterType((*EmbedField)(nil), "publishing.EmbedField")
	proto.RegisterType((*Embed)(nil), "publishing.Embed")
	proto.RegisterType((*File)(nil), "publishing.File")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// PublishingClient is the client API for Publishing service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type PublishingClient interface {
	SendMessage(ctx context.Context, in *SendMessageRequest, opts ...grpc.CallOption) (*SendMessageReply, error)
}

type publishingClient struct {
	cc *grpc.ClientConn
}

func NewPublishingClient(cc *grpc.ClientConn) PublishingClient {
	return &publishingClient{cc}
}

func (c *publishingClient) SendMessage(ctx context.Context, in *SendMessageRequest, opts ...grpc.CallOption) (*SendMessageReply, error) {
	out := new(SendMessageReply)
	err := c.cc.Invoke(ctx, "/publishing.Publishing/SendMessage", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PublishingServer is the server API for Publishing service.
type PublishingServer interface {
	SendMessage(context.Context, *SendMessageRequest) (*SendMessageReply, error)
}

func RegisterPublishingServer(s *grpc.Server, srv PublishingServer) {
	s.RegisterService(&_Publishing_serviceDesc, srv)
}

func _Publishing_SendMessage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendMessageRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PublishingServer).SendMessage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/publishing.Publishing/SendMessage",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PublishingServer).SendMessage(ctx, req.(*SendMessageRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Publishing_serviceDesc = grpc.ServiceDesc{
	ServiceName: "publishing.Publishing",
	HandlerType: (*PublishingServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SendMessage",
			Handler:    _Publishing_SendMessage_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "publishing/publishing.proto",
}

func init() {
	proto.RegisterFile("publishing/publishing.proto", fileDescriptor_publishing_805a4042085b2f89)
}

var fileDescriptor_publishing_805a4042085b2f89 = []byte{
	// 737 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xb4, 0x55, 0xc1, 0x4e, 0x1b, 0x3d,
	0x10, 0xfe, 0x43, 0xb2, 0x09, 0x99, 0x00, 0xca, 0xef, 0x22, 0xea, 0x52, 0x84, 0xa2, 0x15, 0xa2,
	0xa8, 0x6a, 0x41, 0xa2, 0x42, 0xea, 0xb5, 0x07, 0x90, 0x72, 0xa0, 0x42, 0x2e, 0xf4, 0xd2, 0x43,
	0xb5, 0xc9, 0x1a, 0xd6, 0x92, 0xb3, 0x5e, 0xbc, 0x5e, 0x4a, 0x5e, 0xa0, 0x2f, 0xd5, 0x17, 0xeb,
	0xb1, 0xf2, 0xd8, 0x9b, 0x5d, 0xd8, 0x50, 0xf5, 0xc2, 0xcd, 0x33, 0xfe, 0x66, 0xbe, 0xb1, 0xe7,
	0x1b, 0x1b, 0x5e, 0x67, 0xc5, 0x44, 0x8a, 0x3c, 0x11, 0xe9, 0xcd, 0x51, 0xb5, 0x3c, 0xcc, 0xb4,
	0x32, 0x8a, 0x40, 0xe5, 0x09, 0x35, 0x90, 0x2f, 0x3c, 0x8d, 0xcf, 0x79, 0x9e, 0x47, 0x37, 0x9c,
	0xf1, 0xdb, 0x82, 0xe7, 0x86, 0xbc, 0x87, 0xde, 0xcc, 0x79, 0x68, 0x6b, 0xd4, 0x3a, 0x18, 0x1c,
	0xbf, 0x38, 0xac, 0x65, 0x29, 0xc1, 0x25, 0xc6, 0xc2, 0xb5, 0x8b, 0xa4, 0x2b, 0x4d, 0xb8, 0x4f,
	0xca, 0x4a, 0x4c, 0x18, 0xc2, 0xf0, 0x01, 0x67, 0x26, 0xe7, 0x64, 0x03, 0x56, 0x44, 0x8c, 0x64,
	0x7d, 0xb6, 0x22, 0xe2, 0xf0, 0x67, 0x0b, 0x7a, 0x1e, 0x40, 0x28, 0xf4, 0xa6, 0x2a, 0x35, 0x3c,
	0x35, 0x1e, 0x50, 0x9a, 0xe4, 0x0d, 0x04, 0x7c, 0x36, 0xe1, 0xb1, 0xa7, 0xfd, 0xbf, 0x4e, 0x7b,
	0x6a, 0x37, 0x98, 0xdb, 0x27, 0x43, 0x68, 0x1b, 0x93, 0xd3, 0xf6, 0xa8, 0x75, 0xb0, 0xca, 0xec,
	0x92, 0xec, 0x43, 0x70, 0x2d, 0x24, 0xcf, 0x69, 0x67, 0xd4, 0x3e, 0x18, 0x1c, 0x0f, 0xeb, 0xa1,
	0x67, 0x42, 0x72, 0xe6, 0xb6, 0xc3, 0x5f, 0x2d, 0xe8, 0x95, 0xd7, 0xb2, 0x07, 0xeb, 0x53, 0xa5,
	0x35, 0x97, 0x91, 0x11, 0x2a, 0x1d, 0x97, 0xf5, 0x3e, 0x74, 0x92, 0x1d, 0xe8, 0x4f, 0x93, 0x28,
	0x4d, 0xb9, 0x1c, 0xbb, 0xc2, 0xfa, 0xac, 0x72, 0x90, 0x3d, 0xe8, 0x14, 0x39, 0xd7, 0x58, 0xca,
	0x23, 0xda, 0xab, 0x9c, 0x6b, 0x86, 0xbb, 0x64, 0x1f, 0x36, 0x32, 0x2d, 0xee, 0x22, 0xc3, 0xfd,
	0x25, 0xd0, 0x0e, 0x96, 0xfe, 0xc8, 0x4b, 0x76, 0x01, 0x44, 0xce, 0x78, 0x2c, 0x34, 0x9f, 0x1a,
	0x1a, 0x20, 0xa6, 0xe6, 0x09, 0xdf, 0x42, 0xc7, 0x66, 0x7d, 0x7c, 0xbd, 0x84, 0x40, 0x27, 0x8d,
	0x66, 0xdc, 0x97, 0x87, 0xeb, 0xf0, 0x3b, 0x0c, 0xf0, 0xce, 0xce, 0x94, 0x32, 0x5c, 0x5b, 0x88,
	0xe1, 0xf7, 0xe5, 0x95, 0xe3, 0xda, 0x76, 0x42, 0x4c, 0x55, 0x7a, 0xa5, 0xa5, 0x8f, 0x2c, 0x4d,
	0x12, 0xc2, 0x5a, 0xa6, 0xd5, 0xfd, 0x7c, 0xec, 0xb7, 0xdb, 0xb8, 0xfd, 0xc0, 0x17, 0x26, 0x00,
	0x48, 0x30, 0x9e, 0xd9, 0xd2, 0x87, 0xd0, 0x2e, 0xb4, 0xf4, 0xe9, 0xed, 0x92, 0x6c, 0xc3, 0x2a,
	0xe2, 0xab, 0xf4, 0x0b, 0x9b, 0x6c, 0x42, 0xf0, 0x43, 0xc4, 0x26, 0xc1, 0xc4, 0x01, 0x73, 0x06,
	0xd9, 0x82, 0x6e, 0xc2, 0xc5, 0x4d, 0x62, 0xf0, 0x7a, 0x02, 0xe6, 0xad, 0x50, 0xc2, 0x06, 0x32,
	0x5d, 0x26, 0xc5, 0x6c, 0x92, 0x46, 0x42, 0x3e, 0x2b, 0x5b, 0x79, 0xae, 0xaf, 0x22, 0xe6, 0xea,
	0x59, 0x99, 0x4e, 0x60, 0x1d, 0x99, 0x2e, 0xb4, 0xba, 0x13, 0x31, 0xd7, 0x4b, 0xc8, 0x96, 0x75,
	0xf6, 0xd6, 0x77, 0xf6, 0x53, 0x61, 0x12, 0xf5, 0x8f, 0x41, 0xf5, 0x5e, 0xb7, 0xff, 0xde, 0xeb,
	0xce, 0x92, 0x5e, 0x7f, 0xf6, 0x77, 0x72, 0x26, 0xb8, 0xac, 0xe4, 0xd6, 0xaa, 0xe5, 0xdf, 0x84,
	0xe0, 0x2e, 0x92, 0x45, 0x49, 0xea, 0x0c, 0x7b, 0x72, 0x91, 0x4a, 0x91, 0x72, 0x3f, 0xab, 0xde,
	0x0a, 0x7f, 0xb7, 0x21, 0x38, 0x2d, 0x47, 0xb9, 0x59, 0xbd, 0x99, 0x67, 0x8b, 0xea, 0xed, 0xda,
	0x66, 0x37, 0xc2, 0x48, 0xee, 0x6b, 0x77, 0x06, 0x19, 0xc1, 0x20, 0xe6, 0xf9, 0x54, 0x8b, 0xcc,
	0xce, 0xaa, 0x2f, 0xbc, 0xee, 0xb2, 0xc3, 0x6b, 0xc4, 0x8c, 0xe7, 0x26, 0x9a, 0x65, 0x38, 0x4f,
	0x7d, 0x56, 0x39, 0x6c, 0xd6, 0xa9, 0x92, 0x4a, 0xd3, 0xae, 0xeb, 0x16, 0x1a, 0xe4, 0x08, 0xba,
	0xd7, 0x38, 0x33, 0xb4, 0x87, 0x43, 0xfd, 0xb2, 0xf1, 0x0c, 0xb9, 0x91, 0x62, 0x1e, 0x46, 0xde,
	0x41, 0x20, 0xec, 0x0c, 0xd0, 0x55, 0xc4, 0x6f, 0x35, 0xf0, 0x38, 0x21, 0xcc, 0x81, 0xc8, 0x47,
	0xe8, 0x9b, 0x52, 0xc7, 0xb4, 0x8f, 0x11, 0xdb, 0x8d, 0x88, 0x85, 0xd2, 0x59, 0x05, 0xb6, 0x3c,
	0x56, 0x26, 0x8a, 0xc2, 0x13, 0x3c, 0xa8, 0x58, 0xe6, 0x40, 0xe4, 0x04, 0x65, 0x8a, 0xba, 0xa2,
	0x03, 0x0c, 0x78, 0xd5, 0x08, 0x28, 0x85, 0xc7, 0x16, 0x50, 0x7b, 0xfa, 0x08, 0x75, 0x45, 0xd7,
	0x9e, 0x38, 0xbd, 0x93, 0x1d, 0xf3, 0x30, 0x72, 0x08, 0xdd, 0x6b, 0xab, 0x8a, 0x9c, 0xae, 0xe3,
	0xd3, 0xdb, 0x2c, 0x0b, 0x45, 0xc3, 0x3c, 0x2a, 0x64, 0xd0, 0xb1, 0x0f, 0xf2, 0x52, 0x11, 0x8d,
	0x60, 0xe0, 0xff, 0x82, 0xcb, 0x4a, 0x01, 0x75, 0x97, 0x6d, 0xd9, 0x64, 0x6e, 0xb8, 0x7b, 0xfb,
	0xd7, 0x98, 0x33, 0x8e, 0xbf, 0x01, 0x5c, 0x2c, 0x48, 0xc9, 0x39, 0x0c, 0x6a, 0x1f, 0x12, 0xd9,
	0xad, 0x17, 0xd4, 0xfc, 0x1d, 0xb7, 0x77, 0x9e, 0xdc, 0xcf, 0xe4, 0x3c, 0xfc, 0x6f, 0xd2, 0xc5,
	0x6f, 0xf6, 0xc3, 0x9f, 0x00, 0x00, 0x00, 0xff, 0xff, 0x8b, 0x12, 0xb2, 0xff, 0x85, 0x07, 0x00,
	0x00,
}
