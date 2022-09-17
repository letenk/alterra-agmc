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
	"time"

	"github.com/google/uuid"
	"github.com/rizkydarmawan-letenk/jabufaker"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

// CreateRandomUser as function create random account with test and will be used by any function
func CreateRandomUser(t *testing.T) models.Users {
	db := config.SetupTestDB()
	// Delete all data from users
	db.Exec("DELETE FROM users")

	// Use router
	router := routes.SetupRouter(db)

	// Create random dummy data
	data := lib.CreateUser{
		Fullname: jabufaker.RandomPerson(),
		Email:    jabufaker.RandomEmail(),
		Password: "password",
	}

	// String json
	dataBody := fmt.Sprintf(`{"fullname": "%s", "email": "%s", "password": "%s"}`, data.Fullname, data.Email, data.Password)

	// Create new reader
	requestBody := strings.NewReader(dataBody)
	// Create request
	request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/v1/users", requestBody)
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
	assert.Equal(t, 201, response.StatusCode)
	// Response body status code
	assert.Equal(t, 201, int(responsebody["code"].(float64)))
	// Response body status
	assert.Equal(t, "success", responsebody["status"])
	// Response body message
	assert.Equal(t, "User has been created", responsebody["message"])

	// Response body data
	assert.NotEmpty(t, responsebody["data"])

	// Context property data on body
	var contextData = responsebody["data"].(map[string]interface{})
	assert.NotEmpty(t, contextData["id"])
	assert.NotEmpty(t, contextData["created_at"])
	assert.NotEmpty(t, contextData["updated_at"])
	assert.Equal(t, data.Fullname, contextData["fullname"])
	assert.Equal(t, data.Email, contextData["email"])

	// Compare password
	err := bcrypt.CompareHashAndPassword([]byte(contextData["password"].(string)), []byte(data.Password))
	assert.Nil(t, err)

	// Return new user, for re-use any function test
	newUser := models.Users{
		ID:        contextData["id"].(string),
		Fullname:  contextData["fullname"].(string),
		Email:     contextData["email"].(string),
		Password:  data.Password,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	return newUser
}

// Create user success
func TestCreateUserSuccess(t *testing.T) {
	CreateRandomUser(t)
}

// Create user failed email already exist
func TestCreateUserEmailExist(t *testing.T) {
	// Create random user
	user := CreateRandomUser(t)

	db := config.SetupTestDB()
	// Use router
	router := routes.SetupRouter(db)

	// Create random dummy data
	data := lib.CreateUser{
		Fullname: jabufaker.RandomPerson(),
		Email:    user.Email, // Insert same email with new random user created
		Password: "password",
	}

	// String json
	dataBody := fmt.Sprintf(`{"fullname": "%s", "email": "%s", "password": "%s"}`, data.Fullname, data.Email, data.Password)

	// Create new reader
	requestBody := strings.NewReader(dataBody)
	// Create request
	request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/v1/users", requestBody)
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
	assert.Equal(t, "Create user failed", responsebody["message"])

	// Response body data error
	assert.Equal(t, "email already exist", responsebody["data"].(map[string]interface{})["errors"])
}

// Create user failed validation error
func TestCreatesUserValidationError(t *testing.T) {

	db := config.SetupTestDB()
	// Use router
	router := routes.SetupRouter(db)

	// String json
	dataBody := fmt.Sprintf(`{"fullname": "%s", "email": "%s", "password": "%s"}`, "", "emailtest.com", "")
	// Create new reader
	requestBody := strings.NewReader(dataBody)
	// Create request
	request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/v1/users", requestBody)
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
	assert.Equal(t, "Create user failed", responsebody["message"])
	// Response body data message
	assert.NotEmpty(t, responsebody["data"].(map[string]interface{})["message"])
}

// Test Get all user success.
func TestGetAllUserSuccess(t *testing.T) {
	// Login with random account for get token
	_, token := LoginRandomAccount(t)

	db := config.SetupTestDB()
	// Use router
	router := routes.SetupRouter(db)

	// Create request
	request := httptest.NewRequest(http.MethodGet, "http://localhost:8080/v1/users", nil)
	// Added header content type
	request.Header.Add("Content-Type", "application/json")
	// Added header Authorization and added the bearer token
	bearerToken := fmt.Sprintf("Bearer %s", token)
	request.Header.Add("Authorization", bearerToken)

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
	assert.Equal(t, "list of users", responsebody["message"])

	// Iteration data
	data := responsebody["data"].([]interface{})

	assert.NotEqual(t, 0, len(data))
	for _, item := range data {
		list := item.(map[string]interface{})
		assert.NotEmpty(t, list["id"])
		assert.NotEmpty(t, list["fullname"])
		assert.NotEmpty(t, list["email"])
		assert.NotEmpty(t, list["created_at"])
		assert.NotEmpty(t, list["updated_at"])
	}
}

// Test Get all user with data empty slice, because data is not available.
func TestGetUserEmptyData(t *testing.T) {
	// Login with random account for get token
	_, token := LoginRandomAccount(t)

	db := config.SetupTestDB()
	db.Exec("DELETE FROM users")
	// Use router
	router := routes.SetupRouter(db)

	// Create request
	request := httptest.NewRequest(http.MethodGet, "http://localhost:8080/v1/users", nil)
	// Added header content type
	request.Header.Add("Content-Type", "application/json")
	// Added header Authorization and added the bearer token
	bearerToken := fmt.Sprintf("Bearer %s", token)
	request.Header.Add("Authorization", bearerToken)

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
	assert.Equal(t, "list of users", responsebody["message"])

	// Iteration data
	data := responsebody["data"].([]interface{})
	assert.Equal(t, 0, len(data))
}

// Test Get all user failed, unauthorized
func TestGetAllUserUnauthorized(t *testing.T) {
	t.Run("Invalid token or expired", func(t *testing.T) {
		// Create some new user
		CreateRandomUser(t)

		db := config.SetupTestDB()
		// Use router
		router := routes.SetupRouter(db)

		// Create request
		request := httptest.NewRequest(http.MethodGet, "http://localhost:8080/v1/users", nil)
		// Added header content type
		request.Header.Add("Content-Type", "application/json")
		// Added header Authorization and added the bearer token
		request.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiZjM4NjA0ODAtZDBhNi00ZTg0LTg2OWItMTAwZjA5MzQ4NTU2IiwiZXhwIjoxNjYzNDY1NDgwfQ.i-vDWixCqoBFy4zqqw1KBFNFL5JqrPKTlgybr6EUAIk")

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
		assert.Equal(t, 401, response.StatusCode)
		// Response message
		assert.Equal(t, "invalid or expired jwt", responsebody["message"])
	})

	t.Run("Do not have tokens", func(t *testing.T) {
		// Create some new user
		CreateRandomUser(t)

		db := config.SetupTestDB()
		// Use router
		router := routes.SetupRouter(db)

		// Create request
		request := httptest.NewRequest(http.MethodGet, "http://localhost:8080/v1/users", nil)
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
		// Response message
		assert.Equal(t, "missing or malformed jwt", responsebody["message"])
	})
}

// Test Get user find by id success (200).
func TestGetUserByIdSuccess(t *testing.T) {
	// Login with random account for get token
	_, token := LoginRandomAccount(t)

	// Create new user
	newUser := CreateRandomUser(t)

	db := config.SetupTestDB()
	// Use router
	router := routes.SetupRouter(db)
	// Create request
	request := httptest.NewRequest(http.MethodGet, "http://localhost:8080/v1/users/"+newUser.ID, nil)
	// Added header content type
	request.Header.Add("Content-Type", "application/json")
	// Added header Authorization and added the bearer token
	bearerToken := fmt.Sprintf("Bearer %s", token)
	request.Header.Add("Authorization", bearerToken)

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
	assert.Equal(t, "Data of user", responsebody["message"])

	// Iteration data
	data := responsebody["data"].(map[string]interface{})
	assert.Equal(t, newUser.ID, data["id"])
	assert.Equal(t, newUser.Fullname, data["fullname"])
	assert.Equal(t, newUser.Email, data["email"])
	assert.NotEmpty(t, newUser.CreatedAt)
	assert.NotEmpty(t, newUser.UpdatedAt)
}

// Test Get user find by id failed, user not found.
func TestGetUserByIdFailedUserNotFound(t *testing.T) {
	// Login with random account for get token
	_, token := LoginRandomAccount(t)

	db := config.SetupTestDB()
	// Use router
	router := routes.SetupRouter(db)

	// Var for wrong id
	wrongId := "0f2ed479-234a-4148-8acc-004ffcf58ffd"
	// Create request
	request := httptest.NewRequest(http.MethodGet, "http://localhost:8080/v1/users/"+wrongId, nil)
	// Added header content type
	request.Header.Add("Content-Type", "application/json")
	// Added header Authorization and added the bearer token
	bearerToken := fmt.Sprintf("Bearer %s", token)
	request.Header.Add("Authorization", bearerToken)

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
	message := fmt.Sprintf("user with id: %s not found", wrongId)
	assert.Equal(t, message, responsebody["message"])
}

// Test Get user find by id failed, unauthorized (400 and 401).
func TestGetUserByIdUnauthorized(t *testing.T) {
	t.Run("Invalid token or expired", func(t *testing.T) {
		// Create some new user
		newUser := CreateRandomUser(t)

		db := config.SetupTestDB()
		// Use router
		router := routes.SetupRouter(db)

		// Create request
		request := httptest.NewRequest(http.MethodGet, "http://localhost:8080/v1/users/"+newUser.ID, nil)
		// Added header content type
		request.Header.Add("Content-Type", "application/json")
		// Added header Authorization and added the bearer token
		request.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiZjM4NjA0ODAtZDBhNi00ZTg0LTg2OWItMTAwZjA5MzQ4NTU2IiwiZXhwIjoxNjYzNDY1NDgwfQ.i-vDWixCqoBFy4zqqw1KBFNFL5JqrPKTlgybr6EUAIk")

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
		assert.Equal(t, 401, response.StatusCode)
		// Response message
		assert.Equal(t, "invalid or expired jwt", responsebody["message"])
	})

	t.Run("Do not have tokens", func(t *testing.T) {
		// Create some new user
		newUser := CreateRandomUser(t)

		db := config.SetupTestDB()
		// Use router
		router := routes.SetupRouter(db)

		// Create request
		request := httptest.NewRequest(http.MethodGet, "http://localhost:8080/v1/users/"+newUser.ID, nil)
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
		// Response message
		assert.Equal(t, "missing or malformed jwt", responsebody["message"])
	})
}

// Test Update user success (200).
func TestUpdateUserSuccess(t *testing.T) {
	// Login with random account for get token
	newUser, token := LoginRandomAccount(t)

	db := config.SetupTestDB()
	// Use router
	router := routes.SetupRouter(db)

	// A data for use payload update
	dataUpdate := lib.UpdateUser{
		Fullname: "updated",
		Password: "updated",
	}

	// Create payload json
	dataBody := fmt.Sprintf(`{"fullname": "%s", "password": "%s"}`, dataUpdate.Fullname, dataUpdate.Password)
	// Create reader
	requestBody := strings.NewReader(dataBody)

	// Create request
	request := httptest.NewRequest(http.MethodPut, "http://localhost:8080/v1/users/"+newUser.ID, requestBody)
	// Added header content type
	request.Header.Add("Content-Type", "application/json")
	// Added string bearer
	bearerToken := fmt.Sprintf("Bearer %s", token)
	// Added header Authorization and added the bearer token
	request.Header.Add("Authorization", bearerToken)

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
	assert.Equal(t, "User has been updated", responsebody["message"])

	// Response body data
	assert.NotEmpty(t, responsebody["data"])

	// Context property data on body
	var contextData = responsebody["data"].(map[string]interface{})
	assert.Equal(t, newUser.ID, contextData["id"])
	assert.Equal(t, dataUpdate.Fullname, contextData["fullname"])
	assert.Equal(t, newUser.Email, contextData["email"])
	assert.NotEmpty(t, contextData["created_at"])
	assert.NotEmpty(t, contextData["updated_at"])
}

// Test Update user failed, not access. Because only can updating the data self him (400).
func TestUpdateUserFailedNotAccess(t *testing.T) {
	// Login with random account for get token
	_, token := LoginRandomAccount(t)

	// Create new user
	newUser := CreateRandomUser(t)

	db := config.SetupTestDB()
	// Use router
	router := routes.SetupRouter(db)

	// A data for use payload update
	dataUpdate := lib.UpdateUser{
		Fullname: "updated",
		Password: "updated",
	}

	// Create payload json
	dataBody := fmt.Sprintf(`{"fullname": "%s", "password": "%s"}`, dataUpdate.Fullname, dataUpdate.Password)
	// Create reader
	requestBody := strings.NewReader(dataBody)

	// Create request
	request := httptest.NewRequest(http.MethodPut, "http://localhost:8080/v1/users/"+newUser.ID, requestBody)
	// Added header content type
	request.Header.Add("Content-Type", "application/json")
	// Added header Authorization and added the bearer token
	bearerToken := fmt.Sprintf("Bearer %s", token)
	request.Header.Add("Authorization", bearerToken)

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
	assert.Equal(t, 401, response.StatusCode)
	// Response body status code
	assert.Equal(t, 401, int(responsebody["code"].(float64)))
	// Response body status
	assert.Equal(t, "error", responsebody["status"])
	// Response body message
	assert.Equal(t, "unauthorized", responsebody["message"])

	// Response body data
	assert.NotEmpty(t, responsebody["data"])

	// Context property data on body
	var contextData = responsebody["data"].(map[string]interface{})
	assert.Equal(t, "not access", contextData["errors"])
}

// Test Update user failed, validation error (400).
func TestUpdateUserFailedValidationError(t *testing.T) {
	// Login with random account for get token
	newUser, token := LoginRandomAccount(t)

	db := config.SetupTestDB()
	// Use router
	router := routes.SetupRouter(db)

	// A data for use payload update
	dataUpdate := lib.UpdateUser{
		Fullname: "",
		Password: "",
	}

	// Create payload json
	dataBody := fmt.Sprintf(`{"fullname": "%s", "password": "%s"}`, dataUpdate.Fullname, dataUpdate.Password)
	// Create reader
	requestBody := strings.NewReader(dataBody)

	// Create request
	request := httptest.NewRequest(http.MethodPut, "http://localhost:8080/v1/users/"+newUser.ID, requestBody)
	// Added header content type
	request.Header.Add("Content-Type", "application/json")
	// Added header Authorization and added the bearer token
	bearerToken := fmt.Sprintf("Bearer %s", token)
	request.Header.Add("Authorization", bearerToken)

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
	assert.Equal(t, "Create user failed", responsebody["message"])
	// Response body data message
	assert.NotEmpty(t, responsebody["data"].(map[string]interface{})["message"])
}

// Test Update user failed, unauthorized (400 and 401).
func TestUpdateUserUnauthorized(t *testing.T) {
	t.Run("Invalid token or expired", func(t *testing.T) {
		// Create some new user
		newUser := CreateRandomUser(t)

		db := config.SetupTestDB()
		// Use router
		router := routes.SetupRouter(db)

		// A data for use payload update
		dataUpdate := lib.UpdateUser{
			Fullname: "updated",
			Password: "updated",
		}
		// Create payload json
		dataBody := fmt.Sprintf(`{"fullname": "%s", "password": "%s"}`, dataUpdate.Fullname, dataUpdate.Password)
		// Create reader
		requestBody := strings.NewReader(dataBody)

		// Create request
		request := httptest.NewRequest(http.MethodPut, "http://localhost:8080/v1/users/"+newUser.ID, requestBody)
		// Added header content type
		request.Header.Add("Content-Type", "application/json")
		// Added header Authorization and added the bearer token
		request.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiZjM4NjA0ODAtZDBhNi00ZTg0LTg2OWItMTAwZjA5MzQ4NTU2IiwiZXhwIjoxNjYzNDY1NDgwfQ.i-vDWixCqoBFy4zqqw1KBFNFL5JqrPKTlgybr6EUAIk")

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
		assert.Equal(t, 401, response.StatusCode)
		// Response message
		assert.Equal(t, "invalid or expired jwt", responsebody["message"])
	})

	t.Run("Do not have tokens", func(t *testing.T) {
		// Create some new user
		newUser := CreateRandomUser(t)

		db := config.SetupTestDB()
		// Use router
		router := routes.SetupRouter(db)

		// A data for use payload update
		dataUpdate := lib.UpdateUser{
			Fullname: "updated",
			Password: "updated",
		}
		// Create payload json
		dataBody := fmt.Sprintf(`{"fullname": "%s", "password": "%s"}`, dataUpdate.Fullname, dataUpdate.Password)
		// Create reader
		requestBody := strings.NewReader(dataBody)

		// Create request
		request := httptest.NewRequest(http.MethodPut, "http://localhost:8080/v1/users/"+newUser.ID, requestBody)
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
		// Response message
		assert.Equal(t, "missing or malformed jwt", responsebody["message"])
	})
}

// Test Delete user success (200).
func TestDeleteUserSuccess(t *testing.T) {
	// Login with random account for get token
	newUser, token := LoginRandomAccount(t)

	db := config.SetupTestDB()
	// Use router
	router := routes.SetupRouter(db)

	// Create request
	request := httptest.NewRequest(http.MethodDelete, "http://localhost:8080/v1/users/"+newUser.ID, nil)
	// Added header content type
	request.Header.Add("Content-Type", "application/json")
	// Added string bearer
	bearerToken := fmt.Sprintf("Bearer %s", token)
	// Added header Authorization and added the bearer token
	request.Header.Add("Authorization", bearerToken)

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
	assert.Equal(t, "User has been deleted", responsebody["message"])
}

// Test Delete user failed, not access. Because only can deleting the data self him (400) .
func TestDeleteUserFailedNotAccess(t *testing.T) {
	// Login with random account for get token
	_, token := LoginRandomAccount(t)

	db := config.SetupTestDB()
	// Use router
	router := routes.SetupRouter(db)

	// Generate uuid, for sample id
	userID := uuid.NewString()
	// Create request
	request := httptest.NewRequest(http.MethodDelete, "http://localhost:8080/v1/users/"+userID, nil)
	// Added header content type
	request.Header.Add("Content-Type", "application/json")
	// Added header Authorization and added the bearer token
	bearerToken := fmt.Sprintf("Bearer %s", token)
	request.Header.Add("Authorization", bearerToken)

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
	assert.Equal(t, 401, response.StatusCode)
	// Response body status code
	assert.Equal(t, 401, int(responsebody["code"].(float64)))
	// Response body status
	assert.Equal(t, "error", responsebody["status"])
	// Response body message
	assert.Equal(t, "unauthorized", responsebody["message"])

	// Response body data
	assert.NotEmpty(t, responsebody["data"])

	// Context property data on body
	var contextData = responsebody["data"].(map[string]interface{})
	assert.Equal(t, "not access", contextData["errors"])
}

// Test Delete user failed, unauthorized (400 and 401).
func TestDeleteUserUnauthorized(t *testing.T) {
	t.Run("Invalid token or expired", func(t *testing.T) {
		db := config.SetupTestDB()
		// Use router
		router := routes.SetupRouter(db)

		// Generate uuid, for sample id
		userID := uuid.NewString()

		// Create request
		request := httptest.NewRequest(http.MethodDelete, "http://localhost:8080/v1/users/"+userID, nil)
		// Added header content type
		request.Header.Add("Content-Type", "application/json")
		// Added header Authorization and added the bearer token
		request.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiZjM4NjA0ODAtZDBhNi00ZTg0LTg2OWItMTAwZjA5MzQ4NTU2IiwiZXhwIjoxNjYzNDY1NDgwfQ.i-vDWixCqoBFy4zqqw1KBFNFL5JqrPKTlgybr6EUAIk")

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
		assert.Equal(t, 401, response.StatusCode)
		// Response message
		assert.Equal(t, "invalid or expired jwt", responsebody["message"])
	})

	t.Run("Do not have tokens", func(t *testing.T) {
		db := config.SetupTestDB()
		// Use router
		router := routes.SetupRouter(db)

		// Generate uuid, for sample id
		userID := uuid.NewString()

		// Create request
		request := httptest.NewRequest(http.MethodDelete, "http://localhost:8080/v1/users/"+userID, nil)
		// Added header content type
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
		// Response message
		assert.Equal(t, "missing or malformed jwt", responsebody["message"])
	})
}
