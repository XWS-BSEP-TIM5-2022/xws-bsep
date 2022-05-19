// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.19.4
// source: auth_service.proto

package auth

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

// AuthServiceClient is the client API for AuthService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AuthServiceClient interface {
	Register(ctx context.Context, in *RegisterRequest, opts ...grpc.CallOption) (*RegisterResponse, error)
	Login(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginResponse, error)
	PasswordlessLogin(ctx context.Context, in *PasswordlessLoginRequest, opts ...grpc.CallOption) (*PasswordlessLoginResponse, error)
	GetAll(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*GetAllResponse, error)
	UpdateUsername(ctx context.Context, in *UpdateUsernameRequest, opts ...grpc.CallOption) (*UpdateUsernameResponse, error)
	ChangePassword(ctx context.Context, in *ChangePasswordRequest, opts ...grpc.CallOption) (*ChangePasswordResponse, error)
	ConfirmEmailLogin(ctx context.Context, in *ConfirmEmailLoginRequest, opts ...grpc.CallOption) (*ConfirmEmailLoginResponse, error)
	ActivateAccount(ctx context.Context, in *ActivationRequest, opts ...grpc.CallOption) (*ActivationResponse, error)
	SendRecoveryCode(ctx context.Context, in *SendRecoveryCodeRequest, opts ...grpc.CallOption) (*SendRecoveryCodeResponse, error)
	VerifyRecoveryCode(ctx context.Context, in *VerifyRecoveryCodeRequest, opts ...grpc.CallOption) (*Response, error)
	ResetForgottenPassword(ctx context.Context, in *ResetForgottenPasswordRequest, opts ...grpc.CallOption) (*Response, error)
}

type authServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewAuthServiceClient(cc grpc.ClientConnInterface) AuthServiceClient {
	return &authServiceClient{cc}
}

