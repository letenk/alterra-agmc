package tests

import (
	"bookmarket/config"
	"bookmarket/lib"
	"bookmarket/models"
	"bookmarket/routes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/rizkydarmawan-letenk/jabufaker"
	"github.com/stretchr/testify/assert"
)

// LoginRandomAccount for login with account random that return the token
func LoginRandomAccount(t *testing.T) (models.Users, string) {
	// Create random data use function CreateRandomUser on file user controller test
	newUser := CreateRandomUser(t)

	// Passing data newUser to model user
	user := models.Users(newUser)

	db := config.SetupTestDB()
	// Use router
	router := routes.SetupRouter(db)

	// Create random dummy data
	data := lib.LoginPayload{
		Email:    newUser.Email,
		Password: newUser.Password,
	}

	// String json
	dataBody := fmt.Sprintf(`{"email": "%s", "password": "%s"}`, data.Email, data.Password)

	// Create new reader
	requestBody := strings.NewReader(dataBody)
	// Create request
	request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/v1/login", requestBody)
	// Added header content type
	request.Header.Add("Content-Type", "application/json")

	// Create recorder
	recorder := httptest.NewRecorder()

	// Run handler
	router.ServeHTTP(recorder, request)

	// Get response
	response := recorder.Result()
	//  Read body response
	body, _ := io.ReadAll(response.Body)
	var responsebody map[string]interface{}
	json.Unmarshal(body, &responsebody)

	// Test pass
	// Response status code
	assert.Equal(t, 200, response.StatusCode)
	// Response body status code
	assert.Equal(t, 200, int(responsebody["code"].(float64)))
	// Response body status
	assert.Equal(t, "success", responsebody["status"])
	// Response body message
	assert.Equal(t, "You are logged", responsebody["message"])
	// Response body data, property token
	token := responsebody["data"].(map[string]interface{})["token"]
	assert.NotEmpty(t, token)

	// Return user and token, for use by any function that requires authentication
	return user, token.(string)
}

// Test login success
func TestLoginSuccess(t *testing.T) {
	LoginRandomAccount(t)
}

// Test Login failed, user or password incorrect.
func TestLoginOrPasswordIncorrect(t *testing.T) {
	db := config.SetupTestDB()
	// Use router
	router := routes.SetupRouter(db)

	// String json
	dataBody := fmt.Sprintf(`{"email": "%s", "password": "%s"}`, jabufaker.RandomEmail(), jabufaker.RandomString(5))

	// Create new reader
	requestBody := strings.NewReader(dataBody)
	// Create request
	request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/v1/login", requestBody)
	// Added header content type
	request.Header.Add("Content-Type", "application/json")

	// Create recorder
	recorder := httptest.NewRecorder()

	// Run handler
	router.ServeHTTP(recorder, request)

	// Get response
	response := recorder.Result()
	//  Read body response
	body, _ := io.ReadAll(response.Body)
	var responsebody map[string]interface{}
	json.Unmarshal(body, &responsebody)

	// Test pass
	// Response status code
	assert.Equal(t, 400, response.StatusCode)
	// Response body status code
	assert.Equal(t, 400, int(responsebody["code"].(float64)))
	// Response body status
	assert.Equal(t, "error", responsebody["status"])
	// Response body message
	assert.Equal(t, "user or password incorrect", responsebody["message"])
}

// Test Login failed, validation error
func TestLoginValidationError(t *testing.T) {
	db := config.SetupTestDB()
	// Use router
	router := routes.SetupRouter(db)

	// String json
	dataBody := fmt.Sprintf(`{"email": "%s", "password": "%s"}`, "wrongtest.com", "")

	// Create new reader
	requestBody := strings.NewReader(dataBody)
	// Create request
	request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/v1/login", requestBody)
	// Added header content type
	request.Header.Add("Content-Type", "application/json")

	// Create recorder
	recorder := httptest.NewRecorder()

	// Run handler
	router.ServeHTTP(recorder, request)

	// Get response
	response := recorder.Result()
	//  Read body response
	body, _ := io.ReadAll(response.Body)
	var responsebody map[string]interface{}
	json.Unmarshal(body, &responsebody)

	// Test pass
	// Response status code
	assert.Equal(t, 400, response.StatusCode)
	// Response body status code
	assert.Equal(t, 400, int(responsebody["code"].(float64)))
	// Response body status
	assert.Equal(t, "error", responsebody["status"])
	// Response body message
	assert.Equal(t, "bad request", responsebody["message"])
	assert.NotEmpty(t, responsebody["data"].(map[string]interface{})["message"])
}
