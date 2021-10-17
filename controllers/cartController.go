package controllers

import (
	"api-ecommerce/models"
	"api-ecommerce/utils/token"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type cartInput struct {
	Qty  int    `json:"qty"`
	Note string `json:"note"`
}

// AddProductOrder godoc
// @Summary Add Product to Order
// @Description Add Product to Order
// @Tags Buyer
// @Produce json
// @Param id path string true "Product id"
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Success 200 {object} models.Order
// @Router /products/{id}/add-cart [post]
func AddProductOrder(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var product models.Product
	if err := db.Where("id = ?", c.Param("id")).First(&product).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	userId, err := token.ExtractTokenID(c)
	if userId == 0 || err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token!"})
		return
	}

	var order models.Order
	if err := db.Where("user_id=? AND status=?", userId, models.OrderStatusInitial).First(&order).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Order not found, Please create initial order!"})
		return
	}

	var input cartInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.Qty > product.Stock {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Product stock not enough"})
		return
	}

	productTotal := input.Qty * product.Price
	orderProduct := models.OrderProduct{
		Qty:       input.Qty,
		Note:      input.Note,
		ProductID: product.ID,
		OrderID:   order.ID,
		Total:     productTotal,
	}
	db.Create(&orderProduct)
	// Update product stock and Grand total of order
	if _, err = product.SubtractStock(input.Qty, db); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Update Product stock failed"})
		return
	}
	if _, err = order.RecalculateOrder(db); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Order Recalculation failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": order})
}

// UpdateProductOrder godoc
// @Summary Update Product from Order
// @Description Update Product from Order
// @Tags Buyer
// @Produce json
// @Param id path string true "OrderProduct id"
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Param Body body cartInput true "the body to update product"
// @Success 200 {object} models.OrderProduct
// @Router /cart/{id}/ [put]
func UpdateProductOrder(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	// Get model if exist
	var orderProduct models.OrderProduct
	if err := db.Where("id = ?", c.Param("id")).First(&orderProduct).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Order Product no found."})
		return
	}
	var product models.Product
	if err := db.Where("id = ?", orderProduct.ProductID).First(&product).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Product not found!"})
		return
	}
	var order models.Order
	if err := db.Where("id=?", orderProduct.OrderID).First(&order).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Order not found!"})
		return
	}

	// Validate input
	var input cartInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updatedQty := input.Qty - orderProduct.Qty
	if updatedQty > product.Stock {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Product stock not enough"})
		return
	}

	productTotal := input.Qty * product.Price
	var updatedInput models.OrderProduct
	updatedInput.Qty = input.Qty
	updatedInput.Note = input.Note
	updatedInput.Total = productTotal
	updatedInput.UpdatedAt = time.Now()

	db.Model(&orderProduct).Updates(updatedInput)
	// Update product stock and Grand total of order
	if _, err := product.SubtractStock(updatedQty, db); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Update Product stock failed"})
		return
	}
	if _, err := order.RecalculateOrder(db); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Order Recalculation failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": orderProduct})
}

// DeleteProductOrder godoc
// @Summary Delete Product from Order
// @Description Delete Product from Order
// @Tags Buyer
// @Produce json
// @Param id path string true "OrderProduct id"
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Success 200 {object} map[string]boolean
// @Router /cart/{id}/ [delete]
func DeleteProductOrder(c *gin.Context) {
	// Get model if exist
	db := c.MustGet("db").(*gorm.DB)
	var orderProduct models.OrderProduct
	if err := db.Where("id = ?", c.Param("id")).First(&orderProduct).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found"})
		return
	}
	var order models.Order
	if err := db.Where("id=?", orderProduct.OrderID).First(&order).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Order not found!"})
		return
	}

	db.Delete(&orderProduct)
	// Recalculate grand total order
	_, err := order.RecalculateOrder(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Order Recalculation failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": true})
}
