package entities

import "github.com/gofiber/fiber/v2"

const (
	SubmissionStatusPending  = "PENDING"
	SubmissionStatusCorrect  = "CORRECT"
	SubmissionStatusWrong    = "WRONG"
	SubmissionStatusNotSolve = "NOTSOLVE"
)

type Submission struct {
	ID                  uint                  `json:"submission_id" gorm:"primaryKey"`
	Language            string                `json:"language"`
	SourceCode          string                `json:"source_code"`
	Status              string                `json:"status" gorm:"default:PENDING"`
	UserID              uint                  `json:"user_id"`
	User                *User                 `json:"user"`
	ChallengeID         uint                  `json:"challenge_id"`
	Challenge           *Challenge            `json:"challenge"`
	SubmissionTestcases []*SubmissionTestcase `json:"submission_testcases" gorm:"foreignKey:SubmissionID"`
}

type SubmissionCreateDTO struct {
	ChallengeID uint   `json:"challenge_id" validate:"required"`
	Language    string `json:"language" validate:"required"`
	SourceCode  string `json:"source_code" validate:"required"`
}

func ValidateSubmissionCreateDTO(c *fiber.Ctx) SubmissionCreateDTO {
	var dto SubmissionCreateDTO

	if err := c.BodyParser(&dto); err != nil {
		panic(err)
	}

	if err := validate.Struct(&dto); err != nil {
		panic(err)
	}

	return dto
}

func (s *Submission) IsCorrect() bool {
	for _, testcase := range s.SubmissionTestcases {
		if testcase.Status == SubmissionStatusWrong || testcase.Status == SubmissionStatusNotSolve {
			return false
		}
	}

	return true
}
