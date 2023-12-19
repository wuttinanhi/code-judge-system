package tests_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/wuttinanhi/code-judge-system/controllers"
	"github.com/wuttinanhi/code-judge-system/entities"
	"github.com/wuttinanhi/code-judge-system/services"
)

func TestUserRoutes(t *testing.T) {
	testServiceKit := services.CreateTestServiceKit()
	app := controllers.SetupWeb(testServiceKit)

	t.Run("/user/register", func(t *testing.T) {
		dto := entities.UserRegisterDTO{
			DisplayName: "Test User",
			Email:       "testuser@example.com",
			Password:    "testpassword",
		}
		requestBody, _ := json.Marshal(dto)

		request, _ := http.NewRequest(http.MethodPost, "/user/register", bytes.NewBuffer(requestBody))
		request.Header.Set("Content-Type", "application/json")

		response, err := app.Test(request, -1)
		if err != nil {
			t.Error(err)
		}

		if response.StatusCode != http.StatusOK {
			t.Errorf("Expected status OK, got %v", response.StatusCode)
		}
	})

	t.Run("/user/login", func(t *testing.T) {
		dto := entities.UserLoginDTO{
			Email:    "testuser@example.com",
			Password: "testpassword",
		}
		requestBody, _ := json.Marshal(dto)

		request, _ := http.NewRequest(http.MethodPost, "/user/login", bytes.NewBuffer(requestBody))
		request.Header.Set("Content-Type", "application/json")

		response, err := app.Test(request, -1)
		if err != nil {
			t.Error(err)
		}

		if response.StatusCode != http.StatusOK {
			t.Errorf("Expected status OK, got %v", response.StatusCode)
		}
	})
}
