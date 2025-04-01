package moysklad

import (
	"back/internal/models"
	"fmt"
	"gorm.io/gorm"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)


func GetSaveDownloadProductImages(city string, db *gorm.DB) {
	fmt.Println("GetSaveDownloadProductImages running.")
	identifier := 1
	headers := GetToken(city)
	moyskladId := GetMoyskladID(city, db)

	// Create a ticker that ticks every 66.67 milliseconds (3000ms / 45 requests)
	ticker := time.NewTicker(66 * time.Millisecond)
	defer ticker.Stop()

	for _, id := range moyskladId {
		// Wait for the next tick before making a request
		<-ticker.C

		endpoint := fmt.Sprintf("https://api.moysklad.ru/api/remap/1.2/entity/product/%s/images", id)
		// endpoint := fmt.Sprintf("https://api.moysklad.ru/api/remap/1.2/entity/variant?filter=modsid={%s}", moyskladId)
		productsImages, err := GetImages(headers, endpoint)
		if err != nil {
			fmt.Println(err)
			continue // Skip to the next iteration if there's an error
		}
		if productsImages == nil || len(productsImages) == 0 {
			fmt.Println("Images not found. Skipping.")
			continue // Skip to the next iteration if productsImages is empty
		}
		maps := SaveImages(city, productsImages, id, db, identifier)
		downloadImages(city, headers, maps, db, identifier)
	}
}

func GetSaveDownloadModImages(city string, db *gorm.DB) {
	fmt.Println("GetSaveDownloadModImages running.")
	identifier := 2
	headers := GetToken(city)
	modId := GetModID(city, db)

	// Create a ticker that ticks every 66.67 milliseconds (3000ms / 45 requests)
	ticker := time.NewTicker(66 * time.Millisecond)
	defer ticker.Stop()

	for _, id := range modId {
		// Wait for the next tick before making a request
		<-ticker.C

		endpoint := fmt.Sprintf("https://api.moysklad.ru/api/remap/1.2/entity/variant/%s/images", id)
		variantsImages, err := GetImages(headers, endpoint)
		if err != nil {
			fmt.Println(err)
			continue // Skip to the next iteration if there's an error
		}
		if variantsImages == nil || len(variantsImages) == 0 {
			fmt.Println("Images not found. Skipping.")
			continue // Skip to the next iteration if productsImages is empty
		}
		maps := SaveImages(city, variantsImages, id, db, identifier)
		downloadImages(city, headers, maps, db, identifier)
	}
}

func GetImages(headers http.Header, endpoint string) ([]interface{}, error) {
	result, err := GetEssence(headers, endpoint)
	if err != nil {
		fmt.Println(err)
		return nil, err
	} else {
		return result, err
	}
}

func SaveImages(city string, images []interface{}, ID string, db *gorm.DB, identifier int) map[string]string {
	var imgMap = make(map[string]string, len(images))
	for _, item := range images {
		data := item.(map[string]interface{})

		// Access the meta field
		meta := data["meta"].(map[string]interface{})
		imgId := extractImageURL(meta["href"].(string))
		downloadHref := meta["downloadHref"].(string)
		imgMap[imgId] = downloadHref

		if identifier == 1 {
			if city == "moscow" {
				images := models.ProductImages{
					MoyskladID: ID,
					ImgID:      imgId,
				}
				err := db.Where("img_id = ?", imgId).FirstOrCreate(&images).Error
				if err != nil {
					log.Printf("Error: %v", err)
				} else {
					fmt.Println("Image found")
					fmt.Printf("ID\t\tImgID\t\t\n")
					fmt.Printf("%s\t%s\t\n", ID, imgId)
				}
			} else if city == "saratov" {
				images := models.ProductImagesSaratov{
					MoyskladID: ID,
					ImgID:      imgId,
				}
				err := db.Where("img_id = ?", imgId).FirstOrCreate(&images).Error
				if err != nil {
					log.Printf("Error: %v", err)
				} else {
					fmt.Println("Image found")
					fmt.Printf("ID\t\tImgID\t\t\n")
					fmt.Printf("%s\t%s\t\n", ID, imgId)
				}

			}
		} else {
			if city == "moscow" {
				images := models.ModificationImages{
					ModID: ID,
					ImgID: imgId,
				}
				err := db.Where("img_id = ?", imgId).FirstOrCreate(&images).Error
				if err != nil {
					log.Printf("Error finding/creating category: %v", err)
				}
			} else if city == "saratov" {
				images := models.ModificationImagesSaratov{
					ModID: ID,
					ImgID: imgId,
				}
				err := db.Where("img_id = ?", imgId).FirstOrCreate(&images).Error
				if err != nil {
					log.Printf("Error finding/creating category: %v", err)
				}
			}
		}
	}
	// fmt.Println(imgMap)
	return imgMap
}

