package controllers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/wuttinanhi/code-judge-system/entities"
	"github.com/wuttinanhi/code-judge-system/services"
)

func CreateChallenge(c *fiber.Ctx) error {
	user := entities.GetUserFromRequest(c)
	dto := entities.ValidateChallengeCreateDTO(c)

	challenge, err := services.GetServiceKit().ChallengeService.CreateChallenge(&entities.Challenge{
		Name:        dto.Name,
		Description: dto.Description,
		UserID:      user.UserID,
	})
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(entities.HttpError{Message: err.Error()})
	}

	return c.Status(http.StatusOK).JSON(challenge)
}

func CreateChallengeWithTestcase(c *fiber.Ctx) error {
	user := entities.GetUserFromRequest(c)
	dto := entities.ValidateChallengeCreateWithTestcaseDTO(c)

	testcases := make([]entities.ChallengeTestcase, len(dto.Testcases))
	for i, testcase := range dto.Testcases {
		testcases[i] = entities.ChallengeTestcase{
			Input:          testcase.Input,
			ExpectedOutput: testcase.ExpectedOutput,
		}
	}

	challenge, err := services.GetServiceKit().ChallengeService.CreateChallenge(&entities.Challenge{
		Name:        dto.Name,
		Description: dto.Description,
		UserID:      user.UserID,
		Testcases:   testcases,
	})
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(entities.HttpError{Message: err.Error()})
	}

	return c.Status(http.StatusOK).JSON(challenge)
}

func UpdateChallenge(c *fiber.Ctx) error {
	dto := entities.ValidateChallengeUpdateDTO(c)

	challenge, err := services.GetServiceKit().ChallengeService.FindChallengeByID(dto.ChallengeID)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(entities.HttpError{Message: err.Error()})
	}

	err = services.GetServiceKit().ChallengeService.UpdateChallenge(&entities.Challenge{
		ChallengeID: challenge.ChallengeID,
		Name:        dto.Name,
		Description: dto.Description,
		Testcases:   challenge.Testcases,
	})
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(entities.HttpError{Message: err.Error()})
	}

	return c.Status(http.StatusOK).JSON(challenge)
}

func DeleteChallenge(c *fiber.Ctx) error {
	id := ParseIntParam(c, "id")

	challenge, err := services.GetServiceKit().ChallengeService.FindChallengeByID(uint(id))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(entities.HttpError{Message: err.Error()})
	}

	err = services.GetServiceKit().ChallengeService.DeleteChallenge(challenge)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(entities.HttpError{Message: err.Error()})
	}

	return c.SendStatus(http.StatusOK)
}

func GetChallengeByID(c *fiber.Ctx) error {
	id := ParseIntParam(c, "id")

	challenges, err := services.GetServiceKit().ChallengeService.FindChallengeByID(uint(id))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(entities.HttpError{Message: err.Error()})
	}

	return c.Status(http.StatusOK).JSON(challenges)
}

func GetAllChallenges(c *fiber.Ctx) error {
	challenges, err := services.GetServiceKit().ChallengeService.AllChallenges()
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(entities.HttpError{Message: err.Error()})
	}

	return c.Status(http.StatusOK).JSON(challenges)
}
