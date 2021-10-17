package routes

import (
	"api-ecommerce/controllers"
	"api-ecommerce/middlewares"
	"api-ecommerce/models"
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	// set db to gin context
	r.Use(func(c *gin.Context) {
		c.Set("db", db)
	})

	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/swagger/index.html")
	})

	// PUBLIC routes
	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)
	r.GET("/search", controllers.SearchKeyword)

	userAdminRoute := r.Group("/users")
	userAdminRoute.Use(middlewares.AuthMiddleware(models.UserRoleAdmin))
	userAdminRoute.GET("/", controllers.GetAllUser)
	userAdminRoute.POST("/", controllers.CreateUser)
	userAdminRoute.DELETE("/:id", controllers.DeleteUser)

	userRoute := r.Group("/users")
	userRoute.Use(middlewares.AuthMiddleware(models.UserRoleBuyer, models.UserRoleSeller, models.UserRoleAdmin))
	userRoute.GET("/:id", controllers.GetUserById)
	userRoute.PUT("/:id", controllers.UpdateUser)
	userRoute.POST("/change-password", controllers.ChangePassword)

	addressRoute := r.Group("/address")
	addressRoute.Use(middlewares.AuthMiddleware(models.UserRoleBuyer, models.UserRoleSeller))
	addressRoute.GET("/", controllers.GetAllAddress)
	addressRoute.POST("/", controllers.CreateAddress)
	addressRoute.GET("/:id", controllers.GetAddressById)
	addressRoute.PUT("/:id", controllers.UpdateAddress)
	addressRoute.DELETE("/:id", controllers.DeleteAddress)

	categoryRoute := r.Group("/categories")
	categoryRoute.GET("/", controllers.GetAllCategory)
	categoryRoute.GET("/:id/products", controllers.GetProductByCategoryId)
	categoryRoute.Use(middlewares.AuthMiddleware(models.UserRoleAdmin))
	categoryRoute.POST("/", controllers.CreateCategory)
	categoryRoute.PUT("/:id", controllers.UpdateCategory)
	categoryRoute.DELETE("/:id", controllers.DeleteCategory)

	shopRoute := r.Group("/shops")
	shopRoute.GET("/", controllers.GetAllShop)
	shopRoute.GET("/:id", controllers.GetShopById)
	shopRoute.GET("/:id/products", controllers.GetProductByShopId)
	shopRoute.Use(middlewares.AuthMiddleware(models.UserRoleAdmin, models.UserRoleSeller))
	shopRoute.POST("/", controllers.CreateShop)
	shopRoute.PUT("/:id", controllers.UpdateShop)
	shopRoute.DELETE("/:id", controllers.DeleteShop)

	productRoute := r.Group("/products")
	productRoute.GET("/", controllers.GetAllProduct)
	productRoute.GET("/:id", controllers.GetProductById)
	productRoute.Use(middlewares.AuthMiddleware(models.UserRoleAdmin, models.UserRoleSeller))
	productRoute.POST("/", controllers.CreateProduct)
	productRoute.PUT("/:id", controllers.UpdateProduct)
	productRoute.DELETE("/:id", controllers.DeleteProduct)

	addCartRoute := r.Group("/products")
	addCartRoute.Use(middlewares.AuthMiddleware(models.UserRoleBuyer))
	addCartRoute.POST("/products/:id/add-cart", controllers.AddProductOrder)

	cartRoute := r.Group("/cart")
	cartRoute.Use(middlewares.AuthMiddleware(models.UserRoleBuyer))
	cartRoute.PUT("/:id", controllers.UpdateProductOrder)
	cartRoute.DELETE("/:id", controllers.DeleteProductOrder)

	orderRoute := r.Group("/orders")
	orderRoute.Use(middlewares.AuthMiddleware(models.UserRoleBuyer, models.UserRoleSeller, models.UserRoleAdmin))
	orderRoute.GET("/", controllers.GetAllOrder)
	orderRoute.GET("/:id", controllers.GetOrderById)
	orderRoute.POST("/", controllers.CreateOrder)
	orderRoute.PUT("/:id", controllers.UpdateOrder)
	orderRoute.DELETE("/:id", controllers.DeleteOrder)

	orderRoute.POST("/:id/pay", controllers.PayOrder)
	orderRoute.POST("/:id/complete", controllers.CompleteOrder)
	orderRoute.POST("/:id/send", controllers.SendOrder)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}
