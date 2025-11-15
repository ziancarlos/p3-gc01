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

// @title Payment Service API
// @version 1.0
// @description API documentation for Payment Service
// @host localhost:9061
// @BasePath /
func main() {
	cfg, err := config.LoadPaymentConfig()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Swagger docs route
	e.GET("/swagger/*any", echoSwagger.WrapHandler)

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

	db := client.Database(cfg.Database.DBName)
	paymentRepo := repository.NewPaymentRepository(db)
	paymentService := service.NewPaymentService(paymentRepo)
	paymentController := controllers.NewPaymentController(paymentService)

	e.POST("/payments", paymentController.CreatePayment)

	log.Printf("✓ Payment Service running on port %s", cfg.Server.Port)
	e.Logger.Fatal(e.Start(":" + cfg.Server.Port))
}

