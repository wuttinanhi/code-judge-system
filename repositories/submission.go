package repositories

import (
	"fmt"
	"strings"

	"github.com/wuttinanhi/code-judge-system/entities"
	"gorm.io/gorm"
)

type SubmissionRepository interface {
	CreateSubmission(submission *entities.Submission) (*entities.Submission, error)
	DeleteSubmission(submission *entities.Submission) error
	GetSubmissionByID(submissionID uint) (*entities.Submission, error)
	GetSubmissionByUser(user *entities.User) ([]*entities.Submission, error)
	GetSubmissionByChallenge(challenge *entities.Challenge) ([]*entities.Submission, error)
	CreateSubmissionTestcase(submissionTestcase *entities.SubmissionTestcase) (*entities.SubmissionTestcase, error)
	GetSubmissionTestcaseBySubmission(submission *entities.Submission) ([]*entities.SubmissionTestcase, error)
	// CreateSubmissionWithTestcase(submission *entities.Submission, submissionTestcases []entities.SubmissionTestcase) (*entities.Submission, error)
	UpdateSubmission(submission *entities.Submission) (*entities.Submission, error)
	UpdateSubmissionTestcase(submissionTestcase *entities.SubmissionTestcase) (*entities.SubmissionTestcase, error)
	Pagination(options *entities.SubmissionPaginationOptions) (result *entities.PaginationResult[*entities.Submission], err error)
}

type submissionRepository struct {
	db *gorm.DB
}

func (r *submissionRepository) Pagination(options *entities.SubmissionPaginationOptions) (result *entities.PaginationResult[*entities.Submission], err error) {
	result = &entities.PaginationResult[*entities.Submission]{
		Items: make([]*entities.Submission, 0),
		Total: 0,
	}

	// calculate offset for pagination
	offset := (options.Page - 1) * options.Limit

	// convert options.Order to uppercase
	options.Order = strings.ToUpper(options.Order)

	// if options.Order is not ASC or DESC then throw error
	if options.Order != "ASC" && options.Order != "DESC" {
		err = fmt.Errorf("invalid order option")
		return
	}

	// if user or challenge is nil then set it to 0
	if options.User == nil {
		options.User = &entities.User{ID: 0}
	}
	if options.Challenge == nil {
		options.Challenge = &entities.Challenge{ID: 0}
	}

	var submissions []*entities.Submission
	submissionQuery := r.db.Model(&entities.Submission{}).
		Preload("SubmissionTestcases").
		Where(&entities.Submission{
			UserID:      options.User.ID,
			ChallengeID: options.Challenge.ID,
		}).
		Limit(options.Limit).
		Offset(offset).
		Order("id " + options.Order).
		Find(&submissions)
	if submissionQuery.Error != nil {
		err = submissionQuery.Error
		return
	}

	// Query to count total submissions
	var totalCount int64
	r.db.Model(&entities.Submission{}).
		Where(&entities.Submission{
			UserID:      options.User.ID,
			ChallengeID: options.Challenge.ID,
		}).
		Count(&totalCount)

	result.Items = submissions
	result.Total = int(totalCount)

	return
}

// CreateSubmission implements SubmissionRepository.
func (r *submissionRepository) CreateSubmission(submission *entities.Submission) (*entities.Submission, error) {
	result := r.db.Create(submission)
	return submission, result.Error
}

// UpdateSubmissionTestcase implements SubmissionRepository.
func (r *submissionRepository) UpdateSubmissionTestcase(submissionTestcase *entities.SubmissionTestcase) (*entities.SubmissionTestcase, error) {
	result := r.db.Save(submissionTestcase)
	return submissionTestcase, result.Error
}

// UpdateSubmission implements SubmissionRepository.
func (r *submissionRepository) UpdateSubmission(submission *entities.Submission) (*entities.Submission, error) {
	result := r.db.Save(submission)
	return submission, result.Error
}

// CreateSubmissionWithTestcase implements SubmissionRepository.
func (r *submissionRepository) CreateSubmissionWithTestcase(submission *entities.Submission, testcaes []entities.SubmissionTestcase) (*entities.Submission, error) {
	result := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(submission).Error; err != nil {
			return err
		}

		for _, submissionTestcase := range testcaes {
			submissionTestcase.SubmissionID = submission.ID
			if err := tx.Create(&submissionTestcase).Error; err != nil {
				return err
			}
		}

		return nil
	})

	return submission, result
}

// CreateSubmissionTestcase implements SubmissionRepository.
func (r *submissionRepository) CreateSubmissionTestcase(submissionTestcase *entities.SubmissionTestcase) (*entities.SubmissionTestcase, error) {
	result := r.db.Create(submissionTestcase)
	return submissionTestcase, result.Error
}

// DeleteSubmission implements SubmissionRepository.
func (r *submissionRepository) DeleteSubmission(submission *entities.Submission) error {
	err := r.db.Delete(submission).Error
	return err
}

// GetSubmissionByChallenge implements SubmissionRepository.
func (r *submissionRepository) GetSubmissionByChallenge(challenge *entities.Challenge) ([]*entities.Submission, error) {
	var submissions []*entities.Submission
	result := r.db.Where("challenge_id = ?", challenge.ID).Find(&submissions)
	return submissions, result.Error
}

// GetSubmissionByID implements SubmissionRepository.
func (r *submissionRepository) GetSubmissionByID(submissionID uint) (*entities.Submission, error) {
	var submission *entities.Submission
	result := r.db.Model(&entities.Submission{}).
		Preload("SubmissionTestcases").
		Find(&submission, submissionID)
	return submission, result.Error
}

// GetSubmissionByUser implements SubmissionRepository.
func (r *submissionRepository) GetSubmissionByUser(user *entities.User) ([]*entities.Submission, error) {
	var submissions []*entities.Submission
	result := r.db.Model(&entities.Submission{}).
		Preload("SubmissionTestcases").
		Where(&entities.Submission{UserID: user.ID}).
		Find(&submissions)
	return submissions, result.Error
}

// GetSubmissionTestcaseBySubmission implements SubmissionRepository.
func (r *submissionRepository) GetSubmissionTestcaseBySubmission(submission *entities.Submission) ([]*entities.SubmissionTestcase, error) {
	var submissionTestcases []*entities.SubmissionTestcase
	result := r.db.
		Where(&entities.SubmissionTestcase{SubmissionID: submission.ID}).
		Find(&submissionTestcases)
	return submissionTestcases, result.Error
}

func NewSubmissionRepository(db *gorm.DB) SubmissionRepository {
	return &submissionRepository{db: db}
}
