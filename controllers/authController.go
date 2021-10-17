package controllers

import (
	"api-ecommerce/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type LoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required"`
}

// LoginUser godoc
// @Summary Login user
// @Description Logging in to get jwt token to access admin or user api by roles
// @Tags Public
// @Param Body body LoginInput true "the body to login user"
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /login [post]
func Login(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var input LoginInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u := models.User{}

	u.Username = input.Username
	u.Password = input.Password

	token, err := models.LoginCheck(u.Username, u.Password, db)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username or Password is incorrect."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Login Success", "token": token})
}

// Register godoc
// @Summary Register a user
// @Description registering user from public access
// @Tags Public
// @Param Body body RegisterInput true "the body to register user"
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /register [post]
func Register(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var input RegisterInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u := models.User{}
	u.Username = input.Username
	u.Email = input.Email
	u.Password = input.Password
	u.Role = models.UserRoleBuyer

	_, err := u.SaveUser(db)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Registation success", "user": u})
}
