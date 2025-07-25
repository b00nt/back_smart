// internal/moysklad/product.go
package moysklad

import (
	"fmt"
	"log"

	"back/internal/models"

	"gorm.io/gorm"
)

func GetProducts(token, city string) ([]any, error) {
	endpoint := "https://api.moysklad.ru/api/remap/1.2/entity/product"

	// Initial request to get the first page and total count
	const limit = 1000
	offset := 0

	firstPage, totalCount, err := GetEssence(token, endpoint, offset, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get products: %w", err)
	}

	// If total count is less than or equal to limit, return the results
	if totalCount <= limit {
		return firstPage, nil
	}

	// Otherwise, we need to make additional requests to get all products
	allProducts := make([]any, 0, totalCount)
	allProducts = append(allProducts, firstPage...)

	// Make additional requests until we get all products
	for offset += limit; offset < totalCount; offset += limit {
		pageProducts, _, err := GetEssence(token, endpoint, offset, limit)
		if err != nil {
			return nil, fmt.Errorf("failed to get products at offset %d: %w", offset, err)
		}
		allProducts = append(allProducts, pageProducts...)
	}

	return allProducts, nil
}

func SaveProducts(db *gorm.DB, city string, goods []any) error {
	if len(goods) == 0 {
		return nil // Early return if no products
	}

	// Begin a transaction for better performance with multiple records
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for i, item := range goods {
		productData, ok := item.(map[string]interface{})
		if !ok {
			tx.Rollback()
			return fmt.Errorf("item at index %d is not a map[string]interface{}", i)
		}

		// Safely extract data with error handling
		moyskladID, err := extractString(productData, "id")
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("extracting id: %w", err)
		}

		name, err := extractString(productData, "name")
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("extracting name: %w", err)
		}

		code, err := extractString(productData, "code")
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("extracting code: %w", err)
		}

		// Check if salePrices exists and is an array
		salePrices, ok := productData["salePrices"].([]interface{})
		if !ok || len(salePrices) == 0 {
			tx.Rollback()
			return fmt.Errorf("invalid salePrices for product %s", moyskladID)
		}

		// Extract price with proper error handling
		priceData, ok := salePrices[0].(map[string]interface{})
		if !ok {
			tx.Rollback()
			return fmt.Errorf("invalid price format for product %s", moyskladID)
		}

		priceValue, ok := priceData["value"].(float64)
		if !ok {
			tx.Rollback()
			return fmt.Errorf("invalid price value for product %s", moyskladID)
		}

		price := priceValue / 100

		categoryName, err := extractString(productData, "pathName")
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("extracting pathName: %w", err)
		}

		product := models.Product{
			MoyskladID: moyskladID,
			Name:       name,
			Code:       code,
			Price:      price,
			Category:   categoryName,
			City:       city,
		}

		// Use FirstOrCreate with better error handling
		result := tx.Where(models.Product{MoyskladID: moyskladID, City: city}).
			Assign(product).
			FirstOrCreate(&product)

		if result.Error != nil {
			tx.Rollback()
			return fmt.Errorf("database error with product %s: %w", moyskladID, result.Error)
		}

		if result.RowsAffected > 0 {
			log.Printf("Product created/updated: %s", moyskladID)
		}
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// Helper function to safely extract string values
func extractString(data map[string]interface{}, key string) (string, error) {
	value, exists := data[key]
	if !exists {
		return "", fmt.Errorf("key %s not found", key)
	}

	str, ok := value.(string)
	if !ok {
		return "", fmt.Errorf("key %s is not a string", key)
	}

	return str, nil
}
