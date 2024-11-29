// https://dev.moysklad.ru/doc/api/remap/1.2/dictionaries/#suschnosti-modifikaciq

package moysklad

import (
	"back/internal/models"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func GetModifications() []interface{} {
	fmt.Println("GetModifications running.")
	headers := GetToken()
	// endpoint := fmt.Sprintf("https://api.moysklad.ru/api/remap/1.2/entity/variant?filter=modsid={%s}", moyskladId)
	endpoint := "https://api.moysklad.ru/api/remap/1.2/entity/variant"
	result, err := GetEssence(headers, endpoint)
	if err != nil {
		fmt.Println(err)
		return nil
	} else {
		fmt.Println("GetModifications done.")
		return result
	}
}

// SaveModifications saves modifications and associates them with the corresponding Products in the database
func SaveModifications(mods []interface{}, db *gorm.DB) error {
	fmt.Println("SaveModifications done.")
	for _, mod := range mods {
		modMap, ok := mod.(map[string]interface{})
		if !ok {
			return fmt.Errorf("modification data format invalid")
		}

		// Extract ModID, Code, SalePrices, Image URL
		modID, ok := modMap["id"].(string)
		if !ok {
			return fmt.Errorf("modification ID not found or invalid")
		}
		code, _ := modMap["code"].(string)
		//imageURL, _ := modMap["images"].(string) // adjust parsing if necessary
		var salePrice float64
		if prices, ok := modMap["salePrices"].([]interface{}); ok && len(prices) > 0 {
			if firstPrice, ok := prices[0].(map[string]interface{}); ok {
				salePrice, _ = firstPrice["value"].(float64)
			}
		}

		// Extract MoyskladID from product's href URL
		var moyskladID string
		if productInfo, ok := modMap["product"].(map[string]interface{}); ok {
			if meta, ok := productInfo["meta"].(map[string]interface{}); ok {
				if href, ok := meta["href"].(string); ok {
					moyskladID = extractMoyskladIDFromURL(href)
				}
			}
		}

		// Prepare the modification record to insert or update
		modification := models.Modification{
			ModID:      modID,
			MoyskladID: moyskladID,
			Code:       code,
			//Images:      imageURL,
			SalePrices: salePrice / 100,
		}

		// Use Upsert to insert or update the modification
		if err := db.Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "mod_id"}}, // Specify the unique column
			DoUpdates: clause.Assignments(map[string]interface{}{
				"updated_at":  gorm.Expr("CURRENT_TIMESTAMP"), // Update the timestamp
				"moysklad_id": modification.MoyskladID,
				"code":        modification.Code,
				"sale_prices": modification.SalePrices,
			}),
		}).Create(&modification).Error; err != nil {
			return fmt.Errorf("failed to save modification: %v", err)
		} else {
			fmt.Printf("Modification %s is done\n", modID)
		}

		// Collect characteristics for batch insertion
		var modChars []models.ModificationCharacteristics
		var nameChars string
		if characteristics, ok := modMap["characteristics"].([]interface{}); ok {
			for _, char := range characteristics {
				charMap, ok := char.(map[string]interface{})
				if !ok {
					return fmt.Errorf("characteristic format invalid")
				}

				nameChars, _ := charMap["name"].(string)
				value, _ := charMap["value"].(string)

				// Append each characteristic to the slice
				modChars = append(modChars, models.ModificationCharacteristics{
					ModID: modification.ModID,
					Name:  nameChars,
					Value: value,
				})
			}
		}

		// Batch insert characteristics
		if len(modChars) > 0 {
			if err := db.Clauses(clause.OnConflict{
				DoNothing: true, // Avoid duplicating characteristics if they already exist
			}).Create(&modChars).Error; err != nil {
				return fmt.Errorf("failed to batch insert characteristics: %v", err)
			} else {
				fmt.Printf("characteristic %s of modification %s is done\n", nameChars, modID)
			}
		}
	}

	return nil
}
