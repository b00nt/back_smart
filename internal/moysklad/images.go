package moysklad

import (
	"back/internal/models"
	"fmt"
	"gorm.io/gorm"
	"io"
	"log"
	"net/http"
	"os"
)

func GetSaveProductImages(db *gorm.DB) {
	identifier := 1
	headers := GetToken()
	moyskladId := GetMoyskladID(db)
	for _, id := range moyskladId {
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
		SaveImages(productsImages, id, db, identifier)
	}
}

func GetSaveDownloadProductImages(db *gorm.DB) {
	identifier := 1
	headers := GetToken()
	moyskladId := GetMoyskladID(db)
	for _, id := range moyskladId {
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
		maps := SaveImages(productsImages, id, db, identifier)
		downloadImages(headers, maps, db, identifier)
	}
}

func GetSaveDownloadModImages(db *gorm.DB) {
	identifier := 2
	headers := GetToken()
	modId := GetModID(db)
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
		maps := SaveImages(variantsImages, id, db, identifier)
		downloadImages(headers, maps, db, identifier)
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

func SaveImages(images []interface{}, ID string, db *gorm.DB, identifier int) map[string]string {
	var imgMap = make(map[string]string, len(images))
	for _, item := range images {
		data := item.(map[string]interface{})

		// Access the meta field
		meta := data["meta"].(map[string]interface{})
		imgId := extractImageURL(meta["href"].(string))
		downloadHref := meta["downloadHref"].(string)
		imgMap[imgId] = downloadHref

		if identifier == 1 {
			images := models.ProductImages{
				MoyskladID: ID,
				ImgID:      imgId,
			}
			err := db.Where("img_id = ?", imgId).FirstOrCreate(&images).Error
			if err != nil {
				log.Printf("Error: %v", err)
			} else {
				fmt.Println("Image found")
				fmt.Printf("ID\tImgID\t\n")
				fmt.Printf("%s\t%s\t\n", ID, imgId)
			}
		} else {
			images := models.ModificationImages{
				ModID: ID,
				ImgID: imgId,
			}
			err := db.Where("img_id = ?", imgId).FirstOrCreate(&images).Error
			if err != nil {
				log.Printf("Error finding/creating category: %v", err)
			}
		}
	}
	return imgMap
}

func downloadImages(headers http.Header, maps map[string]string, db *gorm.DB, identifier int) error {
	if maps == nil {
		return nil
	}
	// Create the directory if it doesn't exist
	err := os.MkdirAll("static/images", os.ModePerm)
	if err != nil {
		return err
	}

	for key, val := range maps {

		filepath := fmt.Sprintf("static/images/%s.png", key)
		out, err := os.Create(filepath)
		if err != nil {
			return err
		}
		defer out.Close()

		client := &http.Client{}
		req, err := http.NewRequest("GET", val, nil)
		if err != nil {
			return err
		}

		// Set headers
		for key, values := range headers {
			for _, value := range values {
				req.Header.Add(key, value)
			}
		}

		// Make the request
		fmt.Printf("Requesting image from: %s\n", val)
		response, err := client.Do(req)
		if err != nil {
			return err
		}
		defer response.Body.Close()

		// Check if the response status is OK
		if response.StatusCode != http.StatusOK {
			fmt.Printf("Failed to download image: %s\n", response.Status)
			return fmt.Errorf("failed to download image: %s", response.Status)
		}

		// Write the body to file
		_, err = io.Copy(out, response.Body)
		if err != nil {
			return err
		}

		if identifier == 1 {
			images := models.ProductImages{
				ImgID:    key,
				ImageURL: filepath,
			}
			result := db.Where(models.ProductImages{ImgID: key}).
				Assign(images).
				FirstOrCreate(&images)
			if result != nil {
				return err
			}
		} else {
			images := models.ModificationImages{
				ImgID:    key,
				ImageURL: filepath,
			}
			result := db.Where(models.ModificationImages{ImgID: key}).
				Assign(images).
				FirstOrCreate(&images)
			if result != nil {
				return err
			}
		}
	}
	return nil
}
