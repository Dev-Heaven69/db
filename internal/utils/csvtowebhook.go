package utils

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func SendFileToWebhook(url, filePath, requesteeEmail, discordUsername string) error {
	// Determine the file type based on its extension
	var data interface{}
	var jsonData []byte
	var err error

	if strings.HasSuffix(filePath, ".json") {
		// Read JSON file
		jsonData, err = ioutil.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to read JSON file: %w", err)
		}
		// Unmarshal into a generic interface; adjust to specific type as needed
		err = json.Unmarshal(jsonData, &data)
		if err != nil {
			return fmt.Errorf("failed to unmarshal JSON: %w", err)
		}
	} else if strings.HasSuffix(filePath, ".csv") {
		// Read CSV file and convert to JSON
		file, err := os.Open(filePath)
		if err != nil {
			return fmt.Errorf("failed to open CSV file: %w", err)
		}
		defer file.Close()
		reader := csv.NewReader(file)
		records, err := reader.ReadAll()
		if err != nil {
			return fmt.Errorf("failed to read CSV records: %w", err)
		}

		// Convert CSV records to JSON
		data = records 
		jsonData, err = json.Marshal(data)
		if err != nil {
			return fmt.Errorf("failed to marshal records to JSON: %w", err)
		}
	} else {
		return fmt.Errorf("unsupported file type")
	}

	// Create an io.Reader from the JSON data
	requestBody := bytes.NewReader(jsonData)

	// Create the HTTP request
	request, err := http.NewRequest("POST", url, requestBody)
	if err != nil {
		return fmt.Errorf("cannot create request: %w", err)
	}

	// Set the necessary headers
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-REQUESTEE-EMAIL", requesteeEmail)
	request.Header.Set("X-DISCORD-USERNAME", discordUsername)

	// Send the request using a new HTTP client
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer response.Body.Close() // Ensure the response body is closed

	// Optional: Read and log the response body for debugging
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	fmt.Println("API Call Made, Response Status:", response.Status)
	fmt.Println("Response Body:", string(body))

	return nil
}
