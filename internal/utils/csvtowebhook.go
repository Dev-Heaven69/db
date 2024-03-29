package utils

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

func SendCSVToWebhook(url string) {
	// Open the CSV file
	file, err := os.Open("data/response.csv")
	if err != nil {
		fmt.Println("Cannot open file:", err)
		return
	}
	defer file.Close()

	// Create a buffer to store our request body as bytes
	var requestBody bytes.Buffer

	// Create a multipart writer
	multiPartWriter := multipart.NewWriter(&requestBody)

	// Create a new form-data header
	fileWriter, err := multiPartWriter.CreateFormFile("file", "response.csv")
	if err != nil {
		fmt.Println("Cannot create form file:", err)
		return
	}

	// Copy the file into the fileWriter
	_, err = io.Copy(fileWriter, file)
	if err != nil {
		fmt.Println("Cannot write to file:", err)
		return
	}

	// Close the multipart writer to get the terminating boundary.
	multiPartWriter.Close()

	// Create a new http request with the requestBody
	request, err := http.NewRequest("POST", url, &requestBody)
	if err != nil {
		fmt.Println("Cannot create request:", err)
		return
	}

	// Set the content type, this is very important
	request.Header.Set("Content-Type", multiPartWriter.FormDataContentType())

	// Send the request
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println("Failed to send request:", err)
	} else {
		fmt.Println("File upload response status:", response.Status)
	}
}
