package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wuttinanhi/code-judge-system/services"
)

type userHandler struct {
	serviceKit *services.ServiceKit
}

func (h *userHandler) Me(c *fiber.Ctx) error {
	user := GetUserFromRequest(c)

	return c.Status(fiber.StatusOK).JSON(user)
}

func NewUserHandler(serviceKit *services.ServiceKit) *userHandler {
	return &userHandler{
		serviceKit: serviceKit,
	}
}
