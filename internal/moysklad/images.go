package moysklad

import (
	"back/internal/models"
	"fmt"
	//"gorm.io/gorm"
	"io"
	//"log"
	"net/http"
	"os"
	"strings"
)

func GetProductImages(token string, moyskladID string) ([]models.Image, error) {
	endpoint := fmt.Sprintf("https://api.moysklad.ru/api/remap/1.2/entity/product/%s/images", moyskladID)

	// Get image metadata from API
	imageData, _, err := GetEssence(token, endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to get images for product %s: %w", moyskladID, err)
	}

	var images []models.Image
	for _, imgData := range imageData {
		imgMap, ok := imgData.(map[string]interface{})
		if !ok {
			continue
		}

		// Extract image information
		meta, ok := imgMap["meta"].(map[string]interface{})
		if !ok {
			continue
		}

		// Get download URL
		downloadHref, ok := meta["downloadHref"].(string)
		if !ok {
			continue
		}

		// Extract image ID from meta href
		imgHref, ok := meta["href"].(string)
		if !ok {
			continue
		}
		parts := strings.Split(imgHref, "/")
		imgID := parts[len(parts)-1]

		// Get filename
		filename, ok := imgMap["filename"].(string)
		if !ok {
			filename = fmt.Sprintf("%s_%s.jpg", moyskladID, imgID)
		}

		// Make sure filename is URL-safe
		filename = strings.ReplaceAll(filename, " ", "_")

		// Create local file path
		localPath := fmt.Sprintf("static/images/%s", filename)

		// Download image
		imgBytes, err := downloadImage(token, downloadHref)
		if err != nil {
			return nil, err
		}

		// Save to local filesystem
		err = saveImageToStatic(imgBytes, filename)
		if err != nil {
			return nil, err
		}

		// Add to result
		images = append(images, models.Image{
			ImgID:      imgID,
			MoyskladID: moyskladID,
			ImageURL:   "/static/images/" + filename,
			ImagePath:  localPath,
		})
	}

	return images, nil
}

func downloadImage(token string, url string) ([]byte, error) {
	// Create request with authorization
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp != nil {
		defer func() {
			if closeErr := resp.Body.Close(); closeErr != nil {
				if err == nil {
					err = closeErr
				}
			}
		}()
	}

	// Check status
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to download image, status: %d", resp.StatusCode)
	}

	// Read image data
	return io.ReadAll(resp.Body)
}

func saveImageToStatic(imgData []byte, filename string) error {
	// Create directory if not exists
	err := os.MkdirAll("static/images", 0755)
	if err != nil {
		return err
	}

	// Save file
	return os.WriteFile(fmt.Sprintf("static/images/%s", filename), imgData, 0644)
}
