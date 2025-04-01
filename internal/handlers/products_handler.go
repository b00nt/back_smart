// back/internal/handlers/products_handler.go
package handlers

import (
	"back/internal/models"
	"fmt"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"net/http"
	"net/url"
)

type Handler struct {
	DB *gorm.DB
}

// NewHandler initializes a new Handler with DB
func NewHandler(db *gorm.DB) *Handler {
	return &Handler{DB: db}
}

func (h *Handler) GetProductsByCategory(c echo.Context) error {
	var products []models.Products
	var productsSaratov []models.ProductsSaratov
	var query *gorm.DB
	city := c.QueryParam("city")
	// Get the category ID from the request URL
	categoryID := c.Param("categoryID")
	fmt.Println("City:", city)
	fmt.Println("Category ID:", categoryID)

	if city == "saratov" {
		query = h.DB.Table("products_saratovs")
		if err := query.Preload("Category").
			Preload("ProductImages").
			Preload("Modification").
			Preload("Modification.ModificationCharacteristics").
			Preload("Modification.ModificationImages").
			Where("price > ?", 0).
			Where("category_id = ?", categoryID).
			Find(&productsSaratov).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Failed to fetch products",
			})
		} else {
			// Send back the products as JSON
			return c.JSON(http.StatusOK, productsSaratov)
		}
	} else {
		query = h.DB.Table("products")
		if err := query.Preload("Category").
			Preload("ProductImages").
			Preload("Modification").
			Preload("Modification.ModificationCharacteristics").
			Preload("Modification.ModificationImages").
			Where("price > ?", 0).
			Where("category_id = ?", categoryID).
			Find(&products).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Failed to fetch products",
			})
		} else {
			// Send back the products as JSON
			return c.JSON(http.StatusOK, products)
		}
	}
}

func (h *Handler) GetProducts(c echo.Context) error {
	var products []models.Products
	city := c.QueryParam("city")
	searchTerm := c.QueryParam("query")

	fmt.Println("City:", city)
	fmt.Println("Search Term:", searchTerm)

	var query *gorm.DB

	if city == "saratov" {
		query = h.DB.Table("products_saratovs").
			Preload("Category").
			Preload("ProductImages").
			Preload("Modification").
			Preload("Modification.ModificationCharacteristics").
			Preload("Modification.ModificationImages")
	} else {
		query = h.DB.Table("products").
			Preload("Category").
			Preload("ProductImages").
			Preload("Modification").
			Preload("Modification.ModificationCharacteristics").
			Preload("Modification.ModificationImages")
	}

	if searchTerm != "" {
		decodedSearchTerm, err := url.QueryUnescape(searchTerm)
		if err != nil {
			fmt.Println("Error decoding search term:", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Failed to decode search term",
			})
		}
		fmt.Println("Decoded Search Term:", decodedSearchTerm)

		// Use MySQL case-insensitive collation
		query = query.Where("name LIKE ? COLLATE utf8mb4_unicode_ci", "%"+decodedSearchTerm+"%")

		// Log the final query and number of matched products
		fmt.Println("Final Query:", query)
	}

	//
	query = query.Where("price > ?", 0)

	// Ensure no default limit is applied
	query = query.Limit(5000)

	// Execute the query and fetch the products
	result := query.Find(&products)
	if result.Error != nil {
		fmt.Println("Error fetching products:", result.Error)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to fetch products",
		})
	}

	// Print the final number of products fetched
	fmt.Println("Filtered products count:", len(products))
	return c.JSON(http.StatusOK, products)
}

func (h *Handler) GetProductByID(c echo.Context) error {
	var products []models.Products
	var productsSaratov []models.ProductsSaratov
	var query *gorm.DB

	productID := c.Param("moysklad_id")
	city := c.QueryParam("city")

	if productID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Missing product ID",
		})
	}

	//	fmt.Println("City:", city, "\nProduct ID:", productID)

	if city == "saratov" {
		query = h.DB.Table("products_saratovs")
		if err := query.Preload("Category").
			Preload("ProductImages").
			Preload("Modification").
			Preload("Modification.ModificationCharacteristics").
			Preload("Modification.ModificationImages").
			Where("price > ?", 0).
			Where("moysklad_id = ?", productID).
			Find(&productsSaratov).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Failed to fetch products",
			})
		}
		if len(productsSaratov) == 0 {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "Product not found in Saratov",
			})
		}
		return c.JSON(http.StatusOK, productsSaratov)
	}

	query = h.DB.Table("products")
	if err := query.Preload("Category").
		Preload("ProductImages").
		Preload("Modification").
		Preload("Modification.ModificationCharacteristics").
		Preload("Modification.ModificationImages").
		Where("price > ?", 0).
		Where("moysklad_id = ?", productID).
		Find(&products).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to fetch products",
		})
	}
	if len(products) == 0 {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Product not found",
		})
	}
	return c.JSON(http.StatusOK, products)
}

// func (h *Handler) GetProductByID(c echo.Context) error {
// 	var products []models.Products
// 	var productsSaratov []models.ProductsSaratov
// 	var query *gorm.DB
// 	// Extract the product ID from the request parameters
// 	productID := c.Param("moysklad_id")
// 	city := c.QueryParam("city") // Optional: to determine which table to query
//
// 	// fmt.Println("City:", city)
//
// 	if city == "saratov" {
// 		query = h.DB.Table("products_saratovs")
// 		if err := query.Preload("Category").
// 			Preload("ProductImages").
// 			Preload("Modification").
// 			Preload("Modification.ModificationCharacteristics").
// 			Preload("Modification.ModificationImages").
// 			Where("price > ?", 0).
// 			Where("moysklad_id = ?", productID).
// 			Find(&productsSaratov).Error; err != nil {
// 			return c.JSON(http.StatusInternalServerError, map[string]string{
// 				"error": "Failed to fetch products",
// 			})
// 		} else {
// 			// Send back the products as JSON
// 			return c.JSON(http.StatusOK, productsSaratov)
// 		}
// 	} else {
// 		query = h.DB.Table("products")
// 		if err := query.Preload("Category").
// 			Preload("ProductImages").
// 			Preload("Modification").
// 			Preload("Modification.ModificationCharacteristics").
// 			Preload("Modification.ModificationImages").
// 			Where("price > ?", 0).
// 			Where("moysklad_id = ?", productID).
// 			Find(&products).Error; err != nil {
// 			return c.JSON(http.StatusInternalServerError, map[string]string{
// 				"error": "Failed to fetch products",
// 			})
// 		} else {
// 			// Send back the products as JSON
// 			return c.JSON(http.StatusOK, products)
// 		}
// 	}
// }
