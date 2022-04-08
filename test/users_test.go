package test

import (
	"REST-echo-gorm/config"
	"REST-echo-gorm/controllers"
	"REST-echo-gorm/helpers"
	"REST-echo-gorm/models"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func CleanTable(table []string) {
	for _, tableName := range table {
		config.ConnectDB().Exec("DELETE FROM " + tableName)
	}
}

func insertUserDb(name, email, password string) models.Users {
	user := models.Users{
		Name:     name,
		Email:    email,
		Password: password,
	}

	if err := config.ConnectDB().Create(&user).Error; err != nil {
		panic(err)
	}

	return user
}

func TestCreateUserController(t *testing.T) {
	t.Run("Test create user with valid payload", func(t *testing.T) {
		config.InitialMigration()

		requestBody := strings.NewReader(`{"name":"user1","email":"user1@gmail.com","password":"user123"}`)

		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/users", requestBody)
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		e.Validator = &helpers.CustomValidator{Validator: validator.New()}

		if assert.NoError(t, controllers.CreateUserController(c)) {
			// Decode rec.body and exctract Data.id
			var response struct {
				Data models.UsersResponse `json:"data"`
			}

			if err := json.Unmarshal([]byte(rec.Body.String()), &response); err != nil {
				t.Fatal(err)
			}

			expectedReturns := &controllers.ResponseFormat{
				Status:   http.StatusCreated,
				Messages: "Success create user",
				Data: models.UsersResponse{
					ID:       response.Data.ID,
					Name:     "user1",
					Email:    "user1@gmail.com",
					Password: "user123",
				},
			}

			var expectedReturnsJson bytes.Buffer
			if err := json.NewEncoder(&expectedReturnsJson).Encode(expectedReturns); err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, http.StatusCreated, rec.Code)
			assert.Equal(t, expectedReturnsJson.String(), rec.Body.String())
		}

		// Clean table
		CleanTable([]string{"users"})
	})

	t.Run("Test create user with bad payload", func(t *testing.T) {
		requestBody := strings.NewReader(`{"name":"user1","email":"user1","password":"user123"}`)

		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/users", requestBody)
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		e.Validator = &helpers.CustomValidator{Validator: validator.New()}

		if assert.NoError(t, controllers.CreateUserController(c)) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		}

		// Clean table
		CleanTable([]string{"users"})
	})
}

func TestGetUsersController(t *testing.T) {
	t.Run("Test get users", func(t *testing.T) {
		insertUsers := []models.Users{
			{
				Name:     "user1",
				Email:    "user1@gmail.com",
				Password: "user123",
			},
			{
				Name:     "user2",
				Email:    "user2@gmail.com",
				Password: "user123",
			},
			{
				Name:     "user3",
				Email:    "user3@gmail.com",
				Password: "user123",
			},
		}

		for _, user := range insertUsers {
			insertUserDb(user.Name, user.Email, user.Password)
		}

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/users", nil)
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		if assert.NoError(t, controllers.GetUsersController(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)

			body := rec.Body.String()

			var response struct {
				Users []models.Users `json:"data"`
			}

			err := json.Unmarshal([]byte(body), &response)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, len(insertUsers), len(response.Users))
		}

		// Clean table
		CleanTable([]string{"users"})
	})
}

func TestUpdateUserController(t *testing.T) {
	t.Run("Update user with valid payload", func(t *testing.T) {
		// First insert user into DB
		user := insertUserDb("user1", "user1@gmail.com", "user123")

		// Create requestBody for update user
		requestBody := map[string]interface{}{
			"name":     "user2Changed",
			"email":    "user2change@gmail.com",
			"password": "user2",
		}
		// Convert requestBody to json
		requestBodyJson, _ := json.Marshal(requestBody)

		// Initial echo
		e := echo.New()
		// Create request
		req := httptest.NewRequest(http.MethodPut, "/", bytes.NewReader(requestBodyJson))
		req.Header.Set("Content-Type", "application/json")
		// Create recorder
		rec := httptest.NewRecorder()
		// Create context
		c := e.NewContext(req, rec)
		// Set user id
		c.SetPath("/users/:id")
		c.SetParamNames("id")
		c.SetParamValues(helpers.ConvertUintToString(user.ID))

		// Validate request
		e.Validator = &helpers.CustomValidator{Validator: validator.New()}

		// Create response expected
		expectedReturns := &controllers.ResponseFormat{
			Status:   http.StatusOK,
			Messages: "Success",
			Data: models.UsersResponse{
				ID:       user.ID,
				Name:     "user2Changed",
				Email:    "user2change@gmail.com",
				Password: "user2",
			},
		}

		// Check if no error on controller
		if assert.NoError(t, controllers.UpdateUserController(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)

			// Convert struct ResponseFormat to json
			var expectedReturnsJson bytes.Buffer
			if err := json.NewEncoder(&expectedReturnsJson).Encode(expectedReturns); err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, expectedReturnsJson.String(), rec.Body.String())

		}
	})

	t.Run("Update user with invalid payload", func(t *testing.T) {
		// First insert user into DB
		user := insertUserDb("user1", "user1@gmail.com", "user123")

		// Create requestBody for update user
		requestBody := map[string]interface{}{
			"name":     "user2Changed",
			"email":    "user2changed",
			"password": "user2",
		}
		// Convert requestBody to json
		requestBodyJson, _ := json.Marshal(requestBody)

		// Initial echo
		e := echo.New()
		// Create request
		req := httptest.NewRequest(http.MethodPut, "/users", bytes.NewReader(requestBodyJson))
		req.Header.Set("Content-Type", "application/json")
		// Create recorder
		rec := httptest.NewRecorder()
		// Create context
		c := e.NewContext(req, rec)
		// Set user id
		c.SetPath("/users/:id")
		c.SetParamNames("id")
		c.SetParamValues(helpers.ConvertUintToString(user.ID))

		// Validate request
		e.Validator = &helpers.CustomValidator{Validator: validator.New()}

		// Create response expected
		expectedReturns := &controllers.ResponseFormat{
			Status:   http.StatusBadRequest,
			Messages: "Failed",
			Data:     "code=400, message=Key: 'UsersResponse.Email' Error:Field validation for 'Email' failed on the 'email' tag",
		}

		// Check if no error on controller
		if assert.NoError(t, controllers.UpdateUserController(c)) {
			assert.Equal(t, expectedReturns.Status, rec.Code)

			// Convert struct ResponseFormat to json
			var expectedReturnsJson bytes.Buffer
			if err := json.NewEncoder(&expectedReturnsJson).Encode(expectedReturns); err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, expectedReturnsJson.String(), rec.Body.String())

		}
	})
}
