// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.23.0
// 	protoc        v3.12.3
// source: proto/crypto.proto

package config

import (
	proto "github.com/golang/protobuf/proto"
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

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type CryptoConfig struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// A string defining what key type is being stored. Eg. Secp256k1.Type()
	KeyType string `protobuf:"bytes,1,opt,name=KeyType,proto3" json:"KeyType,omitempty"`
	// Private Keys
	PvtKey []byte `protobuf:"bytes,2,opt,name=PvtKey,proto3" json:"PvtKey,omitempty"`
	// Mapping between a node and its public key
	NodeKeyMap map[uint64][]byte `protobuf:"bytes,3,rep,name=NodeKeyMap,proto3" json:"NodeKeyMap,omitempty" protobuf_key:"varint,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *CryptoConfig) Reset() {
	*x = CryptoConfig{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_crypto_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CryptoConfig) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CryptoConfig) ProtoMessage() {}

func (x *CryptoConfig) ProtoReflect() protoreflect.Message {
	mi := &file_proto_crypto_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CryptoConfig.ProtoReflect.Descriptor instead.
func (*CryptoConfig) Descriptor() ([]byte, []int) {
	return file_proto_crypto_proto_rawDescGZIP(), []int{0}
}

func (x *CryptoConfig) GetKeyType() string {
	if x != nil {
		return x.KeyType
	}
	return ""
}

func (x *CryptoConfig) GetPvtKey() []byte {
	if x != nil {
		return x.PvtKey
	}
	return nil
}

func (x *CryptoConfig) GetNodeKeyMap() map[uint64][]byte {
	if x != nil {
		return x.NodeKeyMap
	}
	return nil
}

var File_proto_crypto_proto protoreflect.FileDescriptor

var file_proto_crypto_proto_rawDesc = []byte{
	0x0a, 0x12, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x63, 0x72, 0x79, 0x70, 0x74, 0x6f, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x22, 0xc5, 0x01, 0x0a,
	0x0c, 0x43, 0x72, 0x79, 0x70, 0x74, 0x6f, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0x18, 0x0a,
	0x07, 0x4b, 0x65, 0x79, 0x54, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07,
	0x4b, 0x65, 0x79, 0x54, 0x79, 0x70, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x50, 0x76, 0x74, 0x4b, 0x65,
	0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x06, 0x50, 0x76, 0x74, 0x4b, 0x65, 0x79, 0x12,
	0x44, 0x0a, 0x0a, 0x4e, 0x6f, 0x64, 0x65, 0x4b, 0x65, 0x79, 0x4d, 0x61, 0x70, 0x18, 0x03, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x24, 0x2e, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x43, 0x72, 0x79,
	0x70, 0x74, 0x6f, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x4e, 0x6f, 0x64, 0x65, 0x4b, 0x65,
	0x79, 0x4d, 0x61, 0x70, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x0a, 0x4e, 0x6f, 0x64, 0x65, 0x4b,
	0x65, 0x79, 0x4d, 0x61, 0x70, 0x1a, 0x3d, 0x0a, 0x0f, 0x4e, 0x6f, 0x64, 0x65, 0x4b, 0x65, 0x79,
	0x4d, 0x61, 0x70, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x3a, 0x02, 0x38, 0x01, 0x42, 0x2c, 0x5a, 0x2a, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63,
	0x6f, 0x6d, 0x2f, 0x61, 0x64, 0x69, 0x74, 0x68, 0x79, 0x61, 0x62, 0x68, 0x61, 0x74, 0x6b, 0x61,
	0x6a, 0x61, 0x6b, 0x65, 0x2f, 0x6c, 0x69, 0x62, 0x65, 0x32, 0x63, 0x2f, 0x63, 0x6f, 0x6e, 0x66,
	0x69, 0x67, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_crypto_proto_rawDescOnce sync.Once
	file_proto_crypto_proto_rawDescData = file_proto_crypto_proto_rawDesc
)

func file_proto_crypto_proto_rawDescGZIP() []byte {
	file_proto_crypto_proto_rawDescOnce.Do(func() {
		file_proto_crypto_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_crypto_proto_rawDescData)
	})
	return file_proto_crypto_proto_rawDescData
}

var file_proto_crypto_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_proto_crypto_proto_goTypes = []interface{}{
	(*CryptoConfig)(nil), // 0: config.CryptoConfig
	nil,                  // 1: config.CryptoConfig.NodeKeyMapEntry
}
var file_proto_crypto_proto_depIdxs = []int32{
	1, // 0: config.CryptoConfig.NodeKeyMap:type_name -> config.CryptoConfig.NodeKeyMapEntry
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_proto_crypto_proto_init() }
func file_proto_crypto_proto_init() {
	if File_proto_crypto_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_crypto_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CryptoConfig); i {
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
			RawDescriptor: file_proto_crypto_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_proto_crypto_proto_goTypes,
		DependencyIndexes: file_proto_crypto_proto_depIdxs,
		MessageInfos:      file_proto_crypto_proto_msgTypes,
	}.Build()
	File_proto_crypto_proto = out.File
	file_proto_crypto_proto_rawDesc = nil
	file_proto_crypto_proto_goTypes = nil
	file_proto_crypto_proto_depIdxs = nil
}
