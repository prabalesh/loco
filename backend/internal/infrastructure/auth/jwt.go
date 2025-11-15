package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	UserID int    `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

type JWTService struct {
	accessTokenSecretKey  []byte
	refreshTokenSecretKey []byte
	accessTokenExpires    time.Duration
	refreshTokenExpires   time.Duration
}

func NewJWTService(accessSecret, refreshSecret string, accessTokenExpires, refreshTokenExpires time.Duration) *JWTService {
	return &JWTService{
		accessTokenSecretKey:  []byte(accessSecret),
		refreshTokenSecretKey: []byte(refreshSecret),
		accessTokenExpires:    accessTokenExpires,
		refreshTokenExpires:   refreshTokenExpires,
	}
}

func (j *JWTService) GenerateAccessToken(userID int, email, role string) (string, time.Duration, error) {
	expiresAt := time.Now().Add(j.accessTokenExpires)
	claims := JWTClaims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "loco",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(j.accessTokenSecretKey)
	return tokenStr, j.accessTokenExpires, err
}

func (j *JWTService) GenerateRefreshToken(userID int, email string) (string, time.Duration, error) {
	expiresAt := time.Now().Add(j.accessTokenExpires)
	claims := JWTClaims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "loco",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(j.refreshTokenSecretKey)
	return tokenStr, j.refreshTokenExpires, err

}

func (j *JWTService) ValidateToken(tokenString string, isRefresh bool) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		if isRefresh {
			return j.refreshTokenSecretKey, nil
		}
		return j.accessTokenSecretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
