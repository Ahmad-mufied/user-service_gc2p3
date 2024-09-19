package service

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type JWTManager struct {
	secretKey     string
	tokenDuration time.Duration
}

type ClientClaims struct {
	jwt.RegisteredClaims
	ServiceName string `json:"service_name"`
	Role        string `json:"role"`
}

func NewJWTManager(secretKey string, tokenDuration time.Duration) *JWTManager {
	return &JWTManager{secretKey, tokenDuration}
}

func (manager *JWTManager) Generate(serviceName, role string) (string, error) {
	claims := ClientClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(manager.tokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		ServiceName: serviceName,
		Role:        role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(manager.secretKey))
}

func (manager *JWTManager) Verify(accessToken string) (*ClientClaims, error) {
	token, err := jwt.ParseWithClaims(accessToken, &ClientClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(manager.secretKey), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*ClientClaims)
	if !ok {
		return nil, err
	}

	return claims, nil
}
