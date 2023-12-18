package services

import (
	"github.com/wuttinanhi/code-judge-system/entities"
)

type JWTService interface {
	GenerateToken(user entities.User) (string, error)
	ValidateToken(token string) (*entities.User, error)
}
