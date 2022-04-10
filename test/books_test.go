package test

import (
	"bytes"
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"rest-echo-gorm/config"
	"rest-echo-gorm/controllers"
	"rest-echo-gorm/helpers"
	"rest-echo-gorm/models"
	"strings"
	"testing"
)

func insertBooks(title, author string, year int) *models.Books {
	book := &models.Books{
		Title: title,
		Author: author,
		Year: year,
	}

	if err := config.ConnectDB().Create(&book).Error; err != nil {
		panic(err.Error())
	}
	
	return book
}

func TestCreateBookController(t *testing.T) {
	t.Run("Test create book with valid payload", func(t *testing.T) {
		config.InitialMigration()

		requestBody := strings.NewReader(`{"title":"book1","author":"user1","year":2022}`)

		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/books", requestBody)
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		e.Validator = &helpers.CustomValidator{Validator: validator.New()}

		if assert.NoError(t, controllers.CreateBookController(c)) {
			var response struct{
				Data models.BooksResponse `json:"data"`
			}

			if err := json.Unmarshal([]byte(rec.Body.String()), &response); err != nil {
				t.Fatal(err)
			}

			expectedResponse := &controllers.ResponseFormat{
				Status: http.StatusCreated,
				Messages: "Success",
				Data: models.BooksResponse{
					ID: response.Data.ID,
					Title: "book1",
					Author: "user1",
					Year: 2022,
				},
			}

			var expectedResponses bytes.Buffer
			if err := json.NewEncoder(&expectedResponses).Encode(expectedResponse); err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, expectedResponse.Status, rec.Code)
			assert.Equal(t, expectedResponses.String(), rec.Body.String())
		}

		CleanTable([]string{"books"})
	})

	t.Run("Test create book with invalid payload", func(t *testing.T) {
		config.InitialMigration()

		requestBody := strings.NewReader(`{"title":"book1","author":"user1","year":"2022"}`)

		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/books", requestBody)
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		e.Validator = &helpers.CustomValidator{Validator: validator.New()}

		expectedResponse := &controllers.ResponseFormat{
			Status: http.StatusBadRequest,
			Messages: "Fail",
			Data: "code=400, message=Key: 'Books.Year' Error:Field validation for 'Year' failed on the 'required' tag",
		}

		var expectedResponses bytes.Buffer
		if err := json.NewEncoder(&expectedResponses).Encode(expectedResponse); err != nil {
			t.Fatal(err)
		}

		if assert.NoError(t, controllers.CreateBookController(c)) {
			assert.Equal(t, expectedResponse.Status, rec.Code)
			assert.Equal(t, expectedResponses.String(), rec.Body.String())
		}

		// Clean table books
		CleanTable([]string{"books"})
	})
}

func TestGetBookController(t *testing.T) {
	t.Run("Get book with valid id", func(t *testing.T) {
		book := insertBooks("book1", "user1", 2022)

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/books", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		c.SetPath("/books/:id")
		c.SetParamNames("id")
		c.SetParamValues(helpers.ConvertUintToString(book.ID))

		expectedResponse := &controllers.ResponseFormat{
			Status: http.StatusOK,
			Messages: "Success",
			Data: &models.BooksResponse{
				ID: book.ID,
				Title: book.Title,
				Author: book.Author,
				Year: book.Year,
			},
		}

		var expectedResponses bytes.Buffer
		if err := json.NewEncoder(&expectedResponses).Encode(expectedResponse); err != nil {
			t.Fatal(err.Error())
		}

		if assert.NoError(t, controllers.GetBookController(c)) {
			assert.Equal(t, expectedResponse.Status, rec.Code)
			assert.Equal(t, expectedResponses.String(), rec.Body.String())
		}
	})

	t.Run("Get book with invalid id", func(t *testing.T) {
		insertBooks("book1", "user1", 2022)

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/books", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		c.SetPath("/books/:id")
		c.SetParamNames("id")
		c.SetParamValues("1000")

		expectedResponse := &controllers.ResponseFormat{
			Status: http.StatusNotFound,
			Messages: "Fail",
			Data: "record not found",
		}

		var expectedResponses bytes.Buffer
		if err := json.NewEncoder(&expectedResponses).Encode(expectedResponse); err != nil {
			t.Fatal(err.Error())
		}

		if assert.NoError(t, controllers.GetBookController(c)) {
			assert.Equal(t, expectedResponse.Status, rec.Code)
			assert.Equal(t, expectedResponses.String(), rec.Body.String())
		}
	})
}

