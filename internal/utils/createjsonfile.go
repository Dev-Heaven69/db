package utils

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/DevHeaven/db/domain/models"
)

func PayloadToJSON(data []models.Payload, filename string, requesteeEmail string,emails string) (string, error) {
	// Open existing CSV file
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Cannot open file:", err)
		return "", err
	}
	defer file.Close()

	// Read CSV data
	reader := csv.NewReader(file)
	csvLines, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Cannot read file:", err)
		return "", err
	}

	// Convert CSV lines to JSON objects
	var jsonData []map[string]interface{}
	if emails == "personal" {
		jsonData = convertCSVToJSONForPersonalOnly(csvLines, data)
	}
	if emails == "scan"{
		jsonData = convertCSVToJSON(csvLines, data)
	}
	if emails == "professional"{
		jsonData = convertCSVToJSONForProfessionalOnly(csvLines, data)
	}

	// Generate a new filename
	newFilename := fmt.Sprintf("data/response_%s%v.json", requesteeEmail, time.Now().Unix())
	newFile, err := os.Create(newFilename)
	if err != nil {
		fmt.Println("Cannot open file:", err)
		return "", err
	}
	defer newFile.Close()

	// Write JSON data to file
	encoder := json.NewEncoder(newFile)
	encoder.SetIndent("", "  ") // For pretty printing
	if err := encoder.Encode(jsonData); err != nil {
		fmt.Println("Cannot write to file:", err)
		return "", err
	}

	fmt.Printf("Payload appended successfully to %s!\n", newFilename)
	return newFilename, nil
}

func convertCSVToJSON(csvLines [][]string, payloads []models.Payload) []map[string]interface{} {
	// Initialize JSON data array
	var jsonData []map[string]interface{}

	for i, line := range csvLines {
		jsonObj := make(map[string]interface{})
		// Assuming CSV has headers in the first line
		if i == 0 {
			continue // skip header for data population, adapt if you need headers as data keys
		}

		// Map existing CSV columns to JSON
		for j, header := range csvLines[0] {
			jsonObj[header] = line[j]
		}

		// Add payload data to JSON, match the row with payload index
		if i-1 < len(payloads) {
			payload := payloads[i-1]
			if len(payload.Emails) > 0 {
				jsonObj["Email"] = payload.Emails[0]
			} else {
				jsonObj["Email"] = ""
			}
			if len(payload.Telephone) > 0 {
				jsonObj["Telephone"] = payload.Telephone[0]
			} else {
				jsonObj["Telephone"] = ""
			}
		} else {
			jsonObj["Email"] = ""
			jsonObj["Telephone"] = ""
		}
		jsonData = append(jsonData, jsonObj)
	}
	return jsonData
}

func convertCSVToJSONForPersonalOnly(csvLines [][]string, payloads []models.Payload) []map[string]interface{} {
	var jsonData []map[string]interface{}

	// Define personal email suffixes
	personalSuffixes := []string{"@gmail.com", "@yahoo.in", "@hotmail.me", "@outlook.com", "@protonmail.com", "hotmail.com", "yahoo.com"}

	for i, line := range csvLines {
		jsonObj := make(map[string]interface{})
		// Assuming CSV has headers in the first line
		if i == 0 {
			continue // skip header for data population, adapt if you need headers as data keys
		}

		// Map existing CSV columns to JSON
		for j, header := range csvLines[0] {
			jsonObj[header] = line[j]
		}

		// Add payload data to JSON, match the row with payload index
		if i-1 < len(payloads) {
			payload := payloads[i-1]

			var personalEmails []string

			// Filter emails into personal and professional categories
			for _, email := range payload.Emails {
				for _, suffix := range personalSuffixes {
					if strings.HasSuffix(email, suffix) {
						personalEmails = append(personalEmails, email)
						break
					}
				}
			}

			// Append to JSON object
			jsonObj["PersonalEmails"] = personalEmails

			if len(payload.Telephone) > 0 {
				jsonObj["Telephone"] = payload.Telephone[0]
			} else {
				jsonObj["Telephone"] = ""
			}
		} else {
			jsonObj["PersonalEmails"] = []string{}
			jsonObj["ProfessionalEmails"] = []string{}
			jsonObj["Telephone"] = ""
		}
		jsonData = append(jsonData, jsonObj)
	}
	return jsonData
}

