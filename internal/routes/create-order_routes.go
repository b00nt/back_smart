// routes/products_routes.go
package routes

import (
	"back/internal/handlers"
	"github.com/labstack/echo/v4"
)

func InitOrderRoutes(e *echo.Echo, h *handlers.Handler) {
	// Register routes with and without trailing slash
	e.POST("/api/create-order", h.CreateOrder)
	e.POST("/api/create-order/", h.CreateOrder)
}
