package entities

type ChallengeTestcase struct {
	TestcaseID     uint      `json:"testcase_id" gorm:"primaryKey"`
	Input          string    `json:"input"`
	ExpectedOutput string    `json:"expected_output"`
	ChallengeID    uint      `json:"challenge_id"`
	Challenge      Challenge `json:"challenge" gorm:"foreignKey:ChallengeID"`
}
