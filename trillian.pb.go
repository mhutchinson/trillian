// Code generated by protoc-gen-go. DO NOT EDIT.
// source: trillian.proto

package trillian

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	any "github.com/golang/protobuf/ptypes/any"
	duration "github.com/golang/protobuf/ptypes/duration"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
	keyspb "github.com/google/trillian/crypto/keyspb"
	sigpb "github.com/google/trillian/crypto/sigpb"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// LogRootFormat specifies the fields that are covered by the
// SignedLogRoot signature, as well as their ordering and formats.
type LogRootFormat int32

const (
	LogRootFormat_LOG_ROOT_FORMAT_UNKNOWN LogRootFormat = 0
	LogRootFormat_LOG_ROOT_FORMAT_V1      LogRootFormat = 1
)

var LogRootFormat_name = map[int32]string{
	0: "LOG_ROOT_FORMAT_UNKNOWN",
	1: "LOG_ROOT_FORMAT_V1",
}

var LogRootFormat_value = map[string]int32{
	"LOG_ROOT_FORMAT_UNKNOWN": 0,
	"LOG_ROOT_FORMAT_V1":      1,
}

func (x LogRootFormat) String() string {
	return proto.EnumName(LogRootFormat_name, int32(x))
}

func (LogRootFormat) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_364603a4e17a2a56, []int{0}
}

// MapRootFormat specifies the fields that are covered by the
// SignedMapRoot signature, as well as their ordering and formats.
type MapRootFormat int32

const (
	MapRootFormat_MAP_ROOT_FORMAT_UNKNOWN MapRootFormat = 0
	MapRootFormat_MAP_ROOT_FORMAT_V1      MapRootFormat = 1
)

var MapRootFormat_name = map[int32]string{
	0: "MAP_ROOT_FORMAT_UNKNOWN",
	1: "MAP_ROOT_FORMAT_V1",
}

var MapRootFormat_value = map[string]int32{
	"MAP_ROOT_FORMAT_UNKNOWN": 0,
	"MAP_ROOT_FORMAT_V1":      1,
}

func (x MapRootFormat) String() string {
	return proto.EnumName(MapRootFormat_name, int32(x))
}

func (MapRootFormat) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_364603a4e17a2a56, []int{1}
}

// Defines the way empty / node / leaf hashes are constructed incorporating
// preimage protection, which can be application specific.
type HashStrategy int32

const (
	// Hash strategy cannot be determined. Included to enable detection of
	// mismatched proto versions being used. Represents an invalid value.
	HashStrategy_UNKNOWN_HASH_STRATEGY HashStrategy = 0
	// Certificate Transparency strategy: leaf hash prefix = 0x00, node prefix =
	// 0x01, empty hash is digest([]byte{}), as defined in the specification.
	HashStrategy_RFC6962_SHA256 HashStrategy = 1
	// Sparse Merkle Tree strategy:  leaf hash prefix = 0x00, node prefix = 0x01,
	// empty branch is recursively computed from empty leaf nodes.
	// NOT secure in a multi tree environment. For testing only.
	HashStrategy_TEST_MAP_HASHER HashStrategy = 2
	// Append-only log strategy where leaf nodes are defined as the ObjectHash.
	// All other properties are equal to RFC6962_SHA256.
	HashStrategy_OBJECT_RFC6962_SHA256 HashStrategy = 3
	// The CONIKS sparse tree hasher with SHA512_256 as the hash algorithm.
	HashStrategy_CONIKS_SHA512_256 HashStrategy = 4
	// The CONIKS sparse tree hasher with SHA256 as the hash algorithm.
	HashStrategy_CONIKS_SHA256 HashStrategy = 5
)

var HashStrategy_name = map[int32]string{
	0: "UNKNOWN_HASH_STRATEGY",
	1: "RFC6962_SHA256",
	2: "TEST_MAP_HASHER",
	3: "OBJECT_RFC6962_SHA256",
	4: "CONIKS_SHA512_256",
	5: "CONIKS_SHA256",
}

var HashStrategy_value = map[string]int32{
	"UNKNOWN_HASH_STRATEGY": 0,
	"RFC6962_SHA256":        1,
	"TEST_MAP_HASHER":       2,
	"OBJECT_RFC6962_SHA256": 3,
	"CONIKS_SHA512_256":     4,
	"CONIKS_SHA256":         5,
}

func (x HashStrategy) String() string {
	return proto.EnumName(HashStrategy_name, int32(x))
}

func (HashStrategy) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_364603a4e17a2a56, []int{2}
}

// State of the tree.
type TreeState int32

const (
	// Tree state cannot be determined. Included to enable detection of
	// mismatched proto versions being used. Represents an invalid value.
	TreeState_UNKNOWN_TREE_STATE TreeState = 0
	// Active trees are able to respond to both read and write requests.
	TreeState_ACTIVE TreeState = 1
	// Frozen trees are only able to respond to read requests, writing to a frozen
	// tree is forbidden. Trees should not be frozen when there are entries
	// in the queue that have not yet been integrated. See the DRAINING
	// state for this case.
	TreeState_FROZEN TreeState = 2
	// Deprecated: now tracked in Tree.deleted.
	TreeState_DEPRECATED_SOFT_DELETED TreeState = 3 // Deprecated: Do not use.
	// Deprecated: now tracked in Tree.deleted.
	TreeState_DEPRECATED_HARD_DELETED TreeState = 4 // Deprecated: Do not use.
	// A tree that is draining will continue to integrate queued entries.
	// No new entries should be accepted.
	TreeState_DRAINING TreeState = 5
)

