package route

import (
	"github.com/gin-gonic/gin"

	"example.com/my-api/config"
	"example.com/my-api/internal/handler"
	"example.com/my-api/internal/middleware"
)

func SetupRouter(userHandler *handler.UserHandler, authHandler *handler.AuthHandler, vehicleBrandHandler *handler.VehicleBrandHandler, vehicleModelHandler *handler.VehicleModelHandler, vehicleHandler *handler.VehicleHandler, vehicleTaxHandler *handler.VehicleTaxHandler, vehicleQRHandler *handler.VehicleQRHandler, cfg config.Config) *gin.Engine {
	router := gin.Default()

	router.Static("/uploads", cfg.UploadDir)

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
	
	vehicleBrands := api.Group("/vehicle-brands", middleware.Auth(cfg))
	{
		vehicleBrands.POST("", vehicleBrandHandler.CreateVehicleBrand)
		vehicleBrands.GET("", vehicleBrandHandler.GetVehicleBrands)
		vehicleBrands.GET("/:id", vehicleBrandHandler.GetVehicleBrandById)
		vehicleBrands.PUT("/:id", vehicleBrandHandler.UpdateVehicleBrand)
	}

	vehicleModels := api.Group("/vehicle-models", middleware.Auth(cfg))
	{
		vehicleModels.POST("", vehicleModelHandler.CreateVehicleModel)
		vehicleModels.GET("", vehicleModelHandler.GetVehicleModels)
		vehicleModels.GET("/:id", vehicleModelHandler.GetVehicleModelById)
		vehicleModels.PUT("/:id", vehicleModelHandler.UpdateVehicleModel)
	}
	
	vehicles := api.Group("/vehicles", middleware.Auth(cfg))
	{
		vehicles.POST("", vehicleHandler.CreateVehicle)
		vehicles.GET("", vehicleHandler.GetVehicles)
		vehicles.GET("/:id", vehicleHandler.GetVehicleByID)
		vehicles.PUT("/:id", vehicleHandler.UpdateVehicle)
		vehicles.DELETE("/:id", vehicleHandler.DeleteVehicle)
		vehicles.GET("/qr/:id", vehicleQRHandler.GetVehicleByQRCode) // qr code vehicle
	}

	vehicleTaxes := api.Group("/vehicle-taxes", middleware.Auth(cfg))
	{
		vehicleTaxes.POST("", vehicleTaxHandler.CreateVehicleTax)
		vehicleTaxes.GET("", vehicleTaxHandler.GetVehicleTaxes)
		vehicleTaxes.GET("/:id", vehicleTaxHandler.GetVehicleTaxById)
		vehicleTaxes.PUT("/:id", vehicleTaxHandler.UpdateVehicleTax)
		vehicleTaxes.DELETE("/:id", vehicleTaxHandler.DeleteVehicleTax)
	}

	return router
}
