package controllers

import "github.com/gofiber/fiber/v2"

func ParseIntParam(c *fiber.Ctx, paramName string) int {
	value, err := c.ParamsInt(paramName)
	if err != nil {
		panic(err)
	}
	return value
}
