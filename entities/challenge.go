package entities

import "github.com/gofiber/fiber/v2"

type Challenge struct {
	ID          uint                 `json:"challenge_id" gorm:"primaryKey"`
	Name        string               `json:"name"`
	Description string               `json:"description"`
	UserID      uint                 `json:"user_id"`
	User        *User                `json:"user" gorm:"foreignKey:UserID"`
	Testcases   []*ChallengeTestcase `json:"testcases" gorm:"foreignKey:ChallengeID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE"`
	Submission  []*Submission        `json:"submission" gorm:"foreignKey:ChallengeID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE"`
}

type ChallengeExtended struct {
	Challenge
	SubmissionStatus string
}

type ChallengePaginationOptions struct {
	PaginationOptions
	User *User
}

type ChallengeCreateResponse struct {
	ChallengeID uint   `json:"challenge_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ChallengeCreateWithTestcaseDTO struct {
	Name        string                 `json:"name" validate:"required,min=3,max=255"`
	Description string                 `json:"description" validate:"max=255"`
	Testcases   []ChallengeTestcaseDTO `json:"testcases" validate:"required"`
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

type ChallengeUpdateDTO struct {
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
