package config

import (
	"api-ecommerce/models"
	"api-ecommerce/utils"
	"fmt"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDatabase() *gorm.DB {
	env := utils.Getenv("ENVIRONMENT", "local")
	var db *gorm.DB
	var err error
	if env == "production" {
		username := os.Getenv("DB_USERNAME")
		password := os.Getenv("DB_PASSWORD")
		host := os.Getenv("DB_HOST")
		port := os.Getenv("DB_PORT")
		database := os.Getenv("DB_NAME")

		dsn := "host=" + host + " user=" + username + " password=" + password + " dbname=" + database + " port=" + port + " sslmode=require"
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	} else {
		// local env

		username := "user"
		password := "P@ssw0rd"
		database := "sanber_final"
		host := "tcp(127.0.0.1:3306)"

		dsn := fmt.Sprintf("%v:%v@%v/%v?charset=utf8mb4&parseTime=True&loc=Local", username, password, host, database)
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	}

	if err != nil {
		panic(err.Error())
	}

	db.AutoMigrate(&models.User{}, &models.Shop{}, &models.Category{}, &models.Product{}, &models.Order{}, &models.OrderProduct{}, &models.Address{})

	// Initialize data
	var user models.User
	if err := db.Where("username = ?", utils.Getenv("INITIAL_ADMIN_USERNAME", "admin")).First(&user).Error; err != nil {
		user.FirstName = "Admin"
		user.LastName = "Admin"
		user.Username = utils.Getenv("INITIAL_ADMIN_USERNAME", "admin")
		user.Email = "admin@golang.lo"
		user.Password = utils.Getenv("INITIAL_ADMIN_PASSWORD", "secret")
		user.Role = models.UserRoleAdmin
		_, err := user.SaveUser(db)

		if err != nil {
			panic(err.Error())
		}
	}

	return db
}
