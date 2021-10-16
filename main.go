package main

import (
	"api-ecommerce/config"
	"api-ecommerce/docs"
	"api-ecommerce/routes"
	"api-ecommerce/utils"
	"log"

	"github.com/joho/godotenv"
)

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @termsOfService http://swagger.io/terms/

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
	docs.SwaggerInfo.Host = utils.Getenv("SWAGGER_HOST", "localhost:8080")
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	db := config.ConnectDatabase()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	r := routes.SetupRouter(db)
	r.Run()
}
