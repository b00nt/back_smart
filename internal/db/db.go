package db

import (
	"fmt"
	"log"
	"strings"

	"back/internal/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func connect(cfg *config.Config) (*gorm.DB, error) {
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

func setupDatabaseAndUser(cfg *config.Config) error {
	// connect as superuser (e.g., postgres)
	dsn := fmt.Sprintf("user=%s password=%s host=%s port=%s sslmode=%s",
		cfg.DBRoot, cfg.DBRootPassword, cfg.DBHost, cfg.DBPort, cfg.DBSSLMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("admin connection failed: %w", err)
	}

	// create the database
	createDB := fmt.Sprintf("CREATE DATABASE %s", pqQuoteIdentifier(cfg.DBName))
	if err := db.Exec(createDB).Error; err != nil && !strings.Contains(err.Error(), "already exists") {
		return fmt.Errorf("error creating DB: %w", err)
	}

	// create the user (role)
	createUser := fmt.Sprintf("DO $$ BEGIN IF NOT EXISTS (SELECT FROM pg_catalog.pg_roles WHERE rolname = %s) THEN CREATE ROLE %s WITH LOGIN PASSWORD %s; END IF; END $$;",
		pqQuoteLiteral(cfg.DBUser), pqQuoteIdentifier(cfg.DBUser), pqQuoteLiteral(cfg.DBPassword))
	if err := db.Exec(createUser).Error; err != nil {
		return fmt.Errorf("error creating user: %w", err)
	}

	// grant privileges
	grant := fmt.Sprintf("GRANT ALL PRIVILEGES ON DATABASE %s TO %s", pqQuoteIdentifier(cfg.DBName), pqQuoteIdentifier(cfg.DBUser))
	if err := db.Exec(grant).Error; err != nil {
		return fmt.Errorf("error granting privileges: %w", err)
	}

	// close admin connection
	if sqlDB, err := db.DB(); err == nil {
		if err := sqlDB.Close(); err != nil {
			return fmt.Errorf("failed to close database connection: %v", err)
		}
	} else {
		return fmt.Errorf("failed to get database instance: %v", err)
	}

	return nil
}

// helper functions for quoting
func pqQuoteIdentifier(input string) string {
	return `"` + strings.ReplaceAll(input, `"`, `""`) + `"`
}

func pqQuoteLiteral(input string) string {
	return `'` + strings.ReplaceAll(input, `'`, `''`) + `'`
}
