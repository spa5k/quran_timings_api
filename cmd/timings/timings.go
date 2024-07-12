package timings

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// Define the structure of the JSON response
type Recitation struct {
	ID             int    `json:"id"`
	ReciterName    string `json:"reciter_name"`
	Style          string `json:"style"`
	TranslatedName struct {
		Name         string `json:"name"`
		LanguageName string `json:"language_name"`
	} `json:"translated_name"`
}

type Response struct {
	Recitations []Recitation `json:"recitations"`
}

func FetchQuranTimingReciters() {
	// Define the URL
	url := "https://api.quran.com/api/v4/resources/recitations"

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
	var response Response
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

	fmt.Println("Data successfully written to editions.json")
}
