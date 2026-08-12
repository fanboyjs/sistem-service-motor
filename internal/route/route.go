package route

import (
	"github.com/gin-gonic/gin"

	"example.com/my-api/config"
	"example.com/my-api/internal/handler"
	"example.com/my-api/internal/middleware"
)

func SetupRouter(userHandler *handler.UserHandler, authHandler *handler.AuthHandler, vehicleBrandHandler *handler.VehicleBrandHandler, vehicleModelHandler *handler.VehicleModelHandler, vehicleHandler *handler.VehicleHandler, vehicleTaxHandler *handler.VehicleTaxHandler, vehicleQRHandler *handler.VehicleQRHandler, serviceTypeHandler *handler.ServiceTypeHandler, serviceItemHandler *handler.ServiceItemHandler, serviceRecordHandler *handler.ServiceRecordHandler, cfg config.Config) *gin.Engine {
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

	serviceTypes := api.Group("/service-types", middleware.Auth(cfg))
	{
		serviceTypes.POST("", serviceTypeHandler.CreateServiceType)
		serviceTypes.GET("", serviceTypeHandler.GetServiceTypes)
		serviceTypes.GET("/:id", serviceTypeHandler.GetServiceTypeById)
		serviceTypes.PUT("/:id", serviceTypeHandler.UpdateServiceType)
	}

	serviceItems := api.Group("/service-items", middleware.Auth(cfg))
	{
		serviceItems.POST("", serviceItemHandler.CreateServiceItem)
		serviceItems.GET("", serviceItemHandler.GetServiceItems)
		serviceItems.GET("/:id", serviceItemHandler.GetServiceItemById)
		serviceItems.PUT("/:id", serviceItemHandler.UpdateServiceItem)
	}

	serviceRecords := api.Group("/service-records", middleware.Auth(cfg))
	{
		serviceRecords.POST("", serviceRecordHandler.CreateServiceRecord)
		serviceRecords.GET("", serviceRecordHandler.GetServiceRecords)
		serviceRecords.GET("/:id", serviceRecordHandler.GetServiceRecordById)
		serviceRecords.PUT("/:id", serviceRecordHandler.UpdateServiceRecord)
		serviceRecords.DELETE("/:id", serviceRecordHandler.DeleteServiceRecord)
	}

	return router
}
