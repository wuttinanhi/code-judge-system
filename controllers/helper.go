package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wuttinanhi/code-judge-system/entities"
)

func GetUserFromRequest(c *fiber.Ctx) *entities.User {
	user, ok := c.Locals("user").(*entities.User)
	if !ok {
		panic("User not found")
	}

	return user
}

func ParseIntParam(c *fiber.Ctx, paramName string) int {
	value, err := c.ParamsInt(paramName)
	if err != nil {
		panic(err)
	}
	return value
}
