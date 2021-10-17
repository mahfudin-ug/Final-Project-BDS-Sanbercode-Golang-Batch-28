package controllers

import (
	"api-ecommerce/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SearchKeyword godoc
// @Summary Search by keyword
// @Description Search Product, Shop, Category by keyword
// @Tags Public
// @Param search query string false "keyword of searching"
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /login [post]
func SearchKeyword(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var products []models.Product
	db.Find(&products)
	if search := c.Request.URL.Query().Get("search"); search != "" {
		if err := db.Where("name LIKE ?", "%"+search+"%").Find(&products).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Product not found!"})
			return
		}
	}

	var shops []models.Shop
	db.Find(&shops)
	if search := c.Request.URL.Query().Get("search"); search != "" {
		if err := db.Where("name LIKE ?", "%"+search+"%").Find(&shops).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Shop not found!"})
			return
		}
	}

	var categories []models.Category
	db.Find(&categories)
	if search := c.Request.URL.Query().Get("search"); search != "" {
		if err := db.Where("name LIKE ?", "%"+search+"%").Find(&categories).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Category not found!"})
			return
		}
	}

	result := map[string]interface{}{
		"product":  products,
		"shop":     shops,
		"category": categories,
	}

	c.JSON(http.StatusOK, gin.H{"success": "OK", "data": result})
}
