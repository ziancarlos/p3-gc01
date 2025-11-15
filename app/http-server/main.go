package main

import (
	"context"
	"log"
	"p3-graded-challenge-1-ziancarlos/config"
	"p3-graded-challenge-1-ziancarlos/controllers"
	"p3-graded-challenge-1-ziancarlos/repository"
	"p3-graded-challenge-1-ziancarlos/service"
	"time"

	_ "p3-graded-challenge-1-ziancarlos/docs"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// @title Shopping Service API
// @version 1.0
// @description API documentation for Shopping Service
// @host localhost:9051
// @BasePath /
func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Initialize Echo
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Swagger docs route
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// MongoDB connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.Database.MongoURI))
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}

	defer func() {
		if err = client.Disconnect(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()

	// Test connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Failed to ping MongoDB:", err)
	}
	log.Println("✓ Connected to MongoDB!")

	// Initialize database
	db := client.Database(cfg.Database.DBName)

	// Initialize repositories
	productRepo := repository.NewProductRepository(db)
	transactionRepo := repository.NewTransactionRepository(db, cfg)

	// Initialize services
	productService := service.NewProductService(productRepo)
	transactionService := service.NewTransactionService(transactionRepo, cfg)

	// Initialize controllers
	productController := controllers.NewProductController(productService)
	transactionController := controllers.NewTransactionController(transactionService)

	// Routes - Products
	e.POST("/products", productController.CreateProduct)
	e.GET("/products", productController.GetAllProducts)
	e.GET("/products/:id", productController.GetProductByID)
	e.PUT("/products/:id", productController.UpdateProduct)
	e.DELETE("/products/:id", productController.DeleteProduct)

	// Routes - Transactions
	e.POST("/transactions", transactionController.CreateTransaction)
	e.GET("/transactions", transactionController.GetAllTransactions)
	e.GET("/transactions/:id", transactionController.GetTransactionByID)
	e.PUT("/transactions/:id", transactionController.UpdateTransaction)
	e.DELETE("/transactions/:id", transactionController.DeleteTransaction)

	// Start server
	log.Printf("✓ Shopping Service running on port %s", cfg.Server.Port)

	e.Logger.Fatal(e.Start(":" + cfg.Server.Port))
}