func downloadImages(city string, headers http.Header, maps map[string]string, db *gorm.DB, identifier int) error {
	if maps == nil {
		return nil
	}
	fmt.Println(maps)

	// Create the directory if it doesn't exist
	err := os.MkdirAll("static/images", os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	for key, val := range maps {
		filepath := fmt.Sprintf("static/images/%s.png", key)

		// Open file for writing
		out, err := os.Create(filepath)
		if err != nil {
			return fmt.Errorf("failed to create file: %w", err)
		}

		client := &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return nil
			},
		}
		req, err := http.NewRequest("GET", val, nil)
		if err != nil {
			out.Close()
			return fmt.Errorf("failed to create request: %w", err)
		}

		// Add headers to the request
		for key, values := range headers {
			for _, value := range values {
				req.Header.Add(key, value)
			}
		}
		fmt.Printf("Headers: %+v\n", req.Header)
		fmt.Printf("Requesting image from: %s\n", val)

		// Perform the request
		response, err := client.Do(req)
		if err != nil {
			out.Close()
			return fmt.Errorf("failed to make request: %w", err)
		}
		defer response.Body.Close()

		// Check response status
		if response.StatusCode != http.StatusOK {
			out.Close()
			body, _ := io.ReadAll(response.Body)
			return fmt.Errorf("failed to download image: %s\nResponse Body: %s", response.Status, string(body))
		}

		// Write the response body to the file
		written, err := io.Copy(out, response.Body)
		out.Close()
		if err != nil {
			return fmt.Errorf("failed to write file: %w", err)
		}
		fmt.Printf("Image saved: %s (%d bytes)\n", filepath, written)

		// Save image details to the database
		err = saveImageDetails(city, db, identifier, key, filepath)
		if err != nil {
			return fmt.Errorf("failed to save image details to DB: %w", err)
		}
	}

	return nil
}

// saveImageDetails saves image details to the database based on city and identifier
func saveImageDetails(city string, db *gorm.DB, identifier int, key, filepath string) error {
	var err error
	switch identifier {
	case 1:
		if city == "moscow" {
			images := models.ProductImages{ImgID: key, ImageURL: filepath}
			err = db.Where(models.ProductImages{ImgID: key}).Assign(images).FirstOrCreate(&images).Error
		} else if city == "saratov" {
			images := models.ProductImagesSaratov{ImgID: key, ImageURL: filepath}
			err = db.Where(models.ProductImagesSaratov{ImgID: key}).Assign(images).FirstOrCreate(&images).Error
		}
	default:
		if city == "moscow" {
			images := models.ModificationImages{ImgID: key, ImageURL: filepath}
			err = db.Where(models.ModificationImages{ImgID: key}).Assign(images).FirstOrCreate(&images).Error
		} else if city == "saratov" {
			images := models.ModificationImagesSaratov{ImgID: key, ImageURL: filepath}
			err = db.Where(models.ModificationImagesSaratov{ImgID: key}).Assign(images).FirstOrCreate(&images).Error
		}
	}
	return err
}

// FUNCTIONS FOR TESTING SINGlE IMAGE
func GetSingleSaveDownloadProductImages(city string, db *gorm.DB) {
	identifier := 1
	headers := GetToken(city)
	id := "b437f3cc-b709-11ef-0a80-04420013636d"
	endpoint := fmt.Sprintf("https://api.moysklad.ru/api/remap/1.2/entity/product/%s/images", id)
	// endpoint := fmt.Sprintf("https://api.moysklad.ru/api/remap/1.2/entity/variant?filter=modsid={%s}", moyskladId)
	productsImages, err := GetImages(headers, endpoint)
	if err != nil {
		fmt.Println(err)
	}
	if productsImages == nil || len(productsImages) == 0 {
		fmt.Println("Images not found. Skipping.")
	}
	maps := SaveImages(city, productsImages, id, db, identifier)
	downloadImages(city, headers, maps, db, identifier)
}

func GetSingleSaveDownloadModImages(city string, db *gorm.DB) {
	identifier := 2
	headers := GetToken(city)
	modId := []string{
		"858cf1ce-b70a-11ef-0a80-16a80013ec7c",
		"85910b87-b70a-11ef-0a80-16a80013ec87",
		"8593c1a0-b70a-11ef-0a80-16a80013ec91",
		"85965bca-b70a-11ef-0a80-16a80013ec9b",
		"85990901-b70a-11ef-0a80-16a80013eca5",
		"859c79a6-b70a-11ef-0a80-16a80013ecaf",
		"859f0757-b70a-11ef-0a80-16a80013ecb9",
		"85a18c2d-b70a-11ef-0a80-16a80013ecc3",
		"85a47a25-b70a-11ef-0a80-16a80013eccf",
		"85a74fd0-b70a-11ef-0a80-16a80013ecd9",
		"85aa579d-b70a-11ef-0a80-16a80013ece3",
		"85acedfa-b70a-11ef-0a80-16a80013eced",
	}
	for _, id := range modId {
		endpoint := fmt.Sprintf("https://api.moysklad.ru/api/remap/1.2/entity/variant/%s/images", id)
		variantsImages, err := GetImages(headers, endpoint)
		if err != nil {
			fmt.Println(err)
			continue // Skip to the next iteration if there's an error
		}
		if variantsImages == nil || len(variantsImages) == 0 {
			fmt.Println("Images not found. Skipping.")
			continue // Skip to the next iteration if productsImages is empty
		}
		maps := SaveImages(city, variantsImages, id, db, identifier)
		downloadImages(city, headers, maps, db, identifier)
	}
}
// END TEST FUNCTIONS

// not uses
// func GetSaveProductImages(city string, db *gorm.DB) {
// 	identifier := 1
// 	headers := GetToken(city)
// 	moyskladId := GetMoyskladID(city, db)
// 	for _, id := range moyskladId {
// 		endpoint := fmt.Sprintf("https://api.moysklad.ru/api/remap/1.2/entity/product/%s/images", id)
// 		// endpoint := fmt.Sprintf("https://api.moysklad.ru/api/remap/1.2/entity/variant?filter=modsid={%s}", moyskladId)
// 		productsImages, err := GetImages(headers, endpoint)
// 		if err != nil {
// 			fmt.Println(err)
// 			continue // Skip to the next iteration if there's an error
// 		}
// 		if productsImages == nil || len(productsImages) == 0 {
// 			fmt.Println("Images not found. Skipping.")
// 			continue // Skip to the next iteration if productsImages is empty
// 		}
// 		SaveImages(city, productsImages, id, db, identifier)
// 	}
// }