var TreeState_name = map[int32]string{
	0: "UNKNOWN_TREE_STATE",
	1: "ACTIVE",
	2: "FROZEN",
	3: "DEPRECATED_SOFT_DELETED",
	4: "DEPRECATED_HARD_DELETED",
	5: "DRAINING",
}

var TreeState_value = map[string]int32{
	"UNKNOWN_TREE_STATE":      0,
	"ACTIVE":                  1,
	"FROZEN":                  2,
	"DEPRECATED_SOFT_DELETED": 3,
	"DEPRECATED_HARD_DELETED": 4,
	"DRAINING":                5,
}

func (x TreeState) String() string {
	return proto.EnumName(TreeState_name, int32(x))
}

func (TreeState) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_364603a4e17a2a56, []int{3}
}

// Type of the tree.
type TreeType int32

const (
	// Tree type cannot be determined. Included to enable detection of mismatched
	// proto versions being used. Represents an invalid value.
	TreeType_UNKNOWN_TREE_TYPE TreeType = 0
	// Tree represents a verifiable log.
	TreeType_LOG TreeType = 1
	// Tree represents a verifiable map.
	TreeType_MAP TreeType = 2
	// Tree represents a verifiable pre-ordered log, i.e., a log whose entries are
	// placed according to sequence numbers assigned outside of Trillian.
	TreeType_PREORDERED_LOG TreeType = 3
)

var TreeType_name = map[int32]string{
	0: "UNKNOWN_TREE_TYPE",
	1: "LOG",
	2: "MAP",
	3: "PREORDERED_LOG",
}

var TreeType_value = map[string]int32{
	"UNKNOWN_TREE_TYPE": 0,
	"LOG":               1,
	"MAP":               2,
	"PREORDERED_LOG":    3,
}

func (x TreeType) String() string {
	return proto.EnumName(TreeType_name, int32(x))
}

func (TreeType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_364603a4e17a2a56, []int{4}
}