func TestGetBooksController(t *testing.T) {
	t.Run("Test get books controller", func(t *testing.T) {
		books := []models.Books{
			{
				Title: "Book1",
				Author: "user1",
				Year: 2021,
			},
			{
				Title: "Book2",
				Author: "user2",
				Year: 2022,
			},
		}

		for _, book := range books {
			insertBooks(book.Title, book.Author, book.Year)
		}

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/books", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		if assert.NoError(t, controllers.GetBooksController(c)) {
			body := rec.Body.String()

			var response struct{
				Data []models.BooksResponse `json:"data"`
			}

			err := json.Unmarshal([]byte(body), &response)
			if err != nil {
				t.Fatal(err.Error())
			}

			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, len(books), len(response.Data))
		}
	})
}

func TestUpdateBookController(t *testing.T) {
	t.Run("Test update book with valid payload", func(t *testing.T) {
		book := insertBooks("Book1", "user1", 2022)

		payloadBody := strings.NewReader(`{"title":"Book 2 Changed","author":"User 2 Changed","year":2023}`)

		e := echo.New()
		req := httptest.NewRequest(http.MethodPut, "/books", payloadBody)
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		e.Validator = &helpers.CustomValidator{Validator: validator.New()}

		c.SetPath("/books/:id")
		c.SetParamNames("id")
		c.SetParamValues(helpers.ConvertUintToString(book.ID))

		expectedResponse := &controllers.ResponseFormat{
			Status: http.StatusOK,
			Messages: "Success",
			Data: &models.BooksResponse{
				ID: book.ID,
				Title: "Book 2 Changed",
				Author: "User 2 Changed",
				Year: 2023,
			},
		}

		var expectedResponses bytes.Buffer
		if err := json.NewEncoder(&expectedResponses).Encode(&expectedResponse); err != nil {
			t.Fatal(err.Error())
		}

		if assert.NoError(t, controllers.UpdateBookController(c)) {
			assert.Equal(t, expectedResponse.Status, rec.Code)
			assert.Equal(t, expectedResponses.String(), rec.Body.String())
		}
	})

	t.Run("Test update book with invalid payload", func(t *testing.T) {
		book := insertBooks("Book1", "user1", 2022)

		payloadBody := strings.NewReader(`{"title":"Book 2 Changed","author":"User 2 Changed","year":"2023"}`)

		e := echo.New()
		req := httptest.NewRequest(http.MethodPut, "/books", payloadBody)
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		e.Validator = &helpers.CustomValidator{Validator: validator.New()}

		c.SetPath("/books/:id")
		c.SetParamNames("id")
		c.SetParamValues(helpers.ConvertUintToString(book.ID))

		expectedResponse := &controllers.ResponseFormat{
			Status: http.StatusBadRequest,
			Messages: "Fail",
			Data: "code=400, message=Key: 'BooksResponse.Year' Error:Field validation for 'Year' failed on the 'required' tag",
		}

		var expectedResponses bytes.Buffer
		if err := json.NewEncoder(&expectedResponses).Encode(&expectedResponse); err != nil {
			t.Fatal(err.Error())
		}

		if assert.NoError(t, controllers.UpdateBookController(c)) {
			assert.Equal(t, expectedResponse.Status, rec.Code)
			assert.Equal(t, expectedResponses.String(), rec.Body.String())
		}
	})
}