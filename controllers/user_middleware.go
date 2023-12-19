package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wuttinanhi/code-judge-system/services"
)

func UserMiddleware(serviceKit *services.ServiceKit) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		// Get the token from the Authorization header
		tokenStr := c.Get("Authorization")
		if tokenStr == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Authorization header is missing",
			})
		}

		// Remove the Bearer prefix from the token
		tokenStr = tokenStr[7:]

		// Parse the token
		user, err := serviceKit.JWTService.ValidateToken(tokenStr)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token",
			})
		}

		// Store the user in the context's locals
		c.Locals("user", user)

		// Call the next middleware in the stack
		return c.Next()
	}
}