// Represents a tree, which may be either a verifiable log or map.
// Readonly attributes are assigned at tree creation, after which they may not
// be modified.
//
// Note: Many APIs within the rest of the code require these objects to
// be provided. For safety they should be obtained via Admin API calls and
// not created dynamically.
type Tree struct {
	// ID of the tree.
	// Readonly.
	TreeId int64 `protobuf:"varint,1,opt,name=tree_id,json=treeId,proto3" json:"tree_id,omitempty"`
	// State of the tree.
	// Trees are ACTIVE after creation. At any point the tree may transition
	// between ACTIVE, DRAINING and FROZEN states.
	TreeState TreeState `protobuf:"varint,2,opt,name=tree_state,json=treeState,proto3,enum=trillian.TreeState" json:"tree_state,omitempty"`
	// Type of the tree.
	// Readonly after Tree creation. Exception: Can be switched from
	// PREORDERED_LOG to LOG if the Tree is and remains in the FROZEN state.
	TreeType TreeType `protobuf:"varint,3,opt,name=tree_type,json=treeType,proto3,enum=trillian.TreeType" json:"tree_type,omitempty"`
	// Hash strategy to be used by the tree.
	// Readonly.
	HashStrategy HashStrategy `protobuf:"varint,4,opt,name=hash_strategy,json=hashStrategy,proto3,enum=trillian.HashStrategy" json:"hash_strategy,omitempty"`
	// Hash algorithm to be used by the tree.
	// Readonly.
	HashAlgorithm sigpb.DigitallySigned_HashAlgorithm `protobuf:"varint,5,opt,name=hash_algorithm,json=hashAlgorithm,proto3,enum=sigpb.DigitallySigned_HashAlgorithm" json:"hash_algorithm,omitempty"`
	// Signature algorithm to be used by the tree.
	// Readonly.
	SignatureAlgorithm sigpb.DigitallySigned_SignatureAlgorithm `protobuf:"varint,6,opt,name=signature_algorithm,json=signatureAlgorithm,proto3,enum=sigpb.DigitallySigned_SignatureAlgorithm" json:"signature_algorithm,omitempty"`
	// Display name of the tree.
	// Optional.
	DisplayName string `protobuf:"bytes,8,opt,name=display_name,json=displayName,proto3" json:"display_name,omitempty"`
	// Description of the tree,
	// Optional.
	Description string `protobuf:"bytes,9,opt,name=description,proto3" json:"description,omitempty"`
	// Identifies the private key used for signing tree heads and entry
	// timestamps.
	// This can be any type of message to accommodate different key management
	// systems, e.g. PEM files, HSMs, etc.
	// Private keys are write-only: they're never returned by RPCs.
	// The private_key message can be changed after a tree is created, but the
	// underlying key must remain the same - this is to enable migrating a key
	// from one provider to another.
	PrivateKey *any.Any `protobuf:"bytes,12,opt,name=private_key,json=privateKey,proto3" json:"private_key,omitempty"`
	// Storage-specific settings.
	// Varies according to the storage implementation backing Trillian.
	StorageSettings *any.Any `protobuf:"bytes,13,opt,name=storage_settings,json=storageSettings,proto3" json:"storage_settings,omitempty"`
	// The public key used for verifying tree heads and entry timestamps.
	// Readonly.
	PublicKey *keyspb.PublicKey `protobuf:"bytes,14,opt,name=public_key,json=publicKey,proto3" json:"public_key,omitempty"`
	// Interval after which a new signed root is produced even if there have been
	// no submission.  If zero, this behavior is disabled.
	MaxRootDuration *duration.Duration `protobuf:"bytes,15,opt,name=max_root_duration,json=maxRootDuration,proto3" json:"max_root_duration,omitempty"`
	// Time of tree creation.
	// Readonly.
	CreateTime *timestamp.Timestamp `protobuf:"bytes,16,opt,name=create_time,json=createTime,proto3" json:"create_time,omitempty"`
	// Time of last tree update.
	// Readonly (automatically assigned on updates).
	UpdateTime *timestamp.Timestamp `protobuf:"bytes,17,opt,name=update_time,json=updateTime,proto3" json:"update_time,omitempty"`
	// If true, the tree has been deleted.
	// Deleted trees may be undeleted during a certain time window, after which
	// they're permanently deleted (and unrecoverable).
	// Readonly.
	Deleted bool `protobuf:"varint,19,opt,name=deleted,proto3" json:"deleted,omitempty"`
	// Time of tree deletion, if any.
	// Readonly.
	DeleteTime *timestamp.Timestamp `protobuf:"bytes,20,opt,name=delete_time,json=deleteTime,proto3" json:"delete_time,omitempty"`
	// The number of 8-bit prefix strata in the tree layout. Applies to maps only.
	// The final stratum spans the remaining 256-8*prefix_strata bits.
	// For example, if prefix_strata=3 then there will be 4 tiles in any path:
	//  * 3 prefix tiles of 1 byte each; and
	//  * 1 final leaf tile of 29 bytes height.
	// Readonly.
	PrefixStrata         int32    `protobuf:"varint,21,opt,name=prefix_strata,json=prefixStrata,proto3" json:"prefix_strata,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Tree) Reset()         { *m = Tree{} }
func (m *Tree) String() string { return proto.CompactTextString(m) }
func (*Tree) ProtoMessage()    {}
func (*Tree) Descriptor() ([]byte, []int) {
	return fileDescriptor_364603a4e17a2a56, []int{0}
}

func (m *Tree) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Tree.Unmarshal(m, b)
}
func (m *Tree) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Tree.Marshal(b, m, deterministic)
}
func (m *Tree) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Tree.Merge(m, src)
}
func (m *Tree) XXX_Size() int {
	return xxx_messageInfo_Tree.Size(m)
}
func (m *Tree) XXX_DiscardUnknown() {
	xxx_messageInfo_Tree.DiscardUnknown(m)
}

var xxx_messageInfo_Tree proto.InternalMessageInfo

func (m *Tree) GetTreeId() int64 {
	if m != nil {
		return m.TreeId
	}
	return 0
}

func (m *Tree) GetTreeState() TreeState {
	if m != nil {
		return m.TreeState
	}
	return TreeState_UNKNOWN_TREE_STATE
}

func (m *Tree) GetTreeType() TreeType {
	if m != nil {
		return m.TreeType
	}
	return TreeType_UNKNOWN_TREE_TYPE
}

func (m *Tree) GetHashStrategy() HashStrategy {
	if m != nil {
		return m.HashStrategy
	}
	return HashStrategy_UNKNOWN_HASH_STRATEGY
}

func (m *Tree) GetHashAlgorithm() sigpb.DigitallySigned_HashAlgorithm {
	if m != nil {
		return m.HashAlgorithm
	}
	return sigpb.DigitallySigned_NONE
}

func (m *Tree) GetSignatureAlgorithm() sigpb.DigitallySigned_SignatureAlgorithm {
	if m != nil {
		return m.SignatureAlgorithm
	}
	return sigpb.DigitallySigned_ANONYMOUS
}

func (m *Tree) GetDisplayName() string {
	if m != nil {
		return m.DisplayName
	}
	return ""
}

func (m *Tree) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func (m *Tree) GetPrivateKey() *any.Any {
	if m != nil {
		return m.PrivateKey
	}
	return nil
}

func (m *Tree) GetStorageSettings() *any.Any {
	if m != nil {
		return m.StorageSettings
	}
	return nil
}

func (m *Tree) GetPublicKey() *keyspb.PublicKey {
	if m != nil {
		return m.PublicKey
	}
	return nil
}

func (m *Tree) GetMaxRootDuration() *duration.Duration {
	if m != nil {
		return m.MaxRootDuration
	}
	return nil
}

func (m *Tree) GetCreateTime() *timestamp.Timestamp {
	if m != nil {
		return m.CreateTime
	}
	return nil
}

func (m *Tree) GetUpdateTime() *timestamp.Timestamp {
	if m != nil {
		return m.UpdateTime
	}
	return nil
}

func (m *Tree) GetDeleted() bool {
	if m != nil {
		return m.Deleted
	}
	return false
}

func (m *Tree) GetDeleteTime() *timestamp.Timestamp {
	if m != nil {
		return m.DeleteTime
	}
	return nil
}

func (m *Tree) GetPrefixStrata() int32 {
	if m != nil {
		return m.PrefixStrata
	}
	return 0
}

type SignedEntryTimestamp struct {
	TimestampNanos       int64                  `protobuf:"varint,1,opt,name=timestamp_nanos,json=timestampNanos,proto3" json:"timestamp_nanos,omitempty"`
	LogId                int64                  `protobuf:"varint,2,opt,name=log_id,json=logId,proto3" json:"log_id,omitempty"`
	Signature            *sigpb.DigitallySigned `protobuf:"bytes,3,opt,name=signature,proto3" json:"signature,omitempty"`
	XXX_NoUnkeyedLiteral struct{}               `json:"-"`
	XXX_unrecognized     []byte                 `json:"-"`
	XXX_sizecache        int32                  `json:"-"`
}

func (m *SignedEntryTimestamp) Reset()         { *m = SignedEntryTimestamp{} }
func (m *SignedEntryTimestamp) String() string { return proto.CompactTextString(m) }
func (*SignedEntryTimestamp) ProtoMessage()    {}
func (*SignedEntryTimestamp) Descriptor() ([]byte, []int) {
	return fileDescriptor_364603a4e17a2a56, []int{1}
}

func (m *SignedEntryTimestamp) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SignedEntryTimestamp.Unmarshal(m, b)
}
func (m *SignedEntryTimestamp) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SignedEntryTimestamp.Marshal(b, m, deterministic)
}
func (m *SignedEntryTimestamp) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SignedEntryTimestamp.Merge(m, src)
}
func (m *SignedEntryTimestamp) XXX_Size() int {
	return xxx_messageInfo_SignedEntryTimestamp.Size(m)
}
func (m *SignedEntryTimestamp) XXX_DiscardUnknown() {
	xxx_messageInfo_SignedEntryTimestamp.DiscardUnknown(m)
}

var xxx_messageInfo_SignedEntryTimestamp proto.InternalMessageInfo

func (m *SignedEntryTimestamp) GetTimestampNanos() int64 {
	if m != nil {
		return m.TimestampNanos
	}
	return 0
}

func (m *SignedEntryTimestamp) GetLogId() int64 {
	if m != nil {
		return m.LogId
	}
	return 0
}

func (m *SignedEntryTimestamp) GetSignature() *sigpb.DigitallySigned {
	if m != nil {
		return m.Signature
	}
	return nil
}

// SignedLogRoot represents a commitment by a Log to a particular tree.
type SignedLogRoot struct {
	// key_hint is a hint to identify the public key for signature verification.
	// key_hint is not authenticated and may be incorrect or missing, in which
	// case all known public keys may be used to verify the signature.
	// When directly communicating with a Trillian gRPC server, the key_hint will
	// typically contain the LogID encoded as a big-endian 64-bit integer;
	// however, in other contexts the key_hint is likely to have different
	// contents (e.g. it could be a GUID, a URL + TreeID, or it could be
	// derived from the public key itself).
	KeyHint []byte `protobuf:"bytes,7,opt,name=key_hint,json=keyHint,proto3" json:"key_hint,omitempty"`
	// log_root holds the TLS-serialization of the following structure (described
	// in RFC5246 notation): Clients should validate log_root_signature with
	// VerifySignedLogRoot before deserializing log_root.
	// enum { v1(1), (65535)} Version;
	// struct {
	//   uint64 tree_size;
	//   opaque root_hash<0..128>;
	//   uint64 timestamp_nanos;
	//   uint64 revision;
	//   opaque metadata<0..65535>;
	// } LogRootV1;
	// struct {
	//   Version version;
	//   select(version) {
	//     case v1: LogRootV1;
	//   }
	// } LogRoot;
	//
	// A serialized v1 log root will therefore be laid out as:
	//
	// +---+---+---+---+---+---+---+---+---+---+---+---+---+---+-....--+
	// | ver=1 |          tree_size            |len|    root_hash      |
	// +---+---+---+---+---+---+---+---+---+---+---+---+---+---+-....--+
	//
	// +---+---+---+---+---+---+---+---+---+---+---+---+---+---+---+---+
	// |        timestamp_nanos        |      revision                 |
	// +---+---+---+---+---+---+---+---+---+---+---+---+---+---+---+---+
	//
	// +---+---+---+---+---+-....---+
	// |  len  |    metadata        |
	// +---+---+---+---+---+-....---+
	//
	// (with all integers encoded big-endian).
	LogRoot []byte `protobuf:"bytes,8,opt,name=log_root,json=logRoot,proto3" json:"log_root,omitempty"`
	// log_root_signature is the raw signature over log_root.
	LogRootSignature     []byte   `protobuf:"bytes,9,opt,name=log_root_signature,json=logRootSignature,proto3" json:"log_root_signature,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SignedLogRoot) Reset()         { *m = SignedLogRoot{} }
func (m *SignedLogRoot) String() string { return proto.CompactTextString(m) }
func (*SignedLogRoot) ProtoMessage()    {}
func (*SignedLogRoot) Descriptor() ([]byte, []int) {
	return fileDescriptor_364603a4e17a2a56, []int{2}
}

func (m *SignedLogRoot) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SignedLogRoot.Unmarshal(m, b)
}
func (m *SignedLogRoot) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SignedLogRoot.Marshal(b, m, deterministic)
}
func (m *SignedLogRoot) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SignedLogRoot.Merge(m, src)
}
func (m *SignedLogRoot) XXX_Size() int {
	return xxx_messageInfo_SignedLogRoot.Size(m)
}
func (m *SignedLogRoot) XXX_DiscardUnknown() {
	xxx_messageInfo_SignedLogRoot.DiscardUnknown(m)
}

