// Copyright 2020 New Relic Corporation. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        (unknown)
// source: vendor/github.com/newrelic/go-agent/v3/internal/com_newrelic_trace_v1/v1.proto

package com_newrelic_trace_v1

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

type SpanBatch struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Spans []*Span `protobuf:"bytes,1,rep,name=spans,proto3" json:"spans,omitempty"`
}

func (x *SpanBatch) Reset() {
	*x = SpanBatch{}
	if protoimpl.UnsafeEnabled {
		mi := &file_vendor_github_com_newrelic_go_agent_v3_internal_com_newrelic_trace_v1_v1_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SpanBatch) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SpanBatch) ProtoMessage() {}

func (x *SpanBatch) ProtoReflect() protoreflect.Message {
	mi := &file_vendor_github_com_newrelic_go_agent_v3_internal_com_newrelic_trace_v1_v1_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SpanBatch.ProtoReflect.Descriptor instead.
func (*SpanBatch) Descriptor() ([]byte, []int) {
	return file_vendor_github_com_newrelic_go_agent_v3_internal_com_newrelic_trace_v1_v1_proto_rawDescGZIP(), []int{0}
}

func (x *SpanBatch) GetSpans() []*Span {
	if x != nil {
		return x.Spans
	}
	return nil
}

type Span struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TraceId         string                     `protobuf:"bytes,1,opt,name=trace_id,json=traceId,proto3" json:"trace_id,omitempty"`
	Intrinsics      map[string]*AttributeValue `protobuf:"bytes,2,rep,name=intrinsics,proto3" json:"intrinsics,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	UserAttributes  map[string]*AttributeValue `protobuf:"bytes,3,rep,name=user_attributes,json=userAttributes,proto3" json:"user_attributes,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	AgentAttributes map[string]*AttributeValue `protobuf:"bytes,4,rep,name=agent_attributes,json=agentAttributes,proto3" json:"agent_attributes,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *Span) Reset() {
	*x = Span{}
	if protoimpl.UnsafeEnabled {
		mi := &file_vendor_github_com_newrelic_go_agent_v3_internal_com_newrelic_trace_v1_v1_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Span) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Span) ProtoMessage() {}

func (x *Span) ProtoReflect() protoreflect.Message {
	mi := &file_vendor_github_com_newrelic_go_agent_v3_internal_com_newrelic_trace_v1_v1_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Span.ProtoReflect.Descriptor instead.
func (*Span) Descriptor() ([]byte, []int) {
	return file_vendor_github_com_newrelic_go_agent_v3_internal_com_newrelic_trace_v1_v1_proto_rawDescGZIP(), []int{1}
}

func (x *Span) GetTraceId() string {
	if x != nil {
		return x.TraceId
	}
	return ""
}

func (x *Span) GetIntrinsics() map[string]*AttributeValue {
	if x != nil {
		return x.Intrinsics
	}
	return nil
}

func (x *Span) GetUserAttributes() map[string]*AttributeValue {
	if x != nil {
		return x.UserAttributes
	}
	return nil
}

func (x *Span) GetAgentAttributes() map[string]*AttributeValue {
	if x != nil {
		return x.AgentAttributes
	}
	return nil
}

type AttributeValue struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Value:
	//
	//	*AttributeValue_StringValue
	//	*AttributeValue_BoolValue
	//	*AttributeValue_IntValue
	//	*AttributeValue_DoubleValue
	Value isAttributeValue_Value `protobuf_oneof:"value"`
}

func (x *AttributeValue) Reset() {
	*x = AttributeValue{}
	if protoimpl.UnsafeEnabled {
		mi := &file_vendor_github_com_newrelic_go_agent_v3_internal_com_newrelic_trace_v1_v1_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AttributeValue) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AttributeValue) ProtoMessage() {}

