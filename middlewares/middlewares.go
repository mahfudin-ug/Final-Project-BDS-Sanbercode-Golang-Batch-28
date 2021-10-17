package middlewares

import (
	"api-ecommerce/models"
	"api-ecommerce/utils/token"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AuthMiddleware(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := token.TokenValid(c)
		if err != nil {
			c.String(http.StatusUnauthorized, err.Error())
			c.Abort()
			return
		}

		if cap(roles) > 0 {
			userId, err := token.ExtractTokenID(c)
			if userId == 0 || err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token!"})
				c.Abort()
				return
			}
			db := c.MustGet("db").(*gorm.DB)
			var user models.User
			if err := db.Where("id = ?", userId).First(&user).Error; err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "User not found."})
				c.Abort()
				return
			}

			for _, validRole := range roles {
				if user.Role == validRole {
					c.Next()
					return
				}
			}

			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden!"})
			c.Abort()
			return
		}

		c.Next()
	}
}
