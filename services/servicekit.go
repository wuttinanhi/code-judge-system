package services

import (
	"os"

	"github.com/wuttinanhi/code-judge-system/databases"
	"github.com/wuttinanhi/code-judge-system/repositories"
	"gorm.io/gorm"
)

type ServiceKit struct {
	JWTService        JWTService
	UserService       UserService
	ChallengeService  ChallengeService
	SubmissionService SubmissionService
	SandboxService    SandboxService
}

func newServiceKit(db *gorm.DB) *ServiceKit {
	userRepo := repositories.NewSQLUserRepository(db)
	challengeRepo := repositories.NewChallengeRepository(db)
	submissionRepo := repositories.NewSubmissionRepository(db)

	// read env var "JWT_SECRET" and pass it to JWTService
	// if JWT_SECRET is empty, use default value
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "secret"
	}

	jwtService := NewJWTService(jwtSecret)
	userService := NewUserService(userRepo)
	challengeService := NewChallengeService(challengeRepo)
	submissionService := NewSubmissionService(submissionRepo, challengeService)
	sandboxService := NewSandboxService()

	return &ServiceKit{
		JWTService:        jwtService,
		UserService:       userService,
		ChallengeService:  challengeService,
		SubmissionService: submissionService,
		SandboxService:    sandboxService,
	}
}

func CreateServiceKit() *ServiceKit {
	db := databases.NewSQLiteDatabase()
	return newServiceKit(db)
}

func CreateTestServiceKit() *ServiceKit {
	db := databases.NewTempSQLiteDatabase()
	return newServiceKit(db)
}
