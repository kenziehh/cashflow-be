package main

import (
	"log"
	"os"

	"github.com/kenziehh/cashflow-be/config"
	"github.com/kenziehh/cashflow-be/database/seed"
	_ "github.com/kenziehh/cashflow-be/docs"
	"github.com/kenziehh/cashflow-be/internal/domain/auth/handler/http"
	authRepo "github.com/kenziehh/cashflow-be/internal/domain/auth/repository"
	authService "github.com/kenziehh/cashflow-be/internal/domain/auth/service"
	transactionHandler "github.com/kenziehh/cashflow-be/internal/domain/transaction/handler/http"
	transactionRepo "github.com/kenziehh/cashflow-be/internal/domain/transaction/repository"
	transactionService "github.com/kenziehh/cashflow-be/internal/domain/transaction/service"

	categoryHandler "github.com/kenziehh/cashflow-be/internal/domain/category/handler/http"
	categoryRepo "github.com/kenziehh/cashflow-be/internal/domain/category/repository"
	categoryService "github.com/kenziehh/cashflow-be/internal/domain/category/service"
	"github.com/kenziehh/cashflow-be/internal/middleware"

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

	// Run seeders
	if err := seed.SeedCategoriesIfEmpty(db); err != nil {
		log.Fatal("❌ Seeder failed:", err)
	}

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

	transactionRepository := transactionRepo.NewTransactionRepository(db, redis)
	transactionSvc := transactionService.NewTransactionService(transactionRepository)
	transactionHandler := transactionHandler.NewTransactionHandler(transactionSvc)

	transactions := api.Group("/transactions", middleware.JWTAuth())
	transactions.Post("/", transactionHandler.CreateTransaction)
	transactions.Get("/summary", transactionHandler.GetSummaryTransaction)
	transactions.Get("/:id", transactionHandler.GetTransactionByID)
	transactions.Get("/:id/proof", transactionHandler.GetProofFile)
	transactions.Get("/", transactionHandler.GetTransactionsWithPagination)
	transactions.Put("/:id", transactionHandler.UpdateTransaction)
	transactions.Delete("/:id", transactionHandler.DeleteTransaction)

	categoryRepository := categoryRepo.NewCategoryRepository(db, redis)
	categorySvc := categoryService.NewCategoryService(categoryRepository)
	categoryHandler := categoryHandler.NewCategoryHandler(categorySvc)

	categories := api.Group("/categories", middleware.JWTAuth())
	categories.Get("/", categoryHandler.GetAllCategories)

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
