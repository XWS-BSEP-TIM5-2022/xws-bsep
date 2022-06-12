package api

import (
	"crypto/rsa"
	"time"

	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/auth-service/domain"
	"github.com/dgrijalva/jwt-go"
)

type APITokenService struct {
	accessAPITokenDuration time.Duration
	privateKey             *rsa.PrivateKey
	publicKey              *rsa.PublicKey
}

type APITokenClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func NewAPITokenManager(privateKey, publicKey string) (*APITokenService, error) {
	parsedPrivateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(privateKey))
	if err != nil {
		return nil, err
	}
	parsedPublicKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(publicKey))
	if err != nil {
		return nil, err
	}
	return &APITokenService{
		privateKey:             parsedPrivateKey,
		publicKey:              parsedPublicKey,
		accessAPITokenDuration: 5 * 60 * time.Minute,
	}, nil
}

func (manager *APITokenService) GenerateAPIToken(auth *domain.Authentication) (string, string, error) {
	claims := APITokenClaims{
		StandardClaims: jwt.StandardClaims{
			Subject: auth.Id,
			//ExpiresAt: time.Now().Add(manager.accessAPITokenDuration).Unix(), // TODO: mozda zbog bezbednosti
		},
		Username: auth.Username,
	}
	token := jwt.NewWithClaims(
		jwt.SigningMethodRS256,
		claims,
	)

	signed, _ := token.SignedString(manager.privateKey)
	hashedPassword, err := auth.HashPassword(signed) // hesiranje
	if err != nil {
		return "error", "error", err
	}
	return signed, hashedPassword, nil
}
