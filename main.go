package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	// "back/internal/handlers"
	"back/internal/configs"
	"back/internal/models"
	"back/internal/moysklad"
	"back/internal/routes"
	// "back/internal/services"

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

	// Load configuration
	cfg, err := configs.LoadConfig()
	if err != nil {
		log.Fatal("Error loading configuration:", err)
	}

	// postgres setup
	err = setupDatabaseAndUser(cfg)
	if err != nil {
		log.Fatal("setup failed:", err)
	}

	// Database connection
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

	// Initialize handler with DB instance
	routes.SetupRoutes(e, db)

	// TEST
	headers, err := moysklad.CreateHeader("saratov")
	if err != nil {
		log.Fatal("failed to create headers: ", err)
	}

	token, err := moysklad.GetToken(headers)
	if err != nil {
		log.Fatal("failed to get token: ", err)
	}

	fmt.Println(token)

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

func connect(cfg *configs.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=%s TimeZone=%s",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
		cfg.DBSSLMode,
		cfg.DBTimeZone,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	return db, nil
}

func setupDatabaseAndUser(cfg *configs.Config) error {
	// Connect as superuser (e.g., postgres)
	dsn := fmt.Sprintf("user=%s password=%s host=%s port=%s sslmode=%s",
		cfg.DBRoot, cfg.DBRootPassword, cfg.DBHost, cfg.DBPort, cfg.DBSSLMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("admin connection failed: %w", err)
	}

	// 1. Create the database
	createDB := fmt.Sprintf("CREATE DATABASE %s", pqQuoteIdentifier(cfg.DBName))
	if err := db.Exec(createDB).Error; err != nil && !strings.Contains(err.Error(), "already exists") {
		return fmt.Errorf("error creating DB: %w", err)
	}

	// 2. Create the user (role)
	createUser := fmt.Sprintf("DO $$ BEGIN IF NOT EXISTS (SELECT FROM pg_catalog.pg_roles WHERE rolname = %s) THEN CREATE ROLE %s WITH LOGIN PASSWORD %s; END IF; END $$;",
		pqQuoteLiteral(cfg.DBUser), pqQuoteIdentifier(cfg.DBUser), pqQuoteLiteral(cfg.DBPassword))
	if err := db.Exec(createUser).Error; err != nil {
		return fmt.Errorf("error creating user: %w", err)
	}

	// 3. Grant privileges
	grant := fmt.Sprintf("GRANT ALL PRIVILEGES ON DATABASE %s TO %s", pqQuoteIdentifier(cfg.DBName), pqQuoteIdentifier(cfg.DBUser))
	if err := db.Exec(grant).Error; err != nil {
		return fmt.Errorf("error granting privileges: %w", err)
	}

	// Close admin connection
	if sqlDB, err := db.DB(); err == nil {
		if err := sqlDB.Close(); err != nil {
			return fmt.Errorf("failed to close database connection: %v", err)
		}
	} else {
		return fmt.Errorf("failed to get database instance: %v", err)
	}

	return nil
}

// Helper functions for quoting
func pqQuoteIdentifier(input string) string {
	return `"` + strings.ReplaceAll(input, `"`, `""`) + `"`
}

func pqQuoteLiteral(input string) string {
	return `'` + strings.ReplaceAll(input, `'`, `''`) + `'`
}
