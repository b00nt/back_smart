// https://dev.moysklad.ru/doc/api/remap/1.2/dictionaries/#suschnosti-modifikaciq

package moysklad

import (
	"fmt"
	"strings"

	"back/internal/models"

	"gorm.io/gorm"
	// "gorm.io/gorm/clause"
)

func GetModifications(db *gorm.DB, token, city string) ([]interface{}, error) {
	var allModifications []interface{}
	baseEndpoint := "https://api.moysklad.ru/api/remap/1.2/entity/variant"
	offset := 0
	limit := 1000 // Maximum limit for Moysklad API

	for {
		// Add pagination parameters to the endpoint
		paginatedEndpoint := fmt.Sprintf("%s?limit=%d&offset=%d", baseEndpoint, limit, offset)

		// Get the current page of results
		results, totalCount, err := GetEssence(token, paginatedEndpoint)
		if err != nil {
			return nil, fmt.Errorf("failed to get variants at offset %d: %w", offset, err)
		}

		// Add results to our collection
		allModifications = append(allModifications, results...)

		// Update offset for next iteration
		offset += len(results)

		// Check if we've retrieved all items
		if offset >= totalCount || len(results) == 0 {
			break
		}
	}

	return allModifications, nil
}

func SaveModifications(db *gorm.DB, city string, modifications []interface{}) error {
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

	for _, modData := range modifications {
		mod := modData.(map[string]interface{})

		// Extract the basic modification data
		modID := mod["id"].(string)
		name := mod["name"].(string)

		// Extract code (handle case where it might be missing)
		var code string
		if codeVal, ok := mod["code"]; ok && codeVal != nil {
			code = codeVal.(string)
		}

		// Get the product ID
		var productID string
		if product, ok := mod["product"].(map[string]interface{}); ok {
			if productMeta, ok := product["meta"].(map[string]interface{}); ok {
				if href, ok := productMeta["href"].(string); ok {
					// Extract product ID from the URL
					parts := strings.Split(href, "/")
					if len(parts) > 0 {
						productID = parts[len(parts)-1]
					}
				}
			}
		}

		// Extract the sale price
		var salePrice float64
		if salePrices, ok := mod["salePrices"].([]interface{}); ok && len(salePrices) > 0 {
			if priceData, ok := salePrices[0].(map[string]interface{}); ok {
				if val, ok := priceData["value"].(float64); ok {
					salePrice = val
				}
			}
		}

		// Create or update the modification record
		modification := models.Modification{
			Name:       name,
			ModID:      modID,
			MoyskladID: productID,
			Code:       code,
			Price:      salePrice,
		}

		result := tx.Where("mod_id = ?", modID).FirstOrCreate(&modification)
		if result.Error != nil {
			tx.Rollback()
			return fmt.Errorf("failed to save modification: %w", result.Error)
		}

		// Save characteristics
		if chars, ok := mod["characteristics"].([]interface{}); ok {
			// First, delete existing characteristics for this modification
			if err := tx.Where("mod_id = ?", modID).Delete(&models.Characteristic{}).Error; err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to delete existing characteristics: %w", err)
			}

			for _, charData := range chars {
				char := charData.(map[string]interface{})

				characteristic := models.Characteristic{
					ModID: modID,
					Name:  char["name"].(string),
					Value: char["value"].(string),
				}

				if err := tx.Create(&characteristic).Error; err != nil {
					tx.Rollback()
					return fmt.Errorf("failed to save characteristic: %w", err)
				}
			}
		}

		// Handle images if present
		// if imagesData, ok := mod["images"].(map[string]interface{}); ok {
		// 	if metaData, ok := imagesData["meta"].(map[string]interface{}); ok {
		// 		if href, ok := metaData["href"].(string); ok && metaData["size"].(float64) > 0 {
		// We need to make an additional API call to get the images
		// This is a placeholder for that call
		// You would implement this based on your GetEssence function
		// images, _, err := GetEssence(headers, href)
		// if err == nil {
		//     // Save images logic here
		// }
		//}
		//}
		//}
	}

	// Commit the transaction
	return tx.Commit().Error
}

func UpdateAllModifications(db *gorm.DB, token, city string) error {
	modifications, err := GetModifications(db, token, city)
	if err != nil {
		return fmt.Errorf("failed to get modifications: %w", err)
	}

	err = SaveModifications(db, city, modifications)
	if err != nil {
		return fmt.Errorf("failed to save modifications: %w", err)
	}

	return nil
}
