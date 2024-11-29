package moysklad

import (
	"back/internal/models"
	"fmt"
	"gorm.io/gorm"
)

func GetStock() []interface{} {
	fmt.Println("GetStock running.")
	headers := GetToken()
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

func SaveStock(stock []interface{}, db *gorm.DB) error {
	fmt.Println("SaveStock running.")
	for _, item := range stock {
		stockData := item.(map[string]interface{}) // Type assertion to map
		code := stockData["code"].(string)
		stockValue := int(stockData["stock"].(float64))

		// Check if the product exists based on the code
		var product models.Products
		productExists := db.Where("code = ?", code).First(&product).Error == nil

		// Check if the modification exists based on the code
		var modification models.Modification
		modificationExists := db.Where("code = ?", code).First(&modification).Error == nil

		if productExists {
			// Product exists, update the stock
			product.Stock = stockValue
			db.Save(&product) // Save the updated product
			fmt.Printf("%s has %d\n", code, stockValue)
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

		fmt.Println("##################################################")
	}
	return nil
}
