package controllers

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/wuttinanhi/code-judge-system/entities"
)

func ErrorHandler(c *fiber.Ctx, err error) error {
	// if error is go-validator error
	if err, ok := err.(validator.ValidationErrors); ok {
		return c.Status(fiber.StatusBadRequest).JSON(entities.HttpError{
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusInternalServerError).JSON(entities.HttpError{
		Message: err.Error(),
	})
}
