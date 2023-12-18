package services

import (
	"os"

	"github.com/wuttinanhi/code-judge-system/databases"
	repositories "github.com/wuttinanhi/code-judge-system/repositories/implements"
	services "github.com/wuttinanhi/code-judge-system/services/implements"
	Iservices "github.com/wuttinanhi/code-judge-system/services/interfaces"
	"gorm.io/gorm"
)

var SERVICE_KIT *ServiceKit

type ServiceKit struct {
	JWTService  Iservices.JWTService
	UserService Iservices.UserService
}

func newServiceKit(db *gorm.DB) *ServiceKit {
	userRepo := repositories.NewSQLUserRepository(db)

	// read env var "JWT_SECRET" and pass it to JWTService
	// if JWT_SECRET is empty, use default value
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "secret"
	}

	jwtService := services.NewJWTService(jwtSecret)
	userService := services.NewUserService(userRepo)

	return &ServiceKit{
		JWTService:  jwtService,
		UserService: userService,
	}
}

func InitServiceKit() {
	db := databases.NewSQLiteDatabase()
	SERVICE_KIT = newServiceKit(db)
}

func InitTestServiceKit() {
	db := databases.NewTempSQLiteDatabase()
	SERVICE_KIT = newServiceKit(db)
}
