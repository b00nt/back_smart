// internal/moysklad/product.go
package moysklad

import (
	"back/internal/models"
	"fmt"
	"gorm.io/gorm"
	"log"
)

func GetProducts(city string) []interface{} {
	headers := GetToken(city)
	endpoint := "https://api.moysklad.ru/api/remap/1.2/entity/product"
	result, err := GetEssence(headers, endpoint)
	if err != nil {
		fmt.Println(err)
		return nil
	} else {
		return result
	}
}

func SaveProducts(city string, goods []interface{}, db *gorm.DB) error {
	for _, item := range goods {
		productData := item.(map[string]interface{})
		moyskladID := productData["id"].(string) // moyskladID, ID товара
		name := productData["name"].(string)     // имя товара
		code := productData["code"].(string)
		price := productData["salePrices"].([]interface{})[0].(map[string]interface{})["value"].(float64) / 100 // Делим на 100
		categoryName := productData["pathName"].(string)

		// fmt.Printf("DEBUG\nName: %s\tMoyskladID: %s", name, moyskladID)

		var category models.Category
		err := db.Where("name = ?", categoryName).FirstOrCreate(&category, models.Category{Name: categoryName}).Error
		if err != nil {
			log.Printf("Error finding/creating category: %v", err)
			continue // Skip to the next item if there's an error
		}

		if city == "moscow" {
			product := models.Products{
				MoyskladID: moyskladID,
				Name:       name,
				Code:       code,
				Price:      price,
				CategoryID: category.ID,
			}

			result := db.Where(models.Products{MoyskladID: moyskladID}).Assign(product).FirstOrCreate(&product)

			if result.RowsAffected > 0 {
				fmt.Printf("DONE with product: %s\n", product.MoyskladID)
				fmt.Println("##################################################")
			} else {
				fmt.Printf("WRONG with product: %s\n", product.MoyskladID)
			}
		} else if city == "saratov" {
			product := models.ProductsSaratov{
				MoyskladID: moyskladID,
				Name:       name,
				Code:       code,
				Price:      price,
				CategoryID: category.ID,
			}

			// Upsert: Update or create
			result := db.Where(models.ProductsSaratov{MoyskladID: moyskladID}).Assign(product).FirstOrCreate(&product)

			if result.RowsAffected > 0 {
				fmt.Printf("DONE with product: %s\n", product.MoyskladID)
				fmt.Println("##################################################")
			} else {
				fmt.Printf("WRONG with product: %s\n", product.MoyskladID)
			}
		}

	}
	return nil
}
