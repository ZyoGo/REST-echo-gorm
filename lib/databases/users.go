package databases

import (
	"REST-echo-gorm/config"
	"REST-echo-gorm/models"
)

var DB = config.ConnectDB()

func CreateUser(user models.Users) (models.Users, error) {

	if err := DB.Create(&user).Error; err != nil {
		return models.Users{}, err
	}

	return user, nil
}

func GetUsers() ([]models.Users, error) {
	var users []models.Users

	if err := DB.Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}
