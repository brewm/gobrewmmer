// Code generated by protoc-gen-go. DO NOT EDIT.
// source: api/brewmmer/brew.proto

package brewmmer

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import timestamp "github.com/golang/protobuf/ptypes/timestamp"

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

type StartBrewRequest struct {
	RecipeId             int64    `protobuf:"varint,1,opt,name=recipeId,proto3" json:"recipeId,omitempty"`
	Note                 string   `protobuf:"bytes,2,opt,name=note,proto3" json:"note,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *StartBrewRequest) Reset()         { *m = StartBrewRequest{} }
func (m *StartBrewRequest) String() string { return proto.CompactTextString(m) }
func (*StartBrewRequest) ProtoMessage()    {}
func (*StartBrewRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_brew_1e03cfd08bd125fc, []int{0}
}
func (m *StartBrewRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StartBrewRequest.Unmarshal(m, b)
}
func (m *StartBrewRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StartBrewRequest.Marshal(b, m, deterministic)
}
func (dst *StartBrewRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StartBrewRequest.Merge(dst, src)
}
func (m *StartBrewRequest) XXX_Size() int {
	return xxx_messageInfo_StartBrewRequest.Size(m)
}
func (m *StartBrewRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_StartBrewRequest.DiscardUnknown(m)
}

var xxx_messageInfo_StartBrewRequest proto.InternalMessageInfo

func (m *StartBrewRequest) GetRecipeId() int64 {
	if m != nil {
		return m.RecipeId
	}
	return 0
}

func (m *StartBrewRequest) GetNote() string {
	if m != nil {
		return m.Note
	}
	return ""
}

type StartBrewResponse struct {
	Id                   int64    `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *StartBrewResponse) Reset()         { *m = StartBrewResponse{} }
func (m *StartBrewResponse) String() string { return proto.CompactTextString(m) }
func (*StartBrewResponse) ProtoMessage()    {}
func (*StartBrewResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_brew_1e03cfd08bd125fc, []int{1}
}
func (m *StartBrewResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StartBrewResponse.Unmarshal(m, b)
}
func (m *StartBrewResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StartBrewResponse.Marshal(b, m, deterministic)
}
func (dst *StartBrewResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StartBrewResponse.Merge(dst, src)
}
func (m *StartBrewResponse) XXX_Size() int {
	return xxx_messageInfo_StartBrewResponse.Size(m)
}
func (m *StartBrewResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_StartBrewResponse.DiscardUnknown(m)
}

var xxx_messageInfo_StartBrewResponse proto.InternalMessageInfo

func (m *StartBrewResponse) GetId() int64 {
	if m != nil {
		return m.Id
	}
	return 0
}

