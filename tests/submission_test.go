package tests

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"testing"

	"github.com/wuttinanhi/code-judge-system/cmds"
	"github.com/wuttinanhi/code-judge-system/entities"
	"github.com/wuttinanhi/code-judge-system/services"
)

func TestSubmissionRoute(t *testing.T) {
	serviceKit := services.CreateTestServiceKit()
	app := cmds.SetupWeb(serviceKit)

	// create user
	user, err := serviceKit.UserService.Register("test-submission-route@example.com", "testpassword", "test-submission-route")
	if err != nil {
		t.Error(err)
	}

	// get user access token
	userAccessToken, err := serviceKit.JWTService.GenerateToken(*user)
	if err != nil {
		t.Error(err)
	}

	// create challenge
	challenge, err := serviceKit.ChallengeService.CreateChallenge(&entities.Challenge{
		Name:        "Test Challenge",
		Description: "Test Description",
	})
	if err != nil {
		t.Error(err)
	}

	// create testcases
	challengetestcases := []entities.ChallengeTestcase{
		{Input: "1 2", ExpectedOutput: "3"},
		{Input: "2 3", ExpectedOutput: "5"},
		{Input: "3 4", ExpectedOutput: "7"},
	}
	for _, challengetestcase := range challengetestcases {
		_, err := serviceKit.ChallengeService.AddTestcase(challenge, &challengetestcase)
		if err != nil {
			t.Error(err)
		}
	}

	t.Run("/submission/submit", func(t *testing.T) {
		dto := entities.SubmissionCreateDTO{
			ChallengeID: challenge.ChallengeID,
			Language:    "go",
			SourceCode:  "test source code",
		}
		requestBody, _ := json.Marshal(dto)

		request, _ := http.NewRequest(http.MethodPost, "/submission/submit", bytes.NewBuffer(requestBody))
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Authorization", "Bearer "+userAccessToken)

		response, err := app.Test(request, -1)
		if err != nil {
			t.Error(err)
		}

		if response.StatusCode != http.StatusOK {
			t.Errorf("Expected status OK, got %v", response.StatusCode)
		}

		// get submission in server-side
		submission, err := serviceKit.SubmissionService.GetSubmissionByID(1)
		if err != nil {
			t.Error(err)
		}
		if submission.ChallengeID != challenge.ChallengeID {
			t.Errorf("Expected challenge id %v, got %v", challenge.ChallengeID, submission.ChallengeID)
		}
		if submission.UserID != user.UserID {
			t.Errorf("Expected user id %v, got %v", user.UserID, submission.UserID)
		}
		if submission.Language != dto.Language {
			t.Errorf("Expected language %v, got %v", dto.Language, submission.Language)
		}
		if submission.SourceCode != dto.SourceCode {
			t.Errorf("Expected source code %v, got %v", dto.SourceCode, submission.SourceCode)
		}

		// validate submission testcases
		submissionTestcases, err := serviceKit.SubmissionService.GetSubmissionTestcaseBySubmission(submission)
		if err != nil {
			t.Error(err)
		}
		if len(submissionTestcases) != len(challengetestcases) {
			t.Errorf("Expected %v submission testcases, got %v", len(challengetestcases), len(submissionTestcases))
		}
		for i := range challengetestcases {
			submissionTestcase := submissionTestcases[i]
			if submissionTestcase.SubmissionTestcaseID != uint(i+1) {
				t.Errorf("Expected challenge testcase id %v, got %v", uint(i+1), submissionTestcase.SubmissionTestcaseID)
			}
			if submissionTestcase.Status != entities.SubmissionStatusPending {
				t.Errorf("Expected status %v, got %v", entities.SubmissionStatusPending, submissionTestcase.Status)
			}
			if submissionTestcase.Output != "" {
				t.Errorf("Expected output %v, got %v", "", submissionTestcase.Output)
			}
		}
	})

	t.Run("/submission/get/submission/:id", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/submission/get/submission/1", nil)
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Authorization", "Bearer "+userAccessToken)

		response, err := app.Test(request, -1)
		if err != nil {
			t.Error(err)
		}

		if response.StatusCode != http.StatusOK {
			t.Errorf("Expected status OK, got %v", response.StatusCode)
		}

		bodyBytes := ResponseBodyToBytes(response)

		var submission entities.Submission
		err = json.Unmarshal(bodyBytes, &submission)
		if err != nil {
			t.Error(err)
		}

		if submission.ChallengeID != challenge.ChallengeID {
			t.Errorf("Expected challenge id %v, got %v", challenge.ChallengeID, submission.ChallengeID)
		}
		if submission.UserID != user.UserID {
			t.Errorf("Expected user id %v, got %v", user.UserID, submission.UserID)
		}
		if submission.Language != "go" {
			t.Errorf("Expected language %v, got %v", "go", submission.Language)
		}
		if submission.SourceCode != "test source code" {
			t.Errorf("Expected source code %v, got %v", "test source code", submission.SourceCode)
		}
	})

	t.Run("/submission/get/user", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/submission/get/user", nil)
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Authorization", "Bearer "+userAccessToken)

		response, err := app.Test(request, -1)
		if err != nil {
			t.Error(err)
		}

		bodyBytes, err := io.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}

		if response.StatusCode != http.StatusOK {
			t.Errorf("Expected status OK, got %v", response.StatusCode)
		}

		var submissions []entities.Submission
		err = json.Unmarshal(bodyBytes, &submissions)
		if err != nil {
			t.Error(err)
		}

		if len(submissions) != 1 {
			t.Errorf("Expected %v submissions, got %v", 1, len(submissions))
		}

		submission := submissions[0]
		if submission.ChallengeID != challenge.ChallengeID {
			t.Errorf("Expected challenge id %v, got %v", challenge.ChallengeID, submission.ChallengeID)
		}
		if submission.UserID != user.UserID {
			t.Errorf("Expected user id %v, got %v", user.UserID, submission.UserID)
		}
		if submission.Language != "go" {
			t.Errorf("Expected language %v, got %v", "go", submission.Language)
		}
		if submission.SourceCode != "test source code" {
			t.Errorf("Expected source code %v, got %v", "test source code", submission.SourceCode)
		}
	})

	t.Run("/submission/get/challenge/:id", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/submission/get/challenge/"+strconv.Itoa(int(challenge.ChallengeID)), nil)
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Authorization", "Bearer "+userAccessToken)

		response, err := app.Test(request, -1)
		if err != nil {
			t.Error(err)
		}

		bodyBytes, err := io.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}

		if response.StatusCode != http.StatusOK {
			t.Errorf("Expected status OK, got %v", response.StatusCode)
		}

		var submissions []entities.Submission
		err = json.Unmarshal(bodyBytes, &submissions)
		if err != nil {
			t.Error(err)
		}

		if len(submissions) != 1 {
			t.Errorf("Expected %v submissions, got %v", 1, len(submissions))
		}

		submission := submissions[0]
		if submission.ChallengeID != challenge.ChallengeID {
			t.Errorf("Expected challenge id %v, got %v", challenge.ChallengeID, submission.ChallengeID)
		}
		if submission.UserID != user.UserID {
			t.Errorf("Expected user id %v, got %v", user.UserID, submission.UserID)
		}
		if submission.Language != "go" {
			t.Errorf("Expected language %v, got %v", "go", submission.Language)
		}
		if submission.SourceCode != "test source code" {
			t.Errorf("Expected source code %v, got %v", "test source code", submission.SourceCode)
		}
	})
}
