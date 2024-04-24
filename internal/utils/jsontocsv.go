package utils

import (
	"encoding/csv"
	"fmt"
	"os"
	"time"

	"github.com/DevHeaven/db/domain/models"
)

func PayloadToCSV(data []models.Payload, filename string, requesteeEmail string) (string, error) {
	// Convert Payload data to CSV format
	csvData := convertPayloadToCSV(data)
	// csvData := convertPayloadToCSV(data)

	// Open existing file
	file, err := os.OpenFile(filename, os.O_RDWR, 0755)
	if err != nil {
		fmt.Println("Cannot open file:", err)
		return "", err
	}
	defer file.Close()

	// Read the file
	reader := csv.NewReader(file)
	lines, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Cannot read file:", err)
		return "", err
	}

	// Close the file to clear the file descriptor for the upcoming write operation
	file.Close()

	// Insert new data at specified position and shift other columns
	for i, line := range lines {
		if i == 0 { // For header
			lines[i] = append(line[:5], append([]string{"Email", "Telephone"}, line[5:]...)...)
		} else if i-1 < len(csvData) { // For rows
			lines[i] = append(line[:5], append(csvData[i-1], line[5:]...)...)
		} else { // For empty json
			lines[i] = append(line[:5], append([]string{"", ""}, line[5:]...)...)
		}
	}

	// Overwrite file with the new data
	newFilename := fmt.Sprintf("data/response_%s%v.csv", requesteeEmail, time.Now().Unix())
	newfile, err := os.Create(newFilename)
	if err != nil {
		fmt.Println("Cannot open file:", err)
		return "", err
	}
	defer newfile.Close()

	writer := csv.NewWriter(newfile)
	defer writer.Flush()

	err = writer.WriteAll(lines)
	if err != nil {
		fmt.Println("Cannot write to file:", err)
	}

	fmt.Printf("Payload appended successfully to %s!\n", newFilename)
	return newFilename, nil
}

func convertPayloadToCSV(data []models.Payload) [][]string {
	var rows [][]string
	// Check if data is empty
	if len(data) == 0 {
		rows = append(rows, []string{"", ""}) // append blanks
		return rows
	}

	for _, payload := range data {
		var email, tel string
		// Check if payload is empty
		if payload.Emails == nil && payload.Telephone == nil {
			rows = append(rows, []string{"", ""}) // append blanks
		} else {
			if len(payload.Emails) > 0 {
				email = payload.Emails[0]
			}
			if len(payload.Telephone) > 0 {
				tel = payload.Telephone[0]
			}
			rows = append(rows, []string{email, tel})
		}
	}
	return rows
}
