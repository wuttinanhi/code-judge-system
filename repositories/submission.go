package repositories

import (
	"github.com/wuttinanhi/code-judge-system/entities"
	"gorm.io/gorm"
)

type SubmissionRepository interface {
	CreateSubmission(submission *entities.Submission) (*entities.Submission, error)
	DeleteSubmission(submission *entities.Submission) error
	GetSubmissionByID(submissionID uint) (*entities.Submission, error)
	GetSubmissionByUser(user *entities.User) ([]entities.Submission, error)
	GetSubmissionByChallenge(challenge *entities.Challenge) ([]entities.Submission, error)
	CreateSubmissionTestcase(submissionTestcase *entities.SubmissionTestcase) (*entities.SubmissionTestcase, error)
	GetSubmissionTestcaseBySubmission(submission *entities.Submission) ([]entities.SubmissionTestcase, error)
	CreateSubmissionWithTestcase(submission *entities.Submission, submissionTestcases []entities.SubmissionTestcase) (*entities.Submission, error)
}

type submissionRepository struct {
	db *gorm.DB
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

// CreateSubmission implements SubmissionRepository.
func (r *submissionRepository) CreateSubmission(submission *entities.Submission) (*entities.Submission, error) {
	result := r.db.Create(submission)
	return submission, result.Error
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
func (r *submissionRepository) GetSubmissionByChallenge(challenge *entities.Challenge) ([]entities.Submission, error) {
	var submissions []entities.Submission
	result := r.db.Where("challenge_id = ?", challenge.ID).Find(&submissions)
	return submissions, result.Error
}

// GetSubmissionByID implements SubmissionRepository.
func (r *submissionRepository) GetSubmissionByID(submissionID uint) (*entities.Submission, error) {
	var submission entities.Submission
	result := r.db.Where("submission_id = ?", submissionID).First(&submission)
	return &submission, result.Error
}

// GetSubmissionByUser implements SubmissionRepository.
func (r *submissionRepository) GetSubmissionByUser(user *entities.User) ([]entities.Submission, error) {
	var submissions []entities.Submission
	result := r.db.Where("user_id = ?", user.ID).Find(&submissions)
	return submissions, result.Error
}

// GetSubmissionTestcaseBySubmission implements SubmissionRepository.
func (r *submissionRepository) GetSubmissionTestcaseBySubmission(submission *entities.Submission) ([]entities.SubmissionTestcase, error) {
	var submissionTestcases []entities.SubmissionTestcase
	result := r.db.Where("submission_id = ?", submission.ID).Find(&submissionTestcases)
	return submissionTestcases, result.Error
}

func NewSubmissionRepository(db *gorm.DB) SubmissionRepository {
	return &submissionRepository{db: db}
}
