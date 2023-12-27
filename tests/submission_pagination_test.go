package tests_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/wuttinanhi/code-judge-system/controllers"
	"github.com/wuttinanhi/code-judge-system/databases"
	"github.com/wuttinanhi/code-judge-system/entities"
	"github.com/wuttinanhi/code-judge-system/services"
)

// Quick test to test functionality of PaginationChallengesWithStatus
func TestSubmissionPagination(t *testing.T) {
	db := databases.NewTempSQLiteDatabase()
	testServiceKit := services.CreateServiceKit(db)
	app := controllers.SetupAPI(testServiceKit)

	// register admin
	admin, err := testServiceKit.UserService.Register("admin@example.com", "admin", "admin")
	if err != nil {
		panic(err)
	}

	// register user
	user, err := testServiceKit.UserService.Register("user@example.com", "user", "user")
	if err != nil {
		panic(err)
	}

	// generate user token
	userAccessToken, err := testServiceKit.JWTService.GenerateToken(*user)
	if err != nil {
		panic(err)
	}

	// create challenges
	_, err = testServiceKit.ChallengeService.CreateChallenge(&entities.Challenge{
		Name:        "Test Challenge 1",
		Description: "Test Description 1",
		UserID:      admin.ID,
	})
	if err != nil {
		t.Fatal(err)
	}

	challenge2, err := testServiceKit.ChallengeService.CreateChallenge(&entities.Challenge{
		Name:        "Test Challenge 2",
		Description: "Test Description 2",
		UserID:      admin.ID,
	})
	if err != nil {
		t.Fatal(err)
	}

	// create submissions
	testServiceKit.SubmissionService.SubmitSubmission(&entities.Submission{
		ChallengeID: 1,
		UserID:      admin.ID,
		Language:    "go",
		SourceCode:  "test sourcecode",
		Status:      entities.SubmissionStatusCorrect,
	})

	testServiceKit.SubmissionService.SubmitSubmission(&entities.Submission{
		ChallengeID: 1,
		UserID:      admin.ID,
		Language:    "go",
		SourceCode:  "test sourcecode",
		Status:      entities.SubmissionStatusPending,
	})

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

	testServiceKit.SubmissionService.SubmitSubmission(&entities.Submission{
		ChallengeID: 2,
		UserID:      user.ID,
		Language:    "go",
		SourceCode:  "test sourcecode",
		Status:      entities.SubmissionStatusCorrect,
	})

	t.Run("Test Submission Pagination Normal", func(t *testing.T) {
		submissions, err := testServiceKit.SubmissionService.Pagination(&entities.SubmissionPaginationOptions{
			PaginationOptions: entities.PaginationOptions{
				Page:  1,
				Limit: 10,
				Order: "ASC",
				Sort:  "id",
			},
		})
		if err != nil {
			panic(err)
		}

		if len(submissions.Items) != 5 {
			t.Errorf("expect 5 submissions got %d", len(submissions.Items))
		}
	})

	t.Run("Test Submission Pagination With User", func(t *testing.T) {
		submissions, err := testServiceKit.SubmissionService.Pagination(&entities.SubmissionPaginationOptions{
			PaginationOptions: entities.PaginationOptions{
				Page:  1,
				Limit: 10,
				Order: "ASC",
				Sort:  "id",
			},
			User: admin,
		})
		if err != nil {
			panic(err)
		}

		if len(submissions.Items) != 2 {
			t.Errorf("expect 2 submissions got %d", len(submissions.Items))
		}
	})

	t.Run("Test Submission Pagination With Challenge", func(t *testing.T) {
		submissions, err := testServiceKit.SubmissionService.Pagination(&entities.SubmissionPaginationOptions{
			PaginationOptions: entities.PaginationOptions{
				Page:  1,
				Limit: 10,
				Order: "ASC",
				Sort:  "id",
			},
			User:      user,
			Challenge: challenge2,
		})
		if err != nil {
			panic(err)
		}

		if len(submissions.Items) != 3 {
			t.Errorf("expect 3 submissions got %d", len(submissions.Items))
		}
	})

	t.Run("/submission/pagination", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/submission/pagination", nil)
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Authorization", "Bearer "+userAccessToken)

		response, err := app.Test(request, -1)
		if err != nil {
			t.Error(err)
		}

		if response.StatusCode != http.StatusOK {
			t.Errorf("Expected status OK, got %v", response.StatusCode)
		}

		// try parse json response to pagination result
		var result entities.PaginationResult[entities.Submission]
		err = json.NewDecoder(response.Body).Decode(&result)
		if err != nil {
			t.Error(err)
		}
		if result.Total != 5 {
			t.Errorf("Expected total 5, got %v", result.Total)
		}
	})

	t.Run("/submission/pagination?user_id=1", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/submission/pagination?user_id=1", nil)
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Authorization", "Bearer "+userAccessToken)

		response, err := app.Test(request, -1)
		if err != nil {
			t.Error(err)
		}

		if response.StatusCode != http.StatusOK {
			t.Errorf("Expected status OK, got %v", response.StatusCode)
		}

		// try parse json response to pagination result
		var result entities.PaginationResult[entities.Submission]
		err = json.NewDecoder(response.Body).Decode(&result)
		if err != nil {
			t.Error(err)
		}
		if result.Total != 2 {
			t.Errorf("Expected total 2, got %v", result.Total)
		}
	})

	t.Run("/submission/pagination?challenge_id=2", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/submission/pagination?challenge_id=2", nil)
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Authorization", "Bearer "+userAccessToken)

		response, err := app.Test(request, -1)
		if err != nil {
			t.Error(err)
		}

		if response.StatusCode != http.StatusOK {
			t.Errorf("Expected status OK, got %v", response.StatusCode)
		}

		// try parse json response to pagination result
		var result entities.PaginationResult[entities.Submission]
		err = json.NewDecoder(response.Body).Decode(&result)
		if err != nil {
			t.Error(err)
		}
		if result.Total != 3 {
			t.Errorf("Expected total 3, got %v", result.Total)
		}
	})

	t.Run("/submission/pagination?challenge_id=3", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/submission/pagination?challenge_id=3", nil)
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Authorization", "Bearer "+userAccessToken)

		response, err := app.Test(request, -1)
		if err != nil {
			t.Error(err)
		}

		if response.StatusCode != http.StatusOK {
			t.Errorf("Expected status OK, got %v", response.StatusCode)
		}

		// try parse json response to pagination result
		var result entities.PaginationResult[entities.Submission]
		err = json.NewDecoder(response.Body).Decode(&result)
		if err != nil {
			t.Error(err)
		}
		if result.Total != 0 {
			t.Errorf("Expected total 0, got %v", result.Total)
		}
	})
}
