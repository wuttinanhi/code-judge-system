package services

import "github.com/wuttinanhi/code-judge-system/entities"

type UserService interface {
	Register(email, password, displayname string) (user *entities.User, err error)
	Login(email, password string) (user *entities.User, err error)
}
