package models

import "gorm.io/gorm"

type Users struct {
	gorm.Model
	Name     string `json:"name" form:"name" validate:"required"`
	Email    string `json:"email" form:"email"  validate:"required,email"`
	Password string `json:"password" form:"password"  validate:"required"`
}

type UsersResponse struct {
	Name     string `json:"name" form:"name"`
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
}
