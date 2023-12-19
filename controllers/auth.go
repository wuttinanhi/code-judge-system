package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wuttinanhi/code-judge-system/entities"
	"github.com/wuttinanhi/code-judge-system/services"
)

type authHandler struct {
	serviceKit *services.ServiceKit
}

func (h *authHandler) Register(c *fiber.Ctx) error {
	dto := entities.ValidateUserRegisterDTO(c)

	user, err := h.serviceKit.UserService.Register(dto.Email, dto.Password, dto.DisplayName)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entities.HttpError{
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(entities.UserRegisterResponse{
		UserID:      user.UserID,
		DisplayName: user.DisplayName,
		Email:       user.Email,
	})
}

func (h *authHandler) Login(c *fiber.Ctx) error {
	dto := entities.ValidateUserLoginDTO(c)

	user, err := h.serviceKit.UserService.Login(dto.Email, dto.Password)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entities.HttpError{
			Message: err.Error(),
		})
	}

	token, err := h.serviceKit.JWTService.GenerateToken(*user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(entities.HttpError{
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(entities.UserLoginResponse{Token: token})
}

func NewAuthHandler(serviceKit *services.ServiceKit) *authHandler {
	return &authHandler{
		serviceKit: serviceKit,
	}
}
