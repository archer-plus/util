// Code generated by protoc-gen-go. DO NOT EDIT.
// source: proto/page.proto

package page

import (
	fmt "fmt"

	proto "github.com/golang/protobuf/proto"

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
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Order int32

const (
	Order_ASC  Order = 0
	Order_DESC Order = 1
)

var Order_name = map[int32]string{
	0: "ASC",
	1: "DESC",
}
var Order_value = map[string]int32{
	"ASC":  0,
	"DESC": 1,
}

func (x Order) String() string {
	return proto.EnumName(Order_name, int32(x))
}
func (Order) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_page_80284a56cac7aadc, []int{0}
}

type PageInfo struct {
	// @inject_tag: json:"page_size"
	PageSize int32 `protobuf:"varint,1,opt,name=page_size,json=pageSize,proto3" json:"page_size"`
	// @inject_tag: json:"page"
	Page int32 `protobuf:"varint,2,opt,name=page,proto3" json:"page"`
	// @inject_tag: json:"total"
	Total int32 `protobuf:"varint,3,opt,name=total,proto3" json:"total"`
	// @inject_tag: json:"page_total"
	PageTotal            int32    `protobuf:"varint,4,opt,name=page_total,json=pageTotal,proto3" json:"page_total"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PageInfo) Reset()         { *m = PageInfo{} }
func (m *PageInfo) String() string { return proto.CompactTextString(m) }
func (*PageInfo) ProtoMessage()    {}
func (*PageInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_page_80284a56cac7aadc, []int{0}
}
func (m *PageInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PageInfo.Unmarshal(m, b)
}
func (m *PageInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PageInfo.Marshal(b, m, deterministic)
}
func (dst *PageInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PageInfo.Merge(dst, src)
}
func (m *PageInfo) XXX_Size() int {
	return xxx_messageInfo_PageInfo.Size(m)
}
func (m *PageInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_PageInfo.DiscardUnknown(m)
}

var xxx_messageInfo_PageInfo proto.InternalMessageInfo

func (m *PageInfo) GetPageSize() int32 {
	if m != nil {
		return m.PageSize
	}
	return 0
}

func (m *PageInfo) GetPage() int32 {
	if m != nil {
		return m.Page
	}
	return 0
}

func (m *PageInfo) GetTotal() int32 {
	if m != nil {
		return m.Total
	}
	return 0
}

func (m *PageInfo) GetPageTotal() int32 {
	if m != nil {
		return m.PageTotal
	}
	return 0
}

type PageRequest struct {
	Order                Order    `protobuf:"varint,1,opt,name=order,proto3,enum=core.proto.Order" json:"order,omitempty"`
	Page                 int32    `protobuf:"varint,2,opt,name=page,proto3" json:"page,omitempty"`
	PageSize             int32    `protobuf:"varint,3,opt,name=page_size,json=pageSize,proto3" json:"page_size,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PageRequest) Reset()         { *m = PageRequest{} }
func (m *PageRequest) String() string { return proto.CompactTextString(m) }
func (*PageRequest) ProtoMessage()    {}
func (*PageRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_page_80284a56cac7aadc, []int{1}
}
func (m *PageRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PageRequest.Unmarshal(m, b)
}
func (m *PageRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PageRequest.Marshal(b, m, deterministic)
}
func (dst *PageRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PageRequest.Merge(dst, src)
}
func (m *PageRequest) XXX_Size() int {
	return xxx_messageInfo_PageRequest.Size(m)
}
func (m *PageRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_PageRequest.DiscardUnknown(m)
}

var xxx_messageInfo_PageRequest proto.InternalMessageInfo

func (m *PageRequest) GetOrder() Order {
	if m != nil {
		return m.Order
	}
	return Order_ASC
}

func (m *PageRequest) GetPage() int32 {
	if m != nil {
		return m.Page
	}
	return 0
}

func (m *PageRequest) GetPageSize() int32 {
	if m != nil {
		return m.PageSize
	}
	return 0
}

func init() {
	proto.RegisterType((*PageInfo)(nil), "core.proto.PageInfo")
	proto.RegisterType((*PageRequest)(nil), "core.proto.PageRequest")
	proto.RegisterEnum("core.proto.Order", Order_name, Order_value)
}

func init() { proto.RegisterFile("proto/page.proto", fileDescriptor_page_80284a56cac7aadc) }

var fileDescriptor_page_80284a56cac7aadc = []byte{
	// 232 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x90, 0x41, 0x4b, 0xc4, 0x30,
	0x10, 0x85, 0xad, 0x6d, 0xb4, 0x8e, 0x20, 0x75, 0xf0, 0x50, 0x14, 0x41, 0xf6, 0xa0, 0x8b, 0x60,
	0x02, 0xfa, 0x0b, 0x74, 0xf5, 0xe0, 0x49, 0x69, 0x3d, 0x79, 0x91, 0x6c, 0x1d, 0xb3, 0x85, 0xd5,
	0xc4, 0x34, 0xbd, 0xec, 0xaf, 0x97, 0x4c, 0x85, 0x52, 0xd8, 0xdb, 0x97, 0xf7, 0xc2, 0x7b, 0x8f,
	0x81, 0xc2, 0x79, 0x1b, 0xac, 0x72, 0xda, 0x90, 0x64, 0x44, 0x68, 0xac, 0xff, 0xe7, 0x99, 0x83,
	0xfc, 0x55, 0x1b, 0x7a, 0xfe, 0xf9, 0xb2, 0x78, 0x06, 0x07, 0xf1, 0xd7, 0x47, 0xd7, 0x6e, 0xa8,
	0x4c, 0x2e, 0x92, 0xb9, 0xa8, 0xf2, 0x28, 0xd4, 0xed, 0x86, 0x10, 0x21, 0x8b, 0x5c, 0xee, 0xb2,
	0xce, 0x8c, 0x27, 0x20, 0x82, 0x0d, 0x7a, 0x5d, 0xa6, 0x2c, 0x0e, 0x0f, 0x3c, 0x07, 0xe0, 0x98,
	0xc1, 0xca, 0xd8, 0xe2, 0xe0, 0xb7, 0x28, 0xcc, 0x0c, 0x1c, 0xc6, 0xc6, 0x8a, 0x7e, 0x7b, 0xea,
	0x02, 0x5e, 0x81, 0xb0, 0xfe, 0x93, 0x3c, 0x17, 0x1e, 0xdd, 0x1e, 0xcb, 0x71, 0x9c, 0x7c, 0x89,
	0x46, 0x35, 0xf8, 0x5b, 0x07, 0x4c, 0x16, 0xa7, 0xd3, 0xc5, 0xd7, 0xa7, 0x20, 0x38, 0x00, 0xf7,
	0x21, 0xbd, 0xaf, 0x17, 0xc5, 0x0e, 0xe6, 0x90, 0x3d, 0x3e, 0xd5, 0x8b, 0x22, 0x79, 0x98, 0xbf,
	0x5f, 0x9a, 0x36, 0xac, 0xfa, 0xa5, 0x6c, 0xec, 0xb7, 0xd2, 0xbe, 0x59, 0x91, 0xbf, 0x71, 0xeb,
	0xbe, 0x53, 0xb1, 0x5e, 0x8d, 0x27, 0x5b, 0xee, 0x31, 0xdf, 0xfd, 0x05, 0x00, 0x00, 0xff, 0xff,
	0xa3, 0xcc, 0xfb, 0xaa, 0x47, 0x01, 0x00, 0x00,
}