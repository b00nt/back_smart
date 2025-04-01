// back/internal/routes/routes.go
package routes

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func InitializeRoutes(e *echo.Echo, db *gorm.DB) {
	// Pass the database instance to the context using middleware
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("tom_db", db)
			return next(c)
		}
	})
}
