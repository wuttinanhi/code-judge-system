package tests_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/wuttinanhi/code-judge-system/controllers"
	"github.com/wuttinanhi/code-judge-system/databases"
	"github.com/wuttinanhi/code-judge-system/entities"
	"github.com/wuttinanhi/code-judge-system/services"
)

func TestChallengeRoute(t *testing.T) {
	db := databases.NewTempSQLiteDatabase()
	testServiceKit := services.CreateTestServiceKit(db)
	rateLimitStorage := controllers.GetMemoryStorage()
	app := controllers.SetupAPI(testServiceKit, rateLimitStorage)

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

	t.Run("/challenge/create with sandbox limit", func(t *testing.T) {
		dto := entities.ChallengeCreateWithTestcaseDTO{
			Name:        "Test Challenge",
			Description: "Test Description",
			Testcases: []entities.ChallengeTestcaseDTO{
				{Input: "1 2", ExpectedOutput: "3", LimitMemory: entities.SandboxMemoryGB * 1, LimitTimeMs: 99999},
				{Input: "2 3", ExpectedOutput: "5", LimitMemory: entities.SandboxMemoryGB * 1, LimitTimeMs: 99999},
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

		// expect bad request status code
		if response.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status InternalServerError, got %v", response.StatusCode)
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

func challengeCreateWrapper(app *fiber.App, userCreateToken string) (*http.Response, error) {
	dto := entities.ChallengeCreateWithTestcaseDTO{
		Name:        "Test Challenge",
		Description: "Test Description",
		Testcases: []entities.ChallengeTestcaseDTO{
			{ID: 0,
				Input:          "INPUT",
				ExpectedOutput: "EXPECTED_OUTPUT",
				LimitMemory:    1,
				LimitTimeMs:    1,
				Action:         "create",
			},
		},
	}
	requestBody, _ := json.Marshal(dto)

	request, _ := http.NewRequest(http.MethodPost, "/challenge/create", bytes.NewBuffer(requestBody))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+userCreateToken)

	response, err := app.Test(request, -1)
	return response, err
}

func TestChallengeCreateLimit(t *testing.T) {
	db := databases.NewTempSQLiteDatabase()
	testServiceKit := services.CreateTestServiceKit(db)
	rateLimitStorage := controllers.GetMemoryStorage()
	app := controllers.SetupAPI(testServiceKit, rateLimitStorage)

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

	// generate admin access token
	adminAccessToken, err := testServiceKit.JWTService.GenerateToken(*adminUser)
	if err != nil {
		t.Fatal(err)
	}

	// expect first challenge create 100 to be success
	for i := 0; i < 200; i++ {
		response, err := challengeCreateWrapper(app, adminAccessToken)
		if err != nil {
			t.Error(err)
		}

		if i < 100 {
			if response.StatusCode != http.StatusOK {
				t.Errorf("Expected status OK, got %v", response.StatusCode)
			}
		} else {
			if response.StatusCode != http.StatusTooManyRequests {
				t.Errorf("Expected status TooManyRequests, got %v", response.StatusCode)
			}
		}
	}
}
