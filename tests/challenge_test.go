package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/wuttinanhi/code-judge-system/cmds"
	"github.com/wuttinanhi/code-judge-system/entities"
	"github.com/wuttinanhi/code-judge-system/services"
)

func TestChallengeRoute(t *testing.T) {
	serviceKit := services.CreateTestServiceKit()
	app := cmds.SetupWeb(serviceKit)

	// create user
	user, err := serviceKit.UserService.Register("test-challenge-route@example.com", "testpassword", "testuser")
	if err != nil {
		t.Error(err)
	}

	userAccessToken, err := serviceKit.JWTService.GenerateToken(*user)
	if err != nil {
		t.Error(err)
	}

	t.Run("/challenge/create", func(t *testing.T) {
		dto := entities.ChallengeCreateWithTestcaseDTO{
			ChallengeCreateDTO: entities.ChallengeCreateDTO{
				Name:        "Test Challenge",
				Description: "Test Description",
			},
			Testcases: []entities.ChallengeTestcaseCreateDTO{
				{Input: "1 2", ExpectedOutput: "3"},
				{Input: "2 3", ExpectedOutput: "5"},
			},
		}
		requestBody, _ := json.Marshal(dto)

		request, _ := http.NewRequest(http.MethodPost, "/challenge/create", bytes.NewBuffer(requestBody))
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Authorization", "Bearer "+userAccessToken)

		response, err := app.Test(request, -1)
		if err != nil {
			t.Error(err)
		}

		if response.StatusCode != http.StatusOK {
			t.Errorf("Expected status OK, got %v", response.StatusCode)
		}
	})

	t.Run("/challenge/all", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/challenge/all", nil)
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Authorization", "Bearer "+userAccessToken)

		response, err := app.Test(request, -1)
		if err != nil {
			t.Error(err)
		}

		if response.StatusCode != http.StatusOK {
			t.Errorf("Expected status OK, got %v", response.StatusCode)
		}
	})

	t.Run("/challenge/get/:id", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/challenge/get/1", nil)
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Authorization", "Bearer "+userAccessToken)

		response, err := app.Test(request, -1)
		if err != nil {
			t.Error(err)
		}

		if response.StatusCode != http.StatusOK {
			t.Errorf("Expected status OK, got %v", response.StatusCode)
		}
	})

	t.Run("/challenge/update", func(t *testing.T) {
		dto := entities.ChallengeUpdateDTO{
			ChallengeID: 1,
			Name:        "Test Update Challenge",
			Description: "Test Update Description",
		}
		requestBody, _ := json.Marshal(dto)

		request, _ := http.NewRequest(http.MethodPut, "/challenge/update", bytes.NewBuffer(requestBody))
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Authorization", "Bearer "+userAccessToken)

		response, err := app.Test(request, -1)
		if err != nil {
			t.Error(err)
		}

		if response.StatusCode != http.StatusOK {
			t.Errorf("Expected status OK, got %v", response.StatusCode)
		}

		updatedChallenge, err := serviceKit.ChallengeService.FindChallengeByID(1)
		if err != nil {
			t.Error(err)
		}
		if updatedChallenge.Name != dto.Name {
			t.Errorf("Expected name %v, got %v", dto.Name, updatedChallenge.Name)
		}
		if updatedChallenge.Description != dto.Description {
			t.Errorf("Expected description %v, got %v", dto.Description, updatedChallenge.Description)
		}
	})

	t.Run("/challenge/delete/:id", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodDelete, "/challenge/delete/1", nil)
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Authorization", "Bearer "+userAccessToken)

		response, err := app.Test(request, -1)
		if err != nil {
			t.Error(err)
		}

		if response.StatusCode != http.StatusOK {
			t.Errorf("Expected status OK, got %v", response.StatusCode)
		}

		// get challenge by id must return error
		_, err = serviceKit.ChallengeService.FindChallengeByID(1)
		if err == nil {
			t.Error("Expected error, got nil")
		}
	})
}
