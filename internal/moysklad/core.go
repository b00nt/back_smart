package moysklad

import (
	// "back/internal/models"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	// "gorm.io/gorm"
)

// create header for token
func CreateHeader(city string) (http.Header, error) {
	var username, password string

	switch city {
	case "moscow":
		username = os.Getenv("MOYSKLAD_MOSCOW_USERNAME")
		password = os.Getenv("MOYSKLAD_MOSCOW_PASSWORD")
	case "saratov":
		username = os.Getenv("MOYSKLAD_SARATOV_USERNAME")
		password = os.Getenv("MOYSKLAD_SARATOV_PASSWORD")
	}

	// Base64 encode the credentials
	authString := fmt.Sprintf("%s:%s", username, password)
	b64AuthString := base64.StdEncoding.EncodeToString([]byte(authString))

	// Define the headers
	headers := http.Header{
		"Authorization": []string{fmt.Sprintf("Basic %s", b64AuthString)},
	}

	return headers, nil
}

// returns token string
func GetToken(headers http.Header) (string, error) {
	var tokenResponse struct {
		AccessToken string `json:"access_token"`
	}

	endpoint := "https://api.moysklad.ru/api/remap/1.2/security/token"

	client := &http.Client{}
	req, err := http.NewRequest("POST", endpoint, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Set the headers
	req.Header = headers

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}

	// Error check for resp.Body.Close()
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			if err == nil {
				err = closeErr
			}
		}
	}()

	// Check the status code
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API request failed with status code %d: %s", resp.StatusCode, string(body))
	}

	// Parse the response
	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	return tokenResponse.AccessToken, nil
}

// limit: 1000
func GetEssence(token string, endpoint string) ([]interface{}, int, error) {
	// Define the header
	headers := http.Header{
		"Authorization": []string{fmt.Sprintf("Bearer %s", token)},
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to create request: %w", err)
	}

	// Set the headers
	req.Header = headers

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to send request: %w", err)
	}

	// Error check for resp.Body.Close()
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			if err == nil {
				err = closeErr
			}
		}
	}()

	// Check the status code
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, 0, fmt.Errorf("API request failed with status code %d: %s", resp.StatusCode, string(body))
	}

	var result map[string]interface{}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, fmt.Errorf("io.ReadAll got: %s", err)
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, 0, fmt.Errorf("Unmarshal JSON got: %s", err)
	}

	// Extract meta data for pagination
	meta, metaOk := result["meta"].(map[string]interface{})
	totalCount := 0
	if metaOk {
		if size, ok := meta["size"].(float64); ok {
			totalCount = int(size)
		}
	}

	rows, ok := result["rows"]
	if !ok {
		return nil, totalCount, fmt.Errorf("response doesn't contain 'rows' field")
	}

	// Type assertion to convert to slice
	rowsSlice, ok := rows.([]interface{})
	if !ok {
		return nil, totalCount, fmt.Errorf("'rows' field is not a slice")
	}

	return rowsSlice, totalCount, nil
}
