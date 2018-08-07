// Code generated by protoc-gen-go. DO NOT EDIT.
// source: proto/web.proto

/*
Package server is a generated protocol buffer package.

It is generated from these files:
	proto/web.proto

It has these top-level messages:
	Backend
	X509Cert
	Key
	KV
	ProxyState
	OpResult
	StateRequest
*/
package server

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "github.com/johanbrandhorst/protobuf/proto"

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

type Backend_Protocol int32

const (
	Backend_HTTP1 Backend_Protocol = 0
	Backend_HTTP2 Backend_Protocol = 1
	Backend_GRPC  Backend_Protocol = 3
)

var Backend_Protocol_name = map[int32]string{
	0: "HTTP1",
	1: "HTTP2",
	3: "GRPC",
}
var Backend_Protocol_value = map[string]int32{
	"HTTP1": 0,
	"HTTP2": 1,
	"GRPC":  3,
}

func (x Backend_Protocol) String() string {
	return proto.EnumName(Backend_Protocol_name, int32(x))
}
func (Backend_Protocol) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{0, 0} }

type Backend struct {
	Domain       string            `protobuf:"bytes,1,opt,name=domain" json:"domain,omitempty"`
	Ips          []string          `protobuf:"bytes,2,rep,name=ips" json:"ips,omitempty"`
	HealthCheck  string            `protobuf:"bytes,3,opt,name=health_check,json=healthCheck" json:"health_check,omitempty"`
	HealthStatus string            `protobuf:"bytes,4,opt,name=health_status,json=healthStatus" json:"health_status,omitempty"`
	Protocol     Backend_Protocol  `protobuf:"varint,5,opt,name=protocol,enum=web.Backend_Protocol" json:"protocol,omitempty"`
	InternetCert *X509Cert         `protobuf:"bytes,6,opt,name=internet_cert,json=internetCert" json:"internet_cert,omitempty"`
	BackendCert  *X509Cert         `protobuf:"bytes,7,opt,name=backend_cert,json=backendCert" json:"backend_cert,omitempty"`
	MatchHeaders map[string]string `protobuf:"bytes,8,rep,name=match_headers,json=matchHeaders" json:"match_headers,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
}

func (m *Backend) Reset()                    { *m = Backend{} }
func (m *Backend) String() string            { return proto.CompactTextString(m) }
func (*Backend) ProtoMessage()               {}
func (*Backend) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Backend) GetDomain() string {
	if m != nil {
		return m.Domain
	}
	return ""
}

func (m *Backend) GetIps() []string {
	if m != nil {
		return m.Ips
	}
	return nil
}

func (m *Backend) GetHealthCheck() string {
	if m != nil {
		return m.HealthCheck
	}
	return ""
}

func (m *Backend) GetHealthStatus() string {
	if m != nil {
		return m.HealthStatus
	}
	return ""
}

func (m *Backend) GetProtocol() Backend_Protocol {
	if m != nil {
		return m.Protocol
	}
	return Backend_HTTP1
}

func (m *Backend) GetInternetCert() *X509Cert {
	if m != nil {
		return m.InternetCert
	}
	return nil
}

func (m *Backend) GetBackendCert() *X509Cert {
	if m != nil {
		return m.BackendCert
	}
	return nil
}

func (m *Backend) GetMatchHeaders() map[string]string {
	if m != nil {
		return m.MatchHeaders
	}
	return nil
}

type X509Cert struct {
	Cert []byte `protobuf:"bytes,1,opt,name=cert,proto3" json:"cert,omitempty"`
	Key  []byte `protobuf:"bytes,2,opt,name=key,proto3" json:"key,omitempty"`
}

func (m *X509Cert) Reset()                    { *m = X509Cert{} }
func (m *X509Cert) String() string            { return proto.CompactTextString(m) }
func (*X509Cert) ProtoMessage()               {}
func (*X509Cert) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *X509Cert) GetCert() []byte {
	if m != nil {
		return m.Cert
	}
	return nil
}

func (m *X509Cert) GetKey() []byte {
	if m != nil {
		return m.Key
	}
	return nil
}

type Key struct {
	Prefix []byte `protobuf:"bytes,1,opt,name=prefix,proto3" json:"prefix,omitempty"`
}

func (m *Key) Reset()                    { *m = Key{} }
func (m *Key) String() string            { return proto.CompactTextString(m) }
func (*Key) ProtoMessage()               {}
func (*Key) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *Key) GetPrefix() []byte {
	if m != nil {
		return m.Prefix
	}
	return nil
}

type KV struct {
	Key   []byte `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Value []byte `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
}

func (m *KV) Reset()                    { *m = KV{} }
func (m *KV) String() string            { return proto.CompactTextString(m) }
func (*KV) ProtoMessage()               {}
func (*KV) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *KV) GetKey() []byte {
	if m != nil {
		return m.Key
	}
	return nil
}

func (m *KV) GetValue() []byte {
	if m != nil {
		return m.Value
	}
	return nil
}

type ProxyState struct {
	Backends []*Backend `protobuf:"bytes,1,rep,name=backends" json:"backends,omitempty"`
	// a status message, or an error message.
	Status string `protobuf:"bytes,2,opt,name=status" json:"status,omitempty"`
	// an error code
	Code int32 `protobuf:"varint,3,opt,name=code" json:"code,omitempty"`
}

func (m *ProxyState) Reset()                    { *m = ProxyState{} }
func (m *ProxyState) String() string            { return proto.CompactTextString(m) }
func (*ProxyState) ProtoMessage()               {}
func (*ProxyState) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *ProxyState) GetBackends() []*Backend {
	if m != nil {
		return m.Backends
	}
	return nil
}

func (m *ProxyState) GetStatus() string {
	if m != nil {
		return m.Status
	}
	return ""
}

func (m *ProxyState) GetCode() int32 {
	if m != nil {
		return m.Code
	}
	return 0
}

type OpResult struct {
	Code   int32  `protobuf:"varint,1,opt,name=code" json:"code,omitempty"`
	Status string `protobuf:"bytes,2,opt,name=status" json:"status,omitempty"`
}

func (m *OpResult) Reset()                    { *m = OpResult{} }
func (m *OpResult) String() string            { return proto.CompactTextString(m) }
func (*OpResult) ProtoMessage()               {}
func (*OpResult) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *OpResult) GetCode() int32 {
	if m != nil {
		return m.Code
	}
	return 0
}

func (m *OpResult) GetStatus() string {
	if m != nil {
		return m.Status
	}
	return ""
}

type StateRequest struct {
	// if domain is empty string, return "all" states, otherwise
	// match domain DNS-style, e.g. google.com matches docs.google.com
	Domain string `protobuf:"bytes,1,opt,name=domain" json:"domain,omitempty"`
}

func (m *StateRequest) Reset()                    { *m = StateRequest{} }
func (m *StateRequest) String() string            { return proto.CompactTextString(m) }
func (*StateRequest) ProtoMessage()               {}
func (*StateRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

func (m *StateRequest) GetDomain() string {
	if m != nil {
		return m.Domain
	}
	return ""
}

func init() {
	proto.RegisterType((*Backend)(nil), "web.Backend")
	proto.RegisterType((*X509Cert)(nil), "web.X509Cert")
	proto.RegisterType((*Key)(nil), "web.Key")
	proto.RegisterType((*KV)(nil), "web.KV")
	proto.RegisterType((*ProxyState)(nil), "web.ProxyState")
	proto.RegisterType((*OpResult)(nil), "web.OpResult")
	proto.RegisterType((*StateRequest)(nil), "web.StateRequest")
	proto.RegisterEnum("web.Backend_Protocol", Backend_Protocol_name, Backend_Protocol_value)
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for Proxy service

type ProxyClient interface {
	State(ctx context.Context, in *StateRequest, opts ...grpc.CallOption) (*ProxyState, error)
	Put(ctx context.Context, in *Backend, opts ...grpc.CallOption) (*OpResult, error)
	Remove(ctx context.Context, in *Backend, opts ...grpc.CallOption) (*OpResult, error)
	PutKVStream(ctx context.Context, opts ...grpc.CallOption) (Proxy_PutKVStreamClient, error)
	GetKVStream(ctx context.Context, in *Key, opts ...grpc.CallOption) (Proxy_GetKVStreamClient, error)
}

type proxyClient struct {
	cc *grpc.ClientConn
}

func NewProxyClient(cc *grpc.ClientConn) ProxyClient {
	return &proxyClient{cc}
}

func (c *proxyClient) State(ctx context.Context, in *StateRequest, opts ...grpc.CallOption) (*ProxyState, error) {
	out := new(ProxyState)
	err := grpc.Invoke(ctx, "/web.Proxy/State", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *proxyClient) Put(ctx context.Context, in *Backend, opts ...grpc.CallOption) (*OpResult, error) {
	out := new(OpResult)
	err := grpc.Invoke(ctx, "/web.Proxy/Put", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *proxyClient) Remove(ctx context.Context, in *Backend, opts ...grpc.CallOption) (*OpResult, error) {
	out := new(OpResult)
	err := grpc.Invoke(ctx, "/web.Proxy/Remove", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *proxyClient) PutKVStream(ctx context.Context, opts ...grpc.CallOption) (Proxy_PutKVStreamClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_Proxy_serviceDesc.Streams[0], c.cc, "/web.Proxy/PutKVStream", opts...)
	if err != nil {
		return nil, err
	}
	x := &proxyPutKVStreamClient{stream}
	return x, nil
}

type Proxy_PutKVStreamClient interface {
	Send(*KV) error
	CloseAndRecv() (*OpResult, error)
	grpc.ClientStream
}

type proxyPutKVStreamClient struct {
	grpc.ClientStream
}

func (x *proxyPutKVStreamClient) Send(m *KV) error {
	return x.ClientStream.SendMsg(m)
}

func (x *proxyPutKVStreamClient) CloseAndRecv() (*OpResult, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(OpResult)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *proxyClient) GetKVStream(ctx context.Context, in *Key, opts ...grpc.CallOption) (Proxy_GetKVStreamClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_Proxy_serviceDesc.Streams[1], c.cc, "/web.Proxy/GetKVStream", opts...)
	if err != nil {
		return nil, err
	}
	x := &proxyGetKVStreamClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Proxy_GetKVStreamClient interface {
	Recv() (*KV, error)
	grpc.ClientStream
}

type proxyGetKVStreamClient struct {
	grpc.ClientStream
}

func (x *proxyGetKVStreamClient) Recv() (*KV, error) {
	m := new(KV)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Server API for Proxy service

type ProxyServer interface {
	State(context.Context, *StateRequest) (*ProxyState, error)
	Put(context.Context, *Backend) (*OpResult, error)
	Remove(context.Context, *Backend) (*OpResult, error)
	PutKVStream(Proxy_PutKVStreamServer) error
	GetKVStream(*Key, Proxy_GetKVStreamServer) error
}

func RegisterProxyServer(s *grpc.Server, srv ProxyServer) {
	s.RegisterService(&_Proxy_serviceDesc, srv)
}

func _Proxy_State_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProxyServer).State(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/web.Proxy/State",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProxyServer).State(ctx, req.(*StateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Proxy_Put_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Backend)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProxyServer).Put(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/web.Proxy/Put",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProxyServer).Put(ctx, req.(*Backend))
	}
	return interceptor(ctx, in, info, handler)
}

func _Proxy_Remove_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Backend)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProxyServer).Remove(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/web.Proxy/Remove",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProxyServer).Remove(ctx, req.(*Backend))
	}
	return interceptor(ctx, in, info, handler)
}

func _Proxy_PutKVStream_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(ProxyServer).PutKVStream(&proxyPutKVStreamServer{stream})
}

type Proxy_PutKVStreamServer interface {
	SendAndClose(*OpResult) error
	Recv() (*KV, error)
	grpc.ServerStream
}

type proxyPutKVStreamServer struct {
	grpc.ServerStream
}

func (x *proxyPutKVStreamServer) SendAndClose(m *OpResult) error {
	return x.ServerStream.SendMsg(m)
}

func (x *proxyPutKVStreamServer) Recv() (*KV, error) {
	m := new(KV)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Proxy_GetKVStream_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(Key)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(ProxyServer).GetKVStream(m, &proxyGetKVStreamServer{stream})
}

type Proxy_GetKVStreamServer interface {
	Send(*KV) error
	grpc.ServerStream
}

type proxyGetKVStreamServer struct {
	grpc.ServerStream
}

func (x *proxyGetKVStreamServer) Send(m *KV) error {
	return x.ServerStream.SendMsg(m)
}

var _Proxy_serviceDesc = grpc.ServiceDesc{
	ServiceName: "web.Proxy",
	HandlerType: (*ProxyServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "State",
			Handler:    _Proxy_State_Handler,
		},
		{
			MethodName: "Put",
			Handler:    _Proxy_Put_Handler,
		},
		{
			MethodName: "Remove",
			Handler:    _Proxy_Remove_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "PutKVStream",
			Handler:       _Proxy_PutKVStream_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "GetKVStream",
			Handler:       _Proxy_GetKVStream_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "proto/web.proto",
}

func init() { proto.RegisterFile("proto/web.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 618 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x54, 0xd1, 0x6e, 0xd3, 0x30,
	0x14, 0x6d, 0x9a, 0xb5, 0xcb, 0x6e, 0x53, 0xd6, 0x59, 0x80, 0xa2, 0x4a, 0xa0, 0x12, 0x26, 0x08,
	0x88, 0xb5, 0x5d, 0x11, 0x68, 0xf0, 0x02, 0x5a, 0x85, 0x36, 0xa9, 0x42, 0x54, 0xd9, 0x34, 0x21,
	0x5e, 0x26, 0x27, 0xbd, 0x5b, 0xb2, 0x35, 0x71, 0x71, 0x9c, 0x6d, 0xfd, 0x41, 0x3e, 0x81, 0x8f,
	0xe1, 0x09, 0xd9, 0x71, 0xd6, 0x4e, 0xdb, 0x84, 0x78, 0x3b, 0xf7, 0xde, 0x73, 0xe2, 0x73, 0x6e,
	0x2c, 0xc3, 0xfa, 0x8c, 0x33, 0xc1, 0x7a, 0x97, 0x18, 0x74, 0x15, 0x22, 0xe6, 0x25, 0x06, 0xed,
	0x9d, 0xd3, 0x58, 0x44, 0x79, 0xd0, 0x0d, 0x59, 0xd2, 0x3b, 0x63, 0x11, 0x4d, 0x03, 0x4e, 0xd3,
	0x49, 0xc4, 0x78, 0x26, 0x7a, 0x8a, 0x16, 0xe4, 0x27, 0x05, 0xe8, 0x9d, 0xb2, 0x59, 0x84, 0xfc,
	0x2c, 0x2b, 0xe4, 0xee, 0x2f, 0x13, 0x56, 0x77, 0x69, 0x78, 0x8e, 0xe9, 0x84, 0x3c, 0x86, 0xfa,
	0x84, 0x25, 0x34, 0x4e, 0x1d, 0xa3, 0x63, 0x78, 0x6b, 0xbe, 0xae, 0x48, 0x0b, 0xcc, 0x78, 0x96,
	0x39, 0xd5, 0x8e, 0xe9, 0xad, 0xf9, 0x12, 0x92, 0x67, 0x60, 0x47, 0x48, 0xa7, 0x22, 0x3a, 0x0e,
	0x23, 0x0c, 0xcf, 0x1d, 0x53, 0xf1, 0x1b, 0x45, 0x6f, 0x28, 0x5b, 0xe4, 0x39, 0x34, 0x35, 0x25,
	0x13, 0x54, 0xe4, 0x99, 0xb3, 0xa2, 0x38, 0x5a, 0x77, 0xa0, 0x7a, 0x64, 0x1b, 0x2c, 0x65, 0x23,
	0x64, 0x53, 0xa7, 0xd6, 0x31, 0xbc, 0x07, 0x83, 0x47, 0x5d, 0x19, 0x4d, 0x3b, 0xea, 0x8e, 0xf5,
	0xd0, 0xbf, 0xa6, 0x91, 0x01, 0x34, 0xe3, 0x54, 0x20, 0x4f, 0x51, 0x1c, 0x87, 0xc8, 0x85, 0x53,
	0xef, 0x18, 0x5e, 0x63, 0xd0, 0x54, 0xba, 0xef, 0xef, 0xfa, 0x1f, 0x86, 0xc8, 0x85, 0x6f, 0x97,
	0x1c, 0x59, 0x91, 0x3e, 0xd8, 0x41, 0xf1, 0xc5, 0x42, 0xb2, 0x7a, 0x97, 0xa4, 0xa1, 0x29, 0x4a,
	0x31, 0x84, 0x66, 0x42, 0x45, 0x18, 0x1d, 0x47, 0x48, 0x27, 0xc8, 0x33, 0xc7, 0xea, 0x98, 0x5e,
	0x63, 0xf0, 0xf4, 0x86, 0xbb, 0xaf, 0x92, 0xb1, 0x5f, 0x10, 0xbe, 0xa4, 0x82, 0xcf, 0x7d, 0x3b,
	0x59, 0x6a, 0xb5, 0x3f, 0xc1, 0xc6, 0x2d, 0x8a, 0x5c, 0xe6, 0x39, 0xce, 0xf5, 0x86, 0x25, 0x24,
	0x0f, 0xa1, 0x76, 0x41, 0xa7, 0x39, 0x3a, 0x55, 0xd5, 0x2b, 0x8a, 0x8f, 0xd5, 0x1d, 0xc3, 0x7d,
	0x0d, 0x56, 0xb9, 0x01, 0xb2, 0x06, 0xb5, 0xfd, 0xc3, 0xc3, 0xf1, 0x76, 0xab, 0x52, 0xc2, 0x41,
	0xcb, 0x20, 0x16, 0xac, 0xec, 0xf9, 0xe3, 0x61, 0xcb, 0x74, 0xfb, 0x60, 0x95, 0x51, 0x08, 0x81,
	0x15, 0x95, 0x53, 0x1e, 0x62, 0xfb, 0x0a, 0x97, 0xe7, 0x56, 0x55, 0x4b, 0x42, 0xf7, 0x09, 0x98,
	0x23, 0x9c, 0xcb, 0xbf, 0x3e, 0xe3, 0x78, 0x12, 0x5f, 0x69, 0xba, 0xae, 0xdc, 0x37, 0x50, 0x1d,
	0x1d, 0x2d, 0xdb, 0xb5, 0xef, 0xb0, 0x6b, 0x6b, 0xbb, 0x6e, 0x00, 0x30, 0xe6, 0xec, 0x6a, 0x2e,
	0x7f, 0x2c, 0x12, 0x0f, 0x2c, 0xbd, 0xcd, 0xcc, 0x31, 0xd4, 0xe6, 0xec, 0xe5, 0xcd, 0xf9, 0xd7,
	0x53, 0x79, 0xba, 0xbe, 0x1f, 0x45, 0x7a, 0x5d, 0xa9, 0x08, 0x6c, 0x82, 0xea, 0x66, 0xd5, 0x7c,
	0x85, 0xdd, 0xf7, 0x60, 0x7d, 0x9b, 0xf9, 0x98, 0xe5, 0x53, 0x71, 0x3d, 0x37, 0x16, 0xf3, 0xfb,
	0xbe, 0xe5, 0xbe, 0x00, 0x5b, 0xd9, 0xf2, 0xf1, 0x67, 0x8e, 0x99, 0xb8, 0xef, 0x9e, 0x0f, 0x7e,
	0x1b, 0x50, 0x53, 0x21, 0xc8, 0x16, 0xd4, 0x8a, 0x20, 0x1b, 0xca, 0xf6, 0xb2, 0xba, 0xbd, 0xae,
	0x5a, 0x8b, 0xb0, 0x6e, 0x85, 0x6c, 0x82, 0x39, 0xce, 0x05, 0xb9, 0x91, 0xb1, 0x5d, 0x5c, 0xaf,
	0xd2, 0xb0, 0x5b, 0x21, 0x2f, 0xa1, 0xee, 0x63, 0xc2, 0x2e, 0xf0, 0x5f, 0xc4, 0x57, 0xd0, 0x18,
	0xe7, 0x62, 0x74, 0x74, 0x20, 0x38, 0xd2, 0x84, 0xac, 0xaa, 0xf9, 0xe8, 0xe8, 0x16, 0xd1, 0x33,
	0xc8, 0x26, 0x34, 0xf6, 0x70, 0x41, 0xb5, 0x0a, 0x2a, 0xce, 0xdb, 0xa5, 0xc8, 0xad, 0xf4, 0x8d,
	0xdd, 0xd3, 0x3f, 0x9f, 0xb7, 0x97, 0x5e, 0x08, 0x9a, 0x5e, 0xc5, 0x2c, 0xcf, 0x12, 0x36, 0x41,
	0x9e, 0x26, 0x34, 0xed, 0x85, 0x6c, 0x2b, 0x8c, 0x68, 0xcc, 0xf5, 0x13, 0x11, 0x4e, 0x63, 0x4c,
	0xc5, 0x8f, 0xff, 0x91, 0x64, 0xc8, 0x2f, 0x90, 0x07, 0x75, 0x55, 0xbd, 0xfd, 0x1b, 0x00, 0x00,
	0xff, 0xff, 0x71, 0x24, 0x61, 0xf6, 0xa6, 0x04, 0x00, 0x00,
}
