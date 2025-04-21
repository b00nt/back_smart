package moysklad

import (
	"fmt"
	"log"

	"back/internal/models"

	"gorm.io/gorm"
)

func GetStocks(db *gorm.DB, token, city string) ([]interface{}, error) {
	var productStock []interface{}
	baseEndpoint := "https://api.moysklad.ru/api/remap/1.2/report/stock/all"

	offset := 0
	limit := 1000 // Maximum limit for Moysklad API
	groupBy := "variant"

	for {
		// Add pagination parameters to the endpoint
		paginatedEndpoint := fmt.Sprintf("%s?limit=%d&offset=%d&groupBy=%s", baseEndpoint, limit, offset, groupBy)

		// Get the current page of results
		results, totalCount, err := GetEssence(token, paginatedEndpoint)
		if err != nil {
			return nil, fmt.Errorf("failed to get stocks at offset %d: %w", offset, err)
		}

		// Add results to our collection
		productStock = append(productStock, results...)

		// Update offset for next iteration
		offset += len(results)

		// Check if we've retrieved all items
		if offset >= totalCount || len(results) == 0 {
			break
		}
	}

	return productStock, nil
}

func SaveStocks(db *gorm.DB, city string, stocksData []interface{}) error {
	// Begin a transaction
	tx := db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// Defer rollback in case of error
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for _, item := range stocksData {
		stockData := item.(map[string]interface{})

		// Extract the stock information
		code := stockData["code"].(string)
		stockValue := int(stockData["quantity"].(float64))

		// Update stock for products with matching code
		productResult := tx.Model(&models.Product{}).
			Where("code = ? AND city = ?", code, city).
			Update("stock", stockValue)

		if productResult.Error != nil {
			tx.Rollback()
			return fmt.Errorf("failed to update product stock: %w", productResult.Error)
		}

		// Update stock for modifications with matching code
		modResult := tx.Model(&models.Modification{}).
			Where("code = ?", code).
			Update("stock", stockValue)

		if modResult.Error != nil {
			tx.Rollback()
			return fmt.Errorf("failed to update modification stock: %w", modResult.Error)
		}

		// If no records were updated, log a warning but continue
		if productResult.RowsAffected == 0 && modResult.RowsAffected == 0 {
			log.Printf("Warning: No products or modifications found with code %s in city %s", code, city)
		}
	}

	// Commit the transaction
	return tx.Commit().Error
}

func UpdateAllStocks(db *gorm.DB, token, city string) error {
	stocks, err := GetStocks(db, token, city)
	if err != nil {
		return fmt.Errorf("failed to get stocks: %w", err)
	}

	err = SaveStocks(db, city, stocks)
	if err != nil {
		return fmt.Errorf("failed to save stocks: %w", err)
	}

	return nil
}
