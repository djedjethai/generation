// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.21.9
// source: api/v1/keyvalue/keyvalue.proto

package keyvalue

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

type Empty struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *Empty) Reset() {
	*x = Empty{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_v1_keyvalue_keyvalue_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Empty) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Empty) ProtoMessage() {}

func (x *Empty) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1_keyvalue_keyvalue_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Empty.ProtoReflect.Descriptor instead.
func (*Empty) Descriptor() ([]byte, []int) {
	return file_api_v1_keyvalue_keyvalue_proto_rawDescGZIP(), []int{0}
}

type Records struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key   string `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Value string `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *Records) Reset() {
	*x = Records{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_v1_keyvalue_keyvalue_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Records) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Records) ProtoMessage() {}

func (x *Records) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1_keyvalue_keyvalue_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Records.ProtoReflect.Descriptor instead.
func (*Records) Descriptor() ([]byte, []int) {
	return file_api_v1_keyvalue_keyvalue_proto_rawDescGZIP(), []int{1}
}

func (x *Records) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *Records) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

type GetRecords struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Records *Records `protobuf:"bytes,1,opt,name=records,proto3" json:"records,omitempty"`
}

func (x *GetRecords) Reset() {
	*x = GetRecords{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_v1_keyvalue_keyvalue_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetRecords) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetRecords) ProtoMessage() {}

func (x *GetRecords) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1_keyvalue_keyvalue_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetRecords.ProtoReflect.Descriptor instead.
func (*GetRecords) Descriptor() ([]byte, []int) {
	return file_api_v1_keyvalue_keyvalue_proto_rawDescGZIP(), []int{2}
}

func (x *GetRecords) GetRecords() *Records {
	if x != nil {
		return x.Records
	}
	return nil
}

type GetRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key string `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
}

func (x *GetRequest) Reset() {
	*x = GetRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_v1_keyvalue_keyvalue_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetRequest) ProtoMessage() {}

func (x *GetRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1_keyvalue_keyvalue_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetRequest.ProtoReflect.Descriptor instead.
func (*GetRequest) Descriptor() ([]byte, []int) {
	return file_api_v1_keyvalue_keyvalue_proto_rawDescGZIP(), []int{3}
}

func (x *GetRequest) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

type GetResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Value string `protobuf:"bytes,1,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *GetResponse) Reset() {
	*x = GetResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_v1_keyvalue_keyvalue_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetResponse) ProtoMessage() {}

func (x *GetResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1_keyvalue_keyvalue_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetResponse.ProtoReflect.Descriptor instead.
func (*GetResponse) Descriptor() ([]byte, []int) {
	return file_api_v1_keyvalue_keyvalue_proto_rawDescGZIP(), []int{4}
}

func (x *GetResponse) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

type GetKeysRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *GetKeysRequest) Reset() {
	*x = GetKeysRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_v1_keyvalue_keyvalue_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetKeysRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetKeysRequest) ProtoMessage() {}

func (x *GetKeysRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1_keyvalue_keyvalue_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetKeysRequest.ProtoReflect.Descriptor instead.
func (*GetKeysRequest) Descriptor() ([]byte, []int) {
	return file_api_v1_keyvalue_keyvalue_proto_rawDescGZIP(), []int{5}
}

type GetKeysResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Keys []string `protobuf:"bytes,1,rep,name=keys,proto3" json:"keys,omitempty"`
}

func (x *GetKeysResponse) Reset() {
	*x = GetKeysResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_v1_keyvalue_keyvalue_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetKeysResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetKeysResponse) ProtoMessage() {}

func (x *GetKeysResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1_keyvalue_keyvalue_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetKeysResponse.ProtoReflect.Descriptor instead.
func (*GetKeysResponse) Descriptor() ([]byte, []int) {
	return file_api_v1_keyvalue_keyvalue_proto_rawDescGZIP(), []int{6}
}

func (x *GetKeysResponse) GetKeys() []string {
	if x != nil {
		return x.Keys
	}
	return nil
}

type PutRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Records *Records `protobuf:"bytes,1,opt,name=records,proto3" json:"records,omitempty"`
}

func (x *PutRequest) Reset() {
	*x = PutRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_v1_keyvalue_keyvalue_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PutRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PutRequest) ProtoMessage() {}

func (x *PutRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1_keyvalue_keyvalue_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PutRequest.ProtoReflect.Descriptor instead.
func (*PutRequest) Descriptor() ([]byte, []int) {
	return file_api_v1_keyvalue_keyvalue_proto_rawDescGZIP(), []int{7}
}

func (x *PutRequest) GetRecords() *Records {
	if x != nil {
		return x.Records
	}
	return nil
}

type PutResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *PutResponse) Reset() {
	*x = PutResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_v1_keyvalue_keyvalue_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PutResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PutResponse) ProtoMessage() {}

