package controllers

import (
	"api-ecommerce/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type addressInput struct {
	Address   string `json:"address"`
	City      string `json:"city"`
	Country   string `json:"country"`
	Province  string `json:"province"`
	Zip       uint   `json:"zip"`
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
	UserID    uint   `json:"user_id"`
}

// GetAllAddress godoc
// @Summary Get all Address
// @Description Get a list of Address
// @Tags User
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Param id path string true "User id"
// @Produce json
// @Success 200 {object} []models.Address
// @Router /users/{id}/address [get]
func GetAllAddress(c *gin.Context) {
	// get db from gin context
	db := c.MustGet("db").(*gorm.DB)
	var address []models.Address
	if err := db.Where("user_id = ?", c.Param("id")).Find(&address).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": address})
}

// CreateAddress godoc
// @Summary Create new Address
// @Description Creating new Address
// @Tags User
// @Param id path string true "User id"
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Param Body body addressInput true "the body to create new Address"
// @Produce json
// @Success 200 {object} models.Address
// @Router /users/{id}/address [post]
func CreateAddress(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var user models.User
	if err := db.Where("id = ?", c.Param("id")).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found"})
		return
	}

	// Validate input
	var input addressInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create address
	address := models.Address{
		Address:   input.Address,
		City:      input.City,
		Country:   input.Country,
		Province:  input.Province,
		Zip:       input.Zip,
		Latitude:  input.Latitude,
		Longitude: input.Longitude,
		UserID:    user.ID,
	}
	db.Create(&address)

	c.JSON(http.StatusOK, gin.H{"data": address})
}

// GetAddressById godoc
// @Summary Get Address detail
// @Description Get Address by Id
// @Tags User
// @Produce json
// @Param id path string true "User id"
// @Param address_id path string true "Address id"
// @Success 200 {object} models.Address
// @Router /users/{id}/address/{address_id} [get]
func GetAddressById(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var user models.User
	if err := db.Where("id = ?", c.Param("id")).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found"})
		return
	}

	var address models.Address
	if err := db.Where("id = ? AND user_id = ?", c.Param("address_id"), c.Param("id")).First(&address).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": address})
}

// UpdateAddress godoc
// @Summary Update Address
// @Description Update Address by id
// @Tags User
// @Produce json
// @Param id path string true "User id"
// @Param address_id path string true "Address id"
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Param Body body addressInput true "the body to update address"
// @Success 200 {object} models.Address
// @Router /users/{id}/address/{address_id} [put]
func UpdateAddress(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	// Get model if exist
	var user models.User
	if err := db.Where("id = ?", c.Param("id")).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found"})
		return
	}
	var address models.Address
	if err := db.Where("id = ?", c.Param("address_id")).First(&address).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}
	if user.ID != address.UserID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not the owner of this data"})
		return
	}

	// Validate input
	var input addressInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var updatedInput models.Address
	updatedInput.Address = input.Address
	updatedInput.City = input.City
	updatedInput.Country = input.Country
	updatedInput.Province = input.Province
	updatedInput.Zip = input.Zip
	updatedInput.Latitude = input.Latitude
	updatedInput.Longitude = input.Longitude
	updatedInput.UpdatedAt = time.Now()

	db.Model(&address).Updates(updatedInput)

	c.JSON(http.StatusOK, gin.H{"data": address})
}

// DeleteAddress godoc
// @Summary Delete address
// @Description Delete address by id
// @Tags User
// @Produce json
// @Param id path string true "User id"
// @Param address_id path string true "Address id"
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Success 200 {object} map[string]boolean
// @Router /users/{id}/address/{address_id} [delete]
func DeleteAddress(c *gin.Context) {
	// Get model if exist
	db := c.MustGet("db").(*gorm.DB)
	var user models.User
	if err := db.Where("id = ?", c.Param("id")).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found"})
		return
	}
	var address models.Address
	if err := db.Where("id = ?", c.Param("address_id")).First(&address).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}
	if user.ID != address.UserID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not the owner of this data"})
		return
	}

	db.Delete(&address)
	c.JSON(http.StatusOK, gin.H{"data": true})
}
