package configs

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBRoot                  string
	DBRootPassword          string
	DBUser                  string
	DBPassword              string
	DBHost                  string
	DBPort                  string
	DBName                  string
	DBSSLMode               string
	DBTimeZone              string
	MoyskladUsernameMoscow  string
	MoyskladPasswordMoscow  string
	MoyskladUsernameSaratov string
	MoyskladPasswordSaratov string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	return &Config{
		DBRoot:                  os.Getenv("DB_ROOT"),
		DBRootPassword:          os.Getenv("DB_ROOT_PASSWORD"),
		DBUser:                  os.Getenv("DB_USER"),
		DBPassword:              os.Getenv("DB_PASSWORD"),
		DBHost:                  os.Getenv("DB_HOST"),
		DBPort:                  os.Getenv("DB_PORT"),
		DBName:                  os.Getenv("DB_NAME"),
		DBSSLMode:               os.Getenv("DB_SSLMODE"),
		DBTimeZone:              os.Getenv("DB_TIMEZONE"),
		MoyskladUsernameMoscow:  os.Getenv("MOYSKLAD_USERNAME_MOSCOW"),
		MoyskladPasswordMoscow:  os.Getenv("MOYSKLAD_PASSWORD_MOSCOW"),
		MoyskladUsernameSaratov: os.Getenv("MOYSKLAD_USERNAME_SARATOV"),
		MoyskladPasswordSaratov: os.Getenv("MOYSKLAD_PASSWORD_SARATOV"),
	}, nil
}
