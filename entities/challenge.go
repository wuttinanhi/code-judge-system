package entities

import "github.com/gofiber/fiber/v2"

type Challenge struct {
	ChallengeID uint                `json:"challenge_id" gorm:"primaryKey"`
	Name        string              `json:"name"`
	Description string              `json:"description"`
	UserID      uint                `json:"user_id"`
	User        User                `json:"user" gorm:"foreignKey:UserID"`
	Testcases   []ChallengeTestcase `json:"testcases" gorm:"foreignKey:ChallengeID"`
}

type ChallengeCreateDTO struct {
	Name        string `json:"name" validate:"required,min=3,max=255"`
	Description string `json:"description" validate:"max=255"`
}

func ValidateChallengeCreateDTO(c *fiber.Ctx) ChallengeCreateDTO {
	var dto ChallengeCreateDTO

	if err := c.BodyParser(&dto); err != nil {
		panic(err)
	}

	if err := validate.Struct(&dto); err != nil {
		panic(err)
	}

	return dto
}

type ChallengeUpdateDTO struct {
	ChallengeID uint   `json:"challenge_id" validate:"required"`
	Name        string `json:"name" validate:"required,min=3,max=255"`
	Description string `json:"description" validate:"max=255"`
}

func ValidateChallengeUpdateDTO(c *fiber.Ctx) ChallengeUpdateDTO {
	var dto ChallengeUpdateDTO

	if err := c.BodyParser(&dto); err != nil {
		panic(err)
	}

	if err := validate.Struct(&dto); err != nil {
		panic(err)
	}

	return dto
}

type ChallengeCreateResponse struct {
	ChallengeID uint   `json:"challenge_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ChallengeCreateWithTestcaseDTO struct {
	ChallengeCreateDTO
	Testcases []ChallengeTestcaseCreateDTO `json:"testcases" validate:"required"`
}

func ValidateChallengeCreateWithTestcaseDTO(c *fiber.Ctx) ChallengeCreateWithTestcaseDTO {
	var dto ChallengeCreateWithTestcaseDTO

	if err := c.BodyParser(&dto); err != nil {
		panic(err)
	}

	if err := validate.Struct(&dto); err != nil {
		panic(err)
	}

	return dto
}
