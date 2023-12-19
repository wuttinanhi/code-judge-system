package entities

type Challenge struct {
	ChallengeID uint                `json:"challenge_id" gorm:"primaryKey"`
	Name        string              `json:"name"`
	Description string              `json:"description"`
	UserID      uint                `json:"user_id"`
	User        User                `json:"user" gorm:"foreignKey:UserID"`
	Testcases   []ChallengeTestcase `json:"testcases" gorm:"foreignKey:ChallengeID"`
}
