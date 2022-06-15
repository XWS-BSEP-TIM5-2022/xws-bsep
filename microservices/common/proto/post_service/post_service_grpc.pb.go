// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.20.1
// source: post_service.proto

package post

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

// PostServiceClient is the client API for PostService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PostServiceClient interface {
	Get(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetResponse, error)
	LikePost(ctx context.Context, in *InsertLike, opts ...grpc.CallOption) (*InsertResponse, error)
	DislikePost(ctx context.Context, in *InsertDislike, opts ...grpc.CallOption) (*InsertResponse, error)
	NeutralPost(ctx context.Context, in *InsertNeutralReaction, opts ...grpc.CallOption) (*InsertResponse, error)
	CommentPost(ctx context.Context, in *InsertComment, opts ...grpc.CallOption) (*InsertResponse, error)
	GetAll(ctx context.Context, in *GetAllRequest, opts ...grpc.CallOption) (*GetAllResponse, error)
	Insert(ctx context.Context, in *InsertRequest, opts ...grpc.CallOption) (*InsertResponse, error)
	// TODO: delete
	InsertJobOffer(ctx context.Context, in *InsertJobOfferRequest, opts ...grpc.CallOption) (*InsertResponse, error)
	GetAllByUser(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetAllResponse, error)
	UpdateCompanyInfo(ctx context.Context, in *UpdateCompanyInfoRequest, opts ...grpc.CallOption) (*InsertResponse, error)
}

type postServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewPostServiceClient(cc grpc.ClientConnInterface) PostServiceClient {
	return &postServiceClient{cc}
}

