package repositories

import (
	"fmt"

	"github.com/wuttinanhi/code-judge-system/entities"
	"gorm.io/gorm"
)

type ChallengeRepository interface {
	// CreateChallenge creates a new challenge.
	CreateChallenge(challenge *entities.Challenge) (*entities.Challenge, error)
	// CreateChallengeWithTestcase creates a new challenge with testcases.
	CreateChallengeWithTestcase(challenge *entities.Challenge, testcases []*entities.ChallengeTestcase) (*entities.Challenge, error)
	// UpdateChallenge updates a challenge.
	// UpdateChallenge(challenge *entities.Challenge) error
	// DeleteChallenge deletes a challenge.
	DeleteChallenge(challenge *entities.Challenge) error
	// FindChallengeByID returns a challenge by given ID.
	FindChallengeByID(id uint) (challenge *entities.Challenge, err error)
	// FindChallengeByAuthor returns a challenge by given author.
	FindChallengesByAuthor(author *entities.User) (challenges []*entities.Challenge, err error)
	// AllChallenges returns all challenges.
	AllChallenges() (challenges []*entities.Challenge, err error)
	// AddTestcase adds a testcase to a challenge.
	AddTestcase(challenge *entities.Challenge, testcase *entities.ChallengeTestcase) (*entities.ChallengeTestcase, error)
	// UpdateTestcase updates a testcase.
	UpdateTestcase(testcase *entities.ChallengeTestcase) error
	// DeleteTestcase removes a testcase from a challenge.
	DeleteTestcase(testcase *entities.ChallengeTestcase) error
	// AllTestcases returns all testcases of a challenge.
	AllTestcases(challenge *entities.Challenge) (testcases []*entities.ChallengeTestcase, err error)
	// FindTestcaseByID returns a testcase by given ID.
	FindTestcaseByID(id uint) (testcase *entities.ChallengeTestcase, err error)
	// PaginationChallengesWithStatus returns all challenges with status.
	PaginationChallengesWithStatus(options *entities.ChallengePaginationOptions) (result *entities.PaginationResult[*entities.ChallengeExtended], err error)
	// UpdateChallengeWithTestcase updates a challenge with testcases.
	UpdateChallengeWithTestcase(challenge *entities.Challenge) error
}

type challengeRepository struct {
	db *gorm.DB
}

// UpdateChallengeWithTestcase implements ChallengeRepository.
func (r *challengeRepository) UpdateChallengeWithTestcase(challenge *entities.Challenge) error {
	err := r.db.Transaction(func(tx *gorm.DB) (err error) {
		// Load the Testcases
		oldTestcases := []*entities.ChallengeTestcase{}
		tx.Model(challenge).Association("Testcases").Find(&oldTestcases)

		for _, testcase := range oldTestcases {
			err = tx.Model(testcase).Delete(testcase, testcase.ID).Error
			if err != nil {
				return err
			}
		}

		err = tx.Session(&gorm.Session{FullSaveAssociations: true}).Save(challenge).Error

		return err
	})

	return err
}

// PaginationChallengesWithStatus implements ChallengeRepository.
func (r *challengeRepository) PaginationChallengesWithStatus(options *entities.ChallengePaginationOptions) (result *entities.PaginationResult[*entities.ChallengeExtended], err error) {
	result = &entities.PaginationResult[*entities.ChallengeExtended]{
		Items: make([]*entities.ChallengeExtended, 0),
		Total: 0,
	}

	challengeQuery := fmt.Sprintf(`
SELECT t.*, COALESCE(subq.submission_status, "NOTSOLVE") AS submission_status, (SELECT COUNT(*) FROM challenges) AS total_challenges
FROM challenges AS t
LEFT JOIN (
	SELECT
		challenge_id,
		MAX(id),
		MAX(status) as submission_status
	FROM submissions
	WHERE user_id = ?
	GROUP BY challenge_id
) AS subq ON t.id = subq.challenge_id
ORDER BY ? %s
LIMIT ?
OFFSET ?
	`, options.Order)

	offset := (options.Page - 1) * options.Limit

	var storeVaule []*entities.ChallengeExtended
	err = r.db.Raw(challengeQuery,
		options.User.ID,
		options.Sort,
		options.Limit,
		offset,
	).
		Scan(&storeVaule).Error

	// Query to count total challenges
	var totalChallenges int64
	r.db.Model(&entities.Challenge{}).Count(&totalChallenges)

	result.Items = storeVaule
	result.Total = int(totalChallenges)

	return
}

// CreateChallengeWithTestcase implements ChallengeRepository.
func (r *challengeRepository) CreateChallengeWithTestcase(challenge *entities.Challenge, testcases []*entities.ChallengeTestcase) (*entities.Challenge, error) {
	challenge.Testcases = testcases
	result := r.db.Create(challenge)
	return challenge, result.Error
}

// CreateChallenge implements ChallengeRepository.
func (r *challengeRepository) CreateChallenge(challenge *entities.Challenge) (*entities.Challenge, error) {
	result := r.db.Create(challenge)
	return challenge, result.Error
}

// AllChallenges implements ChallengeRepository.
func (r *challengeRepository) AllChallenges() (challenges []*entities.Challenge, err error) {
	result := r.db.Find(&challenges)
	return challenges, result.Error
}

// DeleteChallenge implements ChallengeRepository.
func (r *challengeRepository) DeleteChallenge(challenge *entities.Challenge) error {
	result := r.db.Delete(challenge)
	return result.Error
}

// FindChallengeByID implements ChallengeRepository.
func (r *challengeRepository) FindChallengeByID(id uint) (challenge *entities.Challenge, err error) {
	result := r.db.Preload("Testcases").First(&challenge, id)
	return challenge, result.Error
}

// FindChallengesByAuthor implements ChallengeRepository.
func (r *challengeRepository) FindChallengesByAuthor(author *entities.User) (challenges []*entities.Challenge, err error) {
	result := r.db.Where("author_id = ?", author.ID).Find(&challenges)
	return challenges, result.Error
}

// UpdateChallenge implements ChallengeRepository.
// func (r *challengeRepository) UpdateChallenge(challenge *entities.Challenge) error {
// 	result := r.db.Save(challenge)
// 	return result.Error
// }

func NewChallengeRepository(db *gorm.DB) ChallengeRepository {
	return &challengeRepository{db}
}