func convertCSVToJSONForProfessionalOnly(csvLines [][]string, payloads []models.Payload) []map[string]interface{} {
	// Initialize JSON data array
	var jsonData []map[string]interface{}

	// Define personal email suffixes
	personalSuffixes := []string{"@gmail.com", "@yahoo.in", "@hotmail.me", "@outlook.com", "@protonmail.com", "hotmail.com", "yahoo.com"}

	for i, line := range csvLines {
		jsonObj := make(map[string]interface{})
		// Assuming CSV has headers in the first line
		if i == 0 {
			continue // skip header for data population, adapt if you need headers as data keys
		}

		// Map existing CSV columns to JSON
		for j, header := range csvLines[0] {
			jsonObj[header] = line[j]
		}

		// Add payload data to JSON, match the row with payload index
		if i-1 < len(payloads) {
			payload := payloads[i-1]

			var personalEmails []string
			var professionalEmails []string

			// Filter emails into personal and professional categories
			for _, email := range payload.Emails {
				isPersonal := false
				for _, suffix := range personalSuffixes {
					if strings.HasSuffix(email, suffix) {
						personalEmails = append(personalEmails, email)
						isPersonal = true
						break
					}
				}
				if !isPersonal {
					professionalEmails = append(professionalEmails, email)
				}
			}

			// Append to JSON object
			jsonObj["ProfessionalEmails"] = professionalEmails

			if len(payload.Telephone) > 0 {
				jsonObj["Telephone"] = payload.Telephone[0]
			} else {
				jsonObj["Telephone"] = ""
			}
		} else {
			jsonObj["ProfessionalEmails"] = []string{}
			jsonObj["Telephone"] = ""
		}
		jsonData = append(jsonData, jsonObj)
	}
	return jsonData
}


func PayloadToJSONforFiltering(data []models.Payload, filename string, requesteeEmail string) (string, error) {
	// Open existing CSV file
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Cannot open file:", err)
		return "", err
	}
	defer file.Close()

	// Read CSV data
	reader := csv.NewReader(file)
	csvLines, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Cannot read file:", err)
		return "", err
	}

	// Convert CSV lines to JSON objects
	jsonData := convertCSVToJSONForFiltering(csvLines, data)

	// Generate a new filename
	newFilename := fmt.Sprintf("data/response_%s_%v.json", strings.ReplaceAll(requesteeEmail, "@", "_"), time.Now().Unix())
	newFile, err := os.Create(newFilename)
	if err != nil {
		fmt.Println("Cannot create file:", err)
		return "", err
	}
	defer newFile.Close()

	// Write JSON data to file
	encoder := json.NewEncoder(newFile)
	encoder.SetIndent("", "  ") // For pretty printing
	if err := encoder.Encode(jsonData); err != nil {
		fmt.Println("Cannot write to file:", err)
		return "", err
	}

	fmt.Printf("Payload appended successfully to %s!\n", newFilename)
	return newFilename, nil
}

func convertCSVToJSONForFiltering(csvLines [][]string, payloads []models.Payload) []map[string]interface{} {
	// Initialize JSON data array
	var jsonData []map[string]interface{}

	// Define personal email suffixes
	personalSuffixes := []string{"@gmail.com", "@yahoo.in", "@hotmail.me"}

	for i, line := range csvLines {
		jsonObj := make(map[string]interface{})
		// Assuming CSV has headers in the first line
		if i == 0 {
			continue // skip header for data population, adapt if you need headers as data keys
		}

		// Map existing CSV columns to JSON
		for j, header := range csvLines[0] {
			jsonObj[header] = line[j]
		}

		// Add payload data to JSON, match the row with payload index
		if i-1 < len(payloads) {
			payload := payloads[i-1]

			var personalEmails []string
			var professionalEmails []string

			// Filter emails into personal and professional categories
			for _, email := range payload.Emails {
				isPersonal := false
				for _, suffix := range personalSuffixes {
					if strings.HasSuffix(email, suffix) {
						personalEmails = append(personalEmails, email)
						isPersonal = true
						break
					}
				}
				if !isPersonal {
					professionalEmails = append(professionalEmails, email)
				}
			}

			// Append to JSON object
			jsonObj["PersonalEmails"] = personalEmails
			jsonObj["ProfessionalEmails"] = professionalEmails

			if len(payload.Telephone) > 0 {
				jsonObj["Telephone"] = payload.Telephone[0]
			} else {
				jsonObj["Telephone"] = ""
			}
		} else {
			jsonObj["PersonalEmails"] = []string{}
			jsonObj["ProfessionalEmails"] = []string{}
			jsonObj["Telephone"] = ""
		}
		jsonData = append(jsonData, jsonObj)
	}
	return jsonData
}