type CompleteBrewStepRequest struct {
	Id                   int64    `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CompleteBrewStepRequest) Reset()         { *m = CompleteBrewStepRequest{} }
func (m *CompleteBrewStepRequest) String() string { return proto.CompactTextString(m) }
func (*CompleteBrewStepRequest) ProtoMessage()    {}
func (*CompleteBrewStepRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_brew_1e03cfd08bd125fc, []int{2}
}
func (m *CompleteBrewStepRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CompleteBrewStepRequest.Unmarshal(m, b)
}
func (m *CompleteBrewStepRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CompleteBrewStepRequest.Marshal(b, m, deterministic)
}
func (dst *CompleteBrewStepRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CompleteBrewStepRequest.Merge(dst, src)
}
func (m *CompleteBrewStepRequest) XXX_Size() int {
	return xxx_messageInfo_CompleteBrewStepRequest.Size(m)
}
func (m *CompleteBrewStepRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_CompleteBrewStepRequest.DiscardUnknown(m)
}

var xxx_messageInfo_CompleteBrewStepRequest proto.InternalMessageInfo

func (m *CompleteBrewStepRequest) GetId() int64 {
	if m != nil {
		return m.Id
	}
	return 0
}

type CompleteBrewStepResponse struct {
	NextStep             *Step    `protobuf:"bytes,1,opt,name=nextStep,proto3" json:"nextStep,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CompleteBrewStepResponse) Reset()         { *m = CompleteBrewStepResponse{} }
func (m *CompleteBrewStepResponse) String() string { return proto.CompactTextString(m) }
func (*CompleteBrewStepResponse) ProtoMessage()    {}
func (*CompleteBrewStepResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_brew_1e03cfd08bd125fc, []int{3}
}
func (m *CompleteBrewStepResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CompleteBrewStepResponse.Unmarshal(m, b)
}
func (m *CompleteBrewStepResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CompleteBrewStepResponse.Marshal(b, m, deterministic)
}
func (dst *CompleteBrewStepResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CompleteBrewStepResponse.Merge(dst, src)
}
func (m *CompleteBrewStepResponse) XXX_Size() int {
	return xxx_messageInfo_CompleteBrewStepResponse.Size(m)
}
func (m *CompleteBrewStepResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_CompleteBrewStepResponse.DiscardUnknown(m)
}

var xxx_messageInfo_CompleteBrewStepResponse proto.InternalMessageInfo

func (m *CompleteBrewStepResponse) GetNextStep() *Step {
	if m != nil {
		return m.NextStep
	}
	return nil
}

type StopBrewRequest struct {
	Id                   int64    `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *StopBrewRequest) Reset()         { *m = StopBrewRequest{} }
func (m *StopBrewRequest) String() string { return proto.CompactTextString(m) }
func (*StopBrewRequest) ProtoMessage()    {}
func (*StopBrewRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_brew_1e03cfd08bd125fc, []int{4}
}
func (m *StopBrewRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StopBrewRequest.Unmarshal(m, b)
}
func (m *StopBrewRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StopBrewRequest.Marshal(b, m, deterministic)
}
func (dst *StopBrewRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StopBrewRequest.Merge(dst, src)
}
func (m *StopBrewRequest) XXX_Size() int {
	return xxx_messageInfo_StopBrewRequest.Size(m)
}
func (m *StopBrewRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_StopBrewRequest.DiscardUnknown(m)
}

var xxx_messageInfo_StopBrewRequest proto.InternalMessageInfo

func (m *StopBrewRequest) GetId() int64 {
	if m != nil {
		return m.Id
	}
	return 0
}

type StopBrewResponse struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *StopBrewResponse) Reset()         { *m = StopBrewResponse{} }
func (m *StopBrewResponse) String() string { return proto.CompactTextString(m) }
func (*StopBrewResponse) ProtoMessage()    {}
func (*StopBrewResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_brew_1e03cfd08bd125fc, []int{5}
}
func (m *StopBrewResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StopBrewResponse.Unmarshal(m, b)
}
func (m *StopBrewResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StopBrewResponse.Marshal(b, m, deterministic)
}
func (dst *StopBrewResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StopBrewResponse.Merge(dst, src)
}
func (m *StopBrewResponse) XXX_Size() int {
	return xxx_messageInfo_StopBrewResponse.Size(m)
}
func (m *StopBrewResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_StopBrewResponse.DiscardUnknown(m)
}

var xxx_messageInfo_StopBrewResponse proto.InternalMessageInfo

type GetBrewRequest struct {
	Id                   int64    `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetBrewRequest) Reset()         { *m = GetBrewRequest{} }
func (m *GetBrewRequest) String() string { return proto.CompactTextString(m) }
func (*GetBrewRequest) ProtoMessage()    {}
func (*GetBrewRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_brew_1e03cfd08bd125fc, []int{6}
}
func (m *GetBrewRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetBrewRequest.Unmarshal(m, b)
}
func (m *GetBrewRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetBrewRequest.Marshal(b, m, deterministic)
}
func (dst *GetBrewRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetBrewRequest.Merge(dst, src)
}
func (m *GetBrewRequest) XXX_Size() int {
	return xxx_messageInfo_GetBrewRequest.Size(m)
}
func (m *GetBrewRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetBrewRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetBrewRequest proto.InternalMessageInfo

func (m *GetBrewRequest) GetId() int64 {
	if m != nil {
		return m.Id
	}
	return 0
}

type GetBrewResponse struct {
	Brew                 *Brew    `protobuf:"bytes,1,opt,name=brew,proto3" json:"brew,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetBrewResponse) Reset()         { *m = GetBrewResponse{} }
func (m *GetBrewResponse) String() string { return proto.CompactTextString(m) }
func (*GetBrewResponse) ProtoMessage()    {}
func (*GetBrewResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_brew_1e03cfd08bd125fc, []int{7}
}
func (m *GetBrewResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetBrewResponse.Unmarshal(m, b)
}
func (m *GetBrewResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetBrewResponse.Marshal(b, m, deterministic)
}
func (dst *GetBrewResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetBrewResponse.Merge(dst, src)
}
func (m *GetBrewResponse) XXX_Size() int {
	return xxx_messageInfo_GetBrewResponse.Size(m)
}
func (m *GetBrewResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_GetBrewResponse.DiscardUnknown(m)
}

var xxx_messageInfo_GetBrewResponse proto.InternalMessageInfo

func (m *GetBrewResponse) GetBrew() *Brew {
	if m != nil {
		return m.Brew
	}
	return nil
}

type GetActiveBrewRequest struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetActiveBrewRequest) Reset()         { *m = GetActiveBrewRequest{} }
func (m *GetActiveBrewRequest) String() string { return proto.CompactTextString(m) }
func (*GetActiveBrewRequest) ProtoMessage()    {}
func (*GetActiveBrewRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_brew_1e03cfd08bd125fc, []int{8}
}
func (m *GetActiveBrewRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetActiveBrewRequest.Unmarshal(m, b)
}
func (m *GetActiveBrewRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetActiveBrewRequest.Marshal(b, m, deterministic)
}
func (dst *GetActiveBrewRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetActiveBrewRequest.Merge(dst, src)
}
func (m *GetActiveBrewRequest) XXX_Size() int {
	return xxx_messageInfo_GetActiveBrewRequest.Size(m)
}
func (m *GetActiveBrewRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetActiveBrewRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetActiveBrewRequest proto.InternalMessageInfo

type GetActiveBrewResponse struct {
	Brew                 *Brew    `protobuf:"bytes,1,opt,name=brew,proto3" json:"brew,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetActiveBrewResponse) Reset()         { *m = GetActiveBrewResponse{} }
func (m *GetActiveBrewResponse) String() string { return proto.CompactTextString(m) }
func (*GetActiveBrewResponse) ProtoMessage()    {}
func (*GetActiveBrewResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_brew_1e03cfd08bd125fc, []int{9}
}
func (m *GetActiveBrewResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetActiveBrewResponse.Unmarshal(m, b)
}
func (m *GetActiveBrewResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetActiveBrewResponse.Marshal(b, m, deterministic)
}
func (dst *GetActiveBrewResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetActiveBrewResponse.Merge(dst, src)
}
func (m *GetActiveBrewResponse) XXX_Size() int {
	return xxx_messageInfo_GetActiveBrewResponse.Size(m)
}
func (m *GetActiveBrewResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_GetActiveBrewResponse.DiscardUnknown(m)
}

var xxx_messageInfo_GetActiveBrewResponse proto.InternalMessageInfo

func (m *GetActiveBrewResponse) GetBrew() *Brew {
	if m != nil {
		return m.Brew
	}
	return nil
}

type ListBrewRequest struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ListBrewRequest) Reset()         { *m = ListBrewRequest{} }
func (m *ListBrewRequest) String() string { return proto.CompactTextString(m) }
func (*ListBrewRequest) ProtoMessage()    {}
func (*ListBrewRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_brew_1e03cfd08bd125fc, []int{10}
}
func (m *ListBrewRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ListBrewRequest.Unmarshal(m, b)
}
func (m *ListBrewRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ListBrewRequest.Marshal(b, m, deterministic)
}
func (dst *ListBrewRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ListBrewRequest.Merge(dst, src)
}
func (m *ListBrewRequest) XXX_Size() int {
	return xxx_messageInfo_ListBrewRequest.Size(m)
}
func (m *ListBrewRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ListBrewRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ListBrewRequest proto.InternalMessageInfo

type ListBrewResponse struct {
	Brews                []*Brew  `protobuf:"bytes,1,rep,name=brews,proto3" json:"brews,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ListBrewResponse) Reset()         { *m = ListBrewResponse{} }
func (m *ListBrewResponse) String() string { return proto.CompactTextString(m) }
func (*ListBrewResponse) ProtoMessage()    {}
func (*ListBrewResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_brew_1e03cfd08bd125fc, []int{11}
}
func (m *ListBrewResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ListBrewResponse.Unmarshal(m, b)
}
func (m *ListBrewResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ListBrewResponse.Marshal(b, m, deterministic)
}
func (dst *ListBrewResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ListBrewResponse.Merge(dst, src)
}
func (m *ListBrewResponse) XXX_Size() int {
	return xxx_messageInfo_ListBrewResponse.Size(m)
}
func (m *ListBrewResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_ListBrewResponse.DiscardUnknown(m)
}

var xxx_messageInfo_ListBrewResponse proto.InternalMessageInfo

func (m *ListBrewResponse) GetBrews() []*Brew {
	if m != nil {
		return m.Brews
	}
	return nil
}

type Brew struct {
	Id                   int64                `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	StartTime            *timestamp.Timestamp `protobuf:"bytes,2,opt,name=startTime,proto3" json:"startTime,omitempty"`
	CompletedTime        *timestamp.Timestamp `protobuf:"bytes,3,opt,name=completedTime,proto3" json:"completedTime,omitempty"`
	Recipe               *Recipe              `protobuf:"bytes,4,opt,name=recipe,proto3" json:"recipe,omitempty"`
	CompletedSteps       []*BrewStep          `protobuf:"bytes,5,rep,name=completedSteps,proto3" json:"completedSteps,omitempty"`
	Note                 string               `protobuf:"bytes,6,opt,name=note,proto3" json:"note,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *Brew) Reset()         { *m = Brew{} }
func (m *Brew) String() string { return proto.CompactTextString(m) }
func (*Brew) ProtoMessage()    {}
func (*Brew) Descriptor() ([]byte, []int) {
	return fileDescriptor_brew_1e03cfd08bd125fc, []int{12}
}
func (m *Brew) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Brew.Unmarshal(m, b)
}
func (m *Brew) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Brew.Marshal(b, m, deterministic)
}
func (dst *Brew) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Brew.Merge(dst, src)
}
func (m *Brew) XXX_Size() int {
	return xxx_messageInfo_Brew.Size(m)
}
func (m *Brew) XXX_DiscardUnknown() {
	xxx_messageInfo_Brew.DiscardUnknown(m)
}

var xxx_messageInfo_Brew proto.InternalMessageInfo

func (m *Brew) GetId() int64 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *Brew) GetStartTime() *timestamp.Timestamp {
	if m != nil {
		return m.StartTime
	}
	return nil
}

func (m *Brew) GetCompletedTime() *timestamp.Timestamp {
	if m != nil {
		return m.CompletedTime
	}
	return nil
}

func (m *Brew) GetRecipe() *Recipe {
	if m != nil {
		return m.Recipe
	}
	return nil
}

func (m *Brew) GetCompletedSteps() []*BrewStep {
	if m != nil {
		return m.CompletedSteps
	}
	return nil
}

func (m *Brew) GetNote() string {
	if m != nil {
		return m.Note
	}
	return ""
}

type BrewStep struct {
	StartTime            *timestamp.Timestamp `protobuf:"bytes,1,opt,name=startTime,proto3" json:"startTime,omitempty"`
	CompletedTime        *timestamp.Timestamp `protobuf:"bytes,2,opt,name=completedTime,proto3" json:"completedTime,omitempty"`
	Step                 *Step                `protobuf:"bytes,4,opt,name=step,proto3" json:"step,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *BrewStep) Reset()         { *m = BrewStep{} }
func (m *BrewStep) String() string { return proto.CompactTextString(m) }
func (*BrewStep) ProtoMessage()    {}
func (*BrewStep) Descriptor() ([]byte, []int) {
	return fileDescriptor_brew_1e03cfd08bd125fc, []int{13}
}
func (m *BrewStep) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_BrewStep.Unmarshal(m, b)
}
func (m *BrewStep) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_BrewStep.Marshal(b, m, deterministic)
}
func (dst *BrewStep) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BrewStep.Merge(dst, src)
}
func (m *BrewStep) XXX_Size() int {
	return xxx_messageInfo_BrewStep.Size(m)
}
func (m *BrewStep) XXX_DiscardUnknown() {
	xxx_messageInfo_BrewStep.DiscardUnknown(m)
}

var xxx_messageInfo_BrewStep proto.InternalMessageInfo

func (m *BrewStep) GetStartTime() *timestamp.Timestamp {
	if m != nil {
		return m.StartTime
	}
	return nil
}

func (m *BrewStep) GetCompletedTime() *timestamp.Timestamp {
	if m != nil {
		return m.CompletedTime
	}
	return nil
}

func (m *BrewStep) GetStep() *Step {
	if m != nil {
		return m.Step
	}
	return nil
}

func init() {
	proto.RegisterType((*StartBrewRequest)(nil), "brewmmer.StartBrewRequest")
	proto.RegisterType((*StartBrewResponse)(nil), "brewmmer.StartBrewResponse")
	proto.RegisterType((*CompleteBrewStepRequest)(nil), "brewmmer.CompleteBrewStepRequest")
	proto.RegisterType((*CompleteBrewStepResponse)(nil), "brewmmer.CompleteBrewStepResponse")
	proto.RegisterType((*StopBrewRequest)(nil), "brewmmer.StopBrewRequest")
	proto.RegisterType((*StopBrewResponse)(nil), "brewmmer.StopBrewResponse")
	proto.RegisterType((*GetBrewRequest)(nil), "brewmmer.GetBrewRequest")
	proto.RegisterType((*GetBrewResponse)(nil), "brewmmer.GetBrewResponse")
	proto.RegisterType((*GetActiveBrewRequest)(nil), "brewmmer.GetActiveBrewRequest")
	proto.RegisterType((*GetActiveBrewResponse)(nil), "brewmmer.GetActiveBrewResponse")
	proto.RegisterType((*ListBrewRequest)(nil), "brewmmer.ListBrewRequest")
	proto.RegisterType((*ListBrewResponse)(nil), "brewmmer.ListBrewResponse")
	proto.RegisterType((*Brew)(nil), "brewmmer.Brew")
	proto.RegisterType((*BrewStep)(nil), "brewmmer.BrewStep")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// BrewServiceClient is the client API for BrewService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type BrewServiceClient interface {
	StartBrew(ctx context.Context, in *StartBrewRequest, opts ...grpc.CallOption) (*StartBrewResponse, error)
	CompleteBrewStep(ctx context.Context, in *CompleteBrewStepRequest, opts ...grpc.CallOption) (*CompleteBrewStepResponse, error)
	StopBrew(ctx context.Context, in *StopBrewRequest, opts ...grpc.CallOption) (*StopBrewResponse, error)
	GetBrew(ctx context.Context, in *GetBrewRequest, opts ...grpc.CallOption) (*GetBrewResponse, error)
	GetActiveBrew(ctx context.Context, in *GetActiveBrewRequest, opts ...grpc.CallOption) (*GetActiveBrewResponse, error)
	ListBrews(ctx context.Context, in *ListBrewRequest, opts ...grpc.CallOption) (*ListBrewResponse, error)
}

type brewServiceClient struct {
	cc *grpc.ClientConn
}

func NewBrewServiceClient(cc *grpc.ClientConn) BrewServiceClient {
	return &brewServiceClient{cc}
}

func (c *brewServiceClient) StartBrew(ctx context.Context, in *StartBrewRequest, opts ...grpc.CallOption) (*StartBrewResponse, error) {
	out := new(StartBrewResponse)
	err := c.cc.Invoke(ctx, "/brewmmer.BrewService/StartBrew", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *brewServiceClient) CompleteBrewStep(ctx context.Context, in *CompleteBrewStepRequest, opts ...grpc.CallOption) (*CompleteBrewStepResponse, error) {
	out := new(CompleteBrewStepResponse)
	err := c.cc.Invoke(ctx, "/brewmmer.BrewService/CompleteBrewStep", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *brewServiceClient) StopBrew(ctx context.Context, in *StopBrewRequest, opts ...grpc.CallOption) (*StopBrewResponse, error) {
	out := new(StopBrewResponse)
	err := c.cc.Invoke(ctx, "/brewmmer.BrewService/StopBrew", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *brewServiceClient) GetBrew(ctx context.Context, in *GetBrewRequest, opts ...grpc.CallOption) (*GetBrewResponse, error) {
	out := new(GetBrewResponse)
	err := c.cc.Invoke(ctx, "/brewmmer.BrewService/GetBrew", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *brewServiceClient) GetActiveBrew(ctx context.Context, in *GetActiveBrewRequest, opts ...grpc.CallOption) (*GetActiveBrewResponse, error) {
	out := new(GetActiveBrewResponse)
	err := c.cc.Invoke(ctx, "/brewmmer.BrewService/GetActiveBrew", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *brewServiceClient) ListBrews(ctx context.Context, in *ListBrewRequest, opts ...grpc.CallOption) (*ListBrewResponse, error) {
	out := new(ListBrewResponse)
	err := c.cc.Invoke(ctx, "/brewmmer.BrewService/ListBrews", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// BrewServiceServer is the server API for BrewService service.
type BrewServiceServer interface {
	StartBrew(context.Context, *StartBrewRequest) (*StartBrewResponse, error)
	CompleteBrewStep(context.Context, *CompleteBrewStepRequest) (*CompleteBrewStepResponse, error)
	StopBrew(context.Context, *StopBrewRequest) (*StopBrewResponse, error)
	GetBrew(context.Context, *GetBrewRequest) (*GetBrewResponse, error)
	GetActiveBrew(context.Context, *GetActiveBrewRequest) (*GetActiveBrewResponse, error)
	ListBrews(context.Context, *ListBrewRequest) (*ListBrewResponse, error)
}

func RegisterBrewServiceServer(s *grpc.Server, srv BrewServiceServer) {
	s.RegisterService(&_BrewService_serviceDesc, srv)
}

func _BrewService_StartBrew_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StartBrewRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BrewServiceServer).StartBrew(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/brewmmer.BrewService/StartBrew",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BrewServiceServer).StartBrew(ctx, req.(*StartBrewRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BrewService_CompleteBrewStep_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CompleteBrewStepRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BrewServiceServer).CompleteBrewStep(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/brewmmer.BrewService/CompleteBrewStep",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BrewServiceServer).CompleteBrewStep(ctx, req.(*CompleteBrewStepRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BrewService_StopBrew_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StopBrewRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BrewServiceServer).StopBrew(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/brewmmer.BrewService/StopBrew",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BrewServiceServer).StopBrew(ctx, req.(*StopBrewRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BrewService_GetBrew_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetBrewRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BrewServiceServer).GetBrew(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/brewmmer.BrewService/GetBrew",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BrewServiceServer).GetBrew(ctx, req.(*GetBrewRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BrewService_GetActiveBrew_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetActiveBrewRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BrewServiceServer).GetActiveBrew(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/brewmmer.BrewService/GetActiveBrew",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BrewServiceServer).GetActiveBrew(ctx, req.(*GetActiveBrewRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BrewService_ListBrews_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListBrewRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BrewServiceServer).ListBrews(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/brewmmer.BrewService/ListBrews",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BrewServiceServer).ListBrews(ctx, req.(*ListBrewRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _BrewService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "brewmmer.BrewService",
	HandlerType: (*BrewServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "StartBrew",
			Handler:    _BrewService_StartBrew_Handler,
		},
		{
			MethodName: "CompleteBrewStep",
			Handler:    _BrewService_CompleteBrewStep_Handler,
		},
		{
			MethodName: "StopBrew",
			Handler:    _BrewService_StopBrew_Handler,
		},
		{
			MethodName: "GetBrew",
			Handler:    _BrewService_GetBrew_Handler,
		},
		{
			MethodName: "GetActiveBrew",
			Handler:    _BrewService_GetActiveBrew_Handler,
		},
		{
			MethodName: "ListBrews",
			Handler:    _BrewService_ListBrews_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/brewmmer/brew.proto",
}

func init() { proto.RegisterFile("api/brewmmer/brew.proto", fileDescriptor_brew_1e03cfd08bd125fc) }

var fileDescriptor_brew_1e03cfd08bd125fc = []byte{
	// 540 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x53, 0x5d, 0x6f, 0xd3, 0x30,
	0x14, 0x55, 0xd2, 0xac, 0xa4, 0xb7, 0x5a, 0xdb, 0x59, 0xc0, 0x32, 0x23, 0xb1, 0xce, 0xf0, 0x50,
	0x78, 0x48, 0xa5, 0x22, 0xa4, 0x09, 0x24, 0xc4, 0x0a, 0x62, 0x42, 0xe2, 0x01, 0xa5, 0x7b, 0xe1,
	0xb1, 0x4d, 0x2f, 0x93, 0xa5, 0xa5, 0x09, 0xb1, 0xb7, 0xf1, 0x23, 0xf8, 0x25, 0xfc, 0x11, 0xfe,
	0x16, 0x8a, 0x63, 0xe7, 0xab, 0xd9, 0x50, 0xc5, 0x5b, 0x72, 0xef, 0x39, 0xc7, 0xf7, 0xda, 0xe7,
	0xc0, 0xe1, 0x32, 0xe1, 0xd3, 0x55, 0x8a, 0xb7, 0x51, 0x84, 0xa9, 0xfa, 0xf0, 0x93, 0x34, 0x96,
	0x31, 0x71, 0x4d, 0x91, 0x1e, 0x5f, 0xc6, 0xf1, 0xe5, 0x15, 0x4e, 0x55, 0x7d, 0x75, 0xfd, 0x7d,
	0x2a, 0x79, 0x84, 0x42, 0x2e, 0xa3, 0x24, 0x87, 0x52, 0x5a, 0xd3, 0x48, 0x31, 0xc4, 0x84, 0x63,
	0xde, 0x63, 0x73, 0x18, 0x2d, 0xe4, 0x32, 0x95, 0xf3, 0x14, 0x6f, 0x03, 0xfc, 0x71, 0x8d, 0x42,
	0x12, 0x0a, 0x6e, 0x8a, 0x21, 0x4f, 0xf0, 0xf3, 0xda, 0xb3, 0xc6, 0xd6, 0xa4, 0x13, 0x14, 0xff,
	0x84, 0x80, 0xb3, 0x89, 0x25, 0x7a, 0xf6, 0xd8, 0x9a, 0xf4, 0x02, 0xf5, 0xcd, 0x9e, 0xc1, 0x41,
	0x45, 0x43, 0x24, 0xf1, 0x46, 0x20, 0x19, 0x80, 0xcd, 0x0d, 0xdd, 0xe6, 0x6b, 0xf6, 0x02, 0x0e,
	0x3f, 0xc4, 0x51, 0x72, 0x85, 0x12, 0x33, 0xdc, 0x42, 0x62, 0x62, 0xce, 0x6b, 0x42, 0x3f, 0x81,
	0xb7, 0x0d, 0xd5, 0xb2, 0x2f, 0xc1, 0xdd, 0xe0, 0x4f, 0x99, 0xd5, 0x14, 0xa3, 0x3f, 0x1b, 0xf8,
	0x66, 0x35, 0x5f, 0x21, 0x8b, 0x3e, 0x3b, 0x81, 0xe1, 0x42, 0xc6, 0x49, 0x75, 0xb5, 0xe6, 0x51,
	0x24, 0x5b, 0xdf, 0x40, 0xf2, 0x23, 0xd8, 0x18, 0x06, 0xe7, 0x28, 0xef, 0x63, 0xbd, 0x86, 0x61,
	0x81, 0xd0, 0x73, 0x31, 0x70, 0xb2, 0x31, 0xb6, 0x67, 0x52, 0x28, 0xd5, 0x63, 0x8f, 0xe1, 0xe1,
	0x39, 0xca, 0xb3, 0x50, 0xf2, 0x1b, 0xac, 0xc8, 0xb3, 0xb7, 0xf0, 0xa8, 0x51, 0xdf, 0x41, 0xf4,
	0x00, 0x86, 0x5f, 0xb8, 0xa8, 0x8e, 0xcb, 0x4e, 0x61, 0x54, 0x96, 0xb4, 0xd4, 0x73, 0xd8, 0xcb,
	0xe0, 0xc2, 0xb3, 0xc6, 0x9d, 0x16, 0xad, 0xbc, 0xc9, 0x7e, 0xd9, 0xe0, 0x64, 0xff, 0xcd, 0x8d,
	0xc9, 0x29, 0xf4, 0x44, 0xf6, 0xc4, 0x17, 0x3c, 0xca, 0xdf, 0xbe, 0x3f, 0xa3, 0x7e, 0xee, 0x3b,
	0xdf, 0xf8, 0xce, 0xbf, 0x30, 0xbe, 0x0b, 0x4a, 0x30, 0x79, 0x0f, 0xfb, 0xa1, 0x7e, 0xcc, 0xb5,
	0x62, 0x77, 0xfe, 0xc9, 0xae, 0x13, 0xc8, 0x04, 0xba, 0xb9, 0xfd, 0x3c, 0x47, 0x51, 0x47, 0xe5,
	0xec, 0x81, 0xaa, 0x07, 0xba, 0x4f, 0xde, 0xc0, 0xa0, 0xa0, 0x66, 0x0e, 0x10, 0xde, 0x9e, 0xda,
	0x96, 0xd4, 0xb7, 0x55, 0x36, 0x69, 0x20, 0x0b, 0x63, 0x77, 0x2b, 0xc6, 0xfe, 0x6d, 0x81, 0x6b,
	0x08, 0xf5, 0x2b, 0xb0, 0xfe, 0xeb, 0x0a, 0xec, 0x5d, 0xaf, 0x80, 0x81, 0x23, 0x32, 0xc7, 0x3b,
	0xad, 0x8e, 0x57, 0xbd, 0xd9, 0x9f, 0x0e, 0xf4, 0xd5, 0xb0, 0x98, 0xde, 0xf0, 0x10, 0xc9, 0x47,
	0xe8, 0x15, 0xa9, 0x24, 0xb4, 0x4a, 0xa9, 0xc7, 0x9d, 0x3e, 0x69, 0xed, 0x69, 0xdf, 0x7c, 0x83,
	0x51, 0x33, 0x8b, 0xe4, 0xa4, 0x24, 0xdc, 0x11, 0x69, 0xca, 0xee, 0x83, 0x68, 0xe9, 0x33, 0x70,
	0x4d, 0xf6, 0xc8, 0x51, 0x75, 0x86, 0x5a, 0x64, 0x29, 0x6d, 0x6b, 0x69, 0x89, 0x77, 0xf0, 0x40,
	0x07, 0x91, 0x78, 0x25, 0xac, 0x9e, 0x5e, 0x7a, 0xd4, 0xd2, 0xd1, 0xfc, 0xaf, 0xb0, 0x5f, 0x4b,
	0x1e, 0x79, 0x5a, 0xc3, 0x6e, 0x45, 0x95, 0x1e, 0xdf, 0xd9, 0xd7, 0x8a, 0x73, 0xe8, 0x99, 0xec,
	0x89, 0xea, 0x56, 0x8d, 0x8c, 0x56, 0xb7, 0x6a, 0x66, 0x75, 0xd5, 0x55, 0x86, 0x78, 0xf5, 0x37,
	0x00, 0x00, 0xff, 0xff, 0xcf, 0x81, 0x28, 0x64, 0xfc, 0x05, 0x00, 0x00,
}