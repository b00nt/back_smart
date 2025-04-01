package main

import (
	"back/internal/database"
	"back/internal/handlers"
	"back/internal/models"
	// "back/internal/moysklad"
	"back/internal/routes"
	"back/internal/services"
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	// "log"
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

	// Database connection
	db, err := database.Connect()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	db.AutoMigrate(&models.Category{})
	db.AutoMigrate(&models.Products{}, &models.Modification{}, &models.ModificationImages{}, &models.ProductImages{}, &models.ModificationCharacteristics{})
	db.AutoMigrate(&models.ProductsSaratov{}, &models.ModificationSaratov{}, &models.ModificationImagesSaratov{}, &models.ProductImagesSaratov{}, &models.ModificationCharacteristicsSaratov{})
	db.AutoMigrate(&models.Feedback{}, &models.OrderItem{}, &models.CustomerInfo{}, &models.Order{}, &models.ModificationCharacteristicsOrder{})

	// Initialize the CronService
	cronService := services.NewCronService()
	fmt.Println(cronService)

	// Start the cron job
	// cronService.Start("saratov", db)
	// cronService.Start("moscow", db)

	// get & save products 
	// resultSaratovProduct := moysklad.GetProducts("saratov")
	// err = moysklad.SaveProducts("saratov", resultSaratovProduct, db)
	// if err != nil {
	// 	log.Errorf("Error updating product:", err)
	// } else {
	// 	fmt.Println("Product update successful")
	// }
	//
	// resultMoscowProduct := moysklad.GetProducts("moscow")
	// err = moysklad.SaveProducts("moscow", resultMoscowProduct, db)
	// if errMoscowProduct != nil {
	// 	log.Errorf("Error updating product:", err)
	// } else {
	// 	fmt.Println("Product update successful")
	// }

	// get & save modifications
	// resultSaratovMod := moysklad.GetModifications("saratov", db)
	// err = moysklad.SaveModifications("saratov", resultSaratovMod, db)
	// if err != nil {
	// 	log.Errorf("Error updating modification:", err)
	// } else {
	// 	fmt.Println("Modification update successful")
	// }
	//
	// resultMoscowMod := moysklad.GetModifications("moscow", db)
	// err = moysklad.SaveModifications("moscow", resultMoscowMod, db)
	// if err != nil {
	// 	log.Errorf("Error updating modification:", err)
	// } else {
	// 	fmt.Println("Modification update successful")
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


	// Initialize handler with DB instance
	handler := handlers.NewHandler(db)
	routes.InitProductsRoutes(e, handler)
	routes.InitFeedbackRoutes(e, handler)
	routes.InitOrderRoutes(e, handler)

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
	fmt.Println("Shutting down server...")

	// Create a deadline to wait for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Attempt to gracefully shut down the server
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}

	// Close the database connection
	sqlDB, err := db.DB()
	if err == nil {
		sqlDB.Close()
		fmt.Println("Database connection closed.")
	} else {
		fmt.Println("Failed to close database connection:", err)
	}
}
