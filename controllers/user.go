package controllers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/wuttinanhi/code-judge-system/entities"
	"github.com/wuttinanhi/code-judge-system/services"
)

type userHandler struct {
	serviceKit *services.ServiceKit
}

func (h *userHandler) Me(c *fiber.Ctx) error {
	user := GetUserFromRequest(c)

	return c.Status(fiber.StatusOK).JSON(user)
}

func (h *userHandler) UpdateRole(c *fiber.Ctx) error {
	user := GetUserFromRequest(c)
	dto := entities.ValidateUserUpdateRoleDTO(c)

	// only user with role admin can update role
	if user.Role != entities.UserRoleAdmin {
		return c.SendStatus(fiber.StatusForbidden)
	}

	// get target user
	targetUser, err := h.serviceKit.UserService.FindUserByID(dto.UserID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entities.HttpError{Message: err.Error()})
	}

	// update role
	err = h.serviceKit.UserService.UpdateRole(targetUser, dto.Role)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entities.HttpError{Message: err.Error()})
	}

	return c.SendStatus(fiber.StatusOK)
}

func (h *challengeHandler) Pagination(c *fiber.Ctx) error {
	user := GetUserFromRequest(c)
	options := ParsePaginationOptions(c)

	// only user with role admin can get all users
	if user.Role != entities.UserRoleAdmin {
		return c.SendStatus(fiber.StatusForbidden)
	}

	result, err := h.serviceKit.UserService.Pagination(options)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(entities.HttpError{Message: err.Error()})
	}

	return c.Status(http.StatusOK).JSON(result)
}

func NewUserHandler(serviceKit *services.ServiceKit) *userHandler {
	return &userHandler{
		serviceKit: serviceKit,
	}
}
