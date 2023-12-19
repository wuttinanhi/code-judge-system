package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wuttinanhi/code-judge-system/entities"
	"github.com/wuttinanhi/code-judge-system/services"
)

func SubmitSubmission(c *fiber.Ctx) error {
	dto := entities.ValidateSubmissionCreateDTO(c)

	user := entities.GetUserFromRequest(c)

	challenge, err := services.GetServiceKit().ChallengeService.FindChallengeByID(dto.ChallengeID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entities.HttpError{Message: err.Error()})
	}

	submission, err := services.GetServiceKit().SubmissionService.SubmitSubmission(&entities.Submission{
		ChallengeID: challenge.ChallengeID,
		UserID:      user.UserID,
		Language:    dto.Language,
		SourceCode:  dto.SourceCode,
	})
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entities.HttpError{Message: err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(submission)
}

func GetSubmissionByID(c *fiber.Ctx) error {
	id := ParseIntParam(c, "id")

	submission, err := services.GetServiceKit().SubmissionService.GetSubmissionByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entities.HttpError{Message: err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(submission)
}

func GetSubmissionByUser(c *fiber.Ctx) error {
	user := entities.GetUserFromRequest(c)

	submissions, err := services.GetServiceKit().SubmissionService.GetSubmissionByUser(user)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entities.HttpError{Message: err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(submissions)
}

func GetSubmissionByChallenge(c *fiber.Ctx) error {
	id := ParseIntParam(c, "id")

	challenge, err := services.GetServiceKit().ChallengeService.FindChallengeByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entities.HttpError{Message: err.Error()})
	}

	submissions, err := services.GetServiceKit().SubmissionService.GetSubmissionByChallenge(challenge)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entities.HttpError{Message: err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(submissions)
}
