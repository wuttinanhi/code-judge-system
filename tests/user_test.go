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

func TestAuthRoutes(t *testing.T) {
	db := databases.NewTempSQLiteDatabase()
	testServiceKit := services.CreateServiceKit(db)
	app := controllers.SetupAPI(testServiceKit)

	t.Run("/auth/register", func(t *testing.T) {
		dto := entities.UserRegisterDTO{
			DisplayName: "Test User",
			Email:       "testuser@example.com",
			Password:    "testpassword",
		}
		requestBody, _ := json.Marshal(dto)

		request, _ := http.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(requestBody))
		request.Header.Set("Content-Type", "application/json")

		response, err := app.Test(request, -1)
		if err != nil {
			t.Error(err)
		}

		if response.StatusCode != http.StatusOK {
			t.Errorf("Expected status OK, got %v", response.StatusCode)
		}
	})

	t.Run("/auth/login", func(t *testing.T) {
		dto := entities.UserLoginDTO{
			Email:    "testuser@example.com",
			Password: "testpassword",
		}
		requestBody, _ := json.Marshal(dto)

		request, _ := http.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(requestBody))
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
