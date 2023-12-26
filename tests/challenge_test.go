package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/wuttinanhi/code-judge-system/controllers"
	"github.com/wuttinanhi/code-judge-system/databases"
	"github.com/wuttinanhi/code-judge-system/entities"
	"github.com/wuttinanhi/code-judge-system/services"
)

func TestChallengeRoute(t *testing.T) {
	db := databases.NewTempSQLiteDatabase()
	testServiceKit := services.CreateServiceKit(db)
	app := controllers.SetupWeb(testServiceKit)

	// create admin user
	adminUser, err := testServiceKit.UserService.Register("admin@example.com", "testpassword", "admin")
	if err != nil {
		t.Fatal(err)
	}

	// set user role to admin
	err = testServiceKit.UserService.UpdateRole(adminUser, entities.UserRoleAdmin)
	if err != nil {
		t.Fatal(err)
	}

	adminAccessToken, err := testServiceKit.JWTService.GenerateToken(*adminUser)
	if err != nil {
		t.Fatal(err)
	}

	// create user
	user, err := testServiceKit.UserService.Register("user@example.com", "testpassword", "user")
	if err != nil {
		t.Fatal(err)
	}

	userAccessToken, err := testServiceKit.JWTService.GenerateToken(*user)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("/challenge/create", func(t *testing.T) {
		dto := entities.ChallengeCreateWithTestcaseDTO{
			Name:        "Test Challenge",
			Description: "Test Description",
			Testcases: []entities.ChallengeTestcaseDTO{
				{Input: "1 2", ExpectedOutput: "3", LimitMemory: 1, LimitTimeMs: 1},
				{Input: "2 3", ExpectedOutput: "5", LimitMemory: 2, LimitTimeMs: 2},
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

	t.Run("/challenge/pagination", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/challenge/pagination", nil)
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Authorization", "Bearer "+adminAccessToken)

		response, err := app.Test(request, -1)
		if err != nil {
			t.Error(err)
		}

		if response.StatusCode != http.StatusOK {
			t.Errorf("Expected status OK, got %v", response.StatusCode)
		}

		// try parse json response to pagination result
		var result entities.PaginationResult[entities.ChallengeExtended]
		err = json.NewDecoder(response.Body).Decode(&result)
		if err != nil {
			t.Error(err)
		}

		// expect total 1
		if result.Total != 1 {
			t.Errorf("Expected total 1, got %v", result.Total)
		}

		// expect challenge name to be Test Challenge
		if result.Items[0].Name != "Test Challenge" {
			t.Errorf("Expected challenge name Test Challenge, got %v", result.Items[0].Name)
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
			Name:        "Test Challenge Updated",
			Description: "Test Description Updated",
		}
		requestBody, _ := json.Marshal(dto)

		request, _ := http.NewRequest(http.MethodPut, "/challenge/update/1", bytes.NewBuffer(requestBody))
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

		// expect 2 testcases
		if len(updatedChallenge.Testcases) != 2 {
			t.Errorf("Expected 2 testcases, got %v", len(updatedChallenge.Testcases))
		}
	})

	t.Run("/challenge/create user should not be able to create challenge", func(t *testing.T) {
		dto := entities.ChallengeCreateWithTestcaseDTO{
			Name:        "Test Challenge",
			Description: "Test Description",
			Testcases: []entities.ChallengeTestcaseDTO{
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

		// challenge count should be 0
		challenges, err := testServiceKit.ChallengeService.AllChallenges()
		if err != nil {
			t.Error(err)
		}
		if len(challenges) != 0 {
			t.Errorf("Expected 0 challenge, got %v", len(challenges))
		}
	})
}
