package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wuttinanhi/code-judge-system/entities"
	"github.com/wuttinanhi/code-judge-system/services"
)

type submissionHandler struct {
	serviceKit *services.ServiceKit
}

func (h *submissionHandler) SubmitSubmission(c *fiber.Ctx) error {
	dto := entities.ValidateSubmissionCreateDTO(c)

	user := GetUserFromRequest(c)

	challenge, err := h.serviceKit.ChallengeService.FindChallengeByID(dto.ChallengeID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entities.HttpError{Message: err.Error()})
	}

	submission, err := h.serviceKit.SubmissionService.SubmitSubmission(&entities.Submission{
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

func (h *submissionHandler) GetSubmissionByID(c *fiber.Ctx) error {
	id := ParseIntParam(c, "id")

	submission, err := h.serviceKit.SubmissionService.GetSubmissionByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entities.HttpError{Message: err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(submission)
}

func (h *submissionHandler) GetSubmissionByUser(c *fiber.Ctx) error {
	user := GetUserFromRequest(c)

	submissions, err := h.serviceKit.SubmissionService.GetSubmissionByUser(user)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entities.HttpError{Message: err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(submissions)
}

func (h *submissionHandler) GetSubmissionByChallenge(c *fiber.Ctx) error {
	id := ParseIntParam(c, "id")

	challenge, err := h.serviceKit.ChallengeService.FindChallengeByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entities.HttpError{Message: err.Error()})
	}

	submissions, err := h.serviceKit.SubmissionService.GetSubmissionByChallenge(challenge)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entities.HttpError{Message: err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(submissions)
}

func NewSubmissionHandler(serviceKit *services.ServiceKit) *submissionHandler {
	return &submissionHandler{
		serviceKit: serviceKit,
	}
}