func (x *PutResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1_keyvalue_keyvalue_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PutResponse.ProtoReflect.Descriptor instead.
func (*PutResponse) Descriptor() ([]byte, []int) {
	return file_api_v1_keyvalue_keyvalue_proto_rawDescGZIP(), []int{8}
}

type DeleteRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key string `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
}

func (x *DeleteRequest) Reset() {
	*x = DeleteRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_v1_keyvalue_keyvalue_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteRequest) ProtoMessage() {}

func (x *DeleteRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1_keyvalue_keyvalue_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteRequest.ProtoReflect.Descriptor instead.
func (*DeleteRequest) Descriptor() ([]byte, []int) {
	return file_api_v1_keyvalue_keyvalue_proto_rawDescGZIP(), []int{9}
}

func (x *DeleteRequest) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

type DeleteResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *DeleteResponse) Reset() {
	*x = DeleteResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_v1_keyvalue_keyvalue_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteResponse) ProtoMessage() {}

func (x *DeleteResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1_keyvalue_keyvalue_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteResponse.ProtoReflect.Descriptor instead.
func (*DeleteResponse) Descriptor() ([]byte, []int) {
	return file_api_v1_keyvalue_keyvalue_proto_rawDescGZIP(), []int{10}
}

var File_api_v1_keyvalue_keyvalue_proto protoreflect.FileDescriptor

var file_api_v1_keyvalue_keyvalue_proto_rawDesc = []byte{
	0x0a, 0x1e, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x6b, 0x65, 0x79, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x2f, 0x6b, 0x65, 0x79, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x22, 0x07, 0x0a, 0x05, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x31, 0x0a, 0x07, 0x52, 0x65, 0x63,
	0x6f, 0x72, 0x64, 0x73, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x30, 0x0a, 0x0a,
	0x47, 0x65, 0x74, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x73, 0x12, 0x22, 0x0a, 0x07, 0x72, 0x65,
	0x63, 0x6f, 0x72, 0x64, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x08, 0x2e, 0x52, 0x65,
	0x63, 0x6f, 0x72, 0x64, 0x73, 0x52, 0x07, 0x72, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x73, 0x22, 0x1e,
	0x0a, 0x0a, 0x47, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x10, 0x0a, 0x03,
	0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x22, 0x23,
	0x0a, 0x0b, 0x47, 0x65, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x14, 0x0a,
	0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x22, 0x10, 0x0a, 0x0e, 0x47, 0x65, 0x74, 0x4b, 0x65, 0x79, 0x73, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x25, 0x0a, 0x0f, 0x47, 0x65, 0x74, 0x4b, 0x65, 0x79, 0x73,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x6b, 0x65, 0x79, 0x73,
	0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x04, 0x6b, 0x65, 0x79, 0x73, 0x22, 0x30, 0x0a, 0x0a,
	0x50, 0x75, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x22, 0x0a, 0x07, 0x72, 0x65,
	0x63, 0x6f, 0x72, 0x64, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x08, 0x2e, 0x52, 0x65,
	0x63, 0x6f, 0x72, 0x64, 0x73, 0x52, 0x07, 0x72, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x73, 0x22, 0x0d,
	0x0a, 0x0b, 0x50, 0x75, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x21, 0x0a,
	0x0d, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x10,
	0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79,
	0x22, 0x10, 0x0a, 0x0e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x32, 0xd5, 0x01, 0x0a, 0x08, 0x4b, 0x65, 0x79, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12,
	0x20, 0x0a, 0x03, 0x47, 0x65, 0x74, 0x12, 0x0b, 0x2e, 0x47, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x0c, 0x2e, 0x47, 0x65, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x20, 0x0a, 0x03, 0x50, 0x75, 0x74, 0x12, 0x0b, 0x2e, 0x50, 0x75, 0x74, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0c, 0x2e, 0x50, 0x75, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x29, 0x0a, 0x06, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x12, 0x0e, 0x2e,
	0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0f, 0x2e,
	0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2c,
	0x0a, 0x07, 0x47, 0x65, 0x74, 0x4b, 0x65, 0x79, 0x73, 0x12, 0x0f, 0x2e, 0x47, 0x65, 0x74, 0x4b,
	0x65, 0x79, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x10, 0x2e, 0x47, 0x65, 0x74,
	0x4b, 0x65, 0x79, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2c, 0x0a, 0x13,
	0x47, 0x65, 0x74, 0x4b, 0x65, 0x79, 0x73, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x53, 0x74, 0x72,
	0x65, 0x61, 0x6d, 0x12, 0x06, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x0b, 0x2e, 0x47, 0x65,
	0x74, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x73, 0x30, 0x01, 0x42, 0x32, 0x5a, 0x30, 0x67, 0x69,
	0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x64, 0x6a, 0x65, 0x64, 0x6a, 0x65, 0x74,
	0x68, 0x61, 0x69, 0x2f, 0x67, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x61,
	0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x6b, 0x65, 0x79, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_v1_keyvalue_keyvalue_proto_rawDescOnce sync.Once
	file_api_v1_keyvalue_keyvalue_proto_rawDescData = file_api_v1_keyvalue_keyvalue_proto_rawDesc
)

func file_api_v1_keyvalue_keyvalue_proto_rawDescGZIP() []byte {
	file_api_v1_keyvalue_keyvalue_proto_rawDescOnce.Do(func() {
		file_api_v1_keyvalue_keyvalue_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_v1_keyvalue_keyvalue_proto_rawDescData)
	})
	return file_api_v1_keyvalue_keyvalue_proto_rawDescData
}

var file_api_v1_keyvalue_keyvalue_proto_msgTypes = make([]protoimpl.MessageInfo, 11)
var file_api_v1_keyvalue_keyvalue_proto_goTypes = []interface{}{
	(*Empty)(nil),           // 0: Empty
	(*Records)(nil),         // 1: Records
	(*GetRecords)(nil),      // 2: GetRecords
	(*GetRequest)(nil),      // 3: GetRequest
	(*GetResponse)(nil),     // 4: GetResponse
	(*GetKeysRequest)(nil),  // 5: GetKeysRequest
	(*GetKeysResponse)(nil), // 6: GetKeysResponse
	(*PutRequest)(nil),      // 7: PutRequest
	(*PutResponse)(nil),     // 8: PutResponse
	(*DeleteRequest)(nil),   // 9: DeleteRequest
	(*DeleteResponse)(nil),  // 10: DeleteResponse
}
var file_api_v1_keyvalue_keyvalue_proto_depIdxs = []int32{
	1,  // 0: GetRecords.records:type_name -> Records
	1,  // 1: PutRequest.records:type_name -> Records
	3,  // 2: KeyValue.Get:input_type -> GetRequest
	7,  // 3: KeyValue.Put:input_type -> PutRequest
	9,  // 4: KeyValue.Delete:input_type -> DeleteRequest
	5,  // 5: KeyValue.GetKeys:input_type -> GetKeysRequest
	0,  // 6: KeyValue.GetKeysValuesStream:input_type -> Empty
	4,  // 7: KeyValue.Get:output_type -> GetResponse
	8,  // 8: KeyValue.Put:output_type -> PutResponse
	10, // 9: KeyValue.Delete:output_type -> DeleteResponse
	6,  // 10: KeyValue.GetKeys:output_type -> GetKeysResponse
	2,  // 11: KeyValue.GetKeysValuesStream:output_type -> GetRecords
	7,  // [7:12] is the sub-list for method output_type
	2,  // [2:7] is the sub-list for method input_type
	2,  // [2:2] is the sub-list for extension type_name
	2,  // [2:2] is the sub-list for extension extendee
	0,  // [0:2] is the sub-list for field type_name
}

func init() { file_api_v1_keyvalue_keyvalue_proto_init() }
func file_api_v1_keyvalue_keyvalue_proto_init() {
	if File_api_v1_keyvalue_keyvalue_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_api_v1_keyvalue_keyvalue_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Empty); i {
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
		file_api_v1_keyvalue_keyvalue_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Records); i {
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
		file_api_v1_keyvalue_keyvalue_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetRecords); i {
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
		file_api_v1_keyvalue_keyvalue_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetRequest); i {
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
		file_api_v1_keyvalue_keyvalue_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetResponse); i {
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
		file_api_v1_keyvalue_keyvalue_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetKeysRequest); i {
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
		file_api_v1_keyvalue_keyvalue_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetKeysResponse); i {
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
		file_api_v1_keyvalue_keyvalue_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PutRequest); i {
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
		file_api_v1_keyvalue_keyvalue_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PutResponse); i {
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
		file_api_v1_keyvalue_keyvalue_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteRequest); i {
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
		file_api_v1_keyvalue_keyvalue_proto_msgTypes[10].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteResponse); i {
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
			RawDescriptor: file_api_v1_keyvalue_keyvalue_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   11,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_api_v1_keyvalue_keyvalue_proto_goTypes,
		DependencyIndexes: file_api_v1_keyvalue_keyvalue_proto_depIdxs,
		MessageInfos:      file_api_v1_keyvalue_keyvalue_proto_msgTypes,
	}.Build()
	File_api_v1_keyvalue_keyvalue_proto = out.File
	file_api_v1_keyvalue_keyvalue_proto_rawDesc = nil
	file_api_v1_keyvalue_keyvalue_proto_goTypes = nil
	file_api_v1_keyvalue_keyvalue_proto_depIdxs = nil
}
