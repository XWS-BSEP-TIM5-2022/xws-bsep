package application

import (
	"crypto/rsa"
	"time"

	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/auth-service/domain"
	"github.com/dgrijalva/jwt-go"
)

type JWTService struct {
	privateKey          *rsa.PrivateKey
	publicKey           *rsa.PublicKey
	accessTokenDuration time.Duration
}

type Claims struct {
	Username    string   `json:"username"`
	Roles       []string `json:"roles"`
	Permissions []string `json:"permissions"`
	jwt.StandardClaims
}

func NewJWTManager(privateKey, publicKey string) (*JWTService, error) {
	parsedPrivateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(privateKey))
	if err != nil {
		return nil, err
	}
	parsedPublicKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(publicKey))
	if err != nil {
		return nil, err
	}
	return &JWTService{
		privateKey:          parsedPrivateKey,
		publicKey:           parsedPublicKey,
		accessTokenDuration: 60 * time.Minute,
	}, nil
}

func (manager *JWTService) GenerateToken(auth *domain.Authentication) (string, error) {
	var roleNames, permissionNames []string
	for _, role := range *auth.Roles {
		roleNames = append(roleNames, role.Name)
		for _, permission := range role.Permissions {
			permissionNames = append(permissionNames, permission.Name)
		}
	}
	claims := Claims{
		StandardClaims: jwt.StandardClaims{
			Subject:   auth.Id,
			ExpiresAt: time.Now().Add(manager.accessTokenDuration).Unix(),
		},
		Username:    auth.Username,
		Roles:       roleNames,
		Permissions: permissionNames,
	}
	token := jwt.NewWithClaims(
		jwt.SigningMethodRS256,
		claims,
	)
	return token.SignedString(manager.privateKey)
}
