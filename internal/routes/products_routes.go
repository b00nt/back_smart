// routes/products_routes.go
package routes

import (
	"back/internal/handlers"
	"github.com/labstack/echo/v4"
)

func InitProductsRoutes(e *echo.Echo, h *handlers.Handler) {
	// Define a group for productsroutes
	productsGroup := e.Group("/api/products")

	// Products routes
	productsGroup.GET("/:categoryID", h.GetProductsByCategory)
	productsGroup.GET("", h.GetProducts)
	productsGroup.GET("/product/:moysklad_id", h.GetProductByID)
}
