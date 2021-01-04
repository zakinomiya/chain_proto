// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package gw

import (
	context "context"
	empty "github.com/golang/protobuf/ptypes/empty"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion7

// HTTPServiceClient is the client API for HTTPService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type HTTPServiceClient interface {
	SendTransaction(ctx context.Context, in *SendTransactionRequest, opts ...grpc.CallOption) (*empty.Empty, error)
	GetTransactionByHash(ctx context.Context, in *GetTransactionByHashRequest, opts ...grpc.CallOption) (*GetTransactionResponse, error)
	SendBlock(ctx context.Context, in *SendBlockRequest, opts ...grpc.CallOption) (*empty.Empty, error)
	GetBlockByHeight(ctx context.Context, in *GetBlockByHeightRequest, opts ...grpc.CallOption) (*GetBlockResponse, error)
	GetBlockByHash(ctx context.Context, in *GetBlockByHashRequest, opts ...grpc.CallOption) (*GetBlockResponse, error)
	GetAccount(ctx context.Context, in *GetAccountRequest, opts ...grpc.CallOption) (*GetAccountResponse, error)
}

type hTTPServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewHTTPServiceClient(cc grpc.ClientConnInterface) HTTPServiceClient {
	return &hTTPServiceClient{cc}
}

func (c *hTTPServiceClient) SendTransaction(ctx context.Context, in *SendTransactionRequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/gw.HTTPService/SendTransaction", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *hTTPServiceClient) GetTransactionByHash(ctx context.Context, in *GetTransactionByHashRequest, opts ...grpc.CallOption) (*GetTransactionResponse, error) {
	out := new(GetTransactionResponse)
	err := c.cc.Invoke(ctx, "/gw.HTTPService/GetTransactionByHash", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *hTTPServiceClient) SendBlock(ctx context.Context, in *SendBlockRequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/gw.HTTPService/SendBlock", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *hTTPServiceClient) GetBlockByHeight(ctx context.Context, in *GetBlockByHeightRequest, opts ...grpc.CallOption) (*GetBlockResponse, error) {
	out := new(GetBlockResponse)
	err := c.cc.Invoke(ctx, "/gw.HTTPService/GetBlockByHeight", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *hTTPServiceClient) GetBlockByHash(ctx context.Context, in *GetBlockByHashRequest, opts ...grpc.CallOption) (*GetBlockResponse, error) {
	out := new(GetBlockResponse)
	err := c.cc.Invoke(ctx, "/gw.HTTPService/GetBlockByHash", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *hTTPServiceClient) GetAccount(ctx context.Context, in *GetAccountRequest, opts ...grpc.CallOption) (*GetAccountResponse, error) {
	out := new(GetAccountResponse)
	err := c.cc.Invoke(ctx, "/gw.HTTPService/GetAccount", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// HTTPServiceServer is the server API for HTTPService service.
// All implementations must embed UnimplementedHTTPServiceServer
// for forward compatibility
type HTTPServiceServer interface {
	SendTransaction(context.Context, *SendTransactionRequest) (*empty.Empty, error)
	GetTransactionByHash(context.Context, *GetTransactionByHashRequest) (*GetTransactionResponse, error)
	SendBlock(context.Context, *SendBlockRequest) (*empty.Empty, error)
	GetBlockByHeight(context.Context, *GetBlockByHeightRequest) (*GetBlockResponse, error)
	GetBlockByHash(context.Context, *GetBlockByHashRequest) (*GetBlockResponse, error)
	GetAccount(context.Context, *GetAccountRequest) (*GetAccountResponse, error)
	mustEmbedUnimplementedHTTPServiceServer()
}

// UnimplementedHTTPServiceServer must be embedded to have forward compatible implementations.
type UnimplementedHTTPServiceServer struct {
}

func (UnimplementedHTTPServiceServer) SendTransaction(context.Context, *SendTransactionRequest) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendTransaction not implemented")
}
func (UnimplementedHTTPServiceServer) GetTransactionByHash(context.Context, *GetTransactionByHashRequest) (*GetTransactionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTransactionByHash not implemented")
}
func (UnimplementedHTTPServiceServer) SendBlock(context.Context, *SendBlockRequest) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendBlock not implemented")
}
func (UnimplementedHTTPServiceServer) GetBlockByHeight(context.Context, *GetBlockByHeightRequest) (*GetBlockResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetBlockByHeight not implemented")
}
func (UnimplementedHTTPServiceServer) GetBlockByHash(context.Context, *GetBlockByHashRequest) (*GetBlockResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetBlockByHash not implemented")
}
func (UnimplementedHTTPServiceServer) GetAccount(context.Context, *GetAccountRequest) (*GetAccountResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAccount not implemented")
}
func (UnimplementedHTTPServiceServer) mustEmbedUnimplementedHTTPServiceServer() {}

// UnsafeHTTPServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to HTTPServiceServer will
// result in compilation errors.
type UnsafeHTTPServiceServer interface {
	mustEmbedUnimplementedHTTPServiceServer()
}

func RegisterHTTPServiceServer(s grpc.ServiceRegistrar, srv HTTPServiceServer) {
	s.RegisterService(&_HTTPService_serviceDesc, srv)
}

func _HTTPService_SendTransaction_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendTransactionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HTTPServiceServer).SendTransaction(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gw.HTTPService/SendTransaction",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HTTPServiceServer).SendTransaction(ctx, req.(*SendTransactionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _HTTPService_GetTransactionByHash_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetTransactionByHashRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HTTPServiceServer).GetTransactionByHash(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gw.HTTPService/GetTransactionByHash",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HTTPServiceServer).GetTransactionByHash(ctx, req.(*GetTransactionByHashRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _HTTPService_SendBlock_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendBlockRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HTTPServiceServer).SendBlock(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gw.HTTPService/SendBlock",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HTTPServiceServer).SendBlock(ctx, req.(*SendBlockRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _HTTPService_GetBlockByHeight_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetBlockByHeightRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HTTPServiceServer).GetBlockByHeight(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gw.HTTPService/GetBlockByHeight",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HTTPServiceServer).GetBlockByHeight(ctx, req.(*GetBlockByHeightRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _HTTPService_GetBlockByHash_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetBlockByHashRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HTTPServiceServer).GetBlockByHash(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gw.HTTPService/GetBlockByHash",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HTTPServiceServer).GetBlockByHash(ctx, req.(*GetBlockByHashRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _HTTPService_GetAccount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAccountRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HTTPServiceServer).GetAccount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gw.HTTPService/GetAccount",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HTTPServiceServer).GetAccount(ctx, req.(*GetAccountRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _HTTPService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "gw.HTTPService",
	HandlerType: (*HTTPServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SendTransaction",
			Handler:    _HTTPService_SendTransaction_Handler,
		},
		{
			MethodName: "GetTransactionByHash",
			Handler:    _HTTPService_GetTransactionByHash_Handler,
		},
		{
			MethodName: "SendBlock",
			Handler:    _HTTPService_SendBlock_Handler,
		},
		{
			MethodName: "GetBlockByHeight",
			Handler:    _HTTPService_GetBlockByHeight_Handler,
		},
		{
			MethodName: "GetBlockByHash",
			Handler:    _HTTPService_GetBlockByHash_Handler,
		},
		{
			MethodName: "GetAccount",
			Handler:    _HTTPService_GetAccount_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "service.proto",
}