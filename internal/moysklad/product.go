package moysklad

import (
	"back/internal/models"
	"fmt"
	"gorm.io/gorm"
	"log"
)

func GetProducts() []interface{} {
	headers := GetToken()
	endpoint := "https://api.moysklad.ru/api/remap/1.2/entity/product"
	result, err := GetEssence(headers, endpoint)
	if err != nil {
		fmt.Println(err)
		return nil
	} else {
		return result
	}
}

func SaveProducts(goods []interface{}, db *gorm.DB) error {
	for _, item := range goods {
		productData := item.(map[string]interface{}) // Type assertion to map
		moyskladID := productData["id"].(string)
		name := productData["name"].(string)
		code := productData["code"].(string)
		price := productData["salePrices"].([]interface{})[0].(map[string]interface{})["value"].(float64) / 100 // Dividing by 100
		categoryName := productData["pathName"].(string)

		// Check if the category exists or create it
		var category models.Category
		err := db.Where("name = ?", categoryName).FirstOrCreate(&category, models.Category{Name: categoryName}).Error
		if err != nil {
			log.Printf("Error finding/creating category: %v", err)
			continue // Skip to the next item if there's an error
		}

		product := models.Products{
			MoyskladID: moyskladID,
			Name:       name,
			Code:       code,
			Price:      price,
			CategoryID: category.ID,
		}

		// Upsert: Update or create
		result := db.Where(models.Products{MoyskladID: moyskladID}).Assign(product).FirstOrCreate(&product)

		if result.RowsAffected > 0 {
			fmt.Printf("DONE with product: %s\n", product.MoyskladID)
			fmt.Println("##################################################")
		} else {
			fmt.Printf("WRONG with product: %s\n", product.MoyskladID)
		}
	}
	return nil
}
