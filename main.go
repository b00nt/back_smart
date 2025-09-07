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
	"back/internal/config"
	"back/internal/db"
	"back/internal/models"
	// "back/internal/moysklad"
	"back/internal/routes"
	// "back/internal/services"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	// middleware configuration
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3000"},
		AllowMethods: []string{echo.GET, echo.POST, echo.OPTIONS},
		AllowHeaders: []string{"Content-Type"},
	}))

	// static setup
	e.Static("/static", "static")

	// set body limit and 20 requests per minute
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(20)))
	e.Use(middleware.BodyLimit("10M"))

	// load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Error loading configuration:", err)
	}

	// postgres setup
	err = db.setupDatabaseAndUser(cfg)
	if err != nil {
		log.Fatal("setup failed:", err)
	}

	// database connection
	db, err := connect(cfg)
	if err != nil {
		log.Fatal("connect failed:", err)
	}

	// db migrate
	err = db.AutoMigrate(
		&models.Product{},
		&models.Modification{},
		&models.ModificationImage{},
		&models.Image{},
		&models.Characteristic{},
		&models.Feedback{},
		&models.OrderItem{},
		&models.CustomerInfo{},
		&models.Order{},
		&models.CharacteristicOrder{},
	)
	if err != nil {
		log.Fatal("Failed to auto-migrate database:", err)
	}

	// initialize handler with DB instance
	routes.SetupRoutes(e, db)

	// channel to listen for interrupt or terminate signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	// start the server in a goroutine
	go func() {
		if err := e.Start(":8080"); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	// wait for an interrupt signal
	<-quit
	log.Println("Shutting down server...")

	// create a deadline to wait for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// attempt to gracefully shut down the server
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}

	// close the database connection
	if posDB, err := db.DB(); err == nil {
		if err := posDB.Close(); err != nil {
			log.Fatal("Failed to close database connection: ", err)
		} else {
			fmt.Println("Database connection closed successfully")
		}
	} else {
		log.Fatal("Failed to get database instance: ", err)
	}
}
