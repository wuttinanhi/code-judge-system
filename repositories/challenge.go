package repositories

import (
	"fmt"
	"strings"
	"time"

	"github.com/wuttinanhi/code-judge-system/entities"
	"gorm.io/gorm"
)

type ChallengeRepository interface {
	// CreateChallenge creates a new challenge.
	CreateChallenge(challenge *entities.Challenge) (*entities.Challenge, error)
	// CreateChallengeWithTestcase creates a new challenge with testcases.
	CreateChallengeWithTestcase(challenge *entities.Challenge, testcases []*entities.ChallengeTestcase) (*entities.Challenge, error)
	// UpdateChallenge updates a challenge.
	UpdateChallenge(challenge *entities.Challenge) error
	// DeleteChallenge deletes a challenge.
	DeleteChallenge(challenge *entities.Challenge) error
	// FindChallengeByID returns a challenge by given ID.
	FindChallengeByID(id uint) (challenge *entities.Challenge, err error)
	// FindChallengeByAuthor returns a challenge by given author.
	// FindChallengesByAuthor(author *entities.User) (challenges []*entities.Challenge, err error)
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
	// CountAllChallengesByUser returns total challenges by given user.
	CountAllChallengesByUser(user *entities.User) (total int64, err error)
}

type challengeRepository struct {
	db *gorm.DB
}

// CountAllChallengesByUser implements ChallengeRepository.
func (r *challengeRepository) CountAllChallengesByUser(user *entities.User) (total int64, err error) {
	result := r.db.Model(&entities.Challenge{}).
		Where(&entities.Challenge{UserID: user.ID}).
		Count(&total)
	return total, result.Error
}

// UpdateChallengeWithTestcase implements ChallengeRepository.
func (r *challengeRepository) UpdateChallengeWithTestcase(challenge *entities.Challenge) error {
	err := r.db.Transaction(func(tx *gorm.DB) (err error) {
		// loop new challenge testcase
		for _, testcase := range challenge.Testcases {
			// if testcase.ActionFlag is "create" then create new testcase
			if testcase.ActionFlag == "create" {
				testcase.ID = 0
				testcase.ChallengeID = challenge.ID
				err = tx.Model(&entities.ChallengeTestcase{}).Create(testcase).Error
				if err != nil {
					return err
				}
			}
			// if testcase.ActionFlag is "update" then update testcase
			if testcase.ActionFlag == "update" {
				testcase.ChallengeID = challenge.ID
				err = tx.Save(testcase).Error
				if err != nil {
					return err
				}
			}
			// if testcase.ActionFlag is "delete" then delete testcase
			if testcase.ActionFlag == "delete" {
				err = tx.Delete(testcase).Error
				if err != nil {
					return err
				}
			}
		}

		err = tx.Session(&gorm.Session{FullSaveAssociations: false}).Omit("Testcases").Save(challenge).Error
		if err != nil {
			return err
		}

		// limit testcases to 100 per challenge
		var totalTestcases int64
		err = tx.
			Model(&entities.ChallengeTestcase{}).
			Where(&entities.ChallengeTestcase{ChallengeID: challenge.ID}).
			Count(&totalTestcases).Error
		if err != nil {
			return err
		}
		if totalTestcases > 100 {
			return fmt.Errorf("testcases limit exceeded")
		}

		return nil
	})
	if err != nil {
		return err
	}

	cleanActionFlag(challenge)

	return err
}

// PaginationChallengesWithStatus implements ChallengeRepository.
func (r *challengeRepository) PaginationChallengesWithStatus(options *entities.ChallengePaginationOptions) (result *entities.PaginationResult[*entities.ChallengeExtended], err error) {
	result = &entities.PaginationResult[*entities.ChallengeExtended]{
		Items: make([]*entities.ChallengeExtended, 0),
		Total: 0,
	}

	// convert options.Order to uppercase
	options.Order = strings.ToUpper(options.Order)

	// if options.Order is not ASC or DESC then throw error
	if options.Order != "ASC" && options.Order != "DESC" {
		err = fmt.Errorf("invalid order option")
		return
	}

	// calculate offset
	offset := (options.Page - 1) * options.Limit

	challengeQuery := fmt.Sprintf(`
SELECT 
	t.id as ORDER_ID,
	t.*, 
	u.*, 
	COALESCE(subq.submission_status, "NOTSOLVE") AS submission_status
FROM challenges AS t
LEFT JOIN users AS u ON t.user_id = u.id
LEFT JOIN (
	SELECT
		MAX(id) as id,
		challenge_id,
		MAX(status) as submission_status
	FROM submissions
	WHERE user_id = ?
	GROUP BY challenge_id
) subq ON t.id = subq.challenge_id
WHERE 
	t.name LIKE ? OR 
	t.description LIKE ? OR
	u.display_name LIKE ?
ORDER BY ORDER_ID %s
LIMIT ?
OFFSET ?
	`, options.Order)

	var storeVaule []*entities.ChallengeExtended
	err = r.db.Raw(
		challengeQuery,
		options.User.ID,
		"%"+options.Search+"%",
		"%"+options.Search+"%",
		"%"+options.Search+"%",
		options.Limit,
		offset,
	).
		Scan(&storeVaule).Error

	// omit user password
	for _, challenge := range storeVaule {
		challenge.User.Password = ""
		challenge.User.CreatedAt = time.Time{}
	}

	// query to count total challenges
	var totalChallenges int64
	r.db.Model(&entities.Challenge{}).
		Preload("User").
		Joins("LEFT JOIN users ON challenges.user_id = users.id").
		Where(
			"name LIKE ? OR description LIKE ? OR users.display_name LIKE ?",
			"%"+options.Search+"%",
			"%"+options.Search+"%",
			"%"+options.Search+"%",
		).
		Count(&totalChallenges)

	// store result
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
	cleanActionFlag(challenge)
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
	cleanActionFlag(challenge)
	return challenge, result.Error
}

// FindChallengesByAuthor implements ChallengeRepository.
func (r *challengeRepository) FindChallengesByAuthor(author *entities.User) (challenges []*entities.Challenge, err error) {
	result := r.db.Where("author_id = ?", author.ID).Find(&challenges)
	return challenges, result.Error
}

// UpdateChallenge implements ChallengeRepository.
func (r *challengeRepository) UpdateChallenge(challenge *entities.Challenge) error {
	result := r.db.Save(challenge)
	return result.Error
}

func cleanActionFlag(challenge *entities.Challenge) {
	for _, testcase := range challenge.Testcases {
		testcase.ActionFlag = ""
	}
}

func NewChallengeRepository(db *gorm.DB) ChallengeRepository {
	return &challengeRepository{db}
}
