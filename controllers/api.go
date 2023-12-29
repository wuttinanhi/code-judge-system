package controllers

import (
	"runtime"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/storage/redis"
	"github.com/wuttinanhi/code-judge-system/services"
)

func SetupAPI(serviceKit *services.ServiceKit) *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler: ErrorHandler,
	})

	redisStore := redis.New(redis.Config{
		Host:      "127.0.0.1",
		Port:      6379,
		Username:  "",
		Password:  "redis",
		Database:  0,
		Reset:     false,
		TLSConfig: nil,
		PoolSize:  10 * runtime.GOMAXPROCS(0),
	})

	app.Use(limiter.New(limiter.Config{
		Max:        100,
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
		Storage: redisStore,
	}))

	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:5173",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))

	app.Use(recover.New())

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	authHandler := NewAuthHandler(serviceKit)
	challengeHandler := NewChallengeHandler(serviceKit)
	submissionHandler := NewSubmissionHandler(serviceKit)
	challengeTestcaseHandler := NewChallengeTestcaseHandler(serviceKit)

	userGroup := app.Group("/user")
	userGroup.Post("/register", authHandler.Register)
	userGroup.Post("/login", authHandler.Login)

	challengeGroup := app.Group("/challenge")
	challengeGroup.Use(UserMiddleware(serviceKit))
	challengeGroup.Post("/create", challengeHandler.CreateChallengeWithTestcase)
	challengeGroup.Get("/pagination", challengeHandler.PaginationChallengesWithStatus)
	challengeGroup.Get("/get/:id", challengeHandler.GetChallengeByID)
	challengeGroup.Put("/update/:id", challengeHandler.UpdateChallenge)
	challengeGroup.Delete("/delete/:id", challengeHandler.DeleteChallenge)

	testcaseGroup := app.Group("/testcase")
	testcaseGroup.Use(UserMiddleware(serviceKit))
	testcaseGroup.Post("/create", challengeTestcaseHandler.CreateTestcase)
	testcaseGroup.Get("/get/:id", challengeTestcaseHandler.GetTestcaseByID)
	testcaseGroup.Put("/update", challengeTestcaseHandler.UpdateTestcase)
	testcaseGroup.Delete("/delete/:id", challengeTestcaseHandler.DeleteTestcase)

	submissionGroup := app.Group("/submission")
	submissionGroup.Use(UserMiddleware(serviceKit))
	submissionGroup.Post("/submit", submissionHandler.SubmitSubmission)
	submissionGroup.Get("/pagination", submissionHandler.Pagination)
	// submissionGroup.Get("/get/user", submissionHandler.GetSubmissionByUser)
	// submissionGroup.Get("/get/challenge/:id", submissionHandler.GetSubmissionByChallenge)
	// submissionGroup.Get("/get/submission/:id", submissionHandler.GetSubmissionByID)

	return app
}
