package entities

const (
	SubmissionStatusPending = "PENDING"
	SubmissionStatusCorrect = "CORRECT"
	SubmissionStatusWrong   = "WRONG"
)

type Submission struct {
	SubmissionID uint      `json:"submission_id" gorm:"primaryKey"`
	ChallengeID  uint      `json:"challenge_id"`
	Challenge    Challenge `json:"challenge" gorm:"foreignKey:ChallengeID"`
	UserID       uint      `json:"user_id"`
	User         User      `json:"user" gorm:"foreignKey:UserID"`
	Language     string    `json:"language"`
	SourceCode   string    `json:"source_code"`
	Status       string    `json:"status" gorm:"default:PENDING"`
}
