package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/wuttinanhi/code-judge-system/controllers"
	"github.com/wuttinanhi/code-judge-system/entities"
	"github.com/wuttinanhi/code-judge-system/services"
)

func TestChallengeRoute(t *testing.T) {
	testServiceKit := services.CreateTestServiceKit()
	app := controllers.SetupWeb(testServiceKit)

	// create admin user
	adminUser, err := testServiceKit.UserService.Register("admin@example.com", "testpassword", "admin")
	if err != nil {
		t.Error(err)
	}

	// set user role to admin
	err = testServiceKit.UserService.UpdateRole(adminUser, entities.UserRoleAdmin)
	if err != nil {
		t.Error(err)
	}

	adminAccessToken, err := testServiceKit.JWTService.GenerateToken(*adminUser)
	if err != nil {
		t.Error(err)
	}

	// create user
	user, err := testServiceKit.UserService.Register("user@example.com", "testpassword", "user")
	if err != nil {
		t.Error(err)
	}

	userAccessToken, err := testServiceKit.JWTService.GenerateToken(*user)
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
		request.Header.Set("Authorization", "Bearer "+adminAccessToken)

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
		request.Header.Set("Authorization", "Bearer "+adminAccessToken)

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
		request.Header.Set("Authorization", "Bearer "+adminAccessToken)

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
		request.Header.Set("Authorization", "Bearer "+adminAccessToken)

		response, err := app.Test(request, -1)
		if err != nil {
			t.Error(err)
		}

		if response.StatusCode != http.StatusOK {
			t.Errorf("Expected status OK, got %v", response.StatusCode)
		}

		updatedChallenge, err := testServiceKit.ChallengeService.FindChallengeByID(1)
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

	t.Run("/challenge/delete/:id user should not be able to delete challenge", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodDelete, "/challenge/delete/1", nil)
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Authorization", "Bearer "+userAccessToken)

		response, err := app.Test(request, -1)
		if err != nil {
			t.Error(err)
		}

		// expect forbidden status code
		if response.StatusCode != http.StatusForbidden {
			t.Errorf("Expected status Forbidden, got %v", response.StatusCode)
		}
	})

	t.Run("/challenge/delete/:id", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodDelete, "/challenge/delete/1", nil)
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Authorization", "Bearer "+adminAccessToken)

		response, err := app.Test(request, -1)
		if err != nil {
			t.Error(err)
		}

		if response.StatusCode != http.StatusOK {
			t.Errorf("Expected status OK, got %v", response.StatusCode)
		}

		// get challenge by id must return error
		_, err = testServiceKit.ChallengeService.FindChallengeByID(1)
		if err == nil {
			t.Error("Expected error, got nil")
		}
	})

	t.Run("/challenge/create user should not be able to create challenge", func(t *testing.T) {
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

		// expect forbidden status code
		if response.StatusCode != http.StatusForbidden {
			t.Errorf("Expected status Forbidden, got %v", response.StatusCode)
		}
	})

}