func (x *AttributeValue) ProtoReflect() protoreflect.Message {
	mi := &file_vendor_github_com_newrelic_go_agent_v3_internal_com_newrelic_trace_v1_v1_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AttributeValue.ProtoReflect.Descriptor instead.
func (*AttributeValue) Descriptor() ([]byte, []int) {
	return file_vendor_github_com_newrelic_go_agent_v3_internal_com_newrelic_trace_v1_v1_proto_rawDescGZIP(), []int{2}
}

func (m *AttributeValue) GetValue() isAttributeValue_Value {
	if m != nil {
		return m.Value
	}
	return nil
}

func (x *AttributeValue) GetStringValue() string {
	if x, ok := x.GetValue().(*AttributeValue_StringValue); ok {
		return x.StringValue
	}
	return ""
}

func (x *AttributeValue) GetBoolValue() bool {
	if x, ok := x.GetValue().(*AttributeValue_BoolValue); ok {
		return x.BoolValue
	}
	return false
}

func (x *AttributeValue) GetIntValue() int64 {
	if x, ok := x.GetValue().(*AttributeValue_IntValue); ok {
		return x.IntValue
	}
	return 0
}

func (x *AttributeValue) GetDoubleValue() float64 {
	if x, ok := x.GetValue().(*AttributeValue_DoubleValue); ok {
		return x.DoubleValue
	}
	return 0
}

type isAttributeValue_Value interface {
	isAttributeValue_Value()
}

type AttributeValue_StringValue struct {
	StringValue string `protobuf:"bytes,1,opt,name=string_value,json=stringValue,proto3,oneof"`
}

type AttributeValue_BoolValue struct {
	BoolValue bool `protobuf:"varint,2,opt,name=bool_value,json=boolValue,proto3,oneof"`
}

type AttributeValue_IntValue struct {
	IntValue int64 `protobuf:"varint,3,opt,name=int_value,json=intValue,proto3,oneof"`
}

type AttributeValue_DoubleValue struct {
	DoubleValue float64 `protobuf:"fixed64,4,opt,name=double_value,json=doubleValue,proto3,oneof"`
}

func (*AttributeValue_StringValue) isAttributeValue_Value() {}

func (*AttributeValue_BoolValue) isAttributeValue_Value() {}

func (*AttributeValue_IntValue) isAttributeValue_Value() {}

func (*AttributeValue_DoubleValue) isAttributeValue_Value() {}

type RecordStatus struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	MessagesSeen uint64 `protobuf:"varint,1,opt,name=messages_seen,json=messagesSeen,proto3" json:"messages_seen,omitempty"`
}

func (x *RecordStatus) Reset() {
	*x = RecordStatus{}
	if protoimpl.UnsafeEnabled {
		mi := &file_vendor_github_com_newrelic_go_agent_v3_internal_com_newrelic_trace_v1_v1_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RecordStatus) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RecordStatus) ProtoMessage() {}

