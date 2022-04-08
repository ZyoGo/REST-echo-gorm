package routes

import (
	"REST-echo-gorm/controllers"
	"REST-echo-gorm/helpers"

	"github.com/go-playground/validator/v10"

	"github.com/labstack/echo/v4"
)

func New() *echo.Echo {
	e := echo.New()

	// Initialaze validator
	e.Validator = &helpers.CustomValidator{Validator: validator.New()}

	e.POST("/users", controllers.CreateUserController)
	e.GET("/users", controllers.GetUsersController)
	e.PUT("/users/:id", controllers.UpdateUserController)

	return e
}
