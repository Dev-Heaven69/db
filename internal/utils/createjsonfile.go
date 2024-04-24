package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/DevHeaven/db/domain/models"
)

func WriteResponseToJson(response models.Response) (string, error) {
	// Marshal the Response struct to JSON
	data, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		return "",fmt.Errorf("error marshaling to JSON: %v", err)
	}

	var filename = fmt.Sprintf("data/response_%s%v.json", response.ResquesteeEmail, time.Now().Unix())

	// Create a file to write the JSON data
	file, err := os.Create(filename)
	if err != nil {
		return "", fmt.Errorf("error creating file: %v", err)
	}
	defer file.Close()

	// Write the JSON data to the file
	_, err = file.Write(data)
	if err != nil {
		return "", fmt.Errorf("error writing JSON to file: %v", err)
	}

	return filename,nil
}
