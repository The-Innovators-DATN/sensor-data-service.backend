// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v3.12.4
// source: commonpb/common.proto

package commonpb

import (
	any1 "github.com/golang/protobuf/ptypes/any"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type StandardResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Status        string                 `protobuf:"bytes,1,opt,name=status,proto3" json:"status,omitempty"`   // "success" | "error"
	Message       string                 `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"` // Thông báo kèm theo
	Data          *any1.Any              `protobuf:"bytes,3,opt,name=data,proto3" json:"data,omitempty"`       // Payload thực sự
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *StandardResponse) Reset() {
	*x = StandardResponse{}
	mi := &file_commonpb_common_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *StandardResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StandardResponse) ProtoMessage() {}

func (x *StandardResponse) ProtoReflect() protoreflect.Message {
	mi := &file_commonpb_common_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StandardResponse.ProtoReflect.Descriptor instead.
func (*StandardResponse) Descriptor() ([]byte, []int) {
	return file_commonpb_common_proto_rawDescGZIP(), []int{0}
}

func (x *StandardResponse) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

func (x *StandardResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

func (x *StandardResponse) GetData() *any1.Any {
	if x != nil {
		return x.Data
	}
	return nil
}

// ENUM WRAPPER
type EnumValue struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Name          string                 `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *EnumValue) Reset() {
	*x = EnumValue{}
	mi := &file_commonpb_common_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *EnumValue) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EnumValue) ProtoMessage() {}

func (x *EnumValue) ProtoReflect() protoreflect.Message {
	mi := &file_commonpb_common_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EnumValue.ProtoReflect.Descriptor instead.
func (*EnumValue) Descriptor() ([]byte, []int) {
	return file_commonpb_common_proto_rawDescGZIP(), []int{1}
}

func (x *EnumValue) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

var File_commonpb_common_proto protoreflect.FileDescriptor

const file_commonpb_common_proto_rawDesc = "" +
	"\n" +
	"\x15commonpb/common.proto\x12\x06common\x1a\x19google/protobuf/any.proto\"n\n" +
	"\x10StandardResponse\x12\x16\n" +
	"\x06status\x18\x01 \x01(\tR\x06status\x12\x18\n" +
	"\amessage\x18\x02 \x01(\tR\amessage\x12(\n" +
	"\x04data\x18\x03 \x01(\v2\x14.google.protobuf.AnyR\x04data\"\x1f\n" +
	"\tEnumValue\x12\x12\n" +
	"\x04name\x18\x01 \x01(\tR\x04nameB6Z4sensor-data-service.backend/api/pb/commonpb;commonpbb\x06proto3"

var (
	file_commonpb_common_proto_rawDescOnce sync.Once
	file_commonpb_common_proto_rawDescData []byte
)

func file_commonpb_common_proto_rawDescGZIP() []byte {
	file_commonpb_common_proto_rawDescOnce.Do(func() {
		file_commonpb_common_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_commonpb_common_proto_rawDesc), len(file_commonpb_common_proto_rawDesc)))
	})
	return file_commonpb_common_proto_rawDescData
}

var file_commonpb_common_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_commonpb_common_proto_goTypes = []any{
	(*StandardResponse)(nil), // 0: common.StandardResponse
	(*EnumValue)(nil),        // 1: common.EnumValue
	(*any1.Any)(nil),         // 2: google.protobuf.Any
}
var file_commonpb_common_proto_depIdxs = []int32{
	2, // 0: common.StandardResponse.data:type_name -> google.protobuf.Any
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_commonpb_common_proto_init() }
func file_commonpb_common_proto_init() {
	if File_commonpb_common_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_commonpb_common_proto_rawDesc), len(file_commonpb_common_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_commonpb_common_proto_goTypes,
		DependencyIndexes: file_commonpb_common_proto_depIdxs,
		MessageInfos:      file_commonpb_common_proto_msgTypes,
	}.Build()
	File_commonpb_common_proto = out.File
	file_commonpb_common_proto_goTypes = nil
	file_commonpb_common_proto_depIdxs = nil
}
