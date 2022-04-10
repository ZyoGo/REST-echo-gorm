package databases

import "rest-echo-gorm/models"

func LoginUser(user models.Users) (*models.Users, error) {
	if err := DB.Where("email = ? AND password = ?", user.Email, user.Password).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}