func (c *postServiceClient) Get(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetResponse, error) {
	out := new(GetResponse)
	err := c.cc.Invoke(ctx, "/post_service.PostService/Get", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *postServiceClient) LikePost(ctx context.Context, in *InsertLike, opts ...grpc.CallOption) (*InsertResponse, error) {
	out := new(InsertResponse)
	err := c.cc.Invoke(ctx, "/post_service.PostService/LikePost", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *postServiceClient) DislikePost(ctx context.Context, in *InsertDislike, opts ...grpc.CallOption) (*InsertResponse, error) {
	out := new(InsertResponse)
	err := c.cc.Invoke(ctx, "/post_service.PostService/DislikePost", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *postServiceClient) NeutralPost(ctx context.Context, in *InsertNeutralReaction, opts ...grpc.CallOption) (*InsertResponse, error) {
	out := new(InsertResponse)
	err := c.cc.Invoke(ctx, "/post_service.PostService/NeutralPost", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *postServiceClient) CommentPost(ctx context.Context, in *InsertComment, opts ...grpc.CallOption) (*InsertResponse, error) {
	out := new(InsertResponse)
	err := c.cc.Invoke(ctx, "/post_service.PostService/CommentPost", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *postServiceClient) GetAll(ctx context.Context, in *GetAllRequest, opts ...grpc.CallOption) (*GetAllResponse, error) {
	out := new(GetAllResponse)
	err := c.cc.Invoke(ctx, "/post_service.PostService/GetAll", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *postServiceClient) Insert(ctx context.Context, in *InsertRequest, opts ...grpc.CallOption) (*InsertResponse, error) {
	out := new(InsertResponse)
	err := c.cc.Invoke(ctx, "/post_service.PostService/Insert", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *postServiceClient) InsertJobOffer(ctx context.Context, in *InsertJobOfferRequest, opts ...grpc.CallOption) (*InsertResponse, error) {
	out := new(InsertResponse)
	err := c.cc.Invoke(ctx, "/post_service.PostService/InsertJobOffer", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *postServiceClient) GetAllByUser(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetAllResponse, error) {
	out := new(GetAllResponse)
	err := c.cc.Invoke(ctx, "/post_service.PostService/GetAllByUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *postServiceClient) UpdateCompanyInfo(ctx context.Context, in *UpdateCompanyInfoRequest, opts ...grpc.CallOption) (*InsertResponse, error) {
	out := new(InsertResponse)
	err := c.cc.Invoke(ctx, "/post_service.PostService/UpdateCompanyInfo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PostServiceServer is the server API for PostService service.
// All implementations must embed UnimplementedPostServiceServer
// for forward compatibility
type PostServiceServer interface {
	Get(context.Context, *GetRequest) (*GetResponse, error)
	LikePost(context.Context, *InsertLike) (*InsertResponse, error)
	DislikePost(context.Context, *InsertDislike) (*InsertResponse, error)
	NeutralPost(context.Context, *InsertNeutralReaction) (*InsertResponse, error)
	CommentPost(context.Context, *InsertComment) (*InsertResponse, error)
	GetAll(context.Context, *GetAllRequest) (*GetAllResponse, error)
	Insert(context.Context, *InsertRequest) (*InsertResponse, error)
	// TODO: delete
	InsertJobOffer(context.Context, *InsertJobOfferRequest) (*InsertResponse, error)
	GetAllByUser(context.Context, *GetRequest) (*GetAllResponse, error)
	UpdateCompanyInfo(context.Context, *UpdateCompanyInfoRequest) (*InsertResponse, error)
	mustEmbedUnimplementedPostServiceServer()
}

// UnimplementedPostServiceServer must be embedded to have forward compatible implementations.
type UnimplementedPostServiceServer struct {
}

func (UnimplementedPostServiceServer) Get(context.Context, *GetRequest) (*GetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (UnimplementedPostServiceServer) LikePost(context.Context, *InsertLike) (*InsertResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LikePost not implemented")
}
func (UnimplementedPostServiceServer) DislikePost(context.Context, *InsertDislike) (*InsertResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DislikePost not implemented")
}
func (UnimplementedPostServiceServer) NeutralPost(context.Context, *InsertNeutralReaction) (*InsertResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method NeutralPost not implemented")
}
func (UnimplementedPostServiceServer) CommentPost(context.Context, *InsertComment) (*InsertResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CommentPost not implemented")
}
func (UnimplementedPostServiceServer) GetAll(context.Context, *GetAllRequest) (*GetAllResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAll not implemented")
}
func (UnimplementedPostServiceServer) Insert(context.Context, *InsertRequest) (*InsertResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Insert not implemented")
}
func (UnimplementedPostServiceServer) InsertJobOffer(context.Context, *InsertJobOfferRequest) (*InsertResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method InsertJobOffer not implemented")
}
func (UnimplementedPostServiceServer) GetAllByUser(context.Context, *GetRequest) (*GetAllResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAllByUser not implemented")
}
func (UnimplementedPostServiceServer) UpdateCompanyInfo(context.Context, *UpdateCompanyInfoRequest) (*InsertResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateCompanyInfo not implemented")
}
func (UnimplementedPostServiceServer) mustEmbedUnimplementedPostServiceServer() {}

// UnsafePostServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PostServiceServer will
// result in compilation errors.
type UnsafePostServiceServer interface {
	mustEmbedUnimplementedPostServiceServer()
}

func RegisterPostServiceServer(s grpc.ServiceRegistrar, srv PostServiceServer) {
	s.RegisterService(&PostService_ServiceDesc, srv)
}

func _PostService_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PostServiceServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/post_service.PostService/Get",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PostServiceServer).Get(ctx, req.(*GetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PostService_LikePost_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(InsertLike)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PostServiceServer).LikePost(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/post_service.PostService/LikePost",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PostServiceServer).LikePost(ctx, req.(*InsertLike))
	}
	return interceptor(ctx, in, info, handler)
}

func _PostService_DislikePost_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(InsertDislike)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PostServiceServer).DislikePost(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/post_service.PostService/DislikePost",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PostServiceServer).DislikePost(ctx, req.(*InsertDislike))
	}
	return interceptor(ctx, in, info, handler)
}

func _PostService_NeutralPost_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(InsertNeutralReaction)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PostServiceServer).NeutralPost(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/post_service.PostService/NeutralPost",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PostServiceServer).NeutralPost(ctx, req.(*InsertNeutralReaction))
	}
	return interceptor(ctx, in, info, handler)
}

func _PostService_CommentPost_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(InsertComment)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PostServiceServer).CommentPost(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/post_service.PostService/CommentPost",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PostServiceServer).CommentPost(ctx, req.(*InsertComment))
	}
	return interceptor(ctx, in, info, handler)
}

func _PostService_GetAll_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAllRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PostServiceServer).GetAll(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/post_service.PostService/GetAll",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PostServiceServer).GetAll(ctx, req.(*GetAllRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PostService_Insert_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(InsertRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PostServiceServer).Insert(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/post_service.PostService/Insert",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PostServiceServer).Insert(ctx, req.(*InsertRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PostService_InsertJobOffer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(InsertJobOfferRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PostServiceServer).InsertJobOffer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/post_service.PostService/InsertJobOffer",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PostServiceServer).InsertJobOffer(ctx, req.(*InsertJobOfferRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PostService_GetAllByUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PostServiceServer).GetAllByUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/post_service.PostService/GetAllByUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PostServiceServer).GetAllByUser(ctx, req.(*GetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PostService_UpdateCompanyInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateCompanyInfoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PostServiceServer).UpdateCompanyInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/post_service.PostService/UpdateCompanyInfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PostServiceServer).UpdateCompanyInfo(ctx, req.(*UpdateCompanyInfoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// PostService_ServiceDesc is the grpc.ServiceDesc for PostService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var PostService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "post_service.PostService",
	HandlerType: (*PostServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Get",
			Handler:    _PostService_Get_Handler,
		},
		{
			MethodName: "LikePost",
			Handler:    _PostService_LikePost_Handler,
		},
		{
			MethodName: "DislikePost",
			Handler:    _PostService_DislikePost_Handler,
		},
		{
			MethodName: "NeutralPost",
			Handler:    _PostService_NeutralPost_Handler,
		},
		{
			MethodName: "CommentPost",
			Handler:    _PostService_CommentPost_Handler,
		},
		{
			MethodName: "GetAll",
			Handler:    _PostService_GetAll_Handler,
		},
		{
			MethodName: "Insert",
			Handler:    _PostService_Insert_Handler,
		},
		{
			MethodName: "InsertJobOffer",
			Handler:    _PostService_InsertJobOffer_Handler,
		},
		{
			MethodName: "GetAllByUser",
			Handler:    _PostService_GetAllByUser_Handler,
		},
		{
			MethodName: "UpdateCompanyInfo",
			Handler:    _PostService_UpdateCompanyInfo_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "post_service.proto",
}
