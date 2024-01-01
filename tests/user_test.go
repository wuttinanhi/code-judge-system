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

func TestUserRoute(t *testing.T) {

	t.Run("admin can update user role", func(t *testing.T) {
		db := databases.NewTempSQLiteDatabase()
		testServiceKit := services.CreateTestServiceKit(db)
		rateLimitStorage := controllers.GetMemoryStorage()
		app := controllers.SetupAPI(testServiceKit, rateLimitStorage)
		adminToken, _ := createUserWrapper(t, testServiceKit)

		dto := entities.UserUpdateRoleDTO{
			UserID: 2,
			Role:   entities.UserRoleStaff,
		}
		requestBody, _ := json.Marshal(dto)

		request, _ := http.NewRequest(http.MethodPut, "/user/update/role", bytes.NewBuffer(requestBody))
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Authorization", "Bearer "+adminToken)

		response, err := app.Test(request, -1)
		if err != nil {
			t.Error(err)
		}

		if response.StatusCode != http.StatusOK {
			t.Errorf("Expected status OK, got %v", response.StatusCode)
		}
	})

	t.Run("normal user cannot update role", func(t *testing.T) {
		db := databases.NewTempSQLiteDatabase()
		testServiceKit := services.CreateTestServiceKit(db)
		rateLimitStorage := controllers.GetMemoryStorage()
		app := controllers.SetupAPI(testServiceKit, rateLimitStorage)
		_, normalToken := createUserWrapper(t, testServiceKit)

		dto := entities.UserUpdateRoleDTO{
			UserID: 2,
			Role:   entities.UserRoleStaff,
		}
		requestBody, _ := json.Marshal(dto)

		request, _ := http.NewRequest(http.MethodPut, "/user/update/role", bytes.NewBuffer(requestBody))
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Authorization", "Bearer "+normalToken)

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

func createUserWrapper(t *testing.T, testServiceKit *services.ServiceKit) (string, string) {
	adminUser, err := testServiceKit.UserService.Register("admin@example.com", "adminpassword", "admin")
	if err != nil {
		t.Fatal(err)
	}

	err = testServiceKit.UserService.UpdateRole(adminUser, entities.UserRoleAdmin)
	if err != nil {
		t.Fatal(err)
	}

	adminToken, err := testServiceKit.JWTService.GenerateToken(*adminUser)
	if err != nil {
		t.Fatal(err)
	}

	normalUser, err := testServiceKit.UserService.Register("test@example.com", "testpassword", "test")
	if err != nil {
		t.Fatal(err)
	}

	normalToken, err := testServiceKit.JWTService.GenerateToken(*normalUser)
	if err != nil {
		t.Fatal(err)
	}

	return adminToken, normalToken
}
