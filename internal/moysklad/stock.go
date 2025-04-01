package moysklad

import (
	"back/internal/models"
	"fmt"
	"gorm.io/gorm"
	"time"
)

func GetVariantStock(city string, db *gorm.DB) []interface{} {
	fmt.Println("GetStock running.")
	headers := GetToken(city)
	moyskladId := GetModID(city, db)

	// Create a ticker that ticks every 66.67 milliseconds (3000ms / 45 requests)
	ticker := time.NewTicker(66 * time.Millisecond)
	defer ticker.Stop()

	for _, v := range moyskladId {
		// Wait for the next tick before making a request
		<-ticker.C

		//		endpoint := fmt.Sprintf("https://api.moysklad.ru/api/remap/1.2/report/stock?filter=productid=%s", v)
		endpoint := fmt.Sprintf("https://api.moysklad.ru/api/remap/1.2/report/stock/all?filter=variant=https://api.moysklad.ru/api/remap/1.2/entity/variant/%s", v)
		// fmt.Println(endpoint)
		result, err := GetEssence(headers, endpoint)
		if result == nil || len(result) == 0 {
			continue
		} else {
			if err != nil {
				fmt.Println(err)
				return nil
			}
			if err != SaveVariantStock(city, result, db) {
				fmt.Println(err)
			}
		}
	}
	return nil
}

func SaveVariantStock(city string, stock []interface{}, db *gorm.DB) error {
	fmt.Println("SaveVariantStock running.")
	for _, item := range stock {
		stockData := item.(map[string]interface{}) // Type assertion to map
		code := stockData["code"].(string)
		stockValue := int(stockData["quantity"].(float64)) // Атрибут "Доступно"
		//		stockValue := int(stockData["stock"].(float64)) // Атрибут "Остаток"

		var modification models.Modification
		var modificationSaratov models.ModificationSaratov
		var modificationExists bool

		if city == "moscow" {
			modificationExists = db.Where("code = ?", code).First(&modification).Error == nil

			if modificationExists {
				modification.Stock = stockValue
				db.Save(&modification) // Save the updated product
				fmt.Printf("Modification with %s code has %d quantity\n", code, stockValue)
			} else {
				fmt.Println("Modification does not exist with code:", code)
			}
			fmt.Println("#######################################################################################################################")
		} else if city == "saratov" {
			modificationExists = db.Where("code = ?", code).First(&modificationSaratov).Error == nil

			if modificationExists {
				modificationSaratov.Stock = stockValue
				db.Save(&modificationSaratov) // Save the updated product
				fmt.Printf("Modification with %s code has %d quantity\n", code, stockValue)
			} else {
				fmt.Println("Modification does not exist with code:", code)
			}
			fmt.Println("#######################################################################################################################")
		}
	}
	return nil
}

func GetStock(city string) []interface{} {
	fmt.Println("GetStock running.")
	headers := GetToken(city)
	endpoint := "https://api.moysklad.ru/api/remap/1.2/report/stock/all"
	result, err := GetEssence(headers, endpoint)
	if err != nil {
		fmt.Println(err)
		return nil
	} else {
		fmt.Println("GetStock done.")
		return result
	}
}

func SaveStock(city string, stock []interface{}, db *gorm.DB) error {
	fmt.Println("SaveStock running.")
	for _, item := range stock {
		stockData := item.(map[string]interface{}) // Type assertion to map
		code := stockData["code"].(string)
		stockValue := int(stockData["quantity"].(float64)) // Атрибут "Доступно"
		//		stockValue := int(stockData["stock"].(float64)) // Атрибут "Остаток"

		var product models.Products
		var productSaratov models.ProductsSaratov
		var modification models.Modification
		var modificationSaratov models.ModificationSaratov
		var productExists, modificationExists bool

		if city == "moscow" {
			productExists = db.Where("code = ?", code).First(&product).Error == nil
			modificationExists = db.Where("code = ?", code).First(&modification).Error == nil
			if productExists {
				// Product exists, update the stock
				product.Stock = stockValue
				db.Save(&product) // Save the updated product
				fmt.Printf("Product with %s code has %d quantity\n", code, stockValue)
			} else {
				fmt.Println("Product does not exist in Products table with code:", code)
			}

			if modificationExists {
				modification.Stock = stockValue
				db.Save(&modification) // Save the updated product
				fmt.Printf("%s has %d\n", code, stockValue)
			} else {
				fmt.Println("Modification does not exist with code:", code)
			}
			fmt.Println("#######################################################################################################################")
		} else if city == "saratov" {
			productExists = db.Where("code = ?", code).First(&productSaratov).Error == nil
			modificationExists = db.Where("code = ?", code).First(&modificationSaratov).Error == nil
			if productExists {
				// Product exists, update the stock
				productSaratov.Stock = stockValue
				db.Save(&productSaratov) // Save the updated product
				fmt.Printf("Product with %s code has %d quantity\n", code, stockValue)
			} else {
				fmt.Println("Product does not exist in ProductsSaratov table with code:", code)
			}

			if modificationExists {
				modificationSaratov.Stock = stockValue
				db.Save(&modificationSaratov) // Save the updated product
				fmt.Printf("%s has %d\n", code, stockValue)
			} else {
				fmt.Println("Modification does not exist with code:", code)
			}
			fmt.Println("#######################################################################################################################")
		}
	}
	return nil
}
