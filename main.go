package main

import (
	"REST-echo-gorm/config"
	"REST-echo-gorm/routes"
)

func main() {
	// Initialize Database
	config.InitialMigration()

	// Initialize Routes
	e := routes.New()

	e.Logger.Fatal(e.Start(":4001"))
}
