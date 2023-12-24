package repositories

import (
	"github.com/wuttinanhi/code-judge-system/entities"
	"gorm.io/gorm"
)

type ChallengeRepository interface {
	// CreateChallenge creates a new challenge.
	CreateChallenge(challenge *entities.Challenge) (*entities.Challenge, error)
	// CreateChallengeWithTestcase creates a new challenge with testcases.
	CreateChallengeWithTestcase(challenge *entities.Challenge, testcases []entities.ChallengeTestcase) (*entities.Challenge, error)
	// UpdateChallenge updates a challenge.
	UpdateChallenge(challenge *entities.Challenge) error
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
	PaginationChallengesWithStatus(page int, limit int, user *entities.User) (result *entities.PaginationResult[*entities.ChallengeExtended], err error)
}

type challengeRepository struct {
	db *gorm.DB
}

// PaginationChallengesWithStatus implements ChallengeRepository.
func (r *challengeRepository) PaginationChallengesWithStatus(page int, limit int, user *entities.User) (result *entities.PaginationResult[*entities.ChallengeExtended], err error) {
	result = &entities.PaginationResult[*entities.ChallengeExtended]{
		Items: make([]*entities.ChallengeExtended, 0),
		Total: 0,
	}

	var storeVaule []*entities.ChallengeExtended
	err = r.db.Raw(`
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
ORDER BY t.id ASC
LIMIT ?
OFFSET ?
`,
		user.ID,
		limit,
		(page-1)*limit,
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
func (r *challengeRepository) CreateChallengeWithTestcase(challenge *entities.Challenge, testcases []entities.ChallengeTestcase) (*entities.Challenge, error) {
	challenge.Testcases = testcases
	result := r.db.Create(challenge)
	return challenge, result.Error
}

// AddTestcase implements ChallengeRepository.
func (r *challengeRepository) AddTestcase(challenge *entities.Challenge, testcase *entities.ChallengeTestcase) (*entities.ChallengeTestcase, error) {
	err := r.db.Model(challenge).Association("Testcases").Append(testcase)
	return testcase, err
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

// AllTestcases implements ChallengeRepository.
func (r *challengeRepository) AllTestcases(challenge *entities.Challenge) (testcases []*entities.ChallengeTestcase, err error) {
	err = r.db.Model(challenge).Association("Testcases").Find(&testcases)
	return testcases, err
}

// DeleteChallenge implements ChallengeRepository.
func (r *challengeRepository) DeleteChallenge(challenge *entities.Challenge) error {
	result := r.db.Delete(challenge)
	return result.Error
}

// FindChallengeByID implements ChallengeRepository.
func (r *challengeRepository) FindChallengeByID(id uint) (challenge *entities.Challenge, err error) {
	result := r.db.First(&challenge, id).Preload("Testcases")
	return challenge, result.Error
}

// FindChallengesByAuthor implements ChallengeRepository.
func (r *challengeRepository) FindChallengesByAuthor(author *entities.User) (challenges []*entities.Challenge, err error) {
	result := r.db.Where("author_id = ?", author.ID).Find(&challenges)
	return challenges, result.Error
}

// FindTestcaseByID implements ChallengeRepository.
func (r *challengeRepository) FindTestcaseByID(id uint) (testcase *entities.ChallengeTestcase, err error) {
	result := r.db.First(&testcase, id)
	return testcase, result.Error
}

// DeleteTestcase implements ChallengeRepository.
func (r *challengeRepository) DeleteTestcase(testcase *entities.ChallengeTestcase) error {
	result := r.db.Delete(testcase)
	return result.Error
}

// UpdateChallenge implements ChallengeRepository.
func (r *challengeRepository) UpdateChallenge(challenge *entities.Challenge) error {
	result := r.db.Save(challenge)
	return result.Error
}

// UpdateTestcase implements ChallengeRepository.
func (r *challengeRepository) UpdateTestcase(testcase *entities.ChallengeTestcase) error {
	result := r.db.Save(testcase)
	return result.Error
}

func NewChallengeRepository(db *gorm.DB) ChallengeRepository {
	return &challengeRepository{db}
}
