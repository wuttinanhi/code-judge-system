package tests_test

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

// Dedicated test for challenge update
func TestChallengeUpdate(t *testing.T) {
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

	// create new challenge
	_, err = testServiceKit.ChallengeService.CreateChallenge(&entities.Challenge{
		Name:        "Test Challenge",
		Description: "Test Description",
		User:        adminUser,
		Testcases: []*entities.ChallengeTestcase{
			{Input: "1", ExpectedOutput: "1", LimitMemory: 1, LimitTimeMs: 1},
			{Input: "2", ExpectedOutput: "2", LimitMemory: 2, LimitTimeMs: 2},
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	dto := entities.ChallengeUpdateDTO{
		Name:        "Test Challenge Updated",
		Description: "Test Description Updated",
		Testcases: []entities.ChallengeTestcaseDTO{
			{ID: 1, Input: "1", ExpectedOutput: "1", LimitMemory: 1, LimitTimeMs: 1, Action: "update"},
			{ID: 2, Input: "2", ExpectedOutput: "2", LimitMemory: 2, LimitTimeMs: 2, Action: "delete"},
			{ID: 0, Input: "3", ExpectedOutput: "3", LimitMemory: 3, LimitTimeMs: 3, Action: "create"},
		},
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

	// expect first testcase to be updated
	if updatedChallenge.Testcases[0].Input != dto.Testcases[0].Input {
		t.Errorf("Expected input %v, got %v", dto.Testcases[0].Input, updatedChallenge.Testcases[0].Input)
	}
	if updatedChallenge.Testcases[0].ExpectedOutput != dto.Testcases[0].ExpectedOutput {
		t.Errorf("Expected expected output %v, got %v", dto.Testcases[0].ExpectedOutput, updatedChallenge.Testcases[0].ExpectedOutput)
	}
	if updatedChallenge.Testcases[0].LimitMemory != dto.Testcases[0].LimitMemory {
		t.Errorf("Expected limit memory %v, got %v", dto.Testcases[0].LimitMemory, updatedChallenge.Testcases[0].LimitMemory)
	}
	if updatedChallenge.Testcases[0].LimitTimeMs != dto.Testcases[0].LimitTimeMs {
		t.Errorf("Expected limit time %v, got %v", dto.Testcases[0].LimitTimeMs, updatedChallenge.Testcases[0].LimitTimeMs)
	}

	// expect testcase 3 to be created
	if updatedChallenge.Testcases[1].Input != dto.Testcases[2].Input {
		t.Errorf("Expected input %v, got %v", dto.Testcases[2].Input, updatedChallenge.Testcases[1].Input)
	}
	if updatedChallenge.Testcases[1].ExpectedOutput != dto.Testcases[2].ExpectedOutput {
		t.Errorf("Expected expected output %v, got %v", dto.Testcases[2].ExpectedOutput, updatedChallenge.Testcases[1].ExpectedOutput)
	}
	if updatedChallenge.Testcases[1].LimitMemory != dto.Testcases[2].LimitMemory {
		t.Errorf("Expected limit memory %v, got %v", dto.Testcases[2].LimitMemory, updatedChallenge.Testcases[1].LimitMemory)
	}
	if updatedChallenge.Testcases[1].LimitTimeMs != dto.Testcases[2].LimitTimeMs {
		t.Errorf("Expected limit time %v, got %v", dto.Testcases[2].LimitTimeMs, updatedChallenge.Testcases[1].LimitTimeMs)
	}
}

func TestChallengeUpdateTestcaseLimit(t *testing.T) {
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

	// create new challenge
	_, err = testServiceKit.ChallengeService.CreateChallenge(&entities.Challenge{
		Name:        "Test Challenge",
		Description: "Test Description",
		User:        adminUser,
		Testcases:   []*entities.ChallengeTestcase{},
	})
	if err != nil {
		t.Fatal(err)
	}

	// create 101 testcases
	testcases := []entities.ChallengeTestcaseDTO{}
	for i := 0; i < 101; i++ {
		testcases = append(testcases, entities.ChallengeTestcaseDTO{
			ID:             0,
			Input:          "INPUT",
			ExpectedOutput: "EXPECTED OUTPUT",
			LimitMemory:    uint(1),
			LimitTimeMs:    uint(1),
			Action:         "create",
		})
	}

	dto := entities.ChallengeUpdateDTO{
		Name:        "Test Challenge",
		Description: "Test Description",
		Testcases:   testcases,
	}
	requestBody, _ := json.Marshal(dto)

	request, _ := http.NewRequest(http.MethodPut, "/challenge/update/1", bytes.NewBuffer(requestBody))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+adminAccessToken)

	response, err := app.Test(request, -1)
	if err != nil {
		t.Error(err)
	}
	if response.StatusCode == http.StatusTooManyRequests {
		t.Errorf("Expected status StatusTooManyRequests, got %v", response.StatusCode)
	}
}

func TestChallengeUpdateSandboxLimit(t *testing.T) {
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

	// create new challenge
	_, err = testServiceKit.ChallengeService.CreateChallenge(&entities.Challenge{
		Name:        "Test Challenge",
		Description: "Test Description",
		User:        adminUser,
		Testcases:   []*entities.ChallengeTestcase{},
	})
	if err != nil {
		t.Fatal(err)
	}

	testcases := []entities.ChallengeTestcaseDTO{}
	testcases = append(testcases, entities.ChallengeTestcaseDTO{
		ID:             0,
		Input:          "INPUT",
		ExpectedOutput: "EXPECTED OUTPUT",
		LimitMemory:    entities.SandboxMemoryGB * 1,
		LimitTimeMs:    99999,
		Action:         "create",
	})
	testcases = append(testcases, entities.ChallengeTestcaseDTO{
		ID:             0,
		Input:          "INPUT",
		ExpectedOutput: "EXPECTED OUTPUT",
		LimitMemory:    entities.SandboxMemoryGB * 1,
		LimitTimeMs:    99999,
		Action:         "create",
	})

	dto := entities.ChallengeUpdateDTO{
		Name:        "Test Update",
		Description: "Test Update",
		Testcases:   testcases,
	}
	requestBody, _ := json.Marshal(dto)

	request, _ := http.NewRequest(http.MethodPut, "/challenge/update/1", bytes.NewBuffer(requestBody))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+adminAccessToken)

	response, err := app.Test(request, -1)
	if err != nil {
		t.Error(err)
	}
	if response.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status %v, got %v", http.StatusBadRequest, response.StatusCode)
	}
}
