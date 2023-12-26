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
	value, err := c.ParamsInt(paramName, 0)
	if err != nil {
		panic(err)
	}
	return value
}

func ParseIntQuery(c *fiber.Ctx, paramName string) int {
	value := c.QueryInt(paramName, 0)
	return value
}

func ParsePaginationOptions(c *fiber.Ctx) *entities.PaginationOptions {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	sort := c.Query("sort", "id")
	order := c.Query("order", "asc")

	return &entities.PaginationOptions{
		Page:  page,
		Limit: limit,
		Sort:  sort,
		Order: order,
	}
}
