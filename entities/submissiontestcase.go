package entities

type SubmissionTestcase struct {
	ID                  uint               `json:"submission_testcase_id" gorm:"primaryKey"`
	Status              string             `json:"status" gorm:"default:PENDING"`
	Output              string             `json:"output"`
	SubmissionID        uint               `json:"submission_id"`
	Submission          *Submission        `json:"submission"`
	ChallengeTestcaseID uint               `json:"challenge_testcase_id"`
	ChallengeTestcase   *ChallengeTestcase `json:"challenge_testcase" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Note                string             `json:"note"`
}
