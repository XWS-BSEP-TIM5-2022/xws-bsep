package application

import (
	"crypto/rsa"
	"time"

	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/auth-service/domain"
	"github.com/dgrijalva/jwt-go"
)

type JWTManager struct {
	privateKey          *rsa.PrivateKey
	publicKey           *rsa.PublicKey
	accessTokenDuration time.Duration
}

type Claims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.StandardClaims
}

func (manager *JWTManager) GenerateToken(auth *domain.Authentication) (string, error) {
	claims := Claims{
		StandardClaims: jwt.StandardClaims{
			Subject:   auth.Id,
			ExpiresAt: time.Now().Add(manager.accessTokenDuration).Unix(),
		},
		Username: auth.Username,
		Role:     auth.Role,
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodRS256,
		claims,
	)

	return token.SignedString(manager.privateKey)
}
