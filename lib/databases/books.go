package databases

import "rest-echo-gorm/models"

func CreateBook(book models.Books) (*models.Books, error) {
	if err := DB.Create(&book).Error; err != nil {
		return nil, err
	}

	return &book, nil
}

func GetBook(reqId string) (*models.Books, error) {
	var book models.Books
	if err := DB.Where("id = ?", reqId).First(&book).Error; err != nil {
		return nil, err
	}

	return &book, nil
}

func GetBooks() (*[]models.Books, error) {
	var books []models.Books

	if err := DB.Find(&books).Error; err != nil {
		return nil, err
	}

	return &books, nil
}

func UpdateBook(book *models.Books) (error) {
	if err := DB.Save(&book).Error; err != nil {
		return err
	}

	return nil
}