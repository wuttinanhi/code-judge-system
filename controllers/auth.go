package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wuttinanhi/code-judge-system/entities"
	"github.com/wuttinanhi/code-judge-system/services"
)

func Register(c *fiber.Ctx) error {
	dto := entities.ValidateUserRegisterDTO(c)

	user, err := services.SERVICE_KIT.UserService.Register(dto.Email, dto.Password, dto.DisplayName)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entities.HttpError{
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(entities.UserRegisterResponse{
		UserID:      user.ID,
		DisplayName: user.DisplayName,
		Email:       user.Email,
	})
}

func Login(c *fiber.Ctx) error {
	dto := entities.ValidateUserLoginDTO(c)

	user, err := services.SERVICE_KIT.UserService.Login(dto.Email, dto.Password)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entities.HttpError{
			Message: err.Error(),
		})
	}

	token, err := services.SERVICE_KIT.JWTService.GenerateToken(*user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(entities.HttpError{
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(entities.UserLoginResponse{Token: token})
}
