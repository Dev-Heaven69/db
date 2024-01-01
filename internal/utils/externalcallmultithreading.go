package utils

// import (
// 	"bytes"
// 	"encoding/json"
// 	"fmt"
// 	"net/http"
// 	"sync"
// )

// type InputRecord struct {
// 	FirstName string
// 	LastName  string
// 	Domain    string
// }

// type APIResponse struct {
// 	Email string
// }

// type OutputRecord struct {
// 	InputRecord
// 	Email string
// }

// var limiter *rate.Limiter

// func apiCall(input InputRecord, resultsChan chan<- OutputRecord, limiter *rate.Limiter) {
// 	limiter.Wait(nil) // rate limit our request

// 	httpClient := &http.Client{Timeout: 10 * time.Second}

// 	url := "https://punchleads.tech/api-product/incoming-webhook/find-emails-first-last"

// 	jsonBody, _ := json.Marshal(map[string]string{
// 		"api_key":    "R2P9V5Z7-T7Z5P1U3-Z4P7P2D3-U9X6L6B6",
// 		"first_name": input.FirstName,
// 		"last_name":  input.LastName,
// 		"domain":     input.Domain,
// 	})

// 	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
// 	req.Header.Set("Content-Type", "application/json")

// 	resp, err := httpClient.Do(req)
// 	if err != nil {
// 		resultsChan <- OutputRecord{input, "Error: " + err.Error()}
// 		return
// 	}
// 	defer resp.Body.Close()

// 	var apiResp APIResponse
// 	err = json.NewDecoder(resp.Body).Decode(&apiResp)
// 	if err != nil {
// 		resultsChan <- OutputRecord{input, "Error: " + err.Error()}
// 		return
// 	}

// 	resultsChan <- OutputRecord{input, apiResp.Email}
// }

// func FindEmails(inputRecords [][]string, limiter *rate.Limiter) [][]string {
// 	resultsChan := make(chan OutputRecord)

// 	var results [][]string

// 	go func() {
// 		for output := range resultsChan {
// 			results = append(results, []string{output.FirstName, output.LastName, output.Domain, output.Email})
// 		}
// 	}()

// 	for _, record := range inputRecords {
// 		input := InputRecord{record[0], record[1], record[2]}
// 		go apiCall(input, resultsChan, limiter) // start a goroutine for each API call
// 	}

// 	return results
// }

// ...

// type ApiPayload struct {
// 	ApiKey    string `json:"api_key"`
// 	FirstName string `json:"first_name"`
// 	LastName  string `json:"last_name"`
// 	Domain    string `json:"domain"`
// }

// // Function to perform a single POST request
// func postApiRequest(payload ApiPayload, wg *sync.WaitGroup, results chan<- *http.Response, errors chan<- error) {
// 	defer wg.Done()

// 	payloadBytes, err := json.Marshal(payload)
// 	if err != nil {
// 		errors <- err
// 		return
// 	}

// 	body := bytes.NewReader(payloadBytes)
// 	req, err := http.NewRequest("POST", "https://punchleads.tech/api-product/incoming-webhook/find-emails-first-last", body)
// 	if err != nil {
// 		errors <- err
// 		return
// 	}

// 	req.Header.Set("Content-Type", "application/json")

// 	resp, err := http.DefaultClient.Do(req)
// 	if err != nil {
// 		errors <- err
// 		return
// 	}

// 	results <- resp
// }

// Function to handle bulk records and perform API calls concurrently
// func QueryBulkRecords(records []ApiPayload) {
// 	var wg sync.WaitGroup
// 	results := make(chan *http.Response)
// 	errors := make(chan error)

// 	// Set a limit for maximum concurrent goroutines if necessary
// 	maxGoroutines := 20
// 	guard := make(chan struct{}, maxGoroutines)

// 	for _, record := range records {
// 		wg.Add(1)
// 		guard <- struct{}{} // would block if guard channel is already filled

// 		go func(record ApiPayload) {
// 			postApiRequest(record, &wg, results, errors)
// 			<-guard // release a spot in the guard channel
// 		}(record)
// 	}

// 	// Close channels when all goroutines are done
// 	go func() {
// 		wg.Wait()
// 		close(results)
// 		close(errors)
// 	}()

// 	// Process results and errors
// 	for {
// 		select {
// 		case result := <-results:
// 			if result != nil {
// 				// Process successful response
// 				fmt.Println(result)
// 				// ...
// 				result.Body.Close() // Don't forget to close the response body
// 			}
// 		case err := <-errors:
// 			if err != nil {
// 				// Handle error
// 				// ...
// 			}
// 		}
// 	}
// }

// ...

// Example usage:
// Assuming you have a slice of ApiPayload with bulk records
// bulkRecords := []ApiPayload{	}
// Populate with your bulk records
// }
