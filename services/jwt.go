package services

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/wuttinanhi/code-judge-system/entities"
)

type JWTService interface {
	GenerateToken(user entities.User) (string, error)
	ValidateToken(token string) (*entities.User, error)
}

type jwtService struct {
	secret string
}

// GenerateToken implements services.JWTService.
func (s *jwtService) GenerateToken(user entities.User) (string, error) {
	claims := &entities.JWTClaims{
		UserID:      user.UserID,
		DisplayName: user.DisplayName,
		Email:       user.Email,
		Role:        user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secret))
}

// ValidateToken implements services.JWTService.
func (s *jwtService) ValidateToken(token string) (*entities.User, error) {
	parsedToken, err := jwt.ParseWithClaims(token, &entities.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.secret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := parsedToken.Claims.(*entities.JWTClaims)
	if !ok || !parsedToken.Valid {
		return nil, errors.New("invalid token")
	}

	user := &entities.User{
		UserID:      claims.UserID,
		DisplayName: claims.DisplayName,
		Email:       claims.Email,
		Role:        claims.Role,
	}

	return user, nil
}

func NewJWTService(secret string) JWTService {
	return &jwtService{
		secret: secret,
	}
}
