package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/DevHeaven/db/domain/models"
)

func SendToWebhook(url string, resp []models.Payload, requesteeEmail string, discordUsername string) error {
	// Marshal the response data into JSON
	jsonData, err := json.Marshal(resp)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
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

// func SendJSONToWebhook(url string,filename string) {
// 	// Open the JSON file
// 	filepath := fmt.Sprintf("data/%s", filename)
// 	file, err := os.Open(filepath)
// 	if err != nil {
// 		fmt.Println("Cannot open file:", err)
// 		return
// 	}
// 	defer file.Close()

// 	// Create a buffer to store our request body as bytes
// 	var requestBody bytes.Buffer

// 	// Copy the file into the requestBody
// 	_, err = io.Copy(&requestBody, file)
// 	if err != nil {
// 		fmt.Println("Cannot write to file:", err)
// 		return
// 	}

// 	// Create a new http request with the requestBody
// 	request, err := http.NewRequest("POST", url, &requestBody)
// 	if err != nil {
// 		fmt.Println("Cannot create request:", err)
// 		return
// 	}

// 	// Set the content type, this is very important
// 	request.Header.Set("Content-Type", "application/json")

// 	// Send the request
// 	client := &http.Client{}
