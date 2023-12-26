package scripts_test

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

func TestChallengeCreate(t *testing.T) {
	db := databases.NewMySQLDatabase()
	testServiceKit := services.CreateServiceKit(db)
	testServiceKit.KafkaService.OverriddenHost("localhost:9094")
	app := controllers.SetupWeb(testServiceKit)

	db.Migrator().DropTable(
		&entities.Challenge{},
		entities.ChallengeTestcase{},
		entities.Submission{},
		entities.SubmissionTestcase{},
	)

	databases.StartMigration(db)

	// get a user
	user, err := testServiceKit.UserService.FindUserByID(1)
	if err != nil {
		t.Error(err)
	}

	// get user access token
	accessToken, err := testServiceKit.JWTService.GenerateToken(*user)
	if err != nil {
		t.Error(err)
	}

	dto := entities.ChallengeCreateWithTestcaseDTO{
		Name:        "Test Challenge",
		Description: "Test Description",
		Testcases: []entities.ChallengeTestcaseDTO{
			{Input: "1\n2\n", ExpectedOutput: "3\n", LimitMemory: 268435456, LimitTimeMs: 1000},
			{Input: "2\n3\n", ExpectedOutput: "5\n", LimitMemory: 268435456, LimitTimeMs: 1000},
		},
	}
	requestBody, _ := json.Marshal(dto)

	request, _ := http.NewRequest(http.MethodPost, "/challenge/create", bytes.NewBuffer(requestBody))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+accessToken)

	response, err := app.Test(request, -1)
	if err != nil {
		t.Error(err)
	}

	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %v", response.StatusCode)
	}
}
