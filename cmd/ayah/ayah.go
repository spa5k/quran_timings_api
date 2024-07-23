package ayah

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type Reciter struct {
	ID             int    `json:"id"`
	ReciterID      int    `json:"reciter_id"`
	Name           string `json:"name"`
	TranslatedName struct {
		Name         string `json:"name"`
		LanguageName string `json:"language_name"`
	} `json:"translated_name"`
	Style struct {
		Name         string `json:"name"`
		LanguageName string `json:"language_name"`
		Description  string `json:"description"`
	} `json:"style"`
	Qirat struct {
		Name         string `json:"name"`
		LanguageName string `json:"language_name"`
	} `json:"qirat"`
	Slug string `json:"slug"`
}

type Reciters struct {
	Reciters []Reciter `json:"reciters"`
}

type ChapterData struct {
	AudioFiles []struct {
		ID           int     `json:"id"`
		ChapterID    int     `json:"chapter_id"`
		FileSize     float64 `json:"file_size"`
		Format       string  `json:"format"`
		AudioURL     string  `json:"audio_url"`
		Duration     int     `json:"duration"`
		VerseTimings []struct {
			VerseKey      string      `json:"verse_key"`
			TimestampFrom int         `json:"timestamp_from"`
			TimestampTo   int         `json:"timestamp_to"`
			Duration      int         `json:"duration"`
			Segments      [][]float64 `json:"segments"`
		} `json:"verse_timings"`
	} `json:"audio_files"`
}

func AyahTimingsPerReciter() {
	// Read the reciters.json file
	file, err := os.ReadFile("data/reciters.json")
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	// Parse the JSON data
	var reciters Reciters
	err = json.Unmarshal(file, &reciters)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	// Create a wait group for concurrency
	var wg sync.WaitGroup

	// Create a semaphore to limit the number of concurrent requests
	semaphore := make(chan struct{}, 10) // Limit to 10 concurrent requests

	// Fetch chapters for each reciter
	for _, reciter := range reciters.Reciters {
		wg.Add(1)
		go func(reciter Reciter) {
			defer wg.Done()
			for chapter := 1; chapter <= 114; chapter++ {
				semaphore <- struct{}{} // Acquire a token
				go func(reciter Reciter, chapter int) {
					defer func() { <-semaphore }() // Release the token
					fetchAndSaveChapter(reciter, chapter)
				}(reciter, chapter)
			}
		}(reciter)
	}

	wg.Wait()
}

func fetchAndSaveChapter(reciter Reciter, chapter int) {
	// Create directory for the reciter if it doesn't exist
	dir := filepath.Join("data", strings.ToLower(reciter.Style.Name), reciter.Slug)
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		fmt.Println("Error creating directory:", err)
		return
	}

	// Check if the file already exists
	filePath := filepath.Join(dir, fmt.Sprintf("%d.json", chapter))
	if _, err := os.Stat(filePath); err == nil {
		fmt.Printf("File already exists: %s\n", filePath)
		return
	}

	// Fetch chapter data from the API
	url := fmt.Sprintf("https://api.qurancdn.com/api/qdc/audio/reciters/%d/audio_files?chapter=%d&segments=true", reciter.ReciterID, chapter)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error fetching data for reciter %s, chapter %d: %v\n", reciter.Name, chapter, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error: received non-200 response code for reciter %s, chapter %d\n", reciter.Name, chapter)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body for reciter %s, chapter %d: %v\n", reciter.Name, chapter, err)
		return
	}

	// Parse the chapter data
	var chapterData ChapterData
	err = json.Unmarshal(body, &chapterData)
	if err != nil {
		fmt.Printf("Error parsing chapter data for reciter %s, chapter %d: %v\n", reciter.Name, chapter, err)
		return
	}

	// Save the chapter data to a file
	err = os.WriteFile(filePath, body, 0644)
	if err != nil {
		fmt.Printf("Error writing file for reciter %s, chapter %d: %v\n", reciter.Name, chapter, err)
		return
	}

	fmt.Printf("Saved chapter %d for reciter %s\n", chapter, reciter.Name)
}

