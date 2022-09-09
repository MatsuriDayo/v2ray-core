package net

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

// UidList represents a list of uid.
type UidList struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The port that this range starts from.
	Uid []uint32 `protobuf:"varint,1,rep,packed,name=uid,proto3" json:"uid,omitempty"`
}

func (x *UidList) Reset() {
	*x = UidList{}
	if protoimpl.UnsafeEnabled {
		mi := &file_common_net_uid_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UidList) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UidList) ProtoMessage() {}

func (x *UidList) ProtoReflect() protoreflect.Message {
	mi := &file_common_net_uid_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UidList.ProtoReflect.Descriptor instead.
func (*UidList) Descriptor() ([]byte, []int) {
	return file_common_net_uid_proto_rawDescGZIP(), []int{0}
}

func (x *UidList) GetUid() []uint32 {
	if x != nil {
		return x.Uid
	}
	return nil
}

var File_common_net_uid_proto protoreflect.FileDescriptor

var file_common_net_uid_proto_rawDesc = []byte{
	0x0a, 0x14, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2f, 0x6e, 0x65, 0x74, 0x2f, 0x75, 0x69, 0x64,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x15, 0x76, 0x32, 0x72, 0x61, 0x79, 0x2e, 0x63, 0x6f,
	0x72, 0x65, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x6e, 0x65, 0x74, 0x22, 0x1b, 0x0a,
	0x07, 0x55, 0x69, 0x64, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x03, 0x28, 0x0d, 0x52, 0x03, 0x75, 0x69, 0x64, 0x42, 0x60, 0x0a, 0x19, 0x63, 0x6f,
	0x6d, 0x2e, 0x76, 0x32, 0x72, 0x61, 0x79, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x63, 0x6f, 0x6d,
	0x6d, 0x6f, 0x6e, 0x2e, 0x6e, 0x65, 0x74, 0x50, 0x01, 0x5a, 0x29, 0x67, 0x69, 0x74, 0x68, 0x75,
	0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x76, 0x32, 0x66, 0x6c, 0x79, 0x2f, 0x76, 0x32, 0x72, 0x61,
	0x79, 0x2d, 0x63, 0x6f, 0x72, 0x65, 0x2f, 0x76, 0x35, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e,
	0x2f, 0x6e, 0x65, 0x74, 0xaa, 0x02, 0x15, 0x56, 0x32, 0x52, 0x61, 0x79, 0x2e, 0x43, 0x6f, 0x72,
	0x65, 0x2e, 0x43, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x4e, 0x65, 0x74, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_common_net_uid_proto_rawDescOnce sync.Once
	file_common_net_uid_proto_rawDescData = file_common_net_uid_proto_rawDesc
)

func file_common_net_uid_proto_rawDescGZIP() []byte {
	file_common_net_uid_proto_rawDescOnce.Do(func() {
		file_common_net_uid_proto_rawDescData = protoimpl.X.CompressGZIP(file_common_net_uid_proto_rawDescData)
	})
	return file_common_net_uid_proto_rawDescData
}

var file_common_net_uid_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_common_net_uid_proto_goTypes = []interface{}{
	(*UidList)(nil), // 0: v2ray.core.common.net.UidList
}
var file_common_net_uid_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_common_net_uid_proto_init() }
func file_common_net_uid_proto_init() {
	if File_common_net_uid_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_common_net_uid_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UidList); i {
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
			RawDescriptor: file_common_net_uid_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_common_net_uid_proto_goTypes,
		DependencyIndexes: file_common_net_uid_proto_depIdxs,
		MessageInfos:      file_common_net_uid_proto_msgTypes,
	}.Build()
	File_common_net_uid_proto = out.File
	file_common_net_uid_proto_rawDesc = nil
	file_common_net_uid_proto_goTypes = nil
	file_common_net_uid_proto_depIdxs = nil
}
