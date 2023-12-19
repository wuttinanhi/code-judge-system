package services

import (
	"github.com/wuttinanhi/code-judge-system/entities"
	"github.com/wuttinanhi/code-judge-system/repositories"
)

type SubmissionService interface {
	CreateSubmission(submission *entities.Submission) (*entities.Submission, error)
	DeleteSubmission(submission *entities.Submission) error
	GetSubmissionByID(submissionID uint) (*entities.Submission, error)
	GetSubmissionByUser(user *entities.User) ([]entities.Submission, error)
	GetSubmissionByChallenge(challenge *entities.Challenge) ([]entities.Submission, error)
	CreateSubmissionTestcase(submissionTestcase *entities.SubmissionTestcase) (*entities.SubmissionTestcase, error)
	GetSubmissionTestcaseBySubmission(submission *entities.Submission) ([]entities.SubmissionTestcase, error)
	SubmitSubmission(submission *entities.Submission) (*entities.Submission, error)
}

type submissionService struct {
	submissionRepository repositories.SubmissionRepository
	challengeService     ChallengeService
}

// SubmitSubmission implements SubmissionService.
func (s *submissionService) SubmitSubmission(submission *entities.Submission) (*entities.Submission, error) {
	// get challenge
	challenge, err := s.challengeService.FindChallengeByID(submission.ChallengeID)
	if err != nil {
		return nil, err
	}

	// get all challenge testcases
	challengeTestcases, err := s.challengeService.AllTestcases(challenge)
	if err != nil {
		return nil, err
	}

	submissionTestcases := make([]entities.SubmissionTestcase, len(challengeTestcases))
	for i, challengeTestcase := range challengeTestcases {
		submissionTestcases[i] = entities.SubmissionTestcase{
			ChallengeTestcaseID: challengeTestcase.TestcaseID,
			Status:              entities.SubmissionStatusPending,
			Output:              "",
		}
	}

	// create submission
	submission, err = s.submissionRepository.CreateSubmissionWithTestcase(submission, submissionTestcases)
	if err != nil {
		return nil, err
	}

	// finally return submission
	return submission, nil
}

// CreateSubmission implements SubmissionService.
func (s *submissionService) CreateSubmission(submission *entities.Submission) (*entities.Submission, error) {
	submission, err := s.submissionRepository.CreateSubmission(submission)
	return submission, err
}

// CreateSubmissionTestcase implements SubmissionService.
func (s *submissionService) CreateSubmissionTestcase(submissionTestcase *entities.SubmissionTestcase) (*entities.SubmissionTestcase, error) {
	submissionTestcase, err := s.submissionRepository.CreateSubmissionTestcase(submissionTestcase)
	return submissionTestcase, err
}

// DeleteSubmission implements SubmissionService.
func (s *submissionService) DeleteSubmission(submission *entities.Submission) error {
	err := s.submissionRepository.DeleteSubmission(submission)
	return err
}

// GetSubmissionByChallenge implements SubmissionService.
func (s *submissionService) GetSubmissionByChallenge(challenge *entities.Challenge) ([]entities.Submission, error) {
	submissions, err := s.submissionRepository.GetSubmissionByChallenge(challenge)
	return submissions, err
}

// GetSubmissionByID implements SubmissionService.
func (s *submissionService) GetSubmissionByID(submissionID uint) (*entities.Submission, error) {
	submission, err := s.submissionRepository.GetSubmissionByID(submissionID)
	return submission, err
}

// GetSubmissionByUser implements SubmissionService.
func (s *submissionService) GetSubmissionByUser(user *entities.User) ([]entities.Submission, error) {
	submissions, err := s.submissionRepository.GetSubmissionByUser(user)
	return submissions, err
}

// GetSubmissionTestcaseBySubmission implements SubmissionService.
func (s *submissionService) GetSubmissionTestcaseBySubmission(submission *entities.Submission) ([]entities.SubmissionTestcase, error) {
	submissionTestcases, err := s.submissionRepository.GetSubmissionTestcaseBySubmission(submission)
	return submissionTestcases, err
}

func NewSubmissionService(submissionRepository repositories.SubmissionRepository, challengeService ChallengeService) SubmissionService {
	return &submissionService{
		submissionRepository: submissionRepository,
		challengeService:     challengeService,
	}
}
