package config

import (
	"api-ecommerce/models"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func ConnectDatabase() *gorm.DB {
	username := "user"
	password := "P@ssw0rd"
	database := "sanber_final"
	host := "tcp(127.0.0.1:3306)"

	dsn := fmt.Sprintf("%v:%v@%v/%v?charset=utf8mb4&parseTime=True&loc=Local", username, password, host, database)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(err.Error())
	}

	db.AutoMigrate(&models.User{}, &models.Shop{}, &models.Category{}, &models.Product{}, &models.Order{}, &models.OrderProduct{}, &models.Address{})

	return db
}
