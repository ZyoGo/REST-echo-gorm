package config

import (
	"fmt"
	"rest-echo-gorm/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Config struct {
	DB_Username string
	DB_Password string
	DB_Port     string
	DB_Host     string
	DB_Name     string
}

func ConnectDB() *gorm.DB {
	// Development on windows
	// config := &Config{
	// 	DB_Username: "postgres",
	// 	DB_Password: "academy",
	// 	DB_Port:     "5432",
	// 	DB_Host:     "localhost",
	// 	DB_Name:     "rest_test",
	// }

	// Development on mac
	config := &Config{
		DB_Username: "devpostgres",
		DB_Password: "devpostgres",
		DB_Port:     "5432",
		DB_Host:     "localhost",
		DB_Name:     "gorm_test",
	}

	connectionString := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
		config.DB_Host,
		config.DB_Username,
		config.DB_Password,
		config.DB_Name,
		config.DB_Port,
	)

	var err error
	DB, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
		panic(err)
	}

	return DB
}

func InitialMigration() {
	DB := ConnectDB()
	DB.AutoMigrate(&models.Books{}, &models.Users{})
}
