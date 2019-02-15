// Code generated by protoc-gen-go. DO NOT EDIT.
// source: scheduling.proto

package schedulingpb

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

type Schedule struct {
	Minute               string   `protobuf:"bytes,1,opt,name=minute,proto3" json:"minute,omitempty"`
	Hour                 string   `protobuf:"bytes,2,opt,name=hour,proto3" json:"hour,omitempty"`
	DayOfMonth           string   `protobuf:"bytes,3,opt,name=dayOfMonth,proto3" json:"dayOfMonth,omitempty"`
	Month                string   `protobuf:"bytes,4,opt,name=month,proto3" json:"month,omitempty"`
	DayOfWeek            string   `protobuf:"bytes,5,opt,name=dayOfWeek,proto3" json:"dayOfWeek,omitempty"`
	Year                 string   `protobuf:"bytes,6,opt,name=year,proto3" json:"year,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Schedule) Reset()         { *m = Schedule{} }
func (m *Schedule) String() string { return proto.CompactTextString(m) }
func (*Schedule) ProtoMessage()    {}
func (*Schedule) Descriptor() ([]byte, []int) {
	return fileDescriptor_scheduling_ade178e3614cdfac, []int{0}
}
func (m *Schedule) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Schedule.Unmarshal(m, b)
}
func (m *Schedule) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Schedule.Marshal(b, m, deterministic)
}
func (dst *Schedule) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Schedule.Merge(dst, src)
}
func (m *Schedule) XXX_Size() int {
	return xxx_messageInfo_Schedule.Size(m)
}
func (m *Schedule) XXX_DiscardUnknown() {
	xxx_messageInfo_Schedule.DiscardUnknown(m)
}

var xxx_messageInfo_Schedule proto.InternalMessageInfo

func (m *Schedule) GetMinute() string {
	if m != nil {
		return m.Minute
	}
	return ""
}

func (m *Schedule) GetHour() string {
	if m != nil {
		return m.Hour
	}
	return ""
}

func (m *Schedule) GetDayOfMonth() string {
	if m != nil {
		return m.DayOfMonth
	}
	return ""
}

func (m *Schedule) GetMonth() string {
	if m != nil {
		return m.Month
	}
	return ""
}

func (m *Schedule) GetDayOfWeek() string {
	if m != nil {
		return m.DayOfWeek
	}
	return ""
}

func (m *Schedule) GetYear() string {
	if m != nil {
		return m.Year
	}
	return ""
}

type Job struct {
	Id                   string    `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	UserId               string    `protobuf:"bytes,2,opt,name=userId,proto3" json:"userId,omitempty"`
	Source               string    `protobuf:"bytes,3,opt,name=source,proto3" json:"source,omitempty"`
	Content              string    `protobuf:"bytes,4,opt,name=content,proto3" json:"content,omitempty"`
	Endpoint             string    `protobuf:"bytes,5,opt,name=endpoint,proto3" json:"endpoint,omitempty"`
	Destination          string    `protobuf:"bytes,6,opt,name=destination,proto3" json:"destination,omitempty"`
	Schedule             *Schedule `protobuf:"bytes,7,opt,name=schedule,proto3" json:"schedule,omitempty"`
	Enabled              bool      `protobuf:"varint,8,opt,name=enabled,proto3" json:"enabled,omitempty"`
	Created              int32     `protobuf:"varint,9,opt,name=created,proto3" json:"created,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *Job) Reset()         { *m = Job{} }
func (m *Job) String() string { return proto.CompactTextString(m) }
func (*Job) ProtoMessage()    {}
func (*Job) Descriptor() ([]byte, []int) {
	return fileDescriptor_scheduling_ade178e3614cdfac, []int{1}
}
func (m *Job) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Job.Unmarshal(m, b)
}
func (m *Job) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Job.Marshal(b, m, deterministic)
}
func (dst *Job) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Job.Merge(dst, src)
}
func (m *Job) XXX_Size() int {
	return xxx_messageInfo_Job.Size(m)
}
func (m *Job) XXX_DiscardUnknown() {
	xxx_messageInfo_Job.DiscardUnknown(m)
}

var xxx_messageInfo_Job proto.InternalMessageInfo

func (m *Job) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Job) GetUserId() string {
	if m != nil {
		return m.UserId
	}
	return ""
}

func (m *Job) GetSource() string {
	if m != nil {
		return m.Source
	}
	return ""
}

func (m *Job) GetContent() string {
	if m != nil {
		return m.Content
	}
	return ""
}

func (m *Job) GetEndpoint() string {
	if m != nil {
		return m.Endpoint
	}
	return ""
}

func (m *Job) GetDestination() string {
	if m != nil {
		return m.Destination
	}
	return ""
}

func (m *Job) GetSchedule() *Schedule {
	if m != nil {
		return m.Schedule
	}
	return nil
}

func (m *Job) GetEnabled() bool {
	if m != nil {
		return m.Enabled
	}
	return false
}

func (m *Job) GetCreated() int32 {
	if m != nil {
		return m.Created
	}
	return 0
}

type CreateJobRequest struct {
	UserId               string    `protobuf:"bytes,1,opt,name=userId,proto3" json:"userId,omitempty"`
	Source               string    `protobuf:"bytes,2,opt,name=source,proto3" json:"source,omitempty"`
	Content              string    `protobuf:"bytes,3,opt,name=content,proto3" json:"content,omitempty"`
	Endpoint             string    `protobuf:"bytes,4,opt,name=endpoint,proto3" json:"endpoint,omitempty"`
	Destination          string    `protobuf:"bytes,5,opt,name=destination,proto3" json:"destination,omitempty"`
	Schedule             *Schedule `protobuf:"bytes,6,opt,name=schedule,proto3" json:"schedule,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *CreateJobRequest) Reset()         { *m = CreateJobRequest{} }
func (m *CreateJobRequest) String() string { return proto.CompactTextString(m) }
func (*CreateJobRequest) ProtoMessage()    {}
func (*CreateJobRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_scheduling_ade178e3614cdfac, []int{2}
}
func (m *CreateJobRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CreateJobRequest.Unmarshal(m, b)
}
func (m *CreateJobRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CreateJobRequest.Marshal(b, m, deterministic)
}
func (dst *CreateJobRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CreateJobRequest.Merge(dst, src)
}
func (m *CreateJobRequest) XXX_Size() int {
	return xxx_messageInfo_CreateJobRequest.Size(m)
}
func (m *CreateJobRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_CreateJobRequest.DiscardUnknown(m)
}

var xxx_messageInfo_CreateJobRequest proto.InternalMessageInfo

func (m *CreateJobRequest) GetUserId() string {
	if m != nil {
		return m.UserId
	}
	return ""
}

func (m *CreateJobRequest) GetSource() string {
	if m != nil {
		return m.Source
	}
	return ""
}

func (m *CreateJobRequest) GetContent() string {
	if m != nil {
		return m.Content
	}
	return ""
}

func (m *CreateJobRequest) GetEndpoint() string {
	if m != nil {
		return m.Endpoint
	}
	return ""
}

func (m *CreateJobRequest) GetDestination() string {
	if m != nil {
		return m.Destination
	}
	return ""
}

func (m *CreateJobRequest) GetSchedule() *Schedule {
	if m != nil {
		return m.Schedule
	}
	return nil
}

type CreateJobReply struct {
	Job                  *Job     `protobuf:"bytes,1,opt,name=job,proto3" json:"job,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CreateJobReply) Reset()         { *m = CreateJobReply{} }
func (m *CreateJobReply) String() string { return proto.CompactTextString(m) }
func (*CreateJobReply) ProtoMessage()    {}
func (*CreateJobReply) Descriptor() ([]byte, []int) {
	return fileDescriptor_scheduling_ade178e3614cdfac, []int{3}
}
func (m *CreateJobReply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CreateJobReply.Unmarshal(m, b)
}
func (m *CreateJobReply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CreateJobReply.Marshal(b, m, deterministic)
}
func (dst *CreateJobReply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CreateJobReply.Merge(dst, src)
}
func (m *CreateJobReply) XXX_Size() int {
	return xxx_messageInfo_CreateJobReply.Size(m)
}
func (m *CreateJobReply) XXX_DiscardUnknown() {
	xxx_messageInfo_CreateJobReply.DiscardUnknown(m)
}

var xxx_messageInfo_CreateJobReply proto.InternalMessageInfo

func (m *CreateJobReply) GetJob() *Job {
	if m != nil {
		return m.Job
	}
	return nil
}

type JobsRequest struct {
	Time                 int64    `protobuf:"varint,1,opt,name=time,proto3" json:"time,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *JobsRequest) Reset()         { *m = JobsRequest{} }
func (m *JobsRequest) String() string { return proto.CompactTextString(m) }
func (*JobsRequest) ProtoMessage()    {}
func (*JobsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_scheduling_ade178e3614cdfac, []int{4}
}
func (m *JobsRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_JobsRequest.Unmarshal(m, b)
}
func (m *JobsRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_JobsRequest.Marshal(b, m, deterministic)
}
func (dst *JobsRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_JobsRequest.Merge(dst, src)
}
func (m *JobsRequest) XXX_Size() int {
	return xxx_messageInfo_JobsRequest.Size(m)
}
func (m *JobsRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_JobsRequest.DiscardUnknown(m)
}

var xxx_messageInfo_JobsRequest proto.InternalMessageInfo

func (m *JobsRequest) GetTime() int64 {
	if m != nil {
		return m.Time
	}
	return 0
}

type JobsReply struct {
	Jobs                 []*Job   `protobuf:"bytes,1,rep,name=jobs,proto3" json:"jobs,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *JobsReply) Reset()         { *m = JobsReply{} }
func (m *JobsReply) String() string { return proto.CompactTextString(m) }
func (*JobsReply) ProtoMessage()    {}
func (*JobsReply) Descriptor() ([]byte, []int) {
	return fileDescriptor_scheduling_ade178e3614cdfac, []int{5}
}
func (m *JobsReply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_JobsReply.Unmarshal(m, b)
}
func (m *JobsReply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_JobsReply.Marshal(b, m, deterministic)
}
func (dst *JobsReply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_JobsReply.Merge(dst, src)
}
func (m *JobsReply) XXX_Size() int {
	return xxx_messageInfo_JobsReply.Size(m)
}
func (m *JobsReply) XXX_DiscardUnknown() {
	xxx_messageInfo_JobsReply.DiscardUnknown(m)
}

var xxx_messageInfo_JobsReply proto.InternalMessageInfo

func (m *JobsReply) GetJobs() []*Job {
	if m != nil {
		return m.Jobs
	}
	return nil
}

type RunJobsRequest struct {
	Time                 int64    `protobuf:"varint,1,opt,name=time,proto3" json:"time,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RunJobsRequest) Reset()         { *m = RunJobsRequest{} }
func (m *RunJobsRequest) String() string { return proto.CompactTextString(m) }
func (*RunJobsRequest) ProtoMessage()    {}
func (*RunJobsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_scheduling_ade178e3614cdfac, []int{6}
}
func (m *RunJobsRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RunJobsRequest.Unmarshal(m, b)
}
func (m *RunJobsRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RunJobsRequest.Marshal(b, m, deterministic)
}
func (dst *RunJobsRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RunJobsRequest.Merge(dst, src)
}
func (m *RunJobsRequest) XXX_Size() int {
	return xxx_messageInfo_RunJobsRequest.Size(m)
}
func (m *RunJobsRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_RunJobsRequest.DiscardUnknown(m)
}

var xxx_messageInfo_RunJobsRequest proto.InternalMessageInfo

func (m *RunJobsRequest) GetTime() int64 {
	if m != nil {
		return m.Time
	}
	return 0
}

type RunJobsReply struct {
	Jobs                 []*Job   `protobuf:"bytes,1,rep,name=jobs,proto3" json:"jobs,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RunJobsReply) Reset()         { *m = RunJobsReply{} }
func (m *RunJobsReply) String() string { return proto.CompactTextString(m) }
func (*RunJobsReply) ProtoMessage()    {}
func (*RunJobsReply) Descriptor() ([]byte, []int) {
	return fileDescriptor_scheduling_ade178e3614cdfac, []int{7}
}
func (m *RunJobsReply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RunJobsReply.Unmarshal(m, b)
}
func (m *RunJobsReply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RunJobsReply.Marshal(b, m, deterministic)
}
func (dst *RunJobsReply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RunJobsReply.Merge(dst, src)
}
func (m *RunJobsReply) XXX_Size() int {
	return xxx_messageInfo_RunJobsReply.Size(m)
}
func (m *RunJobsReply) XXX_DiscardUnknown() {
	xxx_messageInfo_RunJobsReply.DiscardUnknown(m)
}

var xxx_messageInfo_RunJobsReply proto.InternalMessageInfo

func (m *RunJobsReply) GetJobs() []*Job {
	if m != nil {
		return m.Jobs
	}
	return nil
}

type DisableJobRequest struct {
	Id                   string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DisableJobRequest) Reset()         { *m = DisableJobRequest{} }
func (m *DisableJobRequest) String() string { return proto.CompactTextString(m) }
func (*DisableJobRequest) ProtoMessage()    {}
func (*DisableJobRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_scheduling_ade178e3614cdfac, []int{8}
}
func (m *DisableJobRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DisableJobRequest.Unmarshal(m, b)
}
func (m *DisableJobRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DisableJobRequest.Marshal(b, m, deterministic)
}
func (dst *DisableJobRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DisableJobRequest.Merge(dst, src)
}
func (m *DisableJobRequest) XXX_Size() int {
	return xxx_messageInfo_DisableJobRequest.Size(m)
}
func (m *DisableJobRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_DisableJobRequest.DiscardUnknown(m)
}

var xxx_messageInfo_DisableJobRequest proto.InternalMessageInfo

func (m *DisableJobRequest) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

type DisableJobReply struct {
	Id                   string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Enabled              bool     `protobuf:"varint,2,opt,name=enabled,proto3" json:"enabled,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DisableJobReply) Reset()         { *m = DisableJobReply{} }
func (m *DisableJobReply) String() string { return proto.CompactTextString(m) }
func (*DisableJobReply) ProtoMessage()    {}
func (*DisableJobReply) Descriptor() ([]byte, []int) {
	return fileDescriptor_scheduling_ade178e3614cdfac, []int{9}
}
func (m *DisableJobReply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DisableJobReply.Unmarshal(m, b)
}
func (m *DisableJobReply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DisableJobReply.Marshal(b, m, deterministic)
}
func (dst *DisableJobReply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DisableJobReply.Merge(dst, src)
}
func (m *DisableJobReply) XXX_Size() int {
	return xxx_messageInfo_DisableJobReply.Size(m)
}
func (m *DisableJobReply) XXX_DiscardUnknown() {
	xxx_messageInfo_DisableJobReply.DiscardUnknown(m)
}

var xxx_messageInfo_DisableJobReply proto.InternalMessageInfo

func (m *DisableJobReply) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *DisableJobReply) GetEnabled() bool {
	if m != nil {
		return m.Enabled
	}
	return false
}

func init() {
	proto.RegisterType((*Schedule)(nil), "schedulingpb.Schedule")
	proto.RegisterType((*Job)(nil), "schedulingpb.Job")
	proto.RegisterType((*CreateJobRequest)(nil), "schedulingpb.CreateJobRequest")
	proto.RegisterType((*CreateJobReply)(nil), "schedulingpb.CreateJobReply")
	proto.RegisterType((*JobsRequest)(nil), "schedulingpb.JobsRequest")
	proto.RegisterType((*JobsReply)(nil), "schedulingpb.JobsReply")
	proto.RegisterType((*RunJobsRequest)(nil), "schedulingpb.RunJobsRequest")
	proto.RegisterType((*RunJobsReply)(nil), "schedulingpb.RunJobsReply")
	proto.RegisterType((*DisableJobRequest)(nil), "schedulingpb.DisableJobRequest")
	proto.RegisterType((*DisableJobReply)(nil), "schedulingpb.DisableJobReply")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// SchedulingClient is the client API for Scheduling service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type SchedulingClient interface {
	Create(ctx context.Context, in *CreateJobRequest, opts ...grpc.CallOption) (*CreateJobReply, error)
	Jobs(ctx context.Context, in *JobsRequest, opts ...grpc.CallOption) (*JobsReply, error)
	Run(ctx context.Context, in *RunJobsRequest, opts ...grpc.CallOption) (*RunJobsReply, error)
	Disable(ctx context.Context, in *DisableJobRequest, opts ...grpc.CallOption) (*DisableJobReply, error)
}

type schedulingClient struct {
	cc *grpc.ClientConn
}

func NewSchedulingClient(cc *grpc.ClientConn) SchedulingClient {
	return &schedulingClient{cc}
}

func (c *schedulingClient) Create(ctx context.Context, in *CreateJobRequest, opts ...grpc.CallOption) (*CreateJobReply, error) {
	out := new(CreateJobReply)
	err := c.cc.Invoke(ctx, "/schedulingpb.Scheduling/Create", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *schedulingClient) Jobs(ctx context.Context, in *JobsRequest, opts ...grpc.CallOption) (*JobsReply, error) {
	out := new(JobsReply)
	err := c.cc.Invoke(ctx, "/schedulingpb.Scheduling/Jobs", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *schedulingClient) Run(ctx context.Context, in *RunJobsRequest, opts ...grpc.CallOption) (*RunJobsReply, error) {
	out := new(RunJobsReply)
	err := c.cc.Invoke(ctx, "/schedulingpb.Scheduling/Run", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *schedulingClient) Disable(ctx context.Context, in *DisableJobRequest, opts ...grpc.CallOption) (*DisableJobReply, error) {
	out := new(DisableJobReply)
	err := c.cc.Invoke(ctx, "/schedulingpb.Scheduling/Disable", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SchedulingServer is the server API for Scheduling service.
type SchedulingServer interface {
	Create(context.Context, *CreateJobRequest) (*CreateJobReply, error)
	Jobs(context.Context, *JobsRequest) (*JobsReply, error)
	Run(context.Context, *RunJobsRequest) (*RunJobsReply, error)
	Disable(context.Context, *DisableJobRequest) (*DisableJobReply, error)
}

func RegisterSchedulingServer(s *grpc.Server, srv SchedulingServer) {
	s.RegisterService(&_Scheduling_serviceDesc, srv)
}

func _Scheduling_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateJobRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SchedulingServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/schedulingpb.Scheduling/Create",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SchedulingServer).Create(ctx, req.(*CreateJobRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Scheduling_Jobs_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(JobsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SchedulingServer).Jobs(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/schedulingpb.Scheduling/Jobs",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SchedulingServer).Jobs(ctx, req.(*JobsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Scheduling_Run_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RunJobsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SchedulingServer).Run(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/schedulingpb.Scheduling/Run",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SchedulingServer).Run(ctx, req.(*RunJobsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Scheduling_Disable_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DisableJobRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SchedulingServer).Disable(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/schedulingpb.Scheduling/Disable",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SchedulingServer).Disable(ctx, req.(*DisableJobRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Scheduling_serviceDesc = grpc.ServiceDesc{
	ServiceName: "schedulingpb.Scheduling",
	HandlerType: (*SchedulingServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Create",
			Handler:    _Scheduling_Create_Handler,
		},
		{
			MethodName: "Jobs",
			Handler:    _Scheduling_Jobs_Handler,
		},
		{
			MethodName: "Run",
			Handler:    _Scheduling_Run_Handler,
		},
		{
			MethodName: "Disable",
			Handler:    _Scheduling_Disable_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "scheduling.proto",
}

func init() { proto.RegisterFile("scheduling.proto", fileDescriptor_scheduling_ade178e3614cdfac) }

var fileDescriptor_scheduling_ade178e3614cdfac = []byte{
	// 515 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x94, 0xcf, 0x6e, 0xd3, 0x40,
	0x10, 0xc6, 0xf1, 0x9f, 0xb8, 0xce, 0xa4, 0x0a, 0xed, 0x08, 0x95, 0xc5, 0x2a, 0xc5, 0xb8, 0x20,
	0xe5, 0x94, 0x83, 0x51, 0x4f, 0x70, 0x41, 0x70, 0x80, 0x20, 0x84, 0xe4, 0x1e, 0x38, 0xdb, 0xf1,
	0x42, 0xb6, 0x24, 0xbb, 0xc6, 0x5e, 0x1f, 0xf2, 0x00, 0xbc, 0x06, 0x27, 0x9e, 0x87, 0x67, 0x42,
	0xbb, 0x5e, 0xa7, 0x76, 0x5a, 0x17, 0xb8, 0xed, 0xcc, 0x37, 0xf9, 0xf6, 0xdb, 0x9f, 0x47, 0x81,
	0xa3, 0x6a, 0xb9, 0xa2, 0x79, 0xbd, 0x66, 0xfc, 0xeb, 0xbc, 0x28, 0x85, 0x14, 0x78, 0x78, 0xdd,
	0x29, 0xb2, 0xe8, 0xa7, 0x05, 0xfe, 0x65, 0xd3, 0xa0, 0x78, 0x02, 0xde, 0x86, 0xf1, 0x5a, 0x52,
	0x62, 0x85, 0xd6, 0x6c, 0x9c, 0x98, 0x0a, 0x11, 0xdc, 0x95, 0xa8, 0x4b, 0x62, 0xeb, 0xae, 0x3e,
	0xe3, 0x19, 0x40, 0x9e, 0x6e, 0x3f, 0x7d, 0xf9, 0x28, 0xb8, 0x5c, 0x11, 0x47, 0x2b, 0x9d, 0x0e,
	0x3e, 0x80, 0xd1, 0x46, 0x4b, 0xae, 0x96, 0x9a, 0x02, 0x4f, 0x61, 0xac, 0x67, 0x3e, 0x53, 0xfa,
	0x8d, 0x8c, 0xb4, 0x72, 0xdd, 0x50, 0xf7, 0x6c, 0x69, 0x5a, 0x12, 0xaf, 0xb9, 0x47, 0x9d, 0xa3,
	0x1f, 0x36, 0x38, 0x0b, 0x91, 0xe1, 0x14, 0x6c, 0x96, 0x9b, 0x5c, 0x36, 0xcb, 0x55, 0xd6, 0xba,
	0xa2, 0xe5, 0xfb, 0xdc, 0xa4, 0x32, 0x95, 0xea, 0x57, 0xa2, 0x2e, 0x97, 0xd4, 0x64, 0x32, 0x15,
	0x12, 0x38, 0x58, 0x0a, 0x2e, 0x29, 0x97, 0x26, 0x51, 0x5b, 0x62, 0x00, 0x3e, 0xe5, 0x79, 0x21,
	0x18, 0x97, 0x26, 0xd2, 0xae, 0xc6, 0x10, 0x26, 0x39, 0xad, 0x24, 0xe3, 0xa9, 0x64, 0x82, 0x9b,
	0x60, 0xdd, 0x16, 0xc6, 0xe0, 0x1b, 0xa0, 0x94, 0x1c, 0x84, 0xd6, 0x6c, 0x12, 0x9f, 0xcc, 0xbb,
	0x84, 0xe7, 0x2d, 0xdd, 0x64, 0x37, 0xa7, 0xb2, 0x50, 0x9e, 0x66, 0x6b, 0x9a, 0x13, 0x3f, 0xb4,
	0x66, 0x7e, 0xd2, 0x96, 0x3a, 0x65, 0x49, 0x53, 0x49, 0x73, 0x32, 0x0e, 0xad, 0xd9, 0x28, 0x69,
	0xcb, 0xe8, 0xb7, 0x05, 0x47, 0x6f, 0xf4, 0x79, 0x21, 0xb2, 0x84, 0x7e, 0xaf, 0x69, 0x25, 0x3b,
	0x10, 0xac, 0x01, 0x08, 0xf6, 0x10, 0x04, 0x67, 0x18, 0x82, 0x7b, 0x37, 0x84, 0xd1, 0xdd, 0x10,
	0xbc, 0x7f, 0x83, 0x10, 0x5d, 0xc0, 0xb4, 0xf3, 0x9e, 0x62, 0xbd, 0xc5, 0x73, 0x70, 0xae, 0x44,
	0xa6, 0x9f, 0x32, 0x89, 0x8f, 0xfb, 0x06, 0x6a, 0x48, 0xa9, 0xd1, 0x53, 0x98, 0x2c, 0x44, 0x56,
	0xb5, 0x04, 0x10, 0x5c, 0xc9, 0x36, 0xcd, 0xc2, 0x3a, 0x89, 0x3e, 0x47, 0x31, 0x8c, 0x9b, 0x11,
	0x65, 0xfa, 0x1c, 0xdc, 0x2b, 0x91, 0x55, 0xc4, 0x0a, 0x9d, 0xdb, 0x5d, 0xb5, 0x1c, 0x3d, 0x83,
	0x69, 0x52, 0xf3, 0xbf, 0x39, 0x5f, 0xc0, 0xe1, 0x6e, 0xea, 0x3f, 0xcc, 0xcf, 0xe1, 0xf8, 0x2d,
	0xab, 0xd4, 0x17, 0xee, 0x7c, 0xbb, 0xbd, 0x85, 0x8e, 0x5e, 0xc2, 0xfd, 0xee, 0x90, 0xb2, 0xdf,
	0xdf, 0xf9, 0xce, 0xde, 0xd8, 0xbd, 0xbd, 0x89, 0x7f, 0xd9, 0x00, 0x97, 0xbb, 0xcb, 0xf1, 0x1d,
	0x78, 0x0d, 0x5b, 0x3c, 0xeb, 0x67, 0xda, 0xdf, 0xa0, 0xe0, 0x74, 0x50, 0x2f, 0xd6, 0xdb, 0xe8,
	0x1e, 0xbe, 0x02, 0x57, 0x3d, 0x17, 0x1f, 0xdd, 0x78, 0x5b, 0x0b, 0x2a, 0x78, 0x78, 0x9b, 0xd4,
	0xfc, 0xfa, 0x35, 0x38, 0x49, 0xcd, 0x71, 0xef, 0x92, 0x3e, 0xe8, 0x20, 0x18, 0x50, 0x1b, 0x8b,
	0x0f, 0x70, 0x60, 0xb0, 0xe0, 0x93, 0xfe, 0xe0, 0x0d, 0xa4, 0xc1, 0xe3, 0xe1, 0x01, 0x6d, 0x96,
	0x79, 0xfa, 0x2f, 0xf0, 0xc5, 0x9f, 0x00, 0x00, 0x00, 0xff, 0xff, 0xf6, 0x36, 0xe2, 0x63, 0x16,
	0x05, 0x00, 0x00,
}