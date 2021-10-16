package controllers

import (
	"api-ecommerce/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type userUpdateInput struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Gender    string `json:"gender"`
	Phone     string `json:"phone"`
	PhotoPath string `json:"photo_path"`
	Role      string `json:"role"`
}

type userCreateInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Username string `json:"username"`
	userUpdateInput
}

type changePasswordInput struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

// GetAllUser godoc
// @Summary Get all User
// @Description Get a list of User
// @Tags Admin
// @Produce json
// @Success 200 {object} []models.User
// @Router /users [get]
func GetAllUser(c *gin.Context) {
	// get db from gin context
	db := c.MustGet("db").(*gorm.DB)
	var users []models.User
	db.Find(&users)

	c.JSON(http.StatusOK, gin.H{"data": users})
}

// CreateUser godoc
// @Summary Create new User
// @Description Creating new User
// @Tags Admin
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Param Body body userCreateInput true "the body to create new User"
// @Produce json
// @Success 200 {object} models.User
// @Router /users [post]
func CreateUser(c *gin.Context) {
	// Validate input
	var input userCreateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create user
	user := models.User{
		Email:     input.Email,
		Password:  input.Password,
		Username:  input.Username,
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Gender:    input.Gender,
		Phone:     input.Phone,
		PhotoPath: input.PhotoPath,
		Role:      input.Role,
	}
	db := c.MustGet("db").(*gorm.DB)
	// db.Create(&user)
	_, errUser := user.SaveUser(db)

	if errUser != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}

// GetUserById godoc
// @Summary Get User detail
// @Description Get User by Id
// @Tags Admin
// @Produce json
// @Param id path string true "User id"
// @Success 200 {object} models.User
// @Router /users/{id} [get]
func GetUserById(c *gin.Context) {
	var user models.User

	db := c.MustGet("db").(*gorm.DB)
	if err := db.Where("id =?", c.Param("id")).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}

// UpdateUser godoc
// @Summary Update User
// @Description Update User by id
// @Tags Admin, User
// @Produce json
// @Param id path string true "User id"
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Param Body body userUpdateInput true "the body to update user"
// @Success 200 {object} models.User
// @Router /users/{id} [put]
func UpdateUser(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	// Get model if exist
	var user models.User
	if err := db.Where("id = ?", c.Param("id")).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record no found."})
		return
	}

	// Validate input
	var input userUpdateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var updatedInput models.User
	updatedInput.FirstName = input.FirstName
	updatedInput.LastName = input.LastName
	updatedInput.Gender = input.Gender
	updatedInput.Phone = input.Phone
	updatedInput.PhotoPath = input.PhotoPath
	updatedInput.Role = input.Role
	updatedInput.UpdatedAt = time.Now()

	db.Model(&user).Updates(updatedInput)

	c.JSON(http.StatusOK, gin.H{"data": user})
}

// DeleteUser godoc
// @Summary Delete user
// @Description Delete user by id
// @Tags Admin
// @Produce json
// @Param id path string true "User id"
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Success 200 {object} map[string]boolean
// @Router /users/{id} [delete]
func DeleteUser(c *gin.Context) {
	// Get model if exist
	db := c.MustGet("db").(*gorm.DB)
	var user models.User
	if err := db.Where("id = ?", c.Param("id")).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found"})
		return
	}

	db.Delete(&user)
	c.JSON(http.StatusOK, gin.H{"data": true})
}

// ChangePassword godoc
// @Summary Change Password
// @Description Change Password
// @Tags User
// @Produce json
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Param Body body changePasswordInput true "the body to change password"
// @Success 200 {object} map[string]boolean
// @Router /users/{id}/change-password [put]
func ChangePassword(c *gin.Context) {
	// Get model if exist
	db := c.MustGet("db").(*gorm.DB)
	var user models.User
	if err := db.Where("id = ?", c.Param("id")).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found"})
		return
	}

	// Validate input
	var input changePasswordInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := models.VerifyPassword(input.OldPassword, user.Password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Old password is incorrect"})
		return
	}

	hashedPassword, errPassword := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	if errPassword != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var updatedInput models.User
	updatedInput.Password = string(hashedPassword)
	updatedInput.UpdatedAt = time.Now()
	db.Model(&user).Updates(updatedInput)

	c.JSON(http.StatusOK, gin.H{"data": true})
}
