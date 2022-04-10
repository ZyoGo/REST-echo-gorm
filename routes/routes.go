package routes

import (
	"rest-echo-gorm/controllers"
	"rest-echo-gorm/helpers"

	"github.com/go-playground/validator/v10"

	"github.com/labstack/echo/v4"
)

func New() *echo.Echo {
	e := echo.New()

	// Initialaze validator
	e.Validator = &helpers.CustomValidator{Validator: validator.New()}

	e.POST("/users", controllers.CreateUserController)
	e.GET("/users", controllers.GetUsersController)
	e.GET("/users/:id", controllers.GetUserController)
	e.PUT("/users/:id", controllers.UpdateUserController)
	e.DELETE("/users/:id", controllers.DeleteUserController)

	// Books routes
	e.POST("/books", controllers.CreateBookController)
	e.GET("/books", controllers.GetBooksController)
	e.GET("/books/:id", controllers.GetBookController)
	e.PUT("/books/:id", controllers.UpdateBookController)

	return e
}
