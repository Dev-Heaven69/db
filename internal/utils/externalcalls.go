package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/DevHeaven/db/domain/models"
)

type Request struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Domain    string `json:"domain"`
	ApiKey    string `json:"api_key"`
}

func ExternalApiCall(firstname string, lastname string, domain string) (models.ApiResponse, error) {
	var data models.ApiResponse
	requestData := &Request{
		FirstName: firstname,
		LastName:  lastname,
		Domain:    domain,
		ApiKey:    "R2P9V5Z7-T7Z5P1U3-Z4P7P2D3-U9X6L6B6",
	}

	jsonData, err := json.Marshal(requestData)
	if err != nil {
		fmt.Println("Error preparing request data", err)
		return data, err
	}

	client := &http.Client{
		Timeout: 2 * time.Second,
	}

	resp, err := client.Post("https://punchleads.tech/api-product/incoming-webhook/find-emails-first-last",
		"application/json",
		bytes.NewBuffer(jsonData))

	if err != nil {
		// fmt.Println("Timeout occured. Making a call to second API..")
		// secondUrl := fmt.Sprintf("https://api.leadgo.io/api/search?domain=%s&firstName=%s&lastName=%s&api_key=0f0cde1345882c9c9bfadb592561e4b871428ce1", domain, fs, ls)
		// secondClient := &http.Client{
		// 	Timeout: 4 * time.Second,
		// }
		// resp, err = secondClient.Get(secondUrl)

		// if err != nil {
		// 	fmt.Println("Error making GET call", err)
		// 	return data, err
		// }
		return data, err
	}

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		fmt.Println("Error decoding response body", err)
		return data, err
	}

	return data, nil
}
