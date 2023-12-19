package cmds

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/wuttinanhi/code-judge-system/controllers"
)

func SetupWeb() *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler: controllers.ErrorHandler,
	})

	app.Use(recover.New())

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	userGroup := app.Group("/user")
	userGroup.Post("/register", controllers.Register)
	userGroup.Post("/login", controllers.Login)

	challengeGroup := app.Group("/challenge")
	challengeGroup.Use(controllers.UserMiddleware)
	challengeGroup.Post("/create", controllers.CreateChallengeWithTestcase)
	challengeGroup.Get("/all", controllers.GetAllChallenges)
	challengeGroup.Get("/get/:id", controllers.GetChallengeByID)
	challengeGroup.Put("/update", controllers.UpdateChallenge)
	challengeGroup.Delete("/delete/:id", controllers.DeleteChallenge)

	testcaseGroup := app.Group("/testcase")
	testcaseGroup.Use(controllers.UserMiddleware)
	testcaseGroup.Post("/create", controllers.CreateTestcase)
	testcaseGroup.Get("/get/:id", controllers.GetTestcaseByID)
	testcaseGroup.Put("/update", controllers.UpdateTestcase)
	testcaseGroup.Delete("/delete/:id", controllers.DeleteTestcase)

	submissionGroup := app.Group("/submission")
	submissionGroup.Use(controllers.UserMiddleware)
	submissionGroup.Post("/submit", controllers.SubmitSubmission)
	submissionGroup.Get("/get/:id", controllers.GetSubmissionByID)
	submissionGroup.Get("/get/user", controllers.GetSubmissionByUser)
	submissionGroup.Get("/get/challenge/:id", controllers.GetSubmissionByChallenge)

	return app
}
