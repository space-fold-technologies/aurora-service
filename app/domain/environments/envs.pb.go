// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.21.1
// source: envs.proto

package environments

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

type CreateEnvEntryOrder_Scope int32

const (
	CreateEnvEntryOrder_CLUSTER   CreateEnvEntryOrder_Scope = 0
	CreateEnvEntryOrder_NODE      CreateEnvEntryOrder_Scope = 1
	CreateEnvEntryOrder_CONTAINER CreateEnvEntryOrder_Scope = 2
)

// Enum value maps for CreateEnvEntryOrder_Scope.
var (
	CreateEnvEntryOrder_Scope_name = map[int32]string{
		0: "CLUSTER",
		1: "NODE",
		2: "CONTAINER",
	}
	CreateEnvEntryOrder_Scope_value = map[string]int32{
		"CLUSTER":   0,
		"NODE":      1,
		"CONTAINER": 2,
	}
)

func (x CreateEnvEntryOrder_Scope) Enum() *CreateEnvEntryOrder_Scope {
	p := new(CreateEnvEntryOrder_Scope)
	*p = x
	return p
}

func (x CreateEnvEntryOrder_Scope) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (CreateEnvEntryOrder_Scope) Descriptor() protoreflect.EnumDescriptor {
	return file_envs_proto_enumTypes[0].Descriptor()
}

func (CreateEnvEntryOrder_Scope) Type() protoreflect.EnumType {
	return &file_envs_proto_enumTypes[0]
}

func (x CreateEnvEntryOrder_Scope) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use CreateEnvEntryOrder_Scope.Descriptor instead.
func (CreateEnvEntryOrder_Scope) EnumDescriptor() ([]byte, []int) {
	return file_envs_proto_rawDescGZIP(), []int{1, 0}
}

type EnvEntry struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key   string `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Value string `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *EnvEntry) Reset() {
	*x = EnvEntry{}
	if protoimpl.UnsafeEnabled {
		mi := &file_envs_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EnvEntry) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EnvEntry) ProtoMessage() {}

func (x *EnvEntry) ProtoReflect() protoreflect.Message {
	mi := &file_envs_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EnvEntry.ProtoReflect.Descriptor instead.
func (*EnvEntry) Descriptor() ([]byte, []int) {
	return file_envs_proto_rawDescGZIP(), []int{0}
}

func (x *EnvEntry) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *EnvEntry) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

type CreateEnvEntryOrder struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Scope   CreateEnvEntryOrder_Scope `protobuf:"varint,1,opt,name=scope,proto3,enum=CreateEnvEntryOrder_Scope" json:"scope,omitempty"`
	Target  string                    `protobuf:"bytes,2,opt,name=target,proto3" json:"target,omitempty"`
	Entries []*EnvEntry               `protobuf:"bytes,3,rep,name=entries,proto3" json:"entries,omitempty"`
}

func (x *CreateEnvEntryOrder) Reset() {
	*x = CreateEnvEntryOrder{}
	if protoimpl.UnsafeEnabled {
		mi := &file_envs_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateEnvEntryOrder) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateEnvEntryOrder) ProtoMessage() {}

func (x *CreateEnvEntryOrder) ProtoReflect() protoreflect.Message {
	mi := &file_envs_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateEnvEntryOrder.ProtoReflect.Descriptor instead.
func (*CreateEnvEntryOrder) Descriptor() ([]byte, []int) {
	return file_envs_proto_rawDescGZIP(), []int{1}
}

func (x *CreateEnvEntryOrder) GetScope() CreateEnvEntryOrder_Scope {
	if x != nil {
		return x.Scope
	}
	return CreateEnvEntryOrder_CLUSTER
}

func (x *CreateEnvEntryOrder) GetTarget() string {
	if x != nil {
		return x.Target
	}
	return ""
}

func (x *CreateEnvEntryOrder) GetEntries() []*EnvEntry {
	if x != nil {
		return x.Entries
	}
	return nil
}

type EnvSet struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Entries []*EnvEntry `protobuf:"bytes,1,rep,name=entries,proto3" json:"entries,omitempty"`
}

func (x *EnvSet) Reset() {
	*x = EnvSet{}
	if protoimpl.UnsafeEnabled {
		mi := &file_envs_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EnvSet) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EnvSet) ProtoMessage() {}

func (x *EnvSet) ProtoReflect() protoreflect.Message {
	mi := &file_envs_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EnvSet.ProtoReflect.Descriptor instead.
func (*EnvSet) Descriptor() ([]byte, []int) {
	return file_envs_proto_rawDescGZIP(), []int{2}
}

func (x *EnvSet) GetEntries() []*EnvEntry {
	if x != nil {
		return x.Entries
	}
	return nil
}

type RemoveEnvEntryOrder struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Keys []string `protobuf:"bytes,1,rep,name=keys,proto3" json:"keys,omitempty"`
}

func (x *RemoveEnvEntryOrder) Reset() {
	*x = RemoveEnvEntryOrder{}
	if protoimpl.UnsafeEnabled {
		mi := &file_envs_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RemoveEnvEntryOrder) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RemoveEnvEntryOrder) ProtoMessage() {}

