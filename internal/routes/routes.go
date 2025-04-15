// back/internal/routes/routes.go
package routes

import (
	"back/internal/handlers"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func SetupRoutes(e *echo.Echo, db *gorm.DB) {
	h := handlers.NewHandler(db)

	routesGroup := e.Group("/api")

	routesGroup.POST("/feedback", h.Feedback)
	routesGroup.POST("/create-order", h.CreateOrder)
	routesGroup.GET("/products/:categoryID", h.GetProductsByCategory)

}
