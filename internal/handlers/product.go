// back/internal/handlers/products_handler.go
package handlers

import (
	"back/internal/models"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"net/http"
)

type Handler struct {
	DB *gorm.DB
}

func NewHandler(db *gorm.DB) *Handler {
	return &Handler{DB: db}
}

func (h *Handler) GetProductsByCategory(c echo.Context) error {
	category := c.Param("CATEGORY") // This extracts the CATEGORY parameter from the URL path
	city := c.QueryParam("city")

	var products []models.Product // Note: changed from Products to Product to match your model

	query := h.DB.Where("category = ? AND display = ?", category, true)

	// Add city filter if provided
	if city != "" {
		query = query.Where("city = ?", city)
	}

	// Execute the query with preloaded relationships
	result := query.Preload("Modification", "display = ?", true).
		Preload("Modification.ModificationCharacteristic").
		Preload("Modification.ModificationImage").
		Preload("ProductImages").
		Find(&products)

	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to fetch products",
		})
	}

	return c.JSON(http.StatusOK, products)
}
