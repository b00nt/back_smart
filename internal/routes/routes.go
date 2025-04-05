// back/internal/routes/routes.go
package routes

import (
	"back/internal/handlers"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func SetupRoutes(e *echo.Echo, db *gorm.DB) {
	h := handlers.NewAuthHandler(db)

	e.POST("/feedback", h.Feedback)
	e.POST("/api/create-order", h.CreateOrder)

	// Define a group for productsroutes
	productsGroup := e.Group("/api/products")

	// Products routes
	productsGroup.GET("/:categoryID", h.GetProductsByCategory)
	productsGroup.GET("", h.GetProducts)
	productsGroup.GET("/product/:moysklad_id", h.GetProductByID)

}