var xxx_messageInfo_SignedLogRoot proto.InternalMessageInfo

func (m *SignedLogRoot) GetKeyHint() []byte {
	if m != nil {
		return m.KeyHint
	}
	return nil
}

func (m *SignedLogRoot) GetLogRoot() []byte {
	if m != nil {
		return m.LogRoot
	}
	return nil
}

func (m *SignedLogRoot) GetLogRootSignature() []byte {
	if m != nil {
		return m.LogRootSignature
	}
	return nil
}

// SignedMapRoot represents a commitment by a Map to a particular tree.
type SignedMapRoot struct {
	// map_root holds the TLS-serialization of the following structure (described
	// in RFC5246 notation): Clients should validate signature with
	// VerifySignedMapRoot before deserializing map_root.
	// enum { v1(1), (65535)} Version;
	// struct {
	//   opaque root_hash<0..128>;
	//   uint64 timestamp_nanos;
	//   uint64 revision;
	//   opaque metadata<0..65535>;
	// } MapRootV1;
	// struct {
	//   Version version;
	//   select(version) {
	//     case v1: MapRootV1;
	//   }
	// } MapRoot;
	MapRoot []byte `protobuf:"bytes,9,opt,name=map_root,json=mapRoot,proto3" json:"map_root,omitempty"`
	// Signature is the raw signature over MapRoot.
	Signature            []byte   `protobuf:"bytes,4,opt,name=signature,proto3" json:"signature,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SignedMapRoot) Reset()         { *m = SignedMapRoot{} }
func (m *SignedMapRoot) String() string { return proto.CompactTextString(m) }
func (*SignedMapRoot) ProtoMessage()    {}
func (*SignedMapRoot) Descriptor() ([]byte, []int) {
	return fileDescriptor_364603a4e17a2a56, []int{3}
}

func (m *SignedMapRoot) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SignedMapRoot.Unmarshal(m, b)
}
func (m *SignedMapRoot) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SignedMapRoot.Marshal(b, m, deterministic)
}
func (m *SignedMapRoot) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SignedMapRoot.Merge(m, src)
}
func (m *SignedMapRoot) XXX_Size() int {
	return xxx_messageInfo_SignedMapRoot.Size(m)
}
func (m *SignedMapRoot) XXX_DiscardUnknown() {
	xxx_messageInfo_SignedMapRoot.DiscardUnknown(m)
}

var xxx_messageInfo_SignedMapRoot proto.InternalMessageInfo

func (m *SignedMapRoot) GetMapRoot() []byte {
	if m != nil {
		return m.MapRoot
	}
	return nil
}

func (m *SignedMapRoot) GetSignature() []byte {
	if m != nil {
		return m.Signature
	}
	return nil
}

// Proof holds a consistency or inclusion proof for a Merkle tree, as returned
// by the API.
type Proof struct {
	// leaf_index indicates the requested leaf index when this message is used for
	// a leaf inclusion proof.  This field is set to zero when this message is
	// used for a consistency proof.
	LeafIndex            int64    `protobuf:"varint,1,opt,name=leaf_index,json=leafIndex,proto3" json:"leaf_index,omitempty"`
	Hashes               [][]byte `protobuf:"bytes,3,rep,name=hashes,proto3" json:"hashes,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Proof) Reset()         { *m = Proof{} }
func (m *Proof) String() string { return proto.CompactTextString(m) }
func (*Proof) ProtoMessage()    {}
func (*Proof) Descriptor() ([]byte, []int) {
	return fileDescriptor_364603a4e17a2a56, []int{4}
}

func (m *Proof) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Proof.Unmarshal(m, b)
}
func (m *Proof) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Proof.Marshal(b, m, deterministic)
}
func (m *Proof) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Proof.Merge(m, src)
}
func (m *Proof) XXX_Size() int {
	return xxx_messageInfo_Proof.Size(m)
}
func (m *Proof) XXX_DiscardUnknown() {
	xxx_messageInfo_Proof.DiscardUnknown(m)
}

