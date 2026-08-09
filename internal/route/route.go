package route

import (
	"github.com/gin-gonic/gin"

	"example.com/my-api/internal/handler"
)

func SetupRouter(userHandler *handler.UserHandler) *gin.Engine {
	router := gin.Default()

	api := router.Group("/api")
	users := api.Group("/users")
	{
		users.GET("", userHandler.GetUsers)
		users.GET("/:id", userHandler.GetUserById)
		users.DELETE("/:id", userHandler.DeleteUser)
	}

	return router
}
