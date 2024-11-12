package token

import (
	"github.com/golang-jwt/jwt/v5"
)

type TokenUseCase interface {
	GenerateAccessToken(claims JwtCustomClaims) (string, error)
}

type tokenUseCase struct {
	secretKey string
}

func NewTokenUseCase(secretKey string) TokenUseCase {
	return &tokenUseCase{secretKey}
}

type JwtCustomClaims struct {
	Username string `json:"username"`
	UserID uint `json:"user_id"`
	Role     string `json:"role"`
	FullName string `json:"full_name"`
	jwt.RegisteredClaims
}

func (t *tokenUseCase) GenerateAccessToken(claims JwtCustomClaims) (string, error) {
	plainToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	encodedToken, err := plainToken.SignedString([]byte(t.secretKey))
	if err != nil {
		return "", err
	}

	return encodedToken, nil
}
