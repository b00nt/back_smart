package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	// "back/internal/handlers"
	"back/internal/models"
	// "back/internal/moysklad"
	// "back/internal/routes"
	// "back/internal/services"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	e := echo.New()

	// Middleware configuration
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3000"},
		AllowMethods: []string{echo.GET, echo.POST, echo.OPTIONS},
		AllowHeaders: []string{"Content-Type"},
	}))

	e.Static("/static", "static")

	// set body limit and 20 requests per minute
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(20)))
	e.Use(middleware.BodyLimit("10M"))

	// Database connection
	db, err := Connect()
	if err != nil {
		log.Printf("DB connection error: %s", err)
		os.Exit(1)
	}

	db.AutoMigrate(&models.Product{}, &models.Modification{}, &models.ModificationImage{}, &models.ProductImage{}, &models.ModificationCharacteristic{})
	db.AutoMigrate(&models.Feedback{}, &models.OrderItem{}, &models.CustomerInfo{}, &models.Order{}, &models.ModificationCharacteristicOrder{})

	// setup routes
	// routes.SetupRoutes(e, db)

	// Initialize the CronService
	// cronService := services.NewCronService()
	// fmt.Println(cronService)

	// Start the cron job
	// cronService.Start("saratov", db)
	// cronService.Start("moscow", db)

	// get & save products
	// resultSaratovProduct, err := moysklad.GetProducts("saratov")
	// if err != nil {
	// 	log.Printf("Error get product: %s", err)
	// }
	//
	// fmt.Println(resultSaratovProduct)
	//
	// err = moysklad.SaveProducts("saratov", resultSaratovProduct, db)
	// if err != nil {
	// 	fmt.Errorf("Error updating product:", err)
	// }
	// resultMoscowProduct := moysklad.GetProducts("moscow")
	// err = moysklad.SaveProducts("moscow", resultMoscowProduct, db)
	// if errMoscowProduct != nil {
	// 	log.Errorf("Error updating product:", err)
	// } else {
	// 	fmt.Println("Product update successful")
	// }

	// get & save modifications
	// err = moysklad.UpdateAllStocks("saratov", db)
	// if err != nil {
	// 	log.Printf("Error get stocks: %s", err)
	// }

	// get & save stock
	// resultSaratovStock := moysklad.GetStock("saratov")
	// errSaratovStock := moysklad.SaveStock("saratov", resultSaratovStock, db)
	// if errSaratovStock != nil {
	// 	log.Println("Error updating stock:", err)
	// } else {
	// 	fmt.Println("Stock update successful")
	// }
	//
	// resultMoscowStock := moysklad.GetStock("moscow")
	// errMoscowStock := moysklad.SaveStock("moscow", resultMoscowStock, db)
	// if errMoscowStock != nil {
	// 	log.Println("Error updating stock:", err)
	// } else {
	// 	fmt.Println("Stock update successful")
	// }

	// moysklad.GetSaveDownloadProductImages("saratov", db)
	// moysklad.GetSaveDownloadModImages("saratov", db)

	// Channel to listen for interrupt or terminate signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	// Start the server in a goroutine
	go func() {
		if err := e.Start(":8080"); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	// Wait for an interrupt signal
	<-quit
	log.Println("Shutting down server...")

	// Create a deadline to wait for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Attempt to gracefully shut down the server
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}

	// Close the database connection
	sqlDB, err := db.DB()
	if err != nil {
		log.Println("Failed to close database connection:", err)
		os.Exit(1)
	} else {
		sqlDB.Close()
		log.Println("Database connection closed.")
	}
}

func Connect() (*gorm.DB, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dsn := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=%s TimeZone=%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SSLMODE"),
		os.Getenv("DB_TIMEZONE"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	return db, nil
}
