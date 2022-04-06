package test

import (
	"REST-echo-gorm/controllers"
	"REST-echo-gorm/helpers"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCreateUserController(t *testing.T) {
	t.Run("Test create user with valid payload", func(t *testing.T) {
		requestBody := strings.NewReader(`{"name":"user1","email":"user1@gmail.com","password":"user123"}`)

		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/users", requestBody)
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		e.Validator = &helpers.CustomValidator{Validator: validator.New()}

		if assert.NoError(t, controllers.CreateUserController(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
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
	})
}
