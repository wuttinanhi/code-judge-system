package services

import (
	"errors"
	"log"
	"sync"

	"github.com/wuttinanhi/code-judge-system/entities"
	"github.com/wuttinanhi/code-judge-system/repositories"
)

type SubmissionService interface {
	CreateSubmission(submission *entities.Submission) (*entities.Submission, error)
	DeleteSubmission(submission *entities.Submission) error
	GetSubmissionByID(submissionID uint) (*entities.Submission, error)
	GetSubmissionByUser(user *entities.User) ([]*entities.Submission, error)
	GetSubmissionByChallenge(challenge *entities.Challenge) ([]*entities.Submission, error)
	CreateSubmissionTestcase(submissionTestcase *entities.SubmissionTestcase) (*entities.SubmissionTestcase, error)
	GetSubmissionTestcaseBySubmission(submission *entities.Submission) ([]*entities.SubmissionTestcase, error)
	SubmitSubmission(submission *entities.Submission) (*entities.Submission, error)
	ProcessSubmission(submission *entities.Submission) (*entities.Submission, error)
	Pagination(options *entities.SubmissionPaginationOptions) (result *entities.PaginationResult[*entities.Submission], err error)
}

type submissionService struct {
	submissionRepository repositories.SubmissionRepository
	challengeService     ChallengeService
	sandboxService       SandboxService
}

// Pagination implements SubmissionService.
func (s *submissionService) Pagination(options *entities.SubmissionPaginationOptions) (result *entities.PaginationResult[*entities.Submission], err error) {
	result, err = s.submissionRepository.Pagination(options)
	return
}

// ProcessSubmission implements SubmissionService.
func (s *submissionService) ProcessSubmission(submission *entities.Submission) (*entities.Submission, error) {
	submissionTestcases := submission.SubmissionTestcases

	sandbox, err := s.sandboxService.CreateSandbox(submission.Language, submission.SourceCode)
	if err != nil {
		return nil, errors.New("failed to create sandbox")
	}
	defer s.sandboxService.CleanUp(sandbox)

	wg := sync.WaitGroup{}

	for _, testcase := range submissionTestcases {
		wg.Add(1)

		go func(testcase *entities.SubmissionTestcase) {
			defer wg.Done()

			challengeTestcase, err := s.challengeService.FindTestcaseByID(testcase.ChallengeTestcaseID)
			if err != nil {
				log.Println("failed to get challenge testcase ID:", testcase.ID, "with error:", err)
				return
			}

			result := s.sandboxService.Run(
				sandbox,
				challengeTestcase.Input,
				challengeTestcase.LimitMemory,
				challengeTestcase.LimitTimeMs,
			)
			if result.Err != nil {
				testcase.Status = entities.SubmissionStatusWrong
			}

			testcase.Output = result.Stdout + result.Stderr

			if result.ExitCode != 0 {
				testcase.Status = entities.SubmissionStatusWrong
			}

			if testcase.Output == challengeTestcase.ExpectedOutput {
				testcase.Status = entities.SubmissionStatusCorrect
			} else {
				testcase.Status = entities.SubmissionStatusWrong
			}

			_, err = s.submissionRepository.UpdateSubmissionTestcase(testcase)
			if err != nil {
				log.Println("failed to update submission testcase ID:", testcase.ID, "with error:", err)
			}
		}(testcase)
	}

	// wait for all goroutines to finish
	wg.Wait()

	if submission.IsCorrect() {
		submission.Status = entities.SubmissionStatusCorrect
	} else {
		submission.Status = entities.SubmissionStatusWrong
	}

	submission, err = s.submissionRepository.UpdateSubmission(submission)
	if err != nil {
		return nil, err
	}

	return submission, nil
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

	submissionTestcases := make([]*entities.SubmissionTestcase, len(challengeTestcases))
	for i, challengeTestcase := range challengeTestcases {
		submissionTestcases[i] = &entities.SubmissionTestcase{
			ChallengeTestcaseID: challengeTestcase.ID,
			Status:              entities.SubmissionStatusPending,
			Output:              "",
		}
	}

	submission.SubmissionTestcases = submissionTestcases

	// create submission
	submission, err = s.CreateSubmission(submission)
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
func (s *submissionService) GetSubmissionByChallenge(challenge *entities.Challenge) ([]*entities.Submission, error) {
	submissions, err := s.submissionRepository.GetSubmissionByChallenge(challenge)
	return submissions, err
}

// GetSubmissionByID implements SubmissionService.
func (s *submissionService) GetSubmissionByID(submissionID uint) (*entities.Submission, error) {
	submission, err := s.submissionRepository.GetSubmissionByID(submissionID)
	return submission, err
}

// GetSubmissionByUser implements SubmissionService.
func (s *submissionService) GetSubmissionByUser(user *entities.User) ([]*entities.Submission, error) {
	submissions, err := s.submissionRepository.GetSubmissionByUser(user)
	return submissions, err
}

// GetSubmissionTestcaseBySubmission implements SubmissionService.
func (s *submissionService) GetSubmissionTestcaseBySubmission(submission *entities.Submission) ([]*entities.SubmissionTestcase, error) {
	submissionTestcases, err := s.submissionRepository.GetSubmissionTestcaseBySubmission(submission)
	return submissionTestcases, err
}

func NewSubmissionService(submissionRepository repositories.SubmissionRepository, challengeService ChallengeService, sandboxService SandboxService) SubmissionService {
	return &submissionService{
		submissionRepository: submissionRepository,
		challengeService:     challengeService,
		sandboxService:       sandboxService,
	}
}
