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
)

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
	router := route.SetupRouter(userHandler, authHandler, cfg)

	log.Println("Server running on port", cfg.AppPort)
	if err := router.Run(":" + cfg.AppPort); err != nil {
		log.Fatal(err)
	}
}
