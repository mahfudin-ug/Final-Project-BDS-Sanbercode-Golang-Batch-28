package controllers

import (
	"api-ecommerce/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type shopUpdateInput struct {
	Name      string `json:"name"`
	Bank      string `json:"bank"`
	Phone     string `json:"phone"`
	AddressID uint   `json:"address_id"`
}

type shopInput struct {
	shopUpdateInput
	UserID uint `json:"user_id"`
}

// GetAllShop godoc
// @Summary Get all Shop
// @Description Get a list of Shop
// @Tags Public
// @Param search query string false "name search by keyword"
// @Produce json
// @Success 200 {object} []models.Shop
// @Router /shops [get]
func GetAllShop(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var shops []models.Shop
	db.Find(&shops)

	c.JSON(http.StatusOK, gin.H{"data": shops})
}

// CreateShop godoc
// @Summary Create new Shop
// @Description Creating new Shop
// @Tags User
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Param Body body shopInput true "the body to create new Shop"
// @Produce json
// @Success 200 {object} models.Shop
// @Router /shops [post]
func CreateShop(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var input shopInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: restruct shop

	// Create
	shop := models.Shop{
		Name:      input.Name,
		Bank:      input.Bank,
		Phone:     input.Phone,
		UserID:    input.UserID,
		AddressID: input.AddressID,
	}
	db.Create(&shop)

	c.JSON(http.StatusOK, gin.H{"data": shop})
}

// GetShopById godoc
// @Summary Get Shop detail
// @Description Get Shop by Id
// @Tags Public
// @Produce json
// @Param id path string true "Shop id"
// @Success 200 {object} models.Shop
// @Router /shops/{id} [get]
func GetShopById(c *gin.Context) {
	var shop models.Shop

	db := c.MustGet("db").(*gorm.DB)
	if err := db.Where("id =?", c.Param("id")).First(&shop).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": shop})
}

// GetProductByShopId godoc
// @Summary Get Products
// @Description Get all Products by ShopId
// @Tags Public
// @Produce json
// @Param id path string true "Shop id"
// @Success 200 {object} []models.Product
// @Router /shops/{id}/products [get]
func GetProductByShopId(c *gin.Context) {
	var products []models.Product

	db := c.MustGet("db").(*gorm.DB)

	if err := db.Where("shop_id = ?", c.Param("id")).Find(&products).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": products})
}

// UpdateShop godoc
// @Summary Update Shop
// @Description Update Shop by id
// @Tags User
// @Produce json
// @Param id path string true "Shop id"
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Param Body body shopUpdateInput true "the body to update shop"
// @Success 200 {object} models.Shop
// @Router /shops/{id} [put]
func UpdateShop(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	// Get model if exist
	var shop models.Shop
	if err := db.Where("id = ?", c.Param("id")).First(&shop).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record no found."})
		return
	}

	// Validate input
	var input shopUpdateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var updatedInput models.Shop
	updatedInput.Name = input.Name
	updatedInput.Bank = input.Bank
	updatedInput.Phone = input.Phone
	updatedInput.AddressID = input.AddressID
	updatedInput.UpdatedAt = time.Now()

	db.Model(&shop).Updates(updatedInput)

	c.JSON(http.StatusOK, gin.H{"data": shop})

}

// DeleteShop godoc
// @Summary Delete shop
// @Description Delete shop by id
// @Tags User
// @Produce json
// @Param id path string true "Shop id"
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Success 200 {object} map[string]boolean
// @Router /shops/{id} [delete]
func DeleteShop(c *gin.Context) {
	// Get model if exist
	db := c.MustGet("db").(*gorm.DB)
	var shop models.Shop
	if err := db.Where("id = ?", c.Param("id")).First(&shop).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found"})
		return
	}

	db.Delete(&shop)
	c.JSON(http.StatusOK, gin.H{"data": true})
}
