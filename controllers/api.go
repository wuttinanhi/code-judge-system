package controllers

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/spf13/viper"

	"github.com/wuttinanhi/code-judge-system/services"
)

func SetupAPI(serviceKit *services.ServiceKit, ratelimitStorage fiber.Storage) *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler: ErrorHandler,
	})

	app.Use(limiter.New(limiter.Config{
		Max:        150,
		Expiration: 60 * time.Second,
		KeyGenerator: func(c *fiber.Ctx) string {
			if c.IP() != "" {
				return c.IP()
			}
			return c.Get("x-forwarded-for")
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.SendStatus(fiber.StatusTooManyRequests)
		},
		Storage: ratelimitStorage,
	}))

	allowOrigins := viper.GetStringSlice("APP_API_CORS_ALLOW_ORIGINS")
	app.Use(cors.New(cors.Config{
		AllowOriginsFunc: func(origin string) bool {
			for _, allowOrigin := range allowOrigins {
				if origin == allowOrigin {
					return true
				}
			}
			return false
		},
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))

	log.Println("Allow Origins:", allowOrigins)

	app.Use(recover.New())

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	authHandler := NewAuthHandler(serviceKit)
	userHandler := NewUserHandler(serviceKit)
	challengeHandler := NewChallengeHandler(serviceKit)
	submissionHandler := NewSubmissionHandler(serviceKit)
	// challengeTestcaseHandler := NewChallengeTestcaseHandler(serviceKit)

	authGroup := app.Group("/auth")
	authGroup.Post("/register", authHandler.Register)
	authGroup.Post("/login", authHandler.Login)
	authGroup.Get("/me", authHandler.Me)

	userGroup := app.Group("/user")
	userGroup.Use(UserMiddleware(serviceKit))
	userGroup.Get("/me", userHandler.Me)
	userGroup.Put("/update/role", userHandler.UpdateRole)
	userGroup.Get("/pagination", userHandler.Pagination)

	challengeGroup := app.Group("/challenge")
	challengeGroup.Use(UserMiddleware(serviceKit))
	challengeGroup.Post("/create", challengeHandler.CreateChallengeWithTestcase)
	challengeGroup.Get("/pagination", challengeHandler.PaginationChallengesWithStatus)
	challengeGroup.Get("/get/:id", challengeHandler.GetChallengeByID)
	challengeGroup.Put("/update/:id", challengeHandler.UpdateChallenge)
	challengeGroup.Delete("/delete/:id", challengeHandler.DeleteChallenge)

	// testcaseGroup := app.Group("/testcase")
	// testcaseGroup.Use(UserMiddleware(serviceKit))
	// testcaseGroup.Post("/create", challengeTestcaseHandler.CreateTestcase)
	// testcaseGroup.Get("/get/:id", challengeTestcaseHandler.GetTestcaseByID)
	// testcaseGroup.Put("/update", challengeTestcaseHandler.UpdateTestcase)
	// testcaseGroup.Delete("/delete/:id", challengeTestcaseHandler.DeleteTestcase)

	submissionGroup := app.Group("/submission")
	submissionGroup.Use(UserMiddleware(serviceKit))
	submissionGroup.Post("/submit", submissionHandler.SubmitSubmission)
	submissionGroup.Get("/pagination", submissionHandler.Pagination)
	submissionGroup.Get("/get/:id", submissionHandler.GetSubmissionByID)
	// submissionGroup.Get("/get/user", submissionHandler.GetSubmissionByUser)
	// submissionGroup.Get("/get/challenge/:id", submissionHandler.GetSubmissionByChallenge)

	return app
}
