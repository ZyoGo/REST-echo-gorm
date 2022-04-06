package config

import (
	"REST-echo-gorm/models"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDBTest() *gorm.DB {
	config := &Config{
		DB_Username: "postgres",
		DB_Password: "academy",
		DB_Port:     "5432",
		DB_Host:     "localhost",
		DB_Name:     "rest_test",
	}

	connectionString := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
		config.DB_Host,
		config.DB_Username,
		config.DB_Password,
		config.DB_Name,
		config.DB_Port,
	)

	var err error
	DB, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	DB.AutoMigrate(&models.Books{}, &models.Users{})

	return DB
}
