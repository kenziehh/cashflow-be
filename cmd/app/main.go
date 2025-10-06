package main

import (
	"log"
	"os"

	"devsecops-be/config"
	_ "devsecops-be/docs"
	"devsecops-be/internal/domain/auth/handler/http"
	authRepo "devsecops-be/internal/domain/auth/repository"
	authService "devsecops-be/internal/domain/auth/service"
	"devsecops-be/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

// @title Cash Flow API
// @version 1.0
// @description API untuk Website Cash Flow
// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// Load config
	cfg := config.LoadConfig()

	// Initialize database
	db := config.InitDB(cfg)
	defer db.Close()

	// Initialize Redis
	redis := config.InitRedis(cfg)
	defer redis.Close()

	// Initialize Fiber
	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.ErrorHandler,
	})

	// Middleware
	app.Use(middleware.Logger())

	// Swagger
	app.Get("/docs/*", swagger.HandlerDefault)

	// Routes
	api := app.Group("/api/v1")

	// Auth routes
	authRepository := authRepo.NewAuthRepository(db, redis)
	authSvc := authService.NewAuthService(authRepository)
	authHandler := http.NewAuthHandler(authSvc)

	auth := api.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)
	auth.Post("/logout", middleware.JWTAuth(), authHandler.Logout)
	auth.Get("/me", middleware.JWTAuth(), authHandler.GetProfile)

	// Start server
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server running on port %s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatal(err)
	}
}
