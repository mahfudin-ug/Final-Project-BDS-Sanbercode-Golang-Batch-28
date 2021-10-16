package controllers

import (
	"api-ecommerce/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type orderInput struct {
	Payment string `json:"payment"`
	Courier string `json:"courier"`
	UserID  uint   `json:"user_id"`
}

// GetAllOrder godoc
// @Summary Get all Order
// @Description Get a list of Order
// @Tags Admin
// @Produce json
// @Success 200 {object} []models.Order
// @Router /orders [get]
func GetAllOrder(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var orders []models.Order
	db.Find(&orders)

	c.JSON(http.StatusOK, gin.H{"data": orders})
}

// CreateOrder godoc
// @Summary Create new Order
// @Description Creating new Order
// @Tags User
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Param Body body orderInput true "the body to create new Order"
// @Produce json
// @Success 200 {object} models.Order
// @Router /orders [post]
func CreateOrder(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var input orderInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var user models.User
	if err := db.Where("id = ?", input.UserID).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}

	// Create
	order := models.Order{
		Payment: input.Payment,
		Courier: input.Courier,
		UserID:  user.ID,
		Status:  models.OrderStatusInitial,
	}
	db.Create(&order)
	c.JSON(http.StatusOK, gin.H{"data": order})
}

// GetOrderById godoc
// @Summary Get Order detail
// @Description Get Order by Id
// @Tags User
// @Produce json
// @Param id path string true "Order id"
// @Success 200 {object} models.Order
// @Router /orders/{id} [get]
func GetOrderById(c *gin.Context) {
	var order models.Order

	db := c.MustGet("db").(*gorm.DB)
	if err := db.Where("id = ?", c.Param("id")).First(&order).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Order not found!"})
		return
	}
	_, err := order.RecalculateOrder(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Order Recalculate failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": order})
}

// UpdateOrder godoc
// @Summary Update Order
// @Description Update Order by id
// @Tags User
// @Produce json
// @Param id path string true "Order id"
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Param Body body orderInput true "the body to update order"
// @Success 200 {object} models.Order
// @Router /orders/{id} [put]
func UpdateOrder(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	// Get model if exist
	var order models.Order
	if err := db.Where("id = ?", c.Param("id")).First(&order).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Order not found!"})
		return
	}
	// Validate input
	var input orderInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var user models.User
	if err := db.Where("id = ?", input.UserID).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}

	var updatedInput models.Order
	updatedInput.Payment = input.Payment
	updatedInput.Courier = input.Courier
	updatedInput.UserID = input.UserID
	updatedInput.UpdatedAt = time.Now()

	db.Model(&order).Updates(updatedInput)
	_, err := order.RecalculateOrder(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Order Recalculate failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": order})
}

// DeleteOrder godoc
// @Summary Delete order
// @Description Delete order by id
// @Tags Admin
// @Produce json
// @Param id path string true "Order id"
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Success 200 {object} map[string]boolean
// @Router /orders/{id} [delete]
func DeleteOrder(c *gin.Context) {
	// Get model if exist
	db := c.MustGet("db").(*gorm.DB)
	var order models.Order
	if err := db.Where("id = ?", c.Param("id")).First(&order).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found"})
		return
	}

	// TODO if order has used
	db.Delete(&order)
	c.JSON(http.StatusOK, gin.H{"data": true})
}
