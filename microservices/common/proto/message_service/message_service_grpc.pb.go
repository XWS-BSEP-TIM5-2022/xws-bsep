// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.20.1
// source: common/proto/message_service/message_service.proto

package message

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

// MessageServiceClient is the client API for MessageService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MessageServiceClient interface {
	GetConversationById(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetResponse, error)
	GetConversation(ctx context.Context, in *GetConversationRequest, opts ...grpc.CallOption) (*GetConversationResponse, error)
	GetAllConversationsForUser(ctx context.Context, in *GetAllConversationsForUserRequest, opts ...grpc.CallOption) (*GetAllConversationsForUserResponse, error)
	NewMessage(ctx context.Context, in *NewMessageRequest, opts ...grpc.CallOption) (*NewMessageResponse, error)
}

type messageServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewMessageServiceClient(cc grpc.ClientConnInterface) MessageServiceClient {
	return &messageServiceClient{cc}
}

func (c *messageServiceClient) GetConversationById(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetResponse, error) {
	out := new(GetResponse)
	err := c.cc.Invoke(ctx, "/message_service.MessageService/GetConversationById", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *messageServiceClient) GetConversation(ctx context.Context, in *GetConversationRequest, opts ...grpc.CallOption) (*GetConversationResponse, error) {
	out := new(GetConversationResponse)
	err := c.cc.Invoke(ctx, "/message_service.MessageService/GetConversation", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *messageServiceClient) GetAllConversationsForUser(ctx context.Context, in *GetAllConversationsForUserRequest, opts ...grpc.CallOption) (*GetAllConversationsForUserResponse, error) {
	out := new(GetAllConversationsForUserResponse)
	err := c.cc.Invoke(ctx, "/message_service.MessageService/GetAllConversationsForUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *messageServiceClient) NewMessage(ctx context.Context, in *NewMessageRequest, opts ...grpc.CallOption) (*NewMessageResponse, error) {
	out := new(NewMessageResponse)
	err := c.cc.Invoke(ctx, "/message_service.MessageService/NewMessage", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MessageServiceServer is the server API for MessageService service.
// All implementations must embed UnimplementedMessageServiceServer
// for forward compatibility
type MessageServiceServer interface {
	GetConversationById(context.Context, *GetRequest) (*GetResponse, error)
	GetConversation(context.Context, *GetConversationRequest) (*GetConversationResponse, error)
	GetAllConversationsForUser(context.Context, *GetAllConversationsForUserRequest) (*GetAllConversationsForUserResponse, error)
	NewMessage(context.Context, *NewMessageRequest) (*NewMessageResponse, error)
	mustEmbedUnimplementedMessageServiceServer()
}

// UnimplementedMessageServiceServer must be embedded to have forward compatible implementations.
type UnimplementedMessageServiceServer struct {
}

func (UnimplementedMessageServiceServer) GetConversationById(context.Context, *GetRequest) (*GetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetConversationById not implemented")
}
func (UnimplementedMessageServiceServer) GetConversation(context.Context, *GetConversationRequest) (*GetConversationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetConversation not implemented")
}
func (UnimplementedMessageServiceServer) GetAllConversationsForUser(context.Context, *GetAllConversationsForUserRequest) (*GetAllConversationsForUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAllConversationsForUser not implemented")
}
func (UnimplementedMessageServiceServer) NewMessage(context.Context, *NewMessageRequest) (*NewMessageResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method NewMessage not implemented")
}
func (UnimplementedMessageServiceServer) mustEmbedUnimplementedMessageServiceServer() {}

// UnsafeMessageServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MessageServiceServer will
// result in compilation errors.
type UnsafeMessageServiceServer interface {
	mustEmbedUnimplementedMessageServiceServer()
}

func RegisterMessageServiceServer(s grpc.ServiceRegistrar, srv MessageServiceServer) {
	s.RegisterService(&MessageService_ServiceDesc, srv)
}

func _MessageService_GetConversationById_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MessageServiceServer).GetConversationById(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/message_service.MessageService/GetConversationById",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MessageServiceServer).GetConversationById(ctx, req.(*GetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MessageService_GetConversation_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetConversationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MessageServiceServer).GetConversation(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/message_service.MessageService/GetConversation",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MessageServiceServer).GetConversation(ctx, req.(*GetConversationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MessageService_GetAllConversationsForUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAllConversationsForUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MessageServiceServer).GetAllConversationsForUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/message_service.MessageService/GetAllConversationsForUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MessageServiceServer).GetAllConversationsForUser(ctx, req.(*GetAllConversationsForUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MessageService_NewMessage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NewMessageRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MessageServiceServer).NewMessage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/message_service.MessageService/NewMessage",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MessageServiceServer).NewMessage(ctx, req.(*NewMessageRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// MessageService_ServiceDesc is the grpc.ServiceDesc for MessageService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var MessageService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "message_service.MessageService",
	HandlerType: (*MessageServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetConversationById",
			Handler:    _MessageService_GetConversationById_Handler,
		},
		{
			MethodName: "GetConversation",
			Handler:    _MessageService_GetConversation_Handler,
		},
		{
			MethodName: "GetAllConversationsForUser",
			Handler:    _MessageService_GetAllConversationsForUser_Handler,
		},
		{
			MethodName: "NewMessage",
			Handler:    _MessageService_NewMessage_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "common/proto/message_service/message_service.proto",
}
