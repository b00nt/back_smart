package moysklad

import (
	"fmt"
	"log"

	"back/internal/models"

	"gorm.io/gorm"
)

func GetStocks(city string, db *gorm.DB) ([]interface{}, error) {
	var productStock []interface{}
	baseEndpoint := "https://api.moysklad.ru/api/remap/1.2/report/stock/all"

	headers, err := GetToken(city)
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}

	offset := 0
	limit := 1000 // Maximum limit for Moysklad API
	groupBy := "variant"
	
	for {
		// Add pagination parameters to the endpoint
		paginatedEndpoint := fmt.Sprintf("%s?limit=%d&offset=%d&groupBy=%s", baseEndpoint, limit, offset, groupBy)
		
		// Get the current page of results
		results, totalCount, err := GetEssence(headers, paginatedEndpoint)
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

func SaveStocks(city string, stocksData []interface{}, db *gorm.DB) error {
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

func UpdateAllStocks(city string, db *gorm.DB) error {
	stocks, err := GetStocks(city, db)
	if err != nil {
		return fmt.Errorf("failed to get stocks: %w", err)
	}
	
	err = SaveStocks(city, stocks, db)
	if err != nil {
		return fmt.Errorf("failed to save stocks: %w", err)
	}
	
	return nil
}

//
// func SaveVariantStock(city string, stock []interface{}, db *gorm.DB) error {
// 	fmt.Println("SaveVariantStock running.")
// 	for _, item := range stock {
// 		stockData := item.(map[string]interface{}) // Type assertion to map
// 		code := stockData["code"].(string)
// 		stockValue := int(stockData["quantity"].(float64)) // Атрибут "Доступно"
// 		//		stockValue := int(stockData["stock"].(float64)) // Атрибут "Остаток"
//
// 		var modification models.Modification
// 		var modificationSaratov models.ModificationSaratov
// 		var modificationExists bool
//
// 		if city == "moscow" {
// 			modificationExists = db.Where("code = ?", code).First(&modification).Error == nil
//
// 			if modificationExists {
// 				modification.Stock = stockValue
// 				db.Save(&modification) // Save the updated product
// 				fmt.Printf("Modification with %s code has %d quantity\n", code, stockValue)
// 			} else {
// 				fmt.Println("Modification does not exist with code:", code)
// 			}
// 			fmt.Println("#######################################################################################################################")
// 		} else if city == "saratov" {
// 			modificationExists = db.Where("code = ?", code).First(&modificationSaratov).Error == nil
//
// 			if modificationExists {
// 				modificationSaratov.Stock = stockValue
// 				db.Save(&modificationSaratov) // Save the updated product
// 				fmt.Printf("Modification with %s code has %d quantity\n", code, stockValue)
// 			} else {
// 				fmt.Println("Modification does not exist with code:", code)
// 			}
// 			fmt.Println("#######################################################################################################################")
// 		}
// 	}
// 	return nil
// }
//
// func SaveStock(city string, stock []interface{}, db *gorm.DB) error {
// 	fmt.Println("SaveStock running.")
// 	for _, item := range stock {
// 		stockData := item.(map[string]interface{}) // Type assertion to map
// 		code := stockData["code"].(string)
// 		stockValue := int(stockData["quantity"].(float64)) // Атрибут "Доступно"
// 		//		stockValue := int(stockData["stock"].(float64)) // Атрибут "Остаток"
//
// 		var product models.Products
// 		var productSaratov models.ProductsSaratov
// 		var modification models.Modification
// 		var modificationSaratov models.ModificationSaratov
// 		var productExists, modificationExists bool
//
// 		if city == "moscow" {
// 			productExists = db.Where("code = ?", code).First(&product).Error == nil
// 			modificationExists = db.Where("code = ?", code).First(&modification).Error == nil
// 			if productExists {
// 				// Product exists, update the stock
// 				product.Stock = stockValue
// 				db.Save(&product) // Save the updated product
// 				fmt.Printf("Product with %s code has %d quantity\n", code, stockValue)
// 			} else {
// 				fmt.Println("Product does not exist in Products table with code:", code)
// 			}
//
// 			if modificationExists {
// 				modification.Stock = stockValue
// 				db.Save(&modification) // Save the updated product
// 				fmt.Printf("%s has %d\n", code, stockValue)
// 			} else {
// 				fmt.Println("Modification does not exist with code:", code)
// 			}
// 			fmt.Println("#######################################################################################################################")
// 		} else if city == "saratov" {
// 			productExists = db.Where("code = ?", code).First(&productSaratov).Error == nil
// 			modificationExists = db.Where("code = ?", code).First(&modificationSaratov).Error == nil
// 			if productExists {
// 				// Product exists, update the stock
// 				productSaratov.Stock = stockValue
// 				db.Save(&productSaratov) // Save the updated product
// 				fmt.Printf("Product with %s code has %d quantity\n", code, stockValue)
// 			} else {
// 				fmt.Println("Product does not exist in ProductsSaratov table with code:", code)
// 			}
//
// 			if modificationExists {
// 				modificationSaratov.Stock = stockValue
// 				db.Save(&modificationSaratov) // Save the updated product
// 				fmt.Printf("%s has %d\n", code, stockValue)
// 			} else {
// 				fmt.Println("Modification does not exist with code:", code)
// 			}
// 			fmt.Println("#######################################################################################################################")
// 		}
// 	}
// 	return nil
// }
