package scripts_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/wuttinanhi/code-judge-system/configs"
	"github.com/wuttinanhi/code-judge-system/controllers"
	"github.com/wuttinanhi/code-judge-system/databases"
	"github.com/wuttinanhi/code-judge-system/entities"
	"github.com/wuttinanhi/code-judge-system/services"
	"github.com/wuttinanhi/code-judge-system/tests"
)

func TestSubmissionCreate(t *testing.T) {
	configs.LoadConfig()

	db := databases.NewMySQLDatabase()
	testServiceKit := services.CreateServiceKit(db)
	testServiceKit.KafkaService.OverriddenHost("localhost:9094")

	app := controllers.SetupAPI(testServiceKit)

	// get a challenge
	challenge, err := testServiceKit.ChallengeService.FindChallengeByID(1)
	if err != nil {
		t.Fatal(err)
	}

	// get a user
	user, err := testServiceKit.UserService.FindUserByID(1)
	if err != nil {
		t.Error(err)
	}

	// get user access token
	userAccessToken, err := testServiceKit.JWTService.GenerateToken(*user)
	if err != nil {
		t.Error(err)
	}

	dto := entities.SubmissionCreateDTO{
		ChallengeID: challenge.ID,
		Language:    "python",
		// for testing wrong answer please add
		// + "\nprint('')"
		SourceCode: entities.PythonCodeExample,
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

	fmt.Println(tests.ResponseBodyToString(response))
}
