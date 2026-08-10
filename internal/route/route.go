package route

import (
	"github.com/gin-gonic/gin"

	"example.com/my-api/config"
	"example.com/my-api/internal/handler"
	"example.com/my-api/internal/middleware"
)

func SetupRouter(userHandler *handler.UserHandler, authHandler *handler.AuthHandler, cfg config.Config) *gin.Engine {
	router := gin.Default()

	api := router.Group("/api")
	api.POST("/login", authHandler.Login)
	api.POST("/register", authHandler.Register)

	api.GET("/user-info", middleware.Auth(cfg), userHandler.GetUserInfo)

	users := api.Group("/users", middleware.Auth(cfg))
	{
		users.POST("", userHandler.CreateUser)
		users.GET("", userHandler.GetUsers)
		users.GET("/:id", userHandler.GetUserById)
		users.PUT("/:id", userHandler.UpdateUser)
		users.DELETE("/:id", userHandler.DeleteUser)
	}

	return router
}
