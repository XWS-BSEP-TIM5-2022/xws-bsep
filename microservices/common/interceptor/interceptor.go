package interceptor

import (
	"context"
	"crypto/rsa"
	"fmt"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type AuthInterceptor struct {
	accessiblePermissions map[string]string
	publicKey             *rsa.PublicKey
}

func NewAuthInterceptor(accessiblePermissions map[string]string, publicKey *rsa.PublicKey) *AuthInterceptor {
	return &AuthInterceptor{
		accessiblePermissions: accessiblePermissions,
		publicKey:             publicKey,
	}
}

func (interceptor *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		ctx, err := interceptor.Authorize(ctx, info.FullMethod)
		if err != nil {
			return nil, err
		}

		return handler(ctx, req)
	}
}

func (interceptor *AuthInterceptor) Authorize(ctx context.Context, method string) (context.Context, error) {
	accessiblePermission, ok := interceptor.accessiblePermissions[method]
	if !ok {
		return ctx, nil
	}

	var values []string
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ctx, status.Errorf(codes.Unauthenticated, "Unauthorized")
	}

	values = md.Get("Authorization")
	if len(values) == 0 {
		return ctx, status.Errorf(codes.Unauthenticated, "Unauthorized")
	}

	authHeader := values[0]
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 {
		return ctx, status.Errorf(codes.Unauthenticated, "Unauthorized")
	}

	claims, err := interceptor.verifyToken(parts[1])
	if err != nil {
		return ctx, status.Errorf(codes.Unauthenticated, "Unauthorized")
	}

	for _, jwtPermission := range claims.Permissions {
		if accessiblePermission == jwtPermission {
			return context.WithValue(ctx, LoggedInUserKey{}, claims.Subject), nil
		}
	}

	return ctx, status.Errorf(codes.PermissionDenied, "Unauthorized")
}

func (interceptor *AuthInterceptor) verifyToken(accessToken string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(
		accessToken,
		&UserClaims{},
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodRSA)
			if !ok {
				return nil, fmt.Errorf("Unexpected token signing method")
			}

			return interceptor.publicKey, nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("Invalid token: %w", err)
	}
	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		return nil, fmt.Errorf("Invalid token claims")
	}

	return claims, nil
}

type UserClaims struct {
	Username    string   `json:"username"`
	Roles       []string `json:"roles"`
	Permissions []string `json:"permissions"`
	jwt.StandardClaims
}

type LoggedInUserKey struct{}
