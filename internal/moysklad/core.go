package moysklad

import (
	// "back/internal/models"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	// "gorm.io/gorm"
)

func GetToken(city string) (http.Header, error) {
	var username, password string

	if city == "moscow" {
		username = os.Getenv("MOYSKLAD_MOSCOW_USERNAME")
		password = os.Getenv("MOYSKLAD_MOSCOW_PASSWORD")
	} else if city == "saratov" {
		username = os.Getenv("MOYSKLAD_SARATOV_USERNAME")
		password = os.Getenv("MOYSKLAD_SARATOV_PASSWORD")
	}

	// Base64 encode the credentials
	authString := fmt.Sprintf("%s:%s", username, password)
	b64AuthString := base64.StdEncoding.EncodeToString([]byte(authString))

	// Define the headers
	headers := http.Header{
		"Authorization": []string{fmt.Sprintf("Basic %s", b64AuthString)},
		"Content-Type":  []string{"application/json"},
	}

	return headers, nil
}

func GetEssence(headers http.Header, endpoint string) ([]interface{}, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	// Set headers
	for key, values := range headers {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	// Make the request
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	// Check the response status code
	if response.StatusCode == http.StatusOK {
		var result map[string]interface{}
		body, err := io.ReadAll(response.Body)
		if err != nil {
			return nil, fmt.Errorf("io.ReadAll got: %s", err)
		}

		if err := json.Unmarshal(body, &result); err != nil {
			return nil, fmt.Errorf("Unmarshal JSON got: %s", err)
		}

		if rows, ok := result["rows"]; ok {
			// rowsData, err := json.MarshalIndent(rows, "", "  ")
			// if err != nil {
			// 	return nil, err
			// }
			// fmt.Println("Rows data:", string(rowsData))
			return rows.([]interface{}), nil
		}
		return nil, fmt.Errorf("no rows found in response")
	} else {
		body, _ := io.ReadAll(response.Body)
		return nil, fmt.Errorf("failed to fetch products. Status code: %d, Message: %s", response.StatusCode, string(body))
	}
}

// Helper function to extract MoyskladID from the product URL
func extractMoyskladIDFromURL(href string) string {
	parts := strings.Split(href, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return ""
}

func extractImageURL(href string) string {
	return extractMoyskladIDFromURL(href)
}

// func GetMoyskladID(city string, db *gorm.DB) []string {
// 	var moyskladId []string
// 	var err error
//
// 	if city == "moscow" {
// 		err = db.Model(&models.Product{}).Pluck("MoyskladID", &moyskladId).Error
// 	} else if city == "saratov" {
// 		err = db.Model(&models.ProductsSaratov{}).Pluck("MoyskladID", &moyskladId).Error
// 	}
// 	if err != nil {
// 		fmt.Println("Error fetching moyskladID:", err)
// 		return []string{}
// 	}
// 	return moyskladId
// }
//
// func GetModID(city string, db *gorm.DB) []string {
// 	var modId []string
// 	var err error
//
// 	if city == "moscow" {
// 		err = db.Model(&models.Modification{}).Pluck("ModID", &modId).Error
// 	} else if city == "saratov" {
// 		err = db.Model(&models.ModificationSaratov{}).Pluck("ModID", &modId).Error
// 	}
// 	if err != nil {
// 		fmt.Println("Error fetching modificationID:", err)
// 		return []string{}
// 	}
// 	return modId
// }
