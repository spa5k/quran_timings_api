package ayah

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

// Define the structure of the editions.json file
type Recitation struct {
	ID             int    `json:"id"`
	ReciterName    string `json:"reciter_name"`
	Style          string `json:"style"`
	Slug           string `json:"slug"`
	TranslatedName struct {
		Name         string `json:"name"`
		LanguageName string `json:"language_name"`
	} `json:"translated_name"`
}

type Editions struct {
	Recitations []Recitation `json:"recitations"`
}

// Define the structure of the API response for verses
type Verse struct {
	ID          int    `json:"id"`
	VerseNumber int    `json:"verse_number"`
	VerseKey    string `json:"verse_key"`
	Audio       struct {
		URL      string          `json:"url"`
		Segments [][]interface{} `json:"segments"`
	} `json:"audio"`
}

type VersesResponse struct {
	Verses []Verse `json:"verses"`
}

// Cache to store already fetched data
var cache = make(map[string]bool)
var cacheMutex sync.Mutex

func AyahTimingsPerReciter() {
	// Read the reciters.json file
	editionsFile, err := os.ReadFile("data/reciters.json")
	if err != nil {
		fmt.Printf("Error reading editions.json: %v\n", err)
		return
	}

	// Parse the JSON data
	var editions Editions
	err = json.Unmarshal(editionsFile, &editions)
	if err != nil {
		fmt.Printf("Error parsing editions.json: %v\n", err)
		return
	}

	var wg sync.WaitGroup

	// Iterate over each recitation
	for _, recitation := range editions.Recitations {
		// Create a directory for the recitation inside the data folder
		dir := filepath.Join("data", recitation.Slug)
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			fmt.Printf("Error creating directory %s: %v\n", dir, err)
			continue
		}

		// Fetch and save data for chapters 1-114 concurrently
		for chapter := 1; chapter <= 114; chapter++ {
			wg.Add(1)
			go func(recitation Recitation, chapter int) {
				defer wg.Done()

				cacheKey := fmt.Sprintf("%d-%d", recitation.ID, chapter)
				cacheMutex.Lock()
				if cache[cacheKey] {
					cacheMutex.Unlock()
					fmt.Printf("Data for recitation %d chapter %d already fetched, skipping...\n", recitation.ID, chapter)
					return
				}
				cache[cacheKey] = true
				cacheMutex.Unlock()

				url := fmt.Sprintf("https://api.quran.com/api/v4/verses/by_chapter/%d?audio=%d", chapter, recitation.ID)

				resp, err := http.Get(url)
				if err != nil {
					fmt.Printf("Error making GET request for chapter %d: %v\n", chapter, err)
					return
				}
				defer resp.Body.Close()

				body, err := io.ReadAll(resp.Body)
				if err != nil {
					fmt.Printf("Error reading response body for chapter %d: %v\n", chapter, err)
					return
				}

				var versesResponse VersesResponse
				err = json.Unmarshal(body, &versesResponse)
				if err != nil {
					fmt.Printf("Error parsing JSON response for chapter %d: %v\n and url %s\n", chapter, err, url)
					return
				}

				// Extract only the audio data
				audioData := make([]map[string]interface{}, len(versesResponse.Verses))
				for i, verse := range versesResponse.Verses {
					segments := make([][]int, len(verse.Audio.Segments))
					for j, segment := range verse.Audio.Segments {
						segments[j] = make([]int, len(segment))
						for k, value := range segment {
							switch v := value.(type) {
							case string:
								intValue, err := strconv.Atoi(v)
								if err != nil {
									fmt.Printf("Error converting string to int for segment %v: %v\n", segment, err)
									return
								}
								segments[j][k] = intValue
							case float64:
								segments[j][k] = int(v)
							default:
								fmt.Printf("Unexpected type %T for segment value %v\n", v, v)
								return
							}
						}
					}

					audioData[i] = map[string]interface{}{
						"url":      verse.Audio.URL,
						"verse":    verse.VerseNumber,
						"chapter":  chapter,
						"segments": segments,
					}
				}

				// Marshal the audio data into JSON format
				jsonData, err := json.MarshalIndent(audioData, "", "  ")
				if err != nil {
					fmt.Printf("Error marshalling JSON data for chapter %d: %v\n", chapter, err)
					return
				}

				// Write the JSON data to a file
				filePath := filepath.Join(dir, fmt.Sprintf("%d.json", chapter))
				err = os.WriteFile(filePath, jsonData, 0644)
				if err != nil {
					fmt.Printf("Error writing to file %s: %v\n", filePath, err)
					return
				}

				fmt.Printf("Data for recitation %d chapter %d successfully fetched and saved.\n", recitation.ID, chapter)
			}(recitation, chapter)
		}
	}

	wg.Wait()
	fmt.Println("All data successfully fetched and saved.")
}
