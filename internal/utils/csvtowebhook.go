package utils

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strings"

	"github.com/DevHeaven/db/domain/models"
)

func SendFileToWebhook(url, filePath, requesteeEmail, discordUsername, responseFormat string) error {
	if responseFormat == "file" {
		// Open the file to be sent
		file, err := os.Open(filePath)
		if err != nil {
			return fmt.Errorf("failed to open file: %w", err)
		}
		defer file.Close()

		// Create a buffer to write our multipart form data
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		// Create a form file
		part, err := writer.CreateFormFile("file", file.Name())
		if err != nil {
			return fmt.Errorf("failed to create form file: %w", err)
		}
		_, err = io.Copy(part, file)
		if err != nil {
			return fmt.Errorf("failed to copy file contents: %w", err)
		}

		// Important: Close the writer to write the boundary to the buffer
		err = writer.Close()
		if err != nil {
			return fmt.Errorf("failed to close writer: %w", err)
		}

		// Create the HTTP request with the proper headers
		request, err := http.NewRequest("POST", url, body)
		if err != nil {
			return fmt.Errorf("cannot create request: %w", err)
		}
		request.Header.Set("Content-Type", writer.FormDataContentType())
		request.Header.Set("X-REQUESTEE-EMAIL", requesteeEmail)
		request.Header.Set("X-DISCORD-USERNAME", discordUsername)

		// Execute the request
		return executeRequest(request)

	} else if responseFormat == "data" {
		// Process and send file data as JSON
		var jsonData []byte
		var err error

		if strings.HasSuffix(filePath, ".json") {
			jsonData, err = ioutil.ReadFile(filePath)
		} else if strings.HasSuffix(filePath, ".csv") {
			jsonData, err = convertCSVtoJSON(filePath)
		} else {
			return fmt.Errorf("unsupported file type")
		}
		if err != nil {
			return fmt.Errorf("failed to prepare file data: %w", err)
		}

		// Create the HTTP request
		request, err := http.NewRequest("POST", url, bytes.NewReader(jsonData))
		if err != nil {
			return fmt.Errorf("cannot create request: %w", err)
		}
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("X-REQUESTEE-EMAIL", requesteeEmail)
		request.Header.Set("X-DISCORD-USERNAME", discordUsername)

		// Execute the request
		return executeRequest(request)
	} else {
		return fmt.Errorf("invalid response format")
	}
}

func SendResponseToWebhook(url,requesteeEmail, discordUsername string,response models.CSVFileData) error {
	jsonData, err := json.Marshal(response)
	if err != nil {
		fmt.Println("failed to marshal records to JSON: %w", err)
	}
	request, err := http.NewRequest("POST", url, bytes.NewReader(jsonData))
	if err != nil {
		fmt.Println("cannot create request: %w", err)
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-REQUESTEE-EMAIL", requesteeEmail)
	request.Header.Set("X-DISCORD-USERNAME", discordUsername)
	return executeRequest(request)
}

func convertCSVtoJSON(filePath string) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open CSV file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV records: %w", err)
	}

	jsonData, err := json.Marshal(records)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal records to JSON: %w", err)
	}
	return jsonData, nil
}

func executeRequest(request *http.Request) error {
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	fmt.Println("API Call Made, Response Status:", response.Status)
	fmt.Println("Response Body:", string(body))

	return nil
}
