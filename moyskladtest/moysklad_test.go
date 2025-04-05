package moysklad_test

import (
	"back/internal/moysklad"
	"encoding/json"
	"os"
	"testing"
)

// Test function
func TestGetProducts(t *testing.T) {
	// Get the products
	response, err := moysklad.GetProducts("moscow")
	if err != nil {
		t.Fatalf("Error getting products: %v", err)
	}

	// Check if the response is nil
	if response == nil {
		t.Fatal("Got nil response from GetProducts")
	}

	// Marshal the data to JSON
	jsonData, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		t.Fatalf("Error marshaling JSON: %v", err)
	}

	// Define the file name
	fileName := "moysklad_products.json"

	// Write JSON data to the file
	err = os.WriteFile(fileName, jsonData, 0644)
	if err != nil {
		t.Fatalf("Error writing to file: %v", err)
	}

	// Print success message without removing the file
	t.Logf("JSON data written to %s", fileName)

	// Optional: You can print some basic info about the response
	t.Logf("Retrieved %d products", len(response))
}

func TestSaveProducts(t *testing.T) {
	city := "moscow"
	// Get the products
	response, err := moysklad.GetProducts(city)
	if err != nil {
		t.Fatalf("Error getting products: %v", err)
	}

	// Check if the response is nil
	if response == nil {
		t.Fatal("Got nil response from GetProducts")
	}

	err := moysklad.SaveProducts(city, response)

}
