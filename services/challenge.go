package services

import (
	"fmt"

	"github.com/wuttinanhi/code-judge-system/entities"
	"github.com/wuttinanhi/code-judge-system/repositories"
)

type ChallengeService interface {
	CreateChallenge(challenge *entities.Challenge) (*entities.Challenge, error)
	// UpdateChallenge(challenge *entities.Challenge) (err error)
	DeleteChallenge(challenge *entities.Challenge) (err error)
	FindChallengeByID(challengeID uint) (challenge *entities.Challenge, err error)
	AllChallenges() (challenges []*entities.Challenge, err error)
	AddTestcase(challenge *entities.Challenge, testcase *entities.ChallengeTestcase) (*entities.ChallengeTestcase, error)
	UpdateTestcase(testcase *entities.ChallengeTestcase) (err error)
	DeleteTestcase(testcase *entities.ChallengeTestcase) (err error)
	AllTestcases(challenge *entities.Challenge) (testcases []*entities.ChallengeTestcase, err error)
	FindTestcaseByID(testcaseID uint) (testcase *entities.ChallengeTestcase, err error)
	PaginationChallengesWithStatus(options *entities.ChallengePaginationOptions) (result *entities.PaginationResult[*entities.ChallengeExtended], err error)
	UpdateChallengeWithTestcase(challenge *entities.Challenge) (err error)
	CountAllChallengesByUser(user *entities.User) (total int64, err error)
	ValidateTestcases(testcases []*entities.ChallengeTestcase) (err error)
}

type challengeService struct {
	challengeRepo  repositories.ChallengeRepository
	sandboxService SandboxService
}

// ValidateTestcases implements ChallengeService.
func (s *challengeService) ValidateTestcases(testcases []*entities.ChallengeTestcase) (err error) {
	// loop testcases
	for _, testcase := range testcases {
		maxMemoryErr := s.sandboxService.ValidateMemoryLimit(testcase.LimitMemory)
		if maxMemoryErr != nil {
			err = fmt.Errorf("testcase #%d: max memory exceeded sandbox limit", testcase.ID)
			return
		}

		maxTimeLimitErr := s.sandboxService.ValidateTimeLimit(testcase.LimitTimeMs)
		if maxTimeLimitErr != nil {
			err = fmt.Errorf("testcase #%d: max run time exceeded sandbox limit", testcase.ID)
			return
		}
	}
	return
}

// CountAllChallengesByUser implements ChallengeService.
func (s *challengeService) CountAllChallengesByUser(user *entities.User) (total int64, err error) {
	total, err = s.challengeRepo.CountAllChallengesByUser(user)
	return total, err
}

// UpdateChallengeWithTestcase implements ChallengeService.
func (s *challengeService) UpdateChallengeWithTestcase(challenge *entities.Challenge) (err error) {
	err = s.ValidateTestcases(challenge.Testcases)
	if err != nil {
		return err
	}
	err = s.challengeRepo.UpdateChallengeWithTestcase(challenge)
	return
}

// PaginationChallengesWithStatus implements ChallengeService.
func (s *challengeService) PaginationChallengesWithStatus(options *entities.ChallengePaginationOptions) (result *entities.PaginationResult[*entities.ChallengeExtended], err error) {
	result, err = s.challengeRepo.PaginationChallengesWithStatus(options)
	return result, err
}

// CreateChallenge implements ChallengeService.
func (s *challengeService) CreateChallenge(challenge *entities.Challenge) (*entities.Challenge, error) {
	err := s.ValidateTestcases(challenge.Testcases)
	if err != nil {
		return nil, err
	}
	challenge, err = s.challengeRepo.CreateChallenge(challenge)
	return challenge, err
}

// AddTestcase implements ChallengeService.
func (s *challengeService) AddTestcase(challenge *entities.Challenge, testcase *entities.ChallengeTestcase) (*entities.ChallengeTestcase, error) {
	testcase, err := s.challengeRepo.AddTestcase(challenge, testcase)
	return testcase, err
}

// AllChallenges implements ChallengeService.
func (s *challengeService) AllChallenges() (challenges []*entities.Challenge, err error) {
	challenges, err = s.challengeRepo.AllChallenges()
	return challenges, err
}

// AllTestcases implements ChallengeService.
func (s *challengeService) AllTestcases(challenge *entities.Challenge) (testcases []*entities.ChallengeTestcase, err error) {
	testcases, err = s.challengeRepo.AllTestcases(challenge)
	return testcases, err
}

// DeleteChallenge implements ChallengeService.
func (s *challengeService) DeleteChallenge(challenge *entities.Challenge) (err error) {
	err = s.challengeRepo.DeleteChallenge(challenge)
	return err
}

// DeleteTestcase implements ChallengeService.
func (s *challengeService) DeleteTestcase(testcase *entities.ChallengeTestcase) (err error) {
	err = s.challengeRepo.DeleteTestcase(testcase)
	return err
}

// FindChallengeByID implements ChallengeService.
func (s *challengeService) FindChallengeByID(challengeID uint) (challenge *entities.Challenge, err error) {
	challenge, err = s.challengeRepo.FindChallengeByID(challengeID)
	return challenge, err
}

// FindTestcaseByID implements ChallengeService.
func (s *challengeService) FindTestcaseByID(testcaseID uint) (testcase *entities.ChallengeTestcase, err error) {
	testcase, err = s.challengeRepo.FindTestcaseByID(testcaseID)
	return testcase, err
}

// UpdateChallenge implements ChallengeService.
func (s *challengeService) UpdateChallenge(challenge *entities.Challenge) (err error) {
	err = s.challengeRepo.UpdateChallenge(challenge)
	return err
}

// UpdateTestcase implements ChallengeService.
func (s *challengeService) UpdateTestcase(testcase *entities.ChallengeTestcase) (err error) {
	err = s.challengeRepo.UpdateTestcase(testcase)
	return err
}

func NewChallengeService(challengeRepo repositories.ChallengeRepository, sandboxService SandboxService) ChallengeService {
	return &challengeService{
		challengeRepo:  challengeRepo,
		sandboxService: sandboxService,
	}
}
