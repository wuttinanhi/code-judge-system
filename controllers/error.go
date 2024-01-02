package controllers

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/wuttinanhi/code-judge-system/entities"
)

func ErrorHandler(c *fiber.Ctx, err error) error {
	// if error is "user not found"
	if err.Error() == "user not found" {
		return c.Status(fiber.StatusUnauthorized).JSON(entities.HttpError{
			Message: "Unauthorized",
		})
	}

	// if error is go-validator error
	if errs, ok := err.(validator.ValidationErrors); ok {
		var errors []string
		for _, e := range errs {
			errors = append(errors, fmt.Sprintf("%s %s", e.Field(), e.Tag()))
		}

		return c.Status(fiber.StatusBadRequest).JSON(entities.HttpBadRequest{
			Error:   "Bad Request",
			Message: "Validation error",
			Errors:  errors,
		})
	}

	// default error
	return c.Status(fiber.StatusInternalServerError).JSON(entities.HttpError{
		Message: err.Error(),
	})
}
