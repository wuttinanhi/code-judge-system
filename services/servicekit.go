package services

import (
	"os"

	"github.com/wuttinanhi/code-judge-system/repositories"
	"gorm.io/gorm"
)

type ServiceKit struct {
	JWTService        JWTService
	UserService       UserService
	ChallengeService  ChallengeService
	SubmissionService SubmissionService
	SandboxService    SandboxService
	KafkaService      KafkaService
}

func CreateServiceKit(db *gorm.DB) *ServiceKit {
	userRepo := repositories.NewUserRepository(db)
	challengeRepo := repositories.NewChallengeRepository(db)
	submissionRepo := repositories.NewSubmissionRepository(db)

	// read env var "JWT_SECRET" and pass it to JWTService
	// if JWT_SECRET is empty, use default value
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "secret"
	}

	kafkaHost := os.Getenv("KAFKA_HOST")

	jwtService := NewJWTService(jwtSecret)
	userService := NewUserService(userRepo)
	challengeService := NewChallengeService(challengeRepo)
	sandboxService := NewSandboxService()
	submissionService := NewSubmissionService(submissionRepo, challengeService, sandboxService)
	kafkaService := NewKafkaService(kafkaHost)

	return &ServiceKit{
		JWTService:        jwtService,
		UserService:       userService,
		ChallengeService:  challengeService,
		SubmissionService: submissionService,
		SandboxService:    sandboxService,
		KafkaService:      kafkaService,
	}
}

func CreateTestServiceKit(db *gorm.DB) *ServiceKit {
	userRepo := repositories.NewUserRepository(db)
	challengeRepo := repositories.NewChallengeRepository(db)
	submissionRepo := repositories.NewSubmissionRepository(db)

	jwtService := NewJWTService("test")
	userService := NewUserService(userRepo)
	challengeService := NewChallengeService(challengeRepo)
	sandboxService := NewSandboxService()
	submissionService := NewSubmissionService(submissionRepo, challengeService, sandboxService)
	kafkaService := NewKafkaMockService()

	return &ServiceKit{
		JWTService:        jwtService,
		UserService:       userService,
		ChallengeService:  challengeService,
		SubmissionService: submissionService,
		SandboxService:    sandboxService,
		KafkaService:      kafkaService,
	}
}
