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
type orderSendInput struct {
	Resi string `json:"resi"`
}

// GetAllOrder godoc
// @Summary Get all Order
// @Description Get a list of Order
// @Tags Buyer, Admin
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
// @Tags Buyer
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
// @Tags Buyer, Admin
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
// @Tags Buyer, Admin
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

	var updatedInput models.Order
	updatedInput.Payment = input.Payment
	updatedInput.Courier = input.Courier
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

	db.Delete(&order)
	c.JSON(http.StatusOK, gin.H{"data": true})
}

// PayOrder godoc
// @Summary Pay Order
// @Description Pay Order by id
// @Tags Buyer
// @Produce json
// @Param id path string true "Order id"
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Success 200 {object} models.Order
// @Router /orders/{id}/pay [post]
func PayOrder(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	// Get model if exist
	var order models.Order
	if err := db.Where("id = ?", c.Param("id")).First(&order).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Order not found!"})
		return
	}
	if order.Status != models.OrderStatusInitial {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Order status is not " + models.OrderStatusInitial})
		return
	}

	_, err := order.RecalculateOrder(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Order Recalculate failed"})
		return
	}

	var updatedInput models.Order
	updatedInput.Status = models.OrderStatusPaid
	updatedInput.PaidAt = time.Now()

	db.Model(&order).Updates(updatedInput)

	c.JSON(http.StatusOK, gin.H{"data": order})
}

// CompleteOrder godoc
// @Summary Complete Order
// @Description Complete Order by id
// @Tags Buyer
// @Produce json
// @Param id path string true "Order id"
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Success 200 {object} models.Order
// @Router /orders/{id}/complete [post]
func CompleteOrder(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	// Get model if exist
	var order models.Order
	if err := db.Where("id = ?", c.Param("id")).First(&order).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Order not found!"})
		return
	}
	if order.Status != models.OrderStatusSent {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Order status is not " + models.OrderStatusSent})
		return
	}

	var updatedInput models.Order
	updatedInput.Status = models.OrderStatusComplete
	updatedInput.CompleteAt = time.Now()

	db.Model(&order).Updates(updatedInput)

	c.JSON(http.StatusOK, gin.H{"data": order})
}

// SendOrder godoc
// @Summary Send Order
// @Description Send Order by id
// @Tags Seller, Admin
// @Produce json
// @Param id path string true "Order id"
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Param Body body orderSendInput true "the body to update order"
// @Success 200 {object} models.Order
// @Router /orders/{id}/send [post]
func SendOrder(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	// Get model if exist
	var order models.Order
	if err := db.Where("id = ?", c.Param("id")).First(&order).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Order not found!"})
		return
	}
	if order.Status != models.OrderStatusPaid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Order status is not " + models.OrderStatusPaid})
		return
	}
	// Validate input
	var input orderSendInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var updatedInput models.Order
	updatedInput.Resi = input.Resi
	updatedInput.Status = models.OrderStatusSent
	updatedInput.CompleteAt = time.Now()

	c.JSON(http.StatusOK, gin.H{"data": order})
}
