package controllers

import (
	"api-ecommerce/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type productInput struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	PhotoPath   string `json:"photo_path"`
	Stock       uint   `json:"stock"`
	Price       uint   `json:"price"`
	Weight      uint   `json:"weight"`
	CategoryID  uint   `json:"category_id"`
	ShopID      uint   `json:"shop_id"`
}

// GetAllProduct godoc
// @Summary Get all Product
// @Description Get a list of Product
// @Tags Public
// @Param search query string false "name search by keyword"
// @Produce json
// @Success 200 {object} []models.Product
// @Router /products [get]
func GetAllProduct(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var products []models.Product
	db.Find(&products)
	if search := c.Request.URL.Query().Get("search"); search != "" {
		if err := db.Where("name LIKE ?", "%"+search+"%").Find(&products).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"data": products})
}

// CreateProduct godoc
// @Summary Create new Product
// @Description Creating new Product
// @Tags User
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Param Body body productInput true "the body to create new Product"
// @Produce json
// @Success 200 {object} models.Product
// @Router /products [post]
func CreateProduct(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var input productInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var category models.Category
	if err := db.Where("id = ?", input.CategoryID).First(&category).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Category not found"})
		return
	}
	var shop models.Shop
	if err := db.Where("id = ?", input.ShopID).First(&shop).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Shop not found"})
		return
	}

	// Create
	product := models.Product{
		Name:        input.Name,
		Description: input.Description,
		Stock:       input.Stock,
		PhotoPath:   input.PhotoPath,
		Price:       input.Price,
		Weight:      input.Weight,
		CategoryID:  category.ID,
		ShopID:      shop.ID,
	}
	db.Create(&product)

	c.JSON(http.StatusOK, gin.H{"data": product})
}

// GetProductById godoc
// @Summary Get Product detail
// @Description Get Product by Id
// @Tags Public
// @Produce json
// @Param id path string true "Product id"
// @Success 200 {object} models.Product
// @Router /products/{id} [get]
func GetProductById(c *gin.Context) {
	var product models.Product

	db := c.MustGet("db").(*gorm.DB)
	if err := db.Where("id = ?", c.Param("id")).First(&product).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": product})
}

// UpdateProduct godoc
// @Summary Update Product
// @Description Update Product by id
// @Tags User
// @Produce json
// @Param id path string true "Product id"
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Param Body body productInput true "the body to update product"
// @Success 200 {object} models.Product
// @Router /products/{id} [put]
func UpdateProduct(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	// Get model if exist
	var product models.Product
	if err := db.Where("id = ?", c.Param("id")).First(&product).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record no found."})
		return
	}

	// Validate input
	var input productInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var category models.Category
	if err := db.Where("id = ?", input.CategoryID).First(&category).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Category not found"})
		return
	}
	var shop models.Shop
	if err := db.Where("id = ?", input.ShopID).First(&shop).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Shop not found"})
		return
	}

	var updatedInput models.Product
	updatedInput.Name = input.Name
	updatedInput.Description = input.Description
	updatedInput.Stock = input.Stock
	updatedInput.PhotoPath = input.PhotoPath
	updatedInput.Price = input.Price
	updatedInput.Weight = input.Weight
	updatedInput.CategoryID = category.ID
	updatedInput.ShopID = shop.ID
	updatedInput.UpdatedAt = time.Now()

	db.Model(&product).Updates(updatedInput)

	c.JSON(http.StatusOK, gin.H{"data": product})
}

// DeleteProduct godoc
// @Summary Delete product
// @Description Delete product by id
// @Tags User
// @Produce json
// @Param id path string true "Product id"
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Success 200 {object} map[string]boolean
// @Router /products/{id} [delete]
func DeleteProduct(c *gin.Context) {
	// Get model if exist
	db := c.MustGet("db").(*gorm.DB)
	var product models.Product
	if err := db.Where("id = ?", c.Param("id")).First(&product).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found"})
		return
	}

	// TODO if product has used
	db.Delete(&product)
	c.JSON(http.StatusOK, gin.H{"data": true})
}
