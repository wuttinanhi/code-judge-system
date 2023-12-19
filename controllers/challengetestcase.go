package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wuttinanhi/code-judge-system/entities"
	"github.com/wuttinanhi/code-judge-system/services"
)

type challengeTestcaseHandler struct {
	serviceKit *services.ServiceKit
}

func (h *challengeTestcaseHandler) CreateTestcase(c *fiber.Ctx) error {
	dto := entities.ValidateChallengeTestcaseCreateDTO(c)

	challenge, err := h.serviceKit.ChallengeService.FindChallengeByID(dto.ChallengeID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entities.HttpError{Message: err.Error()})
	}

	testcase, err := h.serviceKit.ChallengeService.AddTestcase(challenge, &entities.ChallengeTestcase{
		Input:          dto.Input,
		ExpectedOutput: dto.ExpectedOutput,
		ChallengeID:    challenge.ChallengeID,
	})
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entities.HttpError{Message: err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(testcase)
}

func (h *challengeTestcaseHandler) UpdateTestcase(c *fiber.Ctx) error {
	dto := entities.ValidateChallengeTestcaseUpdateDTO(c)

	testcase, err := h.serviceKit.ChallengeService.FindTestcaseByID(dto.TestcaseID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entities.HttpError{Message: err.Error()})
	}

	err = h.serviceKit.ChallengeService.UpdateTestcase(&entities.ChallengeTestcase{
		TestcaseID:     testcase.TestcaseID,
		Input:          dto.Input,
		ExpectedOutput: dto.ExpectedOutput,
		ChallengeID:    testcase.ChallengeID,
	})
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entities.HttpError{Message: err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(testcase)
}

func (h *challengeTestcaseHandler) DeleteTestcase(c *fiber.Ctx) error {
	id := ParseIntParam(c, "id")

	testcase, err := h.serviceKit.ChallengeService.FindTestcaseByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entities.HttpError{Message: err.Error()})
	}

	err = h.serviceKit.ChallengeService.DeleteTestcase(testcase)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entities.HttpError{Message: err.Error()})
	}

	return c.SendStatus(fiber.StatusOK)
}

func (h *challengeTestcaseHandler) GetTestcaseByID(c *fiber.Ctx) error {
	id := ParseIntParam(c, "id")

	testcase, err := h.serviceKit.ChallengeService.FindTestcaseByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entities.HttpError{Message: err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(testcase)
}

func NewChallengeTestcaseHandler(serviceKit *services.ServiceKit) *challengeTestcaseHandler {
	return &challengeTestcaseHandler{
		serviceKit: serviceKit,
	}
}
