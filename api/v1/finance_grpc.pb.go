// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v5.26.1
// source: v1/finance.proto

package v1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	Finance_GetStockQuote_FullMethodName     = "/api.v1.Finance/GetStockQuote"
	Finance_GetUSASpending_FullMethodName    = "/api.v1.Finance/GetUSASpending"
	Finance_GetSenateLobbying_FullMethodName = "/api.v1.Finance/GetSenateLobbying"
	Finance_WatchTrades_FullMethodName       = "/api.v1.Finance/WatchTrades"
)

// FinanceClient is the client API for Finance service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type FinanceClient interface {
	GetStockQuote(ctx context.Context, in *GetStockQuoteRequest, opts ...grpc.CallOption) (*GetStockQuoteReply, error)
	GetUSASpending(ctx context.Context, in *GetUSASpendingRequest, opts ...grpc.CallOption) (*GetUSASpendingReply, error)
	GetSenateLobbying(ctx context.Context, in *GetSenateLobbyingRequest, opts ...grpc.CallOption) (*GetSenateLobbyingReply, error)
	WatchTrades(ctx context.Context, in *SyncTradesRequest, opts ...grpc.CallOption) (Finance_WatchTradesClient, error)
}

type financeClient struct {
	cc grpc.ClientConnInterface
}

func NewFinanceClient(cc grpc.ClientConnInterface) FinanceClient {
	return &financeClient{cc}
}

func (c *financeClient) GetStockQuote(ctx context.Context, in *GetStockQuoteRequest, opts ...grpc.CallOption) (*GetStockQuoteReply, error) {
	out := new(GetStockQuoteReply)
	err := c.cc.Invoke(ctx, Finance_GetStockQuote_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *financeClient) GetUSASpending(ctx context.Context, in *GetUSASpendingRequest, opts ...grpc.CallOption) (*GetUSASpendingReply, error) {
	out := new(GetUSASpendingReply)
	err := c.cc.Invoke(ctx, Finance_GetUSASpending_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *financeClient) GetSenateLobbying(ctx context.Context, in *GetSenateLobbyingRequest, opts ...grpc.CallOption) (*GetSenateLobbyingReply, error) {
	out := new(GetSenateLobbyingReply)
	err := c.cc.Invoke(ctx, Finance_GetSenateLobbying_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *financeClient) WatchTrades(ctx context.Context, in *SyncTradesRequest, opts ...grpc.CallOption) (Finance_WatchTradesClient, error) {
	stream, err := c.cc.NewStream(ctx, &Finance_ServiceDesc.Streams[0], Finance_WatchTrades_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &financeWatchTradesClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Finance_WatchTradesClient interface {
	Recv() (*SyncTradesReply, error)
	grpc.ClientStream
}

type financeWatchTradesClient struct {
	grpc.ClientStream
}

func (x *financeWatchTradesClient) Recv() (*SyncTradesReply, error) {
	m := new(SyncTradesReply)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// FinanceServer is the server API for Finance service.
// All implementations must embed UnimplementedFinanceServer
// for forward compatibility
type FinanceServer interface {
	GetStockQuote(context.Context, *GetStockQuoteRequest) (*GetStockQuoteReply, error)
	GetUSASpending(context.Context, *GetUSASpendingRequest) (*GetUSASpendingReply, error)
	GetSenateLobbying(context.Context, *GetSenateLobbyingRequest) (*GetSenateLobbyingReply, error)
	WatchTrades(*SyncTradesRequest, Finance_WatchTradesServer) error
	mustEmbedUnimplementedFinanceServer()
}

// UnimplementedFinanceServer must be embedded to have forward compatible implementations.
type UnimplementedFinanceServer struct {
}

func (UnimplementedFinanceServer) GetStockQuote(context.Context, *GetStockQuoteRequest) (*GetStockQuoteReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetStockQuote not implemented")
}
func (UnimplementedFinanceServer) GetUSASpending(context.Context, *GetUSASpendingRequest) (*GetUSASpendingReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUSASpending not implemented")
}
func (UnimplementedFinanceServer) GetSenateLobbying(context.Context, *GetSenateLobbyingRequest) (*GetSenateLobbyingReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetSenateLobbying not implemented")
}
func (UnimplementedFinanceServer) WatchTrades(*SyncTradesRequest, Finance_WatchTradesServer) error {
	return status.Errorf(codes.Unimplemented, "method WatchTrades not implemented")
}
func (UnimplementedFinanceServer) mustEmbedUnimplementedFinanceServer() {}

// UnsafeFinanceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to FinanceServer will
// result in compilation errors.
type UnsafeFinanceServer interface {
	mustEmbedUnimplementedFinanceServer()
}

func RegisterFinanceServer(s grpc.ServiceRegistrar, srv FinanceServer) {
	s.RegisterService(&Finance_ServiceDesc, srv)
}

func _Finance_GetStockQuote_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetStockQuoteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FinanceServer).GetStockQuote(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Finance_GetStockQuote_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FinanceServer).GetStockQuote(ctx, req.(*GetStockQuoteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Finance_GetUSASpending_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUSASpendingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FinanceServer).GetUSASpending(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Finance_GetUSASpending_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FinanceServer).GetUSASpending(ctx, req.(*GetUSASpendingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Finance_GetSenateLobbying_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetSenateLobbyingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FinanceServer).GetSenateLobbying(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Finance_GetSenateLobbying_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FinanceServer).GetSenateLobbying(ctx, req.(*GetSenateLobbyingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Finance_WatchTrades_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(SyncTradesRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(FinanceServer).WatchTrades(m, &financeWatchTradesServer{stream})
}

type Finance_WatchTradesServer interface {
	Send(*SyncTradesReply) error
	grpc.ServerStream
}

type financeWatchTradesServer struct {
	grpc.ServerStream
}

func (x *financeWatchTradesServer) Send(m *SyncTradesReply) error {
	return x.ServerStream.SendMsg(m)
}

// Finance_ServiceDesc is the grpc.ServiceDesc for Finance service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Finance_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "api.v1.Finance",
	HandlerType: (*FinanceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetStockQuote",
			Handler:    _Finance_GetStockQuote_Handler,
		},
		{
			MethodName: "GetUSASpending",
			Handler:    _Finance_GetUSASpending_Handler,
		},
		{
			MethodName: "GetSenateLobbying",
			Handler:    _Finance_GetSenateLobbying_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "WatchTrades",
			Handler:       _Finance_WatchTrades_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "v1/finance.proto",
}
