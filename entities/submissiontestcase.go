package entities

type SubmissionTestcase struct {
	SubmissionTestcaseID uint              `json:"submission_testcase_id" gorm:"primaryKey"`
	SubmissionID         uint              `json:"submission_id"`
	Submission           Submission        `json:"submission" gorm:"foreignKey:SubmissionID"`
	ChallengeTestcaseID  uint              `json:"challenge_testcase_id"`
	ChallengeTestcase    ChallengeTestcase `json:"challenge_testcase" gorm:"foreignKey:ChallengeTestcaseID"`
	Status               string            `json:"status" gorm:"default:PENDING"`
	Output               string            `json:"output"`
}
