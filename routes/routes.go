package routes

import (
	"api-ecommerce/controllers"
	"api-ecommerce/middlewares"

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

	// PUBLIC routes
	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)

	userRoute := r.Group("/users")
	userRoute.GET("/", controllers.GetAllUser)
	userRoute.GET("/:id", controllers.GetUserById)
	userRoute.Use(middlewares.JwtAuthMiddleware())
	userRoute.POST("/", controllers.CreateUser)
	userRoute.PUT("/:id", controllers.UpdateUser)
	userRoute.DELETE("/:id", controllers.DeleteUser)
	userRoute.POST("/:id/change-password", controllers.ChangePassword)

	userRoute.GET("/:id/address/", controllers.GetAllAddress)
	userRoute.GET("/:id/address/:address_id", controllers.GetAddressById)
	userRoute.POST("/:id/address/", controllers.CreateAddress)
	userRoute.PUT("/:id/address/:address_id", controllers.UpdateAddress)
	userRoute.DELETE("/:id/address/:address_id", controllers.DeleteAddress)

	categoryRoute := r.Group("/categories")
	categoryRoute.GET("/", controllers.GetAllCategory)
	categoryRoute.GET("/:id/products", controllers.GetProductByCategoryId)
	categoryRoute.Use(middlewares.JwtAuthMiddleware())
	categoryRoute.POST("/", controllers.CreateCategory)
	categoryRoute.PUT("/:id", controllers.UpdateCategory)
	categoryRoute.DELETE("/:id", controllers.DeleteCategory)

	shopRoute := r.Group("/shops")
	shopRoute.GET("/", controllers.GetAllShop)
	shopRoute.GET("/:id", controllers.GetShopById)
	shopRoute.GET("/:id/products", controllers.GetProductByShopId)
	shopRoute.Use(middlewares.JwtAuthMiddleware())
	shopRoute.POST("/", controllers.CreateShop)
	shopRoute.PUT("/:id", controllers.UpdateShop)
	shopRoute.DELETE("/:id", controllers.DeleteShop)

	productRoute := r.Group("/products")
	productRoute.GET("/", controllers.GetAllProduct)
	productRoute.GET("/:id", controllers.GetProductById)
	productRoute.Use(middlewares.JwtAuthMiddleware())
	productRoute.POST("/", controllers.CreateProduct)
	productRoute.PUT("/:id", controllers.UpdateProduct)
	productRoute.DELETE("/:id", controllers.DeleteProduct)

	productRoute.POST("/:id/add-cart", controllers.AddProductOrder)

	cartRoute := r.Group("/cart")
	cartRoute.Use(middlewares.JwtAuthMiddleware())
	cartRoute.PUT("/:id", controllers.UpdateProductOrder)
	cartRoute.DELETE("/:id", controllers.DeleteProductOrder)

	orderRoute := r.Group("/orders")
	orderRoute.Use(middlewares.JwtAuthMiddleware())
	orderRoute.GET("/", controllers.GetAllOrder)
	orderRoute.GET("/:id", controllers.GetOrderById)
	orderRoute.POST("/", controllers.CreateOrder)
	orderRoute.PUT("/:id", controllers.UpdateOrder)
	orderRoute.DELETE("/:id", controllers.DeleteOrder)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}
