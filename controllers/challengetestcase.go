package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wuttinanhi/code-judge-system/entities"
	"github.com/wuttinanhi/code-judge-system/services"
)

func CreateTestcase(c *fiber.Ctx) error {
	dto := entities.ValidateChallengeTestcaseCreateDTO(c)

	challenge, err := services.GetServiceKit().ChallengeService.FindChallengeByID(dto.ChallengeID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entities.HttpError{Message: err.Error()})
	}

	testcase, err := services.GetServiceKit().ChallengeService.AddTestcase(challenge, &entities.ChallengeTestcase{
		Input:          dto.Input,
		ExpectedOutput: dto.ExpectedOutput,
		ChallengeID:    challenge.ChallengeID,
	})
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entities.HttpError{Message: err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(testcase)
}

func UpdateTestcase(c *fiber.Ctx) error {
	dto := entities.ValidateChallengeTestcaseUpdateDTO(c)

	testcase, err := services.GetServiceKit().ChallengeService.FindTestcaseByID(dto.TestcaseID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entities.HttpError{Message: err.Error()})
	}

	testcase.Input = dto.Input
	testcase.ExpectedOutput = dto.ExpectedOutput

	err = services.GetServiceKit().ChallengeService.UpdateTestcase(testcase)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entities.HttpError{Message: err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(testcase)
}

func DeleteTestcase(c *fiber.Ctx) error {
	id := ParseIntParam(c, "id")

	testcase, err := services.GetServiceKit().ChallengeService.FindTestcaseByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entities.HttpError{Message: err.Error()})
	}

	err = services.GetServiceKit().ChallengeService.DeleteTestcase(testcase)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entities.HttpError{Message: err.Error()})
	}

	return c.SendStatus(fiber.StatusOK)
}