func (c *authServiceClient) Register(ctx context.Context, in *RegisterRequest, opts ...grpc.CallOption) (*RegisterResponse, error) {
	out := new(RegisterResponse)
	err := c.cc.Invoke(ctx, "/auth_service.AuthService/Register", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authServiceClient) Login(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginResponse, error) {
	out := new(LoginResponse)
	err := c.cc.Invoke(ctx, "/auth_service.AuthService/Login", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authServiceClient) PasswordlessLogin(ctx context.Context, in *PasswordlessLoginRequest, opts ...grpc.CallOption) (*PasswordlessLoginResponse, error) {
	out := new(PasswordlessLoginResponse)
	err := c.cc.Invoke(ctx, "/auth_service.AuthService/PasswordlessLogin", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authServiceClient) GetAll(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*GetAllResponse, error) {
	out := new(GetAllResponse)
	err := c.cc.Invoke(ctx, "/auth_service.AuthService/GetAll", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authServiceClient) UpdateUsername(ctx context.Context, in *UpdateUsernameRequest, opts ...grpc.CallOption) (*UpdateUsernameResponse, error) {
	out := new(UpdateUsernameResponse)
	err := c.cc.Invoke(ctx, "/auth_service.AuthService/UpdateUsername", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authServiceClient) ChangePassword(ctx context.Context, in *ChangePasswordRequest, opts ...grpc.CallOption) (*ChangePasswordResponse, error) {
	out := new(ChangePasswordResponse)
	err := c.cc.Invoke(ctx, "/auth_service.AuthService/ChangePassword", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authServiceClient) ConfirmEmailLogin(ctx context.Context, in *ConfirmEmailLoginRequest, opts ...grpc.CallOption) (*ConfirmEmailLoginResponse, error) {
	out := new(ConfirmEmailLoginResponse)
	err := c.cc.Invoke(ctx, "/auth_service.AuthService/ConfirmEmailLogin", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authServiceClient) ActivateAccount(ctx context.Context, in *ActivationRequest, opts ...grpc.CallOption) (*ActivationResponse, error) {
	out := new(ActivationResponse)
	err := c.cc.Invoke(ctx, "/auth_service.AuthService/ActivateAccount", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authServiceClient) SendRecoveryCode(ctx context.Context, in *SendRecoveryCodeRequest, opts ...grpc.CallOption) (*SendRecoveryCodeResponse, error) {
	out := new(SendRecoveryCodeResponse)
	err := c.cc.Invoke(ctx, "/auth_service.AuthService/SendRecoveryCode", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authServiceClient) VerifyRecoveryCode(ctx context.Context, in *VerifyRecoveryCodeRequest, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/auth_service.AuthService/VerifyRecoveryCode", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authServiceClient) ResetForgottenPassword(ctx context.Context, in *ResetForgottenPasswordRequest, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/auth_service.AuthService/ResetForgottenPassword", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AuthServiceServer is the server API for AuthService service.
// All implementations must embed UnimplementedAuthServiceServer
// for forward compatibility
type AuthServiceServer interface {
	Register(context.Context, *RegisterRequest) (*RegisterResponse, error)
	Login(context.Context, *LoginRequest) (*LoginResponse, error)
	PasswordlessLogin(context.Context, *PasswordlessLoginRequest) (*PasswordlessLoginResponse, error)
	GetAll(context.Context, *Empty) (*GetAllResponse, error)
	UpdateUsername(context.Context, *UpdateUsernameRequest) (*UpdateUsernameResponse, error)
	ChangePassword(context.Context, *ChangePasswordRequest) (*ChangePasswordResponse, error)
	ConfirmEmailLogin(context.Context, *ConfirmEmailLoginRequest) (*ConfirmEmailLoginResponse, error)
	ActivateAccount(context.Context, *ActivationRequest) (*ActivationResponse, error)
	SendRecoveryCode(context.Context, *SendRecoveryCodeRequest) (*SendRecoveryCodeResponse, error)
	VerifyRecoveryCode(context.Context, *VerifyRecoveryCodeRequest) (*Response, error)
	ResetForgottenPassword(context.Context, *ResetForgottenPasswordRequest) (*Response, error)
	mustEmbedUnimplementedAuthServiceServer()
}

// UnimplementedAuthServiceServer must be embedded to have forward compatible implementations.
type UnimplementedAuthServiceServer struct {
}

func (UnimplementedAuthServiceServer) Register(context.Context, *RegisterRequest) (*RegisterResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Register not implemented")
}
func (UnimplementedAuthServiceServer) Login(context.Context, *LoginRequest) (*LoginResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Login not implemented")
}
func (UnimplementedAuthServiceServer) PasswordlessLogin(context.Context, *PasswordlessLoginRequest) (*PasswordlessLoginResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PasswordlessLogin not implemented")
}
func (UnimplementedAuthServiceServer) GetAll(context.Context, *Empty) (*GetAllResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAll not implemented")
}
func (UnimplementedAuthServiceServer) UpdateUsername(context.Context, *UpdateUsernameRequest) (*UpdateUsernameResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateUsername not implemented")
}
func (UnimplementedAuthServiceServer) ChangePassword(context.Context, *ChangePasswordRequest) (*ChangePasswordResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ChangePassword not implemented")
}
func (UnimplementedAuthServiceServer) ConfirmEmailLogin(context.Context, *ConfirmEmailLoginRequest) (*ConfirmEmailLoginResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ConfirmEmailLogin not implemented")
}
func (UnimplementedAuthServiceServer) ActivateAccount(context.Context, *ActivationRequest) (*ActivationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ActivateAccount not implemented")
}
func (UnimplementedAuthServiceServer) SendRecoveryCode(context.Context, *SendRecoveryCodeRequest) (*SendRecoveryCodeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendRecoveryCode not implemented")
}
func (UnimplementedAuthServiceServer) VerifyRecoveryCode(context.Context, *VerifyRecoveryCodeRequest) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method VerifyRecoveryCode not implemented")
}
func (UnimplementedAuthServiceServer) ResetForgottenPassword(context.Context, *ResetForgottenPasswordRequest) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ResetForgottenPassword not implemented")
}
func (UnimplementedAuthServiceServer) mustEmbedUnimplementedAuthServiceServer() {}

// UnsafeAuthServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AuthServiceServer will
// result in compilation errors.
type UnsafeAuthServiceServer interface {
	mustEmbedUnimplementedAuthServiceServer()
}

func RegisterAuthServiceServer(s grpc.ServiceRegistrar, srv AuthServiceServer) {
	s.RegisterService(&AuthService_ServiceDesc, srv)
}

func _AuthService_Register_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServiceServer).Register(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/auth_service.AuthService/Register",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServiceServer).Register(ctx, req.(*RegisterRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthService_Login_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoginRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServiceServer).Login(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/auth_service.AuthService/Login",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServiceServer).Login(ctx, req.(*LoginRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthService_PasswordlessLogin_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PasswordlessLoginRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServiceServer).PasswordlessLogin(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/auth_service.AuthService/PasswordlessLogin",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServiceServer).PasswordlessLogin(ctx, req.(*PasswordlessLoginRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthService_GetAll_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServiceServer).GetAll(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/auth_service.AuthService/GetAll",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServiceServer).GetAll(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthService_UpdateUsername_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateUsernameRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServiceServer).UpdateUsername(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/auth_service.AuthService/UpdateUsername",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServiceServer).UpdateUsername(ctx, req.(*UpdateUsernameRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthService_ChangePassword_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ChangePasswordRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServiceServer).ChangePassword(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/auth_service.AuthService/ChangePassword",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServiceServer).ChangePassword(ctx, req.(*ChangePasswordRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthService_ConfirmEmailLogin_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ConfirmEmailLoginRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServiceServer).ConfirmEmailLogin(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/auth_service.AuthService/ConfirmEmailLogin",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServiceServer).ConfirmEmailLogin(ctx, req.(*ConfirmEmailLoginRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthService_ActivateAccount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ActivationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServiceServer).ActivateAccount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/auth_service.AuthService/ActivateAccount",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServiceServer).ActivateAccount(ctx, req.(*ActivationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthService_SendRecoveryCode_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendRecoveryCodeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServiceServer).SendRecoveryCode(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/auth_service.AuthService/SendRecoveryCode",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServiceServer).SendRecoveryCode(ctx, req.(*SendRecoveryCodeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthService_VerifyRecoveryCode_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(VerifyRecoveryCodeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServiceServer).VerifyRecoveryCode(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/auth_service.AuthService/VerifyRecoveryCode",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServiceServer).VerifyRecoveryCode(ctx, req.(*VerifyRecoveryCodeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthService_ResetForgottenPassword_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ResetForgottenPasswordRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServiceServer).ResetForgottenPassword(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/auth_service.AuthService/ResetForgottenPassword",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServiceServer).ResetForgottenPassword(ctx, req.(*ResetForgottenPasswordRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// AuthService_ServiceDesc is the grpc.ServiceDesc for AuthService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var AuthService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "auth_service.AuthService",
	HandlerType: (*AuthServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Register",
			Handler:    _AuthService_Register_Handler,
		},
		{
			MethodName: "Login",
			Handler:    _AuthService_Login_Handler,
		},
		{
			MethodName: "PasswordlessLogin",
			Handler:    _AuthService_PasswordlessLogin_Handler,
		},
		{
			MethodName: "GetAll",
			Handler:    _AuthService_GetAll_Handler,
		},
		{
			MethodName: "UpdateUsername",
			Handler:    _AuthService_UpdateUsername_Handler,
		},
		{
			MethodName: "ChangePassword",
			Handler:    _AuthService_ChangePassword_Handler,
		},
		{
			MethodName: "ConfirmEmailLogin",
			Handler:    _AuthService_ConfirmEmailLogin_Handler,
		},
		{
			MethodName: "ActivateAccount",
			Handler:    _AuthService_ActivateAccount_Handler,
		},
		{
			MethodName: "SendRecoveryCode",
			Handler:    _AuthService_SendRecoveryCode_Handler,
		},
		{
			MethodName: "VerifyRecoveryCode",
			Handler:    _AuthService_VerifyRecoveryCode_Handler,
		},
		{
			MethodName: "ResetForgottenPassword",
			Handler:    _AuthService_ResetForgottenPassword_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "auth_service.proto",
}
