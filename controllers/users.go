package controllers

import (
	"REST-echo-gorm/lib/databases"
	"REST-echo-gorm/models"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type ResponseFormat struct {
	Status   int         `json:"status"`
	Messages string      `json:"messages"`
	Data     interface{} `json:"data"`
}

func CreateUserController(c echo.Context) (err error) {
	req := models.Users{}

	c.Bind(&req)

	//if err = c.Bind(&req); err != nil {
	//	return c.JSON(http.StatusBadRequest, err.Error())
	//}

	if err = c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	user, _ := databases.CreateUser(req)

	// return userResponses from database
	var userResponses models.UsersResponse
	userResponses.ID = user.ID
	userResponses.Name = user.Name
	userResponses.Email = user.Email
	userResponses.Password = user.Password

	return c.JSON(http.StatusOK, ResponseFormat{
		Status:   http.StatusOK,
		Messages: "Success create user",
		Data:     userResponses,
	})
}

func GetUsersController(c echo.Context) error {
	users, err := databases.GetUsers()

	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseFormat{
			Status:   http.StatusBadRequest,
			Messages: "Failed",
			Data:     err.Error(),
		})
	}

	// make userResponses from database
	var userResponses []models.UsersResponse
	userResponses = make([]models.UsersResponse, len(users))

	for i, user := range users {
		userResponses[i].ID = user.ID
		userResponses[i].Name = user.Name
		userResponses[i].Email = user.Email
		userResponses[i].Password = user.Password
	}

	return c.JSON(http.StatusOK, ResponseFormat{
		Status:   http.StatusOK,
		Messages: "Success",
		Data:     userResponses,
	})
}

func UpdateUserController(c echo.Context) error {
	req := models.Users{}
	c.Bind(&req)

	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	user, err := databases.UpdateUser(req)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseFormat{
			Status:   http.StatusInternalServerError,
			Messages: "Failed",
			Data:     err.Error(),
		})
	}

	// make userResponses from database
	var userResponses models.UsersResponse
	userResponses.ID = user.ID
	userResponses.Name = user.Name
	userResponses.Email = user.Email
	userResponses.Password = user.Password

	return c.JSON(http.StatusOK, ResponseFormat{
		Status:   http.StatusOK,
		Messages: "Success update user with id " + strconv.Itoa(int(user.ID)),
		Data:     userResponses,
	})
}