var xxx_messageInfo_Proof proto.InternalMessageInfo

func (m *Proof) GetLeafIndex() int64 {
	if m != nil {
		return m.LeafIndex
	}
	return 0
}

func (m *Proof) GetHashes() [][]byte {
	if m != nil {
		return m.Hashes
	}
	return nil
}

func init() {
	proto.RegisterEnum("trillian.LogRootFormat", LogRootFormat_name, LogRootFormat_value)
	proto.RegisterEnum("trillian.MapRootFormat", MapRootFormat_name, MapRootFormat_value)
	proto.RegisterEnum("trillian.HashStrategy", HashStrategy_name, HashStrategy_value)
	proto.RegisterEnum("trillian.TreeState", TreeState_name, TreeState_value)
	proto.RegisterEnum("trillian.TreeType", TreeType_name, TreeType_value)
	proto.RegisterType((*Tree)(nil), "trillian.Tree")
	proto.RegisterType((*SignedEntryTimestamp)(nil), "trillian.SignedEntryTimestamp")
	proto.RegisterType((*SignedLogRoot)(nil), "trillian.SignedLogRoot")
	proto.RegisterType((*SignedMapRoot)(nil), "trillian.SignedMapRoot")
	proto.RegisterType((*Proof)(nil), "trillian.Proof")
}

func init() { proto.RegisterFile("trillian.proto", fileDescriptor_364603a4e17a2a56) }

var fileDescriptor_364603a4e17a2a56 = []byte{
	// 1121 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x56, 0xdf, 0x72, 0xda, 0x46,
	0x17, 0x8f, 0x40, 0x80, 0x38, 0x80, 0xbd, 0x5e, 0xc7, 0x8e, 0xcc, 0xf7, 0xb5, 0xa1, 0x6e, 0x67,
	0x4a, 0x33, 0x1d, 0xdc, 0xd0, 0x26, 0x33, 0x9d, 0x5c, 0x74, 0x14, 0x23, 0x1b, 0xb0, 0x0d, 0xcc,
	0xa2, 0xa6, 0x93, 0xdc, 0x68, 0xd6, 0x66, 0x2d, 0x34, 0x16, 0x92, 0x46, 0x5a, 0x77, 0xac, 0x67,
	0x68, 0xef, 0xf3, 0x16, 0x7d, 0xc6, 0xce, 0xae, 0x56, 0xe0, 0x38, 0x49, 0x73, 0x63, 0xef, 0x39,
	0xbf, 0x3f, 0xe7, 0xac, 0xf6, 0xac, 0x04, 0x6c, 0xf1, 0xc4, 0x0f, 0x02, 0x9f, 0x86, 0xbd, 0x38,
	0x89, 0x78, 0x84, 0x8d, 0x22, 0x6e, 0xb7, 0xaf, 0x92, 0x2c, 0xe6, 0xd1, 0xd1, 0x0d, 0xcb, 0xd2,
	0xf8, 0x52, 0xfd, 0xcb, 0x59, 0x6d, 0x53, 0x61, 0xa9, 0xef, 0xc5, 0x97, 0xf9, 0x5f, 0x85, 0x1c,
	0x78, 0x51, 0xe4, 0x05, 0xec, 0x48, 0x46, 0x97, 0xb7, 0xd7, 0x47, 0x34, 0xcc, 0x14, 0xf4, 0xf5,
	0x43, 0x68, 0x71, 0x9b, 0x50, 0xee, 0x47, 0xaa, 0x74, 0xfb, 0xe9, 0x43, 0x9c, 0xfb, 0x2b, 0x96,
	0x72, 0xba, 0x8a, 0x73, 0xc2, 0xe1, 0x3f, 0x35, 0xd0, 0x9d, 0x84, 0x31, 0xfc, 0x04, 0x6a, 0x3c,
	0x61, 0xcc, 0xf5, 0x17, 0xa6, 0xd6, 0xd1, 0xba, 0x65, 0x52, 0x15, 0xe1, 0x68, 0x81, 0xfb, 0x00,
	0x12, 0x48, 0x39, 0xe5, 0xcc, 0x2c, 0x75, 0xb4, 0xee, 0x56, 0x7f, 0xb7, 0xb7, 0xde, 0xa2, 0x10,
	0xcf, 0x05, 0x44, 0xea, 0xbc, 0x58, 0xe2, 0x23, 0x90, 0x81, 0xcb, 0xb3, 0x98, 0x99, 0x65, 0x29,
	0xc1, 0x1f, 0x4a, 0x9c, 0x2c, 0x66, 0xc4, 0xe0, 0x6a, 0x85, 0x5f, 0x41, 0x6b, 0x49, 0xd3, 0xa5,
	0x9b, 0xf2, 0x84, 0x72, 0xe6, 0x65, 0xa6, 0x2e, 0x45, 0xfb, 0x1b, 0xd1, 0x90, 0xa6, 0xcb, 0xb9,
	0x42, 0x49, 0x73, 0x79, 0x2f, 0xc2, 0x67, 0xb0, 0x25, 0xc5, 0x34, 0xf0, 0xa2, 0xc4, 0xe7, 0xcb,
	0x95, 0x59, 0x91, 0xea, 0xef, 0x7a, 0xf9, 0x53, 0x1c, 0xf8, 0x9e, 0xcf, 0x69, 0x10, 0x64, 0x73,
	0xdf, 0x0b, 0xd9, 0x42, 0x5a, 0x59, 0x05, 0x97, 0xc8, 0xc2, 0xeb, 0x10, 0xbf, 0x83, 0xdd, 0xd4,
	0xf7, 0x42, 0xca, 0x6f, 0x13, 0x76, 0xcf, 0xb1, 0x2a, 0x1d, 0x7f, 0xf8, 0x8c, 0xe3, 0xbc, 0x50,
	0x6c, 0x6c, 0x71, 0xfa, 0x51, 0x0e, 0x7f, 0x03, 0xcd, 0x85, 0x9f, 0xc6, 0x01, 0xcd, 0xdc, 0x90,
	0xae, 0x98, 0x69, 0x74, 0xb4, 0x6e, 0x9d, 0x34, 0x54, 0x6e, 0x42, 0x57, 0x0c, 0x77, 0xa0, 0xb1,
	0x60, 0xe9, 0x55, 0xe2, 0xc7, 0xe2, 0x14, 0xcd, 0xba, 0x62, 0x6c, 0x52, 0xf8, 0x05, 0x34, 0xe2,
	0xc4, 0xff, 0x93, 0x72, 0xe6, 0xde, 0xb0, 0xcc, 0x6c, 0x76, 0xb4, 0x6e, 0xa3, 0xff, 0xb8, 0x97,
	0x1f, 0x74, 0xaf, 0x38, 0xe8, 0x9e, 0x15, 0x66, 0x04, 0x14, 0xf1, 0x8c, 0x65, 0xf8, 0x37, 0x40,
	0x29, 0x8f, 0x12, 0xea, 0x31, 0x37, 0x65, 0x9c, 0xfb, 0xa1, 0x97, 0x9a, 0xad, 0xff, 0xd0, 0x6e,
	0x2b, 0xf6, 0x5c, 0x91, 0xf1, 0x4f, 0x00, 0xf1, 0xed, 0x65, 0xe0, 0x5f, 0xc9, 0xb2, 0x5b, 0x52,
	0xba, 0xd3, 0x53, 0x23, 0x3c, 0x93, 0xc8, 0x19, 0xcb, 0x48, 0x3d, 0x2e, 0x96, 0xd8, 0x86, 0x9d,
	0x15, 0xbd, 0x73, 0x93, 0x28, 0xe2, 0x6e, 0x31, 0x97, 0xe6, 0xb6, 0x14, 0x1e, 0x7c, 0x54, 0x73,
	0xa0, 0x08, 0x64, 0x7b, 0x45, 0xef, 0x48, 0x14, 0xf1, 0x22, 0x81, 0x5f, 0x41, 0xe3, 0x2a, 0x61,
	0x62, 0xbf, 0x62, 0x78, 0x4d, 0x24, 0x0d, 0xda, 0x1f, 0x19, 0x38, 0xc5, 0x64, 0x13, 0xc8, 0xe9,
	0x22, 0x21, 0xc4, 0xb7, 0xf1, 0x62, 0x2d, 0xde, 0xf9, 0xb2, 0x38, 0xa7, 0x4b, 0xb1, 0x09, 0xb5,
	0x05, 0x0b, 0x18, 0x67, 0x0b, 0x73, 0xb7, 0xa3, 0x75, 0x0d, 0x52, 0x84, 0xc2, 0x36, 0x5f, 0xe6,
	0xb6, 0x8f, 0xbf, 0x6c, 0x9b, 0xd3, 0xa5, 0xed, 0xb7, 0xd0, 0x8a, 0x13, 0x76, 0xed, 0xdf, 0xe5,
	0xe3, 0x4e, 0xcd, 0xbd, 0x8e, 0xd6, 0xad, 0x90, 0x66, 0x9e, 0x94, 0x63, 0x4d, 0xc7, 0xba, 0x81,
	0xd1, 0xee, 0x58, 0x37, 0x6a, 0xc8, 0x18, 0xeb, 0x06, 0xa0, 0xc6, 0x58, 0x37, 0x1a, 0xa8, 0x79,
	0xf8, 0xb7, 0x06, 0x8f, 0xf3, 0xa9, 0xb3, 0x43, 0x9e, 0x64, 0xeb, 0x0a, 0xf8, 0x7b, 0xd8, 0x5e,
	0x5f, 0x6e, 0x37, 0xa4, 0x61, 0x94, 0xaa, 0x8b, 0xbc, 0xb5, 0x4e, 0x4f, 0x44, 0x16, 0xef, 0x41,
	0x35, 0x88, 0x3c, 0x71, 0xd1, 0x4b, 0x12, 0xaf, 0x04, 0x91, 0x37, 0x5a, 0xe0, 0x5f, 0xa0, 0xbe,
	0x1e, 0x59, 0x79, 0x67, 0x1b, 0xfd, 0xfd, 0x4f, 0x8f, 0x3b, 0xd9, 0x10, 0x0f, 0xdf, 0x6b, 0xd0,
	0xca, 0xb3, 0xe7, 0x91, 0x27, 0x8e, 0x0d, 0x1f, 0x80, 0x71, 0xc3, 0x32, 0x77, 0xe9, 0x87, 0xdc,
	0xac, 0x75, 0xb4, 0x6e, 0x93, 0xd4, 0x6e, 0x58, 0x36, 0xf4, 0x43, 0x09, 0x89, 0xca, 0x62, 0x20,
	0xe4, 0xec, 0x37, 0x49, 0x2d, 0x50, 0xaa, 0x1f, 0x01, 0x17, 0x90, 0xbb, 0x69, 0xa3, 0x2e, 0x49,
	0x48, 0x91, 0xd6, 0xb7, 0x6c, 0xac, 0x1b, 0x1a, 0x2a, 0x8d, 0x75, 0xa3, 0x84, 0xca, 0x63, 0xdd,
	0x28, 0x23, 0x7d, 0xac, 0x1b, 0x3a, 0xaa, 0x8c, 0x75, 0xa3, 0x82, 0xaa, 0x63, 0xdd, 0xa8, 0xa2,
	0xda, 0x61, 0x52, 0x34, 0x76, 0x41, 0xe3, 0xa2, 0xb1, 0x15, 0x8d, 0xf3, 0xea, 0xb9, 0x71, 0x6d,
	0xa5, 0xa0, 0xff, 0xdf, 0xdf, 0xbb, 0x2e, 0xb1, 0x4d, 0xe2, 0x93, 0xd5, 0xd6, 0x75, 0xd6, 0x47,
	0x64, 0xa0, 0xfa, 0xe1, 0x00, 0x2a, 0xb3, 0x24, 0x8a, 0xae, 0xf1, 0x57, 0x00, 0x01, 0xa3, 0xd7,
	0xae, 0x1f, 0x2e, 0xd8, 0x9d, 0x3a, 0x87, 0xba, 0xc8, 0x8c, 0x44, 0x02, 0xef, 0x43, 0x55, 0xbc,
	0x75, 0x58, 0x6a, 0x96, 0x3b, 0xe5, 0x6e, 0x93, 0xa8, 0x28, 0xaf, 0xf1, 0x6c, 0x00, 0x2d, 0xf5,
	0x30, 0x4f, 0xa2, 0x64, 0x45, 0x39, 0xfe, 0x1f, 0x3c, 0x39, 0x9f, 0x9e, 0xba, 0x64, 0x3a, 0x75,
	0xdc, 0x93, 0x29, 0xb9, 0xb0, 0x1c, 0xf7, 0xf7, 0xc9, 0xd9, 0x64, 0xfa, 0xc7, 0x04, 0x3d, 0xc2,
	0xfb, 0x80, 0x1f, 0x82, 0x6f, 0x9e, 0x23, 0x4d, 0xb8, 0xa8, 0x9d, 0x6f, 0x5c, 0x2e, 0xac, 0xd9,
	0xe7, 0x5d, 0x1e, 0x82, 0xd2, 0xe5, 0xbd, 0x06, 0xcd, 0xfb, 0xaf, 0x5e, 0x7c, 0x00, 0x7b, 0x4a,
	0xe5, 0x0e, 0xad, 0xf9, 0xd0, 0x9d, 0x3b, 0xc4, 0x72, 0xec, 0xd3, 0xb7, 0xe8, 0x11, 0xc6, 0xb0,
	0x45, 0x4e, 0x8e, 0x5f, 0xfe, 0xfa, 0xb2, 0xef, 0xce, 0x87, 0x56, 0xff, 0xc5, 0x4b, 0xa4, 0xe1,
	0x5d, 0xd8, 0x76, 0xec, 0xb9, 0xe3, 0x0a, 0x73, 0xc1, 0xb7, 0x09, 0x2a, 0x09, 0x8f, 0xe9, 0xeb,
	0xb1, 0x7d, 0xec, 0xb8, 0x0f, 0xf8, 0x65, 0xbc, 0x07, 0x3b, 0xc7, 0xd3, 0xc9, 0xe8, 0x6c, 0x2e,
	0x52, 0x2f, 0x9e, 0xf7, 0x5d, 0x91, 0xd6, 0xf1, 0x0e, 0xb4, 0x36, 0x69, 0x91, 0xaa, 0x3c, 0xfb,
	0x4b, 0x83, 0xfa, 0xfa, 0xe3, 0x23, 0xfa, 0x2f, 0xda, 0x72, 0x88, 0x6d, 0xbb, 0x73, 0xc7, 0x72,
	0x6c, 0xf4, 0x08, 0x03, 0x54, 0xad, 0x63, 0x67, 0xf4, 0xc6, 0x46, 0x9a, 0x58, 0x9f, 0x90, 0xe9,
	0x3b, 0x7b, 0x82, 0x4a, 0xf8, 0x29, 0x3c, 0x19, 0xd8, 0x33, 0x62, 0x1f, 0x5b, 0x8e, 0x3d, 0x70,
	0xe7, 0xd3, 0x13, 0xc7, 0x1d, 0xd8, 0xe7, 0xb6, 0x63, 0x0f, 0x50, 0xb9, 0x5d, 0x32, 0xb4, 0x07,
	0x84, 0xa1, 0x45, 0x06, 0x6b, 0x82, 0x2e, 0x09, 0x4d, 0x30, 0x06, 0xc4, 0x1a, 0x4d, 0x46, 0x93,
	0x53, 0x54, 0x79, 0x76, 0x0a, 0x46, 0xf1, 0x59, 0x13, 0x7b, 0xf8, 0xa0, 0x17, 0xe7, 0xed, 0x4c,
	0xb4, 0x52, 0x83, 0xf2, 0xf9, 0xf4, 0x14, 0x69, 0x62, 0x71, 0x61, 0xcd, 0x50, 0x49, 0x3c, 0xb0,
	0x19, 0xb1, 0xa7, 0x64, 0x60, 0x13, 0x7b, 0xe0, 0x0a, 0xb0, 0xfc, 0x7a, 0x08, 0x07, 0x57, 0xd1,
	0xaa, 0x78, 0x93, 0x7c, 0xf8, 0x4b, 0xe2, 0x75, 0xcb, 0x51, 0xf1, 0x4c, 0x84, 0x33, 0xed, 0x5d,
	0xdb, 0xf3, 0xf9, 0xf2, 0xf6, 0xb2, 0x77, 0x15, 0xad, 0x8e, 0xd4, 0xa7, 0xbe, 0x90, 0x5c, 0x56,
	0xa5, 0xe6, 0xe7, 0x7f, 0x03, 0x00, 0x00, 0xff, 0xff, 0xfd, 0xbe, 0x46, 0x3e, 0x8f, 0x08, 0x00,
	0x00,
}
