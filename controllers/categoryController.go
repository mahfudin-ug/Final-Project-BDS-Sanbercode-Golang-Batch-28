package controllers

import (
	"api-ecommerce/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type categoryInput struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// GetAllCategory godoc
// @Summary Get all Category
// @Description Get a list of Category
// @Tags Public
// @Param search query string false "name search by keyword"
// @Produce json
// @Success 200 {object} []models.Category
// @Router /categories [get]
func GetAllCategory(c *gin.Context) {
	// get db from gin context
	db := c.MustGet("db").(*gorm.DB)
	var categories []models.Category
	db.Find(&categories)
	if search := c.Request.URL.Query().Get("search"); search != "" {
		if err := db.Where("name LIKE ?", "%"+search+"%").Find(&categories).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"data": categories})
}

// CreateCategory godoc
// @Summary Create new Category
// @Description Creating new Category
// @Tags Admin
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Param Body body categoryInput true "the body to create new Category"
// @Produce json
// @Success 200 {object} models.Category
// @Router /categories [post]
func CreateCategory(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	// Validate input
	var input categoryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create category
	category := models.Category{
		Name:        input.Name,
		Description: input.Description,
	}
	db.Create(&category)

	c.JSON(http.StatusOK, gin.H{"data": category})
}

// GetProductByCategoryId godoc
// @Summary Get Products
// @Description Get all Products by CategoryId
// @Tags Public
// @Produce json
// @Param id path string true "Category id"
// @Success 200 {object} []models.Product
// @Router /categories/{id}/products [get]
func GetProductByCategoryId(c *gin.Context) {
	var products []models.Product

	db := c.MustGet("db").(*gorm.DB)

	if err := db.Where("category_id = ?", c.Param("id")).Find(&products).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": products})
}

// UpdateCategory godoc
// @Summary Update Category
// @Description Update Category by id
// @Tags Admin
// @Produce json
// @Param id path string true "Category id"
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Param Body body categoryInput true "the body to update category"
// @Success 200 {object} models.Category
// @Router /categories/{id} [put]
func UpdateCategory(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	// Get model if exist
	var category models.Category
	if err := db.Where("id = ?", c.Param("id")).First(&category).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record no found."})
		return
	}

	// Validate input
	var input categoryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var updatedInput models.Category
	updatedInput.Name = input.Name
	updatedInput.Description = input.Description
	updatedInput.UpdatedAt = time.Now()

	db.Model(&category).Updates(updatedInput)

	c.JSON(http.StatusOK, gin.H{"data": category})
}

// DeleteCategory godoc
// @Summary Delete category
// @Description Delete category by id
// @Tags Admin
// @Produce json
// @Param id path string true "Category id"
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Success 200 {object} map[string]boolean
// @Router /categories/{id} [delete]
func DeleteCategory(c *gin.Context) {
	// Get model if exist
	db := c.MustGet("db").(*gorm.DB)
	var category models.Category
	if err := db.Where("id = ?", c.Param("id")).First(&category).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found"})
		return
	}

	// TODO if category has used
	db.Delete(&category)
	c.JSON(http.StatusOK, gin.H{"data": true})
}
