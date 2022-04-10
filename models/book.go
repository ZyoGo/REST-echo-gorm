package models

import "gorm.io/gorm"

type Books struct {
	gorm.Model
	Title  string `json:"title" form:"title" validate:"required"`
	Author string `json:"author" form:"author" validate:"required"`
	Year   int    `json:"year" form:"year" validate:"required"`
}

type BooksResponse struct {
	ID     uint   `json:"id"`
	Title  string `json:"title" form:"title" validate:"required"`
	Author string `json:"author" form:"author" validate:"required"`
	Year   int    `json:"year" form:"year" validate:"required"`
}
