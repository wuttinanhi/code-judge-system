package entities

import "github.com/gofiber/fiber/v2"

type ChallengeTestcase struct {
	TestcaseID     uint      `json:"testcase_id" gorm:"primaryKey"`
	Input          string    `json:"input"`
	ExpectedOutput string    `json:"expected_output"`
	ChallengeID    uint      `json:"challenge_id"`
	Challenge      Challenge `json:"challenge" gorm:"foreignKey:ChallengeID"`
}

type ChallengeTestcaseCreateDTO struct {
	Input          string `json:"input" validate:"required,max=1024"`
	ExpectedOutput string `json:"expected_output" validate:"required,max=1024"`
	ChallengeID    uint   `json:"challenge_id" validate:"required"`
}

func ValidateChallengeTestcaseCreateDTO(c *fiber.Ctx) ChallengeTestcaseCreateDTO {
	var dto ChallengeTestcaseCreateDTO

	if err := c.BodyParser(&dto); err != nil {
		panic(err)
	}

	if err := validate.Struct(&dto); err != nil {
		panic(err)
	}

	return dto
}

type ChallengeTestcaseUpdateDTO struct {
	TestcaseID     uint   `json:"testcase_id" validate:"required"`
	Input          string `json:"input" validate:"required,max=1024"`
	ExpectedOutput string `json:"expected_output" validate:"required,max=1024"`
}

func ValidateChallengeTestcaseUpdateDTO(c *fiber.Ctx) ChallengeTestcaseUpdateDTO {
	var dto ChallengeTestcaseUpdateDTO

	if err := c.BodyParser(&dto); err != nil {
		panic(err)
	}

	if err := validate.Struct(&dto); err != nil {
		panic(err)
	}

	return dto
}

type ChallengeTestcaseCreateResponse struct {
	TestcaseID     uint   `json:"testcase_id"`
	Input          string `json:"input"`
	ExpectedOutput string `json:"expected_output"`
	ChallengeID    uint   `json:"challenge_id"`
}