func (x *RecordStatus) ProtoReflect() protoreflect.Message {
	mi := &file_vendor_github_com_newrelic_go_agent_v3_internal_com_newrelic_trace_v1_v1_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RecordStatus.ProtoReflect.Descriptor instead.
func (*RecordStatus) Descriptor() ([]byte, []int) {
	return file_vendor_github_com_newrelic_go_agent_v3_internal_com_newrelic_trace_v1_v1_proto_rawDescGZIP(), []int{3}
}

func (x *RecordStatus) GetMessagesSeen() uint64 {
	if x != nil {
		return x.MessagesSeen
	}
	return 0
}

var File_vendor_github_com_newrelic_go_agent_v3_internal_com_newrelic_trace_v1_v1_proto protoreflect.FileDescriptor

var file_vendor_github_com_newrelic_go_agent_v3_internal_com_newrelic_trace_v1_v1_proto_rawDesc = []byte{
	0x0a, 0x4e, 0x76, 0x65, 0x6e, 0x64, 0x6f, 0x72, 0x2f, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e,
	0x63, 0x6f, 0x6d, 0x2f, 0x6e, 0x65, 0x77, 0x72, 0x65, 0x6c, 0x69, 0x63, 0x2f, 0x67, 0x6f, 0x2d,
	0x61, 0x67, 0x65, 0x6e, 0x74, 0x2f, 0x76, 0x33, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61,
	0x6c, 0x2f, 0x63, 0x6f, 0x6d, 0x5f, 0x6e, 0x65, 0x77, 0x72, 0x65, 0x6c, 0x69, 0x63, 0x5f, 0x74,
	0x72, 0x61, 0x63, 0x65, 0x5f, 0x76, 0x31, 0x2f, 0x76, 0x31, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x15, 0x63, 0x6f, 0x6d, 0x2e, 0x6e, 0x65, 0x77, 0x72, 0x65, 0x6c, 0x69, 0x63, 0x2e, 0x74,
	0x72, 0x61, 0x63, 0x65, 0x2e, 0x76, 0x31, 0x22, 0x3e, 0x0a, 0x09, 0x53, 0x70, 0x61, 0x6e, 0x42,
	0x61, 0x74, 0x63, 0x68, 0x12, 0x31, 0x0a, 0x05, 0x73, 0x70, 0x61, 0x6e, 0x73, 0x18, 0x01, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x6e, 0x65, 0x77, 0x72, 0x65, 0x6c,
	0x69, 0x63, 0x2e, 0x74, 0x72, 0x61, 0x63, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x70, 0x61, 0x6e,
	0x52, 0x05, 0x73, 0x70, 0x61, 0x6e, 0x73, 0x22, 0xe0, 0x04, 0x0a, 0x04, 0x53, 0x70, 0x61, 0x6e,
	0x12, 0x19, 0x0a, 0x08, 0x74, 0x72, 0x61, 0x63, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x07, 0x74, 0x72, 0x61, 0x63, 0x65, 0x49, 0x64, 0x12, 0x4b, 0x0a, 0x0a, 0x69,
	0x6e, 0x74, 0x72, 0x69, 0x6e, 0x73, 0x69, 0x63, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x2b, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x6e, 0x65, 0x77, 0x72, 0x65, 0x6c, 0x69, 0x63, 0x2e, 0x74,
	0x72, 0x61, 0x63, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x70, 0x61, 0x6e, 0x2e, 0x49, 0x6e, 0x74,
	0x72, 0x69, 0x6e, 0x73, 0x69, 0x63, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x0a, 0x69, 0x6e,
	0x74, 0x72, 0x69, 0x6e, 0x73, 0x69, 0x63, 0x73, 0x12, 0x58, 0x0a, 0x0f, 0x75, 0x73, 0x65, 0x72,
	0x5f, 0x61, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x2f, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x6e, 0x65, 0x77, 0x72, 0x65, 0x6c, 0x69, 0x63,
	0x2e, 0x74, 0x72, 0x61, 0x63, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x70, 0x61, 0x6e, 0x2e, 0x55,
	0x73, 0x65, 0x72, 0x41, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x73, 0x45, 0x6e, 0x74,
	0x72, 0x79, 0x52, 0x0e, 0x75, 0x73, 0x65, 0x72, 0x41, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74,
	0x65, 0x73, 0x12, 0x5b, 0x0a, 0x10, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x5f, 0x61, 0x74, 0x74, 0x72,
	0x69, 0x62, 0x75, 0x74, 0x65, 0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x30, 0x2e, 0x63,
	0x6f, 0x6d, 0x2e, 0x6e, 0x65, 0x77, 0x72, 0x65, 0x6c, 0x69, 0x63, 0x2e, 0x74, 0x72, 0x61, 0x63,
	0x65, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x70, 0x61, 0x6e, 0x2e, 0x41, 0x67, 0x65, 0x6e, 0x74, 0x41,
	0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x0f,
	0x61, 0x67, 0x65, 0x6e, 0x74, 0x41, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x73, 0x1a,
	0x64, 0x0a, 0x0f, 0x49, 0x6e, 0x74, 0x72, 0x69, 0x6e, 0x73, 0x69, 0x63, 0x73, 0x45, 0x6e, 0x74,
	0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x03, 0x6b, 0x65, 0x79, 0x12, 0x3b, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x25, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x6e, 0x65, 0x77, 0x72, 0x65, 0x6c,
	0x69, 0x63, 0x2e, 0x74, 0x72, 0x61, 0x63, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x41, 0x74, 0x74, 0x72,
	0x69, 0x62, 0x75, 0x74, 0x65, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x3a, 0x02, 0x38, 0x01, 0x1a, 0x68, 0x0a, 0x13, 0x55, 0x73, 0x65, 0x72, 0x41, 0x74, 0x74,
	0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03,
	0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x3b,
	0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x25, 0x2e,
	0x63, 0x6f, 0x6d, 0x2e, 0x6e, 0x65, 0x77, 0x72, 0x65, 0x6c, 0x69, 0x63, 0x2e, 0x74, 0x72, 0x61,
	0x63, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x41, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x56,
	0x61, 0x6c, 0x75, 0x65, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x1a,
	0x69, 0x0a, 0x14, 0x41, 0x67, 0x65, 0x6e, 0x74, 0x41, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74,
	0x65, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x3b, 0x0a, 0x05, 0x76, 0x61, 0x6c,
	0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x25, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x6e,
	0x65, 0x77, 0x72, 0x65, 0x6c, 0x69, 0x63, 0x2e, 0x74, 0x72, 0x61, 0x63, 0x65, 0x2e, 0x76, 0x31,
	0x2e, 0x41, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52,
	0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0xa3, 0x01, 0x0a, 0x0e, 0x41,
	0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x23, 0x0a,
	0x0c, 0x73, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x5f, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x0b, 0x73, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x56, 0x61, 0x6c,
	0x75, 0x65, 0x12, 0x1f, 0x0a, 0x0a, 0x62, 0x6f, 0x6f, 0x6c, 0x5f, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x48, 0x00, 0x52, 0x09, 0x62, 0x6f, 0x6f, 0x6c, 0x56, 0x61,
	0x6c, 0x75, 0x65, 0x12, 0x1d, 0x0a, 0x09, 0x69, 0x6e, 0x74, 0x5f, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x48, 0x00, 0x52, 0x08, 0x69, 0x6e, 0x74, 0x56, 0x61, 0x6c,
	0x75, 0x65, 0x12, 0x23, 0x0a, 0x0c, 0x64, 0x6f, 0x75, 0x62, 0x6c, 0x65, 0x5f, 0x76, 0x61, 0x6c,
	0x75, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x01, 0x48, 0x00, 0x52, 0x0b, 0x64, 0x6f, 0x75, 0x62,
	0x6c, 0x65, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x42, 0x07, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x22, 0x33, 0x0a, 0x0c, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x12, 0x23, 0x0a, 0x0d, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x5f, 0x73, 0x65, 0x65,
	0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0c, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x73, 0x53, 0x65, 0x65, 0x6e, 0x32, 0xc5, 0x01, 0x0a, 0x0d, 0x49, 0x6e, 0x67, 0x65, 0x73, 0x74,
	0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x54, 0x0a, 0x0a, 0x52, 0x65, 0x63, 0x6f, 0x72,
	0x64, 0x53, 0x70, 0x61, 0x6e, 0x12, 0x1b, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x6e, 0x65, 0x77, 0x72,
	0x65, 0x6c, 0x69, 0x63, 0x2e, 0x74, 0x72, 0x61, 0x63, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x70,
	0x61, 0x6e, 0x1a, 0x23, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x6e, 0x65, 0x77, 0x72, 0x65, 0x6c, 0x69,
	0x63, 0x2e, 0x74, 0x72, 0x61, 0x63, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x65, 0x63, 0x6f, 0x72,
	0x64, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x22, 0x00, 0x28, 0x01, 0x30, 0x01, 0x12, 0x5e, 0x0a,
	0x0f, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x53, 0x70, 0x61, 0x6e, 0x42, 0x61, 0x74, 0x63, 0x68,
	0x12, 0x20, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x6e, 0x65, 0x77, 0x72, 0x65, 0x6c, 0x69, 0x63, 0x2e,
	0x74, 0x72, 0x61, 0x63, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x70, 0x61, 0x6e, 0x42, 0x61, 0x74,
	0x63, 0x68, 0x1a, 0x23, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x6e, 0x65, 0x77, 0x72, 0x65, 0x6c, 0x69,
	0x63, 0x2e, 0x74, 0x72, 0x61, 0x63, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x65, 0x63, 0x6f, 0x72,
	0x64, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x22, 0x00, 0x28, 0x01, 0x30, 0x01, 0x42, 0x40, 0x5a,
	0x3e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6e, 0x65, 0x77, 0x72,
	0x65, 0x6c, 0x69, 0x63, 0x2f, 0x67, 0x6f, 0x2d, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x2f, 0x76, 0x33,
	0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x63, 0x6f, 0x6d, 0x5f, 0x6e, 0x65,
	0x77, 0x72, 0x65, 0x6c, 0x69, 0x63, 0x5f, 0x74, 0x72, 0x61, 0x63, 0x65, 0x5f, 0x76, 0x31, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_vendor_github_com_newrelic_go_agent_v3_internal_com_newrelic_trace_v1_v1_proto_rawDescOnce sync.Once
	file_vendor_github_com_newrelic_go_agent_v3_internal_com_newrelic_trace_v1_v1_proto_rawDescData = file_vendor_github_com_newrelic_go_agent_v3_internal_com_newrelic_trace_v1_v1_proto_rawDesc
)

func file_vendor_github_com_newrelic_go_agent_v3_internal_com_newrelic_trace_v1_v1_proto_rawDescGZIP() []byte {
	file_vendor_github_com_newrelic_go_agent_v3_internal_com_newrelic_trace_v1_v1_proto_rawDescOnce.Do(func() {
		file_vendor_github_com_newrelic_go_agent_v3_internal_com_newrelic_trace_v1_v1_proto_rawDescData = protoimpl.X.CompressGZIP(file_vendor_github_com_newrelic_go_agent_v3_internal_com_newrelic_trace_v1_v1_proto_rawDescData)
	})
	return file_vendor_github_com_newrelic_go_agent_v3_internal_com_newrelic_trace_v1_v1_proto_rawDescData
}

var file_vendor_github_com_newrelic_go_agent_v3_internal_com_newrelic_trace_v1_v1_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_vendor_github_com_newrelic_go_agent_v3_internal_com_newrelic_trace_v1_v1_proto_goTypes = []interface{}{
	(*SpanBatch)(nil),      // 0: com.newrelic.trace.v1.SpanBatch
	(*Span)(nil),           // 1: com.newrelic.trace.v1.Span
	(*AttributeValue)(nil), // 2: com.newrelic.trace.v1.AttributeValue
	(*RecordStatus)(nil),   // 3: com.newrelic.trace.v1.RecordStatus
	nil,                    // 4: com.newrelic.trace.v1.Span.IntrinsicsEntry
	nil,                    // 5: com.newrelic.trace.v1.Span.UserAttributesEntry
	nil,                    // 6: com.newrelic.trace.v1.Span.AgentAttributesEntry
}
var file_vendor_github_com_newrelic_go_agent_v3_internal_com_newrelic_trace_v1_v1_proto_depIdxs = []int32{
	1, // 0: com.newrelic.trace.v1.SpanBatch.spans:type_name -> com.newrelic.trace.v1.Span
	4, // 1: com.newrelic.trace.v1.Span.intrinsics:type_name -> com.newrelic.trace.v1.Span.IntrinsicsEntry
	5, // 2: com.newrelic.trace.v1.Span.user_attributes:type_name -> com.newrelic.trace.v1.Span.UserAttributesEntry
	6, // 3: com.newrelic.trace.v1.Span.agent_attributes:type_name -> com.newrelic.trace.v1.Span.AgentAttributesEntry
	2, // 4: com.newrelic.trace.v1.Span.IntrinsicsEntry.value:type_name -> com.newrelic.trace.v1.AttributeValue
	2, // 5: com.newrelic.trace.v1.Span.UserAttributesEntry.value:type_name -> com.newrelic.trace.v1.AttributeValue
	2, // 6: com.newrelic.trace.v1.Span.AgentAttributesEntry.value:type_name -> com.newrelic.trace.v1.AttributeValue
	1, // 7: com.newrelic.trace.v1.IngestService.RecordSpan:input_type -> com.newrelic.trace.v1.Span
	0, // 8: com.newrelic.trace.v1.IngestService.RecordSpanBatch:input_type -> com.newrelic.trace.v1.SpanBatch
	3, // 9: com.newrelic.trace.v1.IngestService.RecordSpan:output_type -> com.newrelic.trace.v1.RecordStatus
	3, // 10: com.newrelic.trace.v1.IngestService.RecordSpanBatch:output_type -> com.newrelic.trace.v1.RecordStatus
	9, // [9:11] is the sub-list for method output_type
	7, // [7:9] is the sub-list for method input_type
	7, // [7:7] is the sub-list for extension type_name
	7, // [7:7] is the sub-list for extension extendee
	0, // [0:7] is the sub-list for field type_name
}

func init() {
	file_vendor_github_com_newrelic_go_agent_v3_internal_com_newrelic_trace_v1_v1_proto_init()
}
func file_vendor_github_com_newrelic_go_agent_v3_internal_com_newrelic_trace_v1_v1_proto_init() {
	if File_vendor_github_com_newrelic_go_agent_v3_internal_com_newrelic_trace_v1_v1_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_vendor_github_com_newrelic_go_agent_v3_internal_com_newrelic_trace_v1_v1_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SpanBatch); i {
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
		file_vendor_github_com_newrelic_go_agent_v3_internal_com_newrelic_trace_v1_v1_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Span); i {
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
		file_vendor_github_com_newrelic_go_agent_v3_internal_com_newrelic_trace_v1_v1_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AttributeValue); i {
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
		file_vendor_github_com_newrelic_go_agent_v3_internal_com_newrelic_trace_v1_v1_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RecordStatus); i {
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
	file_vendor_github_com_newrelic_go_agent_v3_internal_com_newrelic_trace_v1_v1_proto_msgTypes[2].OneofWrappers = []interface{}{
		(*AttributeValue_StringValue)(nil),
		(*AttributeValue_BoolValue)(nil),
		(*AttributeValue_IntValue)(nil),
		(*AttributeValue_DoubleValue)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_vendor_github_com_newrelic_go_agent_v3_internal_com_newrelic_trace_v1_v1_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_vendor_github_com_newrelic_go_agent_v3_internal_com_newrelic_trace_v1_v1_proto_goTypes,
		DependencyIndexes: file_vendor_github_com_newrelic_go_agent_v3_internal_com_newrelic_trace_v1_v1_proto_depIdxs,
		MessageInfos:      file_vendor_github_com_newrelic_go_agent_v3_internal_com_newrelic_trace_v1_v1_proto_msgTypes,
	}.Build()
	File_vendor_github_com_newrelic_go_agent_v3_internal_com_newrelic_trace_v1_v1_proto = out.File
	file_vendor_github_com_newrelic_go_agent_v3_internal_com_newrelic_trace_v1_v1_proto_rawDesc = nil
	file_vendor_github_com_newrelic_go_agent_v3_internal_com_newrelic_trace_v1_v1_proto_goTypes = nil
	file_vendor_github_com_newrelic_go_agent_v3_internal_com_newrelic_trace_v1_v1_proto_depIdxs = nil
}