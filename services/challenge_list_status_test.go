package services_test

import (
	"testing"

	"github.com/wuttinanhi/code-judge-system/databases"
	"github.com/wuttinanhi/code-judge-system/entities"
	"github.com/wuttinanhi/code-judge-system/services"
)

/*
Quick test to test functionality of PaginationChallengesWithStatus
*/
func TestChallengeWithStatus(t *testing.T) {
	db := databases.NewTempSQLiteDatabase()
	testServiceKit := services.CreateServiceKit(db)

	// register a user
	testServiceKit.UserService.Register("admin@example.com", "admin", "admin")

	user, err := testServiceKit.UserService.FindUserByID(1)
	if err != nil {
		panic(err)
	}

	testServiceKit.ChallengeService.CreateChallenge(&entities.Challenge{
		Name:        "Test Challenge 1",
		Description: "Test Description 1",
		UserID:      user.ID,
	})

	testServiceKit.ChallengeService.CreateChallenge(&entities.Challenge{
		Name:        "Test Challenge 2",
		Description: "Test Description 2",
		UserID:      user.ID,
	})

	testServiceKit.ChallengeService.CreateChallenge(&entities.Challenge{
		Name:        "Test Challenge 3",
		Description: "Test Description 3",
		UserID:      user.ID,
	})

	testServiceKit.ChallengeService.CreateChallenge(&entities.Challenge{
		Name:        "Test Challenge 4",
		Description: "Test Description 4",
		UserID:      user.ID,
	})

	{
		testServiceKit.SubmissionService.SubmitSubmission(&entities.Submission{
			ChallengeID: 1,
			UserID:      user.ID,
			Language:    "go",
			SourceCode:  "test sourcecode",
			Status:      entities.SubmissionStatusCorrect,
		})

		testServiceKit.SubmissionService.SubmitSubmission(&entities.Submission{
			ChallengeID: 1,
			UserID:      user.ID,
			Language:    "go",
			SourceCode:  "test sourcecode",
			Status:      entities.SubmissionStatusPending,
		})
	}

	{
		testServiceKit.SubmissionService.SubmitSubmission(&entities.Submission{
			ChallengeID: 2,
			UserID:      user.ID,
			Language:    "go",
			SourceCode:  "test sourcecode",
			Status:      entities.SubmissionStatusCorrect,
		})

		testServiceKit.SubmissionService.SubmitSubmission(&entities.Submission{
			ChallengeID: 2,
			UserID:      user.ID,
			Language:    "go",
			SourceCode:  "test sourcecode",
			Status:      entities.SubmissionStatusCorrect,
		})
	}

	{
		testServiceKit.SubmissionService.SubmitSubmission(&entities.Submission{
			ChallengeID: 3,
			UserID:      user.ID,
			Language:    "go",
			SourceCode:  "test sourcecode",
			Status:      entities.SubmissionStatusCorrect,
		})

		testServiceKit.SubmissionService.SubmitSubmission(&entities.Submission{
			ChallengeID: 3,
			UserID:      user.ID,
			Language:    "go",
			SourceCode:  "test sourcecode",
			Status:      entities.SubmissionStatusWrong,
		})
	}

	challenges, err := testServiceKit.ChallengeService.PaginationChallengesWithStatus(1, 10, user)
	if err != nil {
		panic(err)
	}

	// Must be
	// 1 Test Challenge 1 PENDING
	// 2 Test Challenge 2 CORRECT
	// 3 Test Challenge 3 WRONG
	// 4 Test Challenge 4 NOTSOLVE

	if challenges[0].SubmissionStatus != entities.SubmissionStatusPending {
		t.Errorf("Expected status PENDING, got %v", challenges[0].SubmissionStatus)
	}

	if challenges[1].SubmissionStatus != entities.SubmissionStatusCorrect {
		t.Errorf("Expected status CORRECT, got %v", challenges[1].SubmissionStatus)
	}

	if challenges[2].SubmissionStatus != entities.SubmissionStatusWrong {
		t.Errorf("Expected status WRONG, got %v", challenges[2].SubmissionStatus)
	}

	if challenges[3].SubmissionStatus != entities.SubmissionStatusNotSolve {
		t.Errorf("Expected status NOTSOLVE, got %v", challenges[3].SubmissionStatus)
	}
}
