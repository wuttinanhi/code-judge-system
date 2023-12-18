package entities

import "github.com/golang-jwt/jwt/v5"

type JWTClaims struct {
	UserID      uint   `json:"user_id"`
	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
	Role        string `json:"role"`
	jwt.RegisteredClaims
}
