package controllers

import (
	"net/http"
	"strings"

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

	// only user with role admin or staff can create challenge
	if user.Role != entities.UserRoleAdmin && user.Role != entities.UserRoleStaff {
		return c.SendStatus(fiber.StatusForbidden)
	}

	// limit challenges created by user to 100
	total, err := h.serviceKit.ChallengeService.CountAllChallengesByUser(user)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(entities.HttpError{Message: err.Error()})
	}
	if total >= 100 {
		return c.Status(fiber.StatusTooManyRequests).
			JSON(entities.HttpError{Message: "You have reached the limit of 100 challenges"})
	}

	// create challenge
	challenge, err := h.serviceKit.ChallengeService.CreateChallenge(&entities.Challenge{
		Name:        dto.Name,
		Description: dto.Description,
		UserID:      user.ID,
		Testcases:   dto.GetTestcases(),
	})
	if err != nil {
		if strings.Contains(err.Error(), "testcase #") {
			return c.Status(http.StatusBadRequest).JSON(CreateValidationError(err))
		}
		return c.Status(http.StatusBadRequest).JSON(entities.HttpError{Message: err.Error()})
	}

	return c.Status(http.StatusOK).JSON(challenge)
}

func (h *challengeHandler) UpdateChallenge(c *fiber.Ctx) error {
	user := GetUserFromRequest(c)
	dto := entities.ValidateChallengeUpdateDTO(c)
	id := ParseIntParam(c, "id")

	// only user with role admin or staff can update challenge
	if user.Role != entities.UserRoleAdmin && user.Role != entities.UserRoleStaff {
		return c.SendStatus(fiber.StatusForbidden)
	}

	// load challenge
	challenge, err := h.serviceKit.ChallengeService.FindChallengeByID(uint(id))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(entities.HttpError{Message: err.Error()})
	}

	// update challenge
	challenge.Name = dto.Name
	challenge.Description = dto.Description
	challenge.Testcases = dto.GetTestcases()
	err = h.serviceKit.ChallengeService.UpdateChallengeWithTestcase(challenge)
	if err != nil {
		if strings.Contains(err.Error(), "testcase #") {
			return c.Status(http.StatusBadRequest).JSON(CreateValidationError(err))
		}
		return c.Status(http.StatusBadRequest).JSON(entities.HttpError{Message: err.Error()})
	}

	return c.Status(http.StatusOK).JSON(challenge)
}

func (h *challengeHandler) DeleteChallenge(c *fiber.Ctx) error {
	user := GetUserFromRequest(c)
	id := ParseIntParam(c, "id")

	// only user with role admin or staff can delete challenge
	if user.Role != entities.UserRoleAdmin && user.Role != entities.UserRoleStaff {
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
