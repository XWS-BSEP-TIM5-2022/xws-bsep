// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.20.1
// source: common/proto/connection_service/connection_service.proto

package connection

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

// ConnectionServiceClient is the client API for ConnectionService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ConnectionServiceClient interface {
	GetConnections(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*Users, error)
	Register(ctx context.Context, in *RegisterRequest, opts ...grpc.CallOption) (*ActionResult, error)
	AddConnection(ctx context.Context, in *AddConnectionRequest, opts ...grpc.CallOption) (*ActionResult, error)
	ApproveConnection(ctx context.Context, in *ApproveConnectionRequest, opts ...grpc.CallOption) (*ActionResult, error)
}

type connectionServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewConnectionServiceClient(cc grpc.ClientConnInterface) ConnectionServiceClient {
	return &connectionServiceClient{cc}
}

func (c *connectionServiceClient) GetConnections(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*Users, error) {
	out := new(Users)
	err := c.cc.Invoke(ctx, "/connection_service.ConnectionService/GetConnections", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *connectionServiceClient) Register(ctx context.Context, in *RegisterRequest, opts ...grpc.CallOption) (*ActionResult, error) {
	out := new(ActionResult)
	err := c.cc.Invoke(ctx, "/connection_service.ConnectionService/Register", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *connectionServiceClient) AddConnection(ctx context.Context, in *AddConnectionRequest, opts ...grpc.CallOption) (*ActionResult, error) {
	out := new(ActionResult)
	err := c.cc.Invoke(ctx, "/connection_service.ConnectionService/AddConnection", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *connectionServiceClient) ApproveConnection(ctx context.Context, in *ApproveConnectionRequest, opts ...grpc.CallOption) (*ActionResult, error) {
	out := new(ActionResult)
	err := c.cc.Invoke(ctx, "/connection_service.ConnectionService/ApproveConnection", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ConnectionServiceServer is the server API for ConnectionService service.
// All implementations must embed UnimplementedConnectionServiceServer
// for forward compatibility
type ConnectionServiceServer interface {
	GetConnections(context.Context, *GetRequest) (*Users, error)
	Register(context.Context, *RegisterRequest) (*ActionResult, error)
	AddConnection(context.Context, *AddConnectionRequest) (*ActionResult, error)
	ApproveConnection(context.Context, *ApproveConnectionRequest) (*ActionResult, error)
	mustEmbedUnimplementedConnectionServiceServer()
}

// UnimplementedConnectionServiceServer must be embedded to have forward compatible implementations.
type UnimplementedConnectionServiceServer struct {
}

func (UnimplementedConnectionServiceServer) GetConnections(context.Context, *GetRequest) (*Users, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetConnections not implemented")
}
func (UnimplementedConnectionServiceServer) Register(context.Context, *RegisterRequest) (*ActionResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Register not implemented")
}
func (UnimplementedConnectionServiceServer) AddConnection(context.Context, *AddConnectionRequest) (*ActionResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddConnection not implemented")
}
func (UnimplementedConnectionServiceServer) ApproveConnection(context.Context, *ApproveConnectionRequest) (*ActionResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ApproveConnection not implemented")
}
func (UnimplementedConnectionServiceServer) mustEmbedUnimplementedConnectionServiceServer() {}

// UnsafeConnectionServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ConnectionServiceServer will
// result in compilation errors.
type UnsafeConnectionServiceServer interface {
	mustEmbedUnimplementedConnectionServiceServer()
}

func RegisterConnectionServiceServer(s grpc.ServiceRegistrar, srv ConnectionServiceServer) {
	s.RegisterService(&ConnectionService_ServiceDesc, srv)
}

func _ConnectionService_GetConnections_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ConnectionServiceServer).GetConnections(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/connection_service.ConnectionService/GetConnections",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ConnectionServiceServer).GetConnections(ctx, req.(*GetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ConnectionService_Register_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ConnectionServiceServer).Register(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/connection_service.ConnectionService/Register",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ConnectionServiceServer).Register(ctx, req.(*RegisterRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ConnectionService_AddConnection_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddConnectionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ConnectionServiceServer).AddConnection(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/connection_service.ConnectionService/AddConnection",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ConnectionServiceServer).AddConnection(ctx, req.(*AddConnectionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ConnectionService_ApproveConnection_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ApproveConnectionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ConnectionServiceServer).ApproveConnection(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/connection_service.ConnectionService/ApproveConnection",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ConnectionServiceServer).ApproveConnection(ctx, req.(*ApproveConnectionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ConnectionService_ServiceDesc is the grpc.ServiceDesc for ConnectionService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ConnectionService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "connection_service.ConnectionService",
	HandlerType: (*ConnectionServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetConnections",
			Handler:    _ConnectionService_GetConnections_Handler,
		},
		{
			MethodName: "Register",
			Handler:    _ConnectionService_Register_Handler,
		},
		{
			MethodName: "AddConnection",
			Handler:    _ConnectionService_AddConnection_Handler,
		},
		{
			MethodName: "ApproveConnection",
			Handler:    _ConnectionService_ApproveConnection_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "common/proto/connection_service/connection_service.proto",
}
