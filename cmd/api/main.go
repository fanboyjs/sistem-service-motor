package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"

	"example.com/my-api/config"
	"example.com/my-api/internal/database"
	"example.com/my-api/internal/handler"
	"example.com/my-api/internal/repository"
	"example.com/my-api/internal/route"
	"example.com/my-api/internal/service"
	"example.com/my-api/internal/storage"
)

// @title Sistem Service Motor API
// @version 1.0
// @description Backend API untuk sistem service motor. Sebagian besar endpoint membutuhkan token JWT (Bearer) yang diperoleh dari login/register.
// @host localhost:8080
// @BasePath /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Masukkan token JWT dengan format: `Bearer <token>`

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println(".env not found")
	}

	if envFile := ".env." + os.Getenv("APP_ENV"); envFile != ".env." {
		if _, err := os.Stat(envFile); err == nil {
			if err := godotenv.Overload(envFile); err != nil {
				log.Println("failed to load", envFile)
			}
		}
	}

	cfg := config.Load()

	db, err := database.NewPostgres(cfg.DatabaseURL())
	if err != nil {
		log.Fatal("database connection failed:", err)
	}
	defer db.Close()

	userRepository := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepository)
	userHandler := handler.NewUserHandler(userService)
	authService := service.NewAuthService(userRepository, cfg)
	authHandler := handler.NewAuthHandler(authService)
	vehicleBrandRepository := repository.NewVehicleBrandRepository(db)
	vehicleBrandStorage := storage.NewLocalStorage(cfg.UploadDir)
	vehicleBrandService := service.NewVehicleBrandService(vehicleBrandRepository, vehicleBrandStorage)
	vehicleBrandHandler := handler.NewVehicleBrandHandler(vehicleBrandService)
	vehicleModelRepository := repository.NewVehicleModelRepository(db)
	vehicleModelService := service.NewVehicleModelService(vehicleModelRepository)
	vehicleModelHandler := handler.NewVehicleModelHandler(vehicleModelService)
	vehicleRepository := repository.NewVehicleRepository(db)
	vehicleStorage := storage.NewLocalStorage(cfg.UploadDir)
	vehicleQRRepository := repository.NewVehicleQRRepository(db)
	vehicleQRService := service.NewVehicleQRService(vehicleQRRepository)
	vehicleService := service.NewVehicleService(vehicleRepository, vehicleStorage)
	vehicleHandler := handler.NewVehicleHandler(vehicleService)
	vehicleQRHandler := handler.NewVehicleQRHandler(vehicleQRService)
	vehicleTaxRepository := repository.NewVehicleTaxRepository(db)
	vehicleTaxService := service.NewVehicleTaxService(vehicleTaxRepository)
	vehicleTaxHandler := handler.NewVehicleTaxHandler(vehicleTaxService)
	serviceTypeRepository := repository.NewServiceTypeRepository(db)
	serviceTypeService := service.NewServiceTypeService(serviceTypeRepository)
	serviceTypeHandler := handler.NewServiceTypeHandler(serviceTypeService)
	serviceItemRepository := repository.NewServiceItemRepository(db)
	serviceItemService := service.NewServiceItemService(serviceItemRepository)
	serviceItemHandler := handler.NewServiceItemHandler(serviceItemService)
	serviceRecordRepository := repository.NewServiceRecordRepository(db)
	serviceRecordService := service.NewServiceRecordService(serviceRecordRepository, serviceItemRepository)
	serviceRecordHandler := handler.NewServiceRecordHandler(serviceRecordService)
	router := route.SetupRouter(userHandler, authHandler, vehicleBrandHandler, vehicleModelHandler, vehicleHandler, vehicleTaxHandler, vehicleQRHandler, serviceTypeHandler, serviceItemHandler, serviceRecordHandler, cfg)

	log.Println("Server running on port", cfg.AppPort)
	if err := router.Run(":" + cfg.AppPort); err != nil {
		log.Fatal(err)
	}
}
