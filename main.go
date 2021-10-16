package main

import (
	"api-ecommerce/config"
	"api-ecommerce/docs"
	"api-ecommerce/routes"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	// for load godotenv
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Programmatically set swagger info
	docs.SwaggerInfo.Title = "Swagger Example API"
	docs.SwaggerInfo.Description = "This is E-commerce server"
	docs.SwaggerInfo.Version = "1.0"
	// docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	db := config.ConnectDatabase()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	r := routes.SetupRouter(db)
	r.Run()
}
