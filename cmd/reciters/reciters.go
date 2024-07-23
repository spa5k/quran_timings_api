package reciters

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type Reciters struct {
	Reciters []Reciter `json:"reciters"`
}

type Reciter struct {
	ID             int64  `json:"id"`
	ReciterID      int64  `json:"reciter_id"`
	Name           string `json:"name"`
	TranslatedName Qirat  `json:"translated_name"`
	Style          Style  `json:"style"`
	Qirat          Qirat  `json:"qirat"`
}

type Qirat struct {
	Name         string `json:"name"`
	LanguageName string `json:"language_name"`
}

type Style struct {
	Name         string `json:"name"`
	LanguageName string `json:"language_name"`
	Description  string `json:"description"`
}

func FetchQuranTimingReciters() {
	// Define the URL
	url := "https://api.qurancdn.com/api/qdc/audio/reciters?locale=en&fields=undefined"

	// Make the HTTP GET request
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error making GET request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %v\n", err)
		return
	}

	// Parse the JSON response
	var response Reciters
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Printf("Error parsing JSON response: %v\n", err)
		return
	}

	// Marshal the data into JSON format
	jsonData, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		fmt.Printf("Error marshalling JSON data: %v\n", err)
		return
	}

	// Write the JSON data to a file
	err = os.WriteFile("data/reciters.json", jsonData, 0644)
	if err != nil {
		fmt.Printf("Error writing to file: %v\n", err)
		return
	}

	fmt.Println("Data successfully written to reciters.json")
}
