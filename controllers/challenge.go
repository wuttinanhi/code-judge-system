package controllers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/wuttinanhi/code-judge-system/entities"
	"github.com/wuttinanhi/code-judge-system/services"
)

type challengeHandler struct {
	serviceKit *services.ServiceKit
}

// func (h *challengeHandler) CreateChallenge(c *fiber.Ctx) error {
// 	user := GetUserFromRequest(c)
// 	dto := entities.ValidateChallengeCreateDTO(c)

// 	// only user with role admin can create challenge
// 	if user.Role != entities.UserRoleAdmin {
// 		return c.SendStatus(fiber.StatusForbidden)
// 	}

// 	challenge, err := h.serviceKit.ChallengeService.CreateChallenge(&entities.Challenge{
// 		Name:        dto.Name,
// 		Description: dto.Description,
// 		UserID:      user.ID,
// 	})
// 	if err != nil {
// 		return c.Status(http.StatusBadRequest).JSON(entities.HttpError{Message: err.Error()})
// 	}

// 	return c.Status(http.StatusOK).JSON(challenge)
// }

func (h *challengeHandler) CreateChallengeWithTestcase(c *fiber.Ctx) error {
	user := GetUserFromRequest(c)
	dto := entities.ValidateChallengeCreateWithTestcaseDTO(c)

	// only user with role admin can create challenge
	if user.Role != entities.UserRoleAdmin {
		return c.SendStatus(fiber.StatusForbidden)
	}

	testcases := make([]*entities.ChallengeTestcase, len(dto.Testcases))
	for i, dtotestcase := range dto.Testcases {
		testcases[i] = &entities.ChallengeTestcase{
			Input:          dtotestcase.Input,
			ExpectedOutput: dtotestcase.ExpectedOutput,
			LimitMemory:    dtotestcase.LimitMemory,
			LimitTimeMs:    dtotestcase.LimitTimeMs,
		}
	}

	challenge, err := h.serviceKit.ChallengeService.CreateChallenge(&entities.Challenge{
		Name:        dto.Name,
		Description: dto.Description,
		UserID:      user.ID,
		Testcases:   testcases,
	})
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(entities.HttpError{Message: err.Error()})
	}

	return c.Status(http.StatusOK).JSON(challenge)
}

func (h *challengeHandler) UpdateChallenge(c *fiber.Ctx) error {
	user := GetUserFromRequest(c)
	dto := entities.ValidateChallengeUpdateDTO(c)
	id := ParseIntParam(c, "id")

	// only user with role admin can update challenge
	if user.Role != entities.UserRoleAdmin {
		return c.SendStatus(fiber.StatusForbidden)
	}

	// load challenge
	challenge, err := h.serviceKit.ChallengeService.FindChallengeByID(uint(id))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(entities.HttpError{Message: err.Error()})
	}

	// update challenge
	err = h.serviceKit.ChallengeService.UpdateChallenge(&entities.Challenge{
		ID:          challenge.ID,
		Name:        dto.Name,
		Description: dto.Description,
		Testcases:   challenge.Testcases,
	})
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(entities.HttpError{Message: err.Error()})
	}

	return c.Status(http.StatusOK).JSON(challenge)
}

func (h *challengeHandler) DeleteChallenge(c *fiber.Ctx) error {
	user := GetUserFromRequest(c)
	id := ParseIntParam(c, "id")

	if user.Role != entities.UserRoleAdmin {
		return c.SendStatus(fiber.StatusForbidden)
	}

	challenge, err := h.serviceKit.ChallengeService.FindChallengeByID(uint(id))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(entities.HttpError{Message: err.Error()})
	}

	err = h.serviceKit.ChallengeService.DeleteChallenge(challenge)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(entities.HttpError{Message: err.Error()})
	}

	return c.SendStatus(http.StatusOK)
}

func (h *challengeHandler) GetChallengeByID(c *fiber.Ctx) error {
	id := ParseIntParam(c, "id")

	challenges, err := h.serviceKit.ChallengeService.FindChallengeByID(uint(id))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(entities.HttpError{Message: err.Error()})
	}

	return c.Status(http.StatusOK).JSON(challenges)
}

func (h *challengeHandler) GetAllChallenges(c *fiber.Ctx) error {
	challenges, err := h.serviceKit.ChallengeService.AllChallenges()
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(entities.HttpError{Message: err.Error()})
	}

	return c.Status(http.StatusOK).JSON(challenges)
}

func (h *challengeHandler) PaginationChallengesWithStatus(c *fiber.Ctx) error {
	user := GetUserFromRequest(c)
	options := ParsePaginationOptions(c)

	challenges, err := h.serviceKit.ChallengeService.PaginationChallengesWithStatus(&entities.ChallengePaginationOptions{
		PaginationOptions: *options,
		User:              user,
	})
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(entities.HttpError{Message: err.Error()})
	}

	return c.Status(http.StatusOK).JSON(challenges)
}

func NewChallengeHandler(serviceKit *services.ServiceKit) *challengeHandler {
	return &challengeHandler{
		serviceKit: serviceKit,
	}
}
