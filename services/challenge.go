package services

import (
	"github.com/wuttinanhi/code-judge-system/entities"
	"github.com/wuttinanhi/code-judge-system/repositories"
)

type ChallengeService interface {
	CreateChallenge(challenge *entities.Challenge) (*entities.Challenge, error)
	UpdateChallenge(challenge *entities.Challenge) (err error)
	DeleteChallenge(challenge *entities.Challenge) (err error)
	FindChallengeByID(challengeID uint) (challenge *entities.Challenge, err error)
	FindChallengesByAuthor(user *entities.User) (challenges []*entities.Challenge, err error)
	AllChallenges() (challenges []*entities.Challenge, err error)
	AddTestcase(challenge *entities.Challenge, testcase *entities.ChallengeTestcase) (*entities.ChallengeTestcase, error)
	UpdateTestcase(testcase *entities.ChallengeTestcase) (err error)
	DeleteTestcase(testcase *entities.ChallengeTestcase) (err error)
	AllTestcases(challenge *entities.Challenge) (testcases []*entities.ChallengeTestcase, err error)
	FindTestcaseByID(testcaseID uint) (testcase *entities.ChallengeTestcase, err error)
	PaginationChallengesWithStatus(options *entities.ChallengePaginationOptions) (result *entities.PaginationResult[*entities.ChallengeExtended], err error)
}

type challengeService struct {
	challengeRepo repositories.ChallengeRepository
}

// PaginationChallengesWithStatus implements ChallengeService.
func (s *challengeService) PaginationChallengesWithStatus(options *entities.ChallengePaginationOptions) (result *entities.PaginationResult[*entities.ChallengeExtended], err error) {
	result, err = s.challengeRepo.PaginationChallengesWithStatus(options)
	return result, err
}

// CreateChallenge implements ChallengeService.
func (s *challengeService) CreateChallenge(challenge *entities.Challenge) (*entities.Challenge, error) {
	challenge, err := s.challengeRepo.CreateChallenge(challenge)
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

// FindChallengesByAuthor implements ChallengeService.
func (s *challengeService) FindChallengesByAuthor(user *entities.User) (challenges []*entities.Challenge, err error) {
	challenges, err = s.challengeRepo.FindChallengesByAuthor(user)
	return challenges, err
}

// FindTestcaseByID implements ChallengeService.
func (s *challengeService) FindTestcaseByID(testcaseID uint) (testcase *entities.ChallengeTestcase, err error) {
	testcase, err = s.challengeRepo.FindTestcaseByID(testcaseID)
	return testcase, err
}

// UpdateChallenge implements ChallengeService.
func (s *challengeService) UpdateChallenge(challenge *entities.Challenge) (err error) {
	err = s.challengeRepo.UpdateChallengeWithTestcase(challenge)
	return err
}

// UpdateTestcase implements ChallengeService.
func (s *challengeService) UpdateTestcase(testcase *entities.ChallengeTestcase) (err error) {
	err = s.challengeRepo.UpdateTestcase(testcase)
	return err
}

func NewChallengeService(challengeRepo repositories.ChallengeRepository) ChallengeService {
	return &challengeService{
		challengeRepo: challengeRepo,
	}
}
