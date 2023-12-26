package entities

import "github.com/gofiber/fiber/v2"

type ChallengeTestcase struct {
	ID                  uint                  `json:"testcase_id" gorm:"primaryKey"`
	Input               string                `json:"input"`
	ExpectedOutput      string                `json:"expected_output"`
	LimitMemory         uint                  `json:"limit_memory"`
	LimitTimeMs         uint                  `json:"limit_time_ms"`
	SubmissionTestcases []*SubmissionTestcase `json:"submission_testcases"`
	ChallengeID         uint                  `json:"challenge_id"`
	Challenge           *Challenge            `json:"challenge"`
}

type ChallengeTestcaseCreateDTO struct {
	Input          string `json:"input" validate:"required,max=1024"`
	ExpectedOutput string `json:"expected_output" validate:"required,max=1024"`
	ChallengeID    uint   `json:"challenge_id" validate:"required"`
	LimitMemory    uint   `json:"limit_memory" validate:"required"`
	LimitTimeMs    uint   `json:"limit_time_ms" validate:"required"`
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

type ChallengeTestcaseDTO struct {
	Input          string `json:"input" validate:"required,max=1024"`
	ExpectedOutput string `json:"expected_output" validate:"required,max=1024"`
	LimitMemory    uint   `json:"limit_memory" validate:"required"`
	LimitTimeMs    uint   `json:"limit_time_ms" validate:"required"`
}

func (t *ChallengeTestcaseDTO) ToTestcase() *ChallengeTestcase {
	return &ChallengeTestcase{
		Input:          t.Input,
		ExpectedOutput: t.ExpectedOutput,
		LimitMemory:    t.LimitMemory,
		LimitTimeMs:    t.LimitTimeMs,
	}
}
