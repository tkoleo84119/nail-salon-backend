package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/tkoleo84119/nail-salon-backend/internal/config"
)

type Store struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type JWTClaims struct {
	UserID    int64   `json:"user_id"`
	Username  string  `json:"username"`
	Role      string  `json:"role"`
	StoreList []Store `json:"store_list"`
	jwt.RegisteredClaims
}

func GenerateJWT(jwtConfig config.JWTConfig, userID int64, username, role string, storeList []Store) (string, error) {
	claims := JWTClaims{
		UserID:    userID,
		Username:  username,
		Role:      role,
		StoreList: storeList,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(jwtConfig.ExpiryHours) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtConfig.Secret))
}

func ValidateJWT(jwtConfig config.JWTConfig, tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtConfig.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