func (x *RemoveEnvEntryOrder) ProtoReflect() protoreflect.Message {
	mi := &file_envs_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RemoveEnvEntryOrder.ProtoReflect.Descriptor instead.
func (*RemoveEnvEntryOrder) Descriptor() ([]byte, []int) {
	return file_envs_proto_rawDescGZIP(), []int{3}
}

func (x *RemoveEnvEntryOrder) GetKeys() []string {
	if x != nil {
		return x.Keys
	}
	return nil
}

var File_envs_proto protoreflect.FileDescriptor

var file_envs_proto_rawDesc = []byte{
	0x0a, 0x0a, 0x65, 0x6e, 0x76, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x32, 0x0a, 0x08,
	0x45, 0x6e, 0x76, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x22, 0xb3, 0x01, 0x0a, 0x13, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x45, 0x6e, 0x76, 0x45, 0x6e,
	0x74, 0x72, 0x79, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x12, 0x30, 0x0a, 0x05, 0x73, 0x63, 0x6f, 0x70,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1a, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65,
	0x45, 0x6e, 0x76, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x2e, 0x53, 0x63,
	0x6f, 0x70, 0x65, 0x52, 0x05, 0x73, 0x63, 0x6f, 0x70, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x74, 0x61,
	0x72, 0x67, 0x65, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x74, 0x61, 0x72, 0x67,
	0x65, 0x74, 0x12, 0x23, 0x0a, 0x07, 0x65, 0x6e, 0x74, 0x72, 0x69, 0x65, 0x73, 0x18, 0x03, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x09, 0x2e, 0x45, 0x6e, 0x76, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x07,
	0x65, 0x6e, 0x74, 0x72, 0x69, 0x65, 0x73, 0x22, 0x2d, 0x0a, 0x05, 0x53, 0x63, 0x6f, 0x70, 0x65,
	0x12, 0x0b, 0x0a, 0x07, 0x43, 0x4c, 0x55, 0x53, 0x54, 0x45, 0x52, 0x10, 0x00, 0x12, 0x08, 0x0a,
	0x04, 0x4e, 0x4f, 0x44, 0x45, 0x10, 0x01, 0x12, 0x0d, 0x0a, 0x09, 0x43, 0x4f, 0x4e, 0x54, 0x41,
	0x49, 0x4e, 0x45, 0x52, 0x10, 0x02, 0x22, 0x2d, 0x0a, 0x06, 0x45, 0x6e, 0x76, 0x53, 0x65, 0x74,
	0x12, 0x23, 0x0a, 0x07, 0x65, 0x6e, 0x74, 0x72, 0x69, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x09, 0x2e, 0x45, 0x6e, 0x76, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x07, 0x65, 0x6e,
	0x74, 0x72, 0x69, 0x65, 0x73, 0x22, 0x29, 0x0a, 0x13, 0x52, 0x65, 0x6d, 0x6f, 0x76, 0x65, 0x45,
	0x6e, 0x76, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x12, 0x12, 0x0a, 0x04,
	0x6b, 0x65, 0x79, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x04, 0x6b, 0x65, 0x79, 0x73,
	0x42, 0x10, 0x5a, 0x0e, 0x2e, 0x3b, 0x65, 0x6e, 0x76, 0x69, 0x72, 0x6f, 0x6e, 0x6d, 0x65, 0x6e,
	0x74, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_envs_proto_rawDescOnce sync.Once
	file_envs_proto_rawDescData = file_envs_proto_rawDesc
)

func file_envs_proto_rawDescGZIP() []byte {
	file_envs_proto_rawDescOnce.Do(func() {
		file_envs_proto_rawDescData = protoimpl.X.CompressGZIP(file_envs_proto_rawDescData)
	})
	return file_envs_proto_rawDescData
}

var file_envs_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_envs_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_envs_proto_goTypes = []interface{}{
	(CreateEnvEntryOrder_Scope)(0), // 0: CreateEnvEntryOrder.Scope
	(*EnvEntry)(nil),               // 1: EnvEntry
	(*CreateEnvEntryOrder)(nil),    // 2: CreateEnvEntryOrder
	(*EnvSet)(nil),                 // 3: EnvSet
	(*RemoveEnvEntryOrder)(nil),    // 4: RemoveEnvEntryOrder
}
var file_envs_proto_depIdxs = []int32{
	0, // 0: CreateEnvEntryOrder.scope:type_name -> CreateEnvEntryOrder.Scope
	1, // 1: CreateEnvEntryOrder.entries:type_name -> EnvEntry
	1, // 2: EnvSet.entries:type_name -> EnvEntry
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_envs_proto_init() }
func file_envs_proto_init() {
	if File_envs_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_envs_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EnvEntry); i {
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
		file_envs_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateEnvEntryOrder); i {
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
		file_envs_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EnvSet); i {
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
		file_envs_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RemoveEnvEntryOrder); i {
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
			RawDescriptor: file_envs_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_envs_proto_goTypes,
		DependencyIndexes: file_envs_proto_depIdxs,
		EnumInfos:         file_envs_proto_enumTypes,
		MessageInfos:      file_envs_proto_msgTypes,
	}.Build()
	File_envs_proto = out.File
	file_envs_proto_rawDesc = nil
	file_envs_proto_goTypes = nil
	file_envs_proto_depIdxs = nil
}
