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
)

// CreateRandomBook as function create random book with test and will be used by any function
func CreateRandomBook(t *testing.T) models.Books {
	// Login Random account
	_, token := LoginRandomAccount(t)

	db := config.SetupTestDB()
	// Use router
	router := routes.SetupRouter(db)

	// Create random dummy data
	data := lib.CreateOrUpdateBook{
		Name:   jabufaker.RandomString(7),
		Author: jabufaker.RandomPerson(),
	}

	// String json
	dataBody := fmt.Sprintf(`{"name": "%s", "author": "%s"}`, data.Name, data.Author)

	// Create new reader
	requestBody := strings.NewReader(dataBody)
	// Create request
	request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/v1/books", requestBody)
	// Added header content type
	request.Header.Add("Content-Type", "application/json")
	// Bearer token
	bearerToken := fmt.Sprintf("Bearer %s", token)
	// Added header Authorization
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
	assert.Equal(t, 201, response.StatusCode)
	// Response body status code
	assert.Equal(t, 201, int(responsebody["code"].(float64)))
	// Response body status
	assert.Equal(t, "success", responsebody["status"])
	// Response body message
	assert.Equal(t, "Book has been created", responsebody["message"])

	// Response body data
	assert.NotEmpty(t, responsebody["data"])

	// Context property data on body
	var contextData = responsebody["data"].(map[string]interface{})
	assert.NotEmpty(t, contextData["id"])
	assert.NotEmpty(t, contextData["created_at"])
	assert.NotEmpty(t, contextData["updated_at"])
	assert.Equal(t, data.Name, contextData["name"])
	assert.Equal(t, data.Author, contextData["author"])

	// Return new books, for re-use any function test
	newBook := models.Books{
		ID:        contextData["id"].(string),
		Name:      contextData["name"].(string),
		Author:    contextData["author"].(string),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	return newBook
}

// Test Create book success (200).
func TestCreateBookSuccess(t *testing.T) {
	CreateRandomBook(t)
}

// - [ ] Create book failed, validation error (400).
func TestCreateBookFailedValidationError(t *testing.T) {
	// Login Random account
	_, token := LoginRandomAccount(t)

	db := config.SetupTestDB()
	// Use router
	router := routes.SetupRouter(db)

	// Create random dummy data
	data := lib.CreateOrUpdateBook{
		Name:   "",
		Author: "",
	}

	// String json
	dataBody := fmt.Sprintf(`{"name": "%s", "author": "%s"}`, data.Name, data.Author)

	// Create new reader
	requestBody := strings.NewReader(dataBody)
	// Create request
	request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/v1/books", requestBody)
	// Added header content type
	request.Header.Add("Content-Type", "application/json")
	// Bearer token
	bearerToken := fmt.Sprintf("Bearer %s", token)
	// Added header Authorization
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

	// Response body data
	assert.NotEmpty(t, responsebody["data"])

	// Response body data message
	assert.NotEmpty(t, responsebody["data"].(map[string]interface{})["message"])
}

// Test Create book failed, unauthorized (400 and 401).
func TestCreateBookUnauthorized(t *testing.T) {
	t.Run("Invalid token or expired", func(t *testing.T) {
		db := config.SetupTestDB()
		// Use router
		router := routes.SetupRouter(db)

		// Create random dummy data
		data := lib.CreateOrUpdateBook{
			Name:   jabufaker.RandomString(7),
			Author: jabufaker.RandomPerson(),
		}

		// String json
		dataBody := fmt.Sprintf(`{"name": "%s", "author": "%s"}`, data.Name, data.Author)

		// Create new reader
		requestBody := strings.NewReader(dataBody)
		// Create request
		request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/v1/books", requestBody)
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

		// Create random dummy data
		data := lib.CreateOrUpdateBook{
			Name:   jabufaker.RandomString(7),
			Author: jabufaker.RandomPerson(),
		}

		// String json
		dataBody := fmt.Sprintf(`{"name": "%s", "author": "%s"}`, data.Name, data.Author)

		// Create new reader
		requestBody := strings.NewReader(dataBody)
		// Create request
		request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/v1/books", requestBody)
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

// Test Get all books success (200).
func TestGetAllBooksSuccess(t *testing.T) {
	// Create some new book
	for i := 0; i < 5; i++ {
		CreateRandomBook(t)
	}

	db := config.SetupTestDB()
	// Use router
	router := routes.SetupRouter(db)

	// Create request
	request := httptest.NewRequest(http.MethodGet, "http://localhost:8080/v1/books", nil)
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
	assert.Equal(t, "List of books", responsebody["message"])

	// Iteration data
	data := responsebody["data"].([]interface{})
	for _, item := range data {
		list := item.(map[string]interface{})
		assert.NotEmpty(t, list["id"])
		assert.NotEmpty(t, list["name"])
		assert.NotEmpty(t, list["author"])
		assert.NotEmpty(t, list["created_at"])
		assert.NotEmpty(t, list["updated_at"])
	}
}

// Test Get book by id success (200).
func TestGetBookByIdSuccess(t *testing.T) {
	// Create a book
	newBook := CreateRandomBook(t)

	db := config.SetupTestDB()
	// Use router
	router := routes.SetupRouter(db)

	// Create request
	request := httptest.NewRequest(http.MethodGet, "http://localhost:8080/v1/books/"+newBook.ID, nil)
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
	assert.Equal(t, "List of books", responsebody["message"])

	// Response body data
	assert.NotEmpty(t, responsebody["data"])

	// Context property data on body
	var contextData = responsebody["data"].(map[string]interface{})
	assert.Equal(t, newBook.ID, contextData["id"])
	assert.Equal(t, newBook.Name, contextData["name"])
	assert.Equal(t, newBook.Author, contextData["author"])
	assert.NotEmpty(t, contextData["created_at"])
	assert.NotEmpty(t, contextData["updated_at"])
}

// Tes Get book by id failed, book not found (400).
func TestGetBookByIdFailedNotFound(t *testing.T) {
	db := config.SetupTestDB()
	// Use router
	router := routes.SetupRouter(db)

	// Generate uuid
	bookId := uuid.NewString()
	// Create request
	request := httptest.NewRequest(http.MethodGet, "http://localhost:8080/v1/books/"+bookId, nil)
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
	message := fmt.Sprintf("Book with id: %s not found", bookId)
	assert.Equal(t, message, responsebody["message"])
}

// Test Update book success (200).
func TestUpdateBookSuccess(t *testing.T) {
	// Login Random account
	_, token := LoginRandomAccount(t)

	// Create a new book
	newBook := CreateRandomBook(t)

	db := config.SetupTestDB()
	// Use router
	router := routes.SetupRouter(db)

	// Create random dummy data
	updatedBook := lib.CreateOrUpdateBook{
		Name:   jabufaker.RandomString(7),
		Author: jabufaker.RandomPerson(),
	}

	// String json
	dataBody := fmt.Sprintf(`{"name": "%s", "author": "%s"}`, updatedBook.Name, updatedBook.Author)

	// Create new reader
	requestBody := strings.NewReader(dataBody)
	// Create request
	request := httptest.NewRequest(http.MethodPut, "http://localhost:8080/v1/books/"+newBook.ID, requestBody)
	// Added header content type
	request.Header.Add("Content-Type", "application/json")
	// Bearer token
	bearerToken := fmt.Sprintf("Bearer %s", token)
	// Added header Authorization
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
	assert.Equal(t, "Book has been updated", responsebody["message"])

	// Response body data
	assert.NotEmpty(t, responsebody["data"])

	// Context property data on body
	var contextData = responsebody["data"].(map[string]interface{})
	assert.Equal(t, newBook.ID, contextData["id"])
	assert.Equal(t, updatedBook.Name, contextData["name"])
	assert.Equal(t, updatedBook.Author, contextData["author"])
	assert.NotEmpty(t, contextData["created_at"])
	assert.NotEmpty(t, contextData["updated_at"])
}

// Test Update book failed, book not found (400).
func TestUpdateBookFailedNotFound(t *testing.T) {
	// Login Random account
	_, token := LoginRandomAccount(t)

	db := config.SetupTestDB()
	// Use router
	router := routes.SetupRouter(db)

	// Create random dummy data
	updatedBook := lib.CreateOrUpdateBook{
		Name:   jabufaker.RandomString(7),
		Author: jabufaker.RandomPerson(),
	}

	// String json
	dataBody := fmt.Sprintf(`{"name": "%s", "author": "%s"}`, updatedBook.Name, updatedBook.Author)

	// Create new reader
	requestBody := strings.NewReader(dataBody)

	// generate uuid
	bookId := uuid.NewString()
	// Create request
	request := httptest.NewRequest(http.MethodPut, "http://localhost:8080/v1/books/"+bookId, requestBody)
	// Added header content type
	request.Header.Add("Content-Type", "application/json")
	// Bearer token
	bearerToken := fmt.Sprintf("Bearer %s", token)
	// Added header Authorization
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
	message := fmt.Sprintf("Book with id: %s not found", bookId)
	assert.Equal(t, message, responsebody["message"])
}

// Test Update book failed, validation error (400).
func TestUpdateBookFailedValidationError(t *testing.T) {
	// Login Random account
	_, token := LoginRandomAccount(t)

	db := config.SetupTestDB()
	// Use router
	router := routes.SetupRouter(db)

	// Create random dummy data
	updatedBook := lib.CreateOrUpdateBook{
		Name:   "",
		Author: "",
	}

	// String json
	dataBody := fmt.Sprintf(`{"name": "%s", "author": "%s"}`, updatedBook.Name, updatedBook.Author)

	// Create new reader
	requestBody := strings.NewReader(dataBody)

	// generate uuid
	bookId := uuid.NewString()
	// Create request
	request := httptest.NewRequest(http.MethodPut, "http://localhost:8080/v1/books/"+bookId, requestBody)
	// Added header content type
	request.Header.Add("Content-Type", "application/json")
	// Bearer token
	bearerToken := fmt.Sprintf("Bearer %s", token)
	// Added header Authorization
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
	assert.Equal(t, "Update book failed", responsebody["message"])
	// Response body data message
	assert.NotEmpty(t, responsebody["data"].(map[string]interface{})["message"])
}

// Test Update book failed, unauthorized (400 and 401).
func TestUpdateBookFailedUnauthorized(t *testing.T) {
	t.Run("Invalid token or expired", func(t *testing.T) {
		db := config.SetupTestDB()
		// Use router
		router := routes.SetupRouter(db)

		// Create random dummy data
		data := lib.CreateOrUpdateBook{
			Name:   jabufaker.RandomString(7),
			Author: jabufaker.RandomPerson(),
		}

		// String json
		dataBody := fmt.Sprintf(`{"name": "%s", "author": "%s"}`, data.Name, data.Author)

		// Create new reader
		requestBody := strings.NewReader(dataBody)

		// Generate uuid
		bookId := uuid.NewString()
		// Create request
		request := httptest.NewRequest(http.MethodPut, "http://localhost:8080/v1/books/"+bookId, requestBody)
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

		// Create random dummy data
		data := lib.CreateOrUpdateBook{
			Name:   jabufaker.RandomString(7),
			Author: jabufaker.RandomPerson(),
		}

		// String json
		dataBody := fmt.Sprintf(`{"name": "%s", "author": "%s"}`, data.Name, data.Author)

		// Create new reader
		requestBody := strings.NewReader(dataBody)

		// Generate uuid
		bookId := uuid.NewString()
		// Create request
		request := httptest.NewRequest(http.MethodPut, "http://localhost:8080/v1/books/"+bookId, requestBody)
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

// Test Delete Book success (200).
func TestDeleteBookSuccess(t *testing.T) {
	// Login Random account
	_, token := LoginRandomAccount(t)

	// Create a new book
	newBook := CreateRandomBook(t)

	db := config.SetupTestDB()
	// Use router
	router := routes.SetupRouter(db)

	// Create request
	request := httptest.NewRequest(http.MethodDelete, "http://localhost:8080/v1/books/"+newBook.ID, nil)
	// Added header content type
	request.Header.Add("Content-Type", "application/json")
	// Bearer token
	bearerToken := fmt.Sprintf("Bearer %s", token)
	// Added header Authorization
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
	assert.Equal(t, "Book has been deleted", responsebody["message"])
}

// - [ ] Delete book failed, book not found (400).
func TestDeleteBookFailedNotFound(t *testing.T) {
	// Login Random account
	_, token := LoginRandomAccount(t)

	db := config.SetupTestDB()
	// Use router
	router := routes.SetupRouter(db)

	// Generate uuid
	bookID := uuid.NewString()
	// Create request
	request := httptest.NewRequest(http.MethodDelete, "http://localhost:8080/v1/books/"+bookID, nil)
	// Added header content type
	request.Header.Add("Content-Type", "application/json")
	// Bearer token
	bearerToken := fmt.Sprintf("Bearer %s", token)
	// Added header Authorization
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
	message := fmt.Sprintf("Book with id: %s not found", bookID)
	assert.Equal(t, message, responsebody["message"])
}

// Test Delete book failed, unauthorized (400 and 401)
func TestDeleteBookFailedUnauthorized(t *testing.T) {
	t.Run("Invalid token or expired", func(t *testing.T) {
		db := config.SetupTestDB()
		// Use router
		router := routes.SetupRouter(db)

		// Generate uuid
		bookId := uuid.NewString()
		// Create request
		request := httptest.NewRequest(http.MethodDelete, "http://localhost:8080/v1/books/"+bookId, nil)
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

		// Added header content typ// Generate uuid
		bookId := uuid.NewString()
		// Create request
		request := httptest.NewRequest(http.MethodDelete, "http://localhost:8080/v1/books/"+bookId, nil)
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
