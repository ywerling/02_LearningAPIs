package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

// ApiResponse represents the expected JSON response structure
type ApiResponse struct {
	Results struct {
		Sunrise string `json:"sunrise"`
		Sunset  string `json:"sunset"`
	} `json:"results"`
	Status string `json:"status"`
	TZID   string `json:"tzid"`
}

func main() {
	lat, lng := readUserInput()
	apiURL := buildAPIURL(lat, lng)
	data, err := fetchSunTimes(apiURL)
	if err != nil {
		log.Fatalf("Error fetching data: %v", err)
	}

	printSunTimes(data)
}

// readUserInput prompts the user for latitude and longitude
func readUserInput() (string, string) {
	var lat, lng string
	fmt.Println("Enter the latitude and longitude (Example: '48.14816 17.10674' for Bratislava or '48.20849 16.37208' for Vienna):")
	_, err := fmt.Scanln(&lat, &lng)
	if err != nil {
		log.Fatalf("Invalid input: %v", err)
	}
	return lat, lng
}

// buildAPIURL constructs the API URL with query parameters
func buildAPIURL(lat, lng string) string {
	baseURL := "https://api.sunrise-sunset.org/json"
	params := url.Values{}
	params.Add("lat", lat)
	params.Add("lng", lng)
	return fmt.Sprintf("%s?%s", baseURL, params.Encode())
}

// fetchSunTimes makes the HTTP request and parses the response into ApiResponse struct
func fetchSunTimes(apiURL string) (*ApiResponse, error) {
	client := &http.Client{}

	// Create a new GET request
	request, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("creating GET request: %w", err)
	}
	request.Header.Set("Content-Type", "application/json")

	// Execute the request
	response, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("making request: %w", err)
	}
	defer response.Body.Close()

	// Read response body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	// Parse JSON response
	var data ApiResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("parsing JSON: %w", err)
	}

	// Ensure response status is OK
	if data.Status != "OK" {
		return nil, fmt.Errorf("API returned error status: %s", data.Status)
	}

	return &data, nil
}

// printSunTimes displays the sunrise and sunset times
func printSunTimes(data *ApiResponse) {
	fmt.Println("Sunrise:", data.Results.Sunrise, data.TZID)
	fmt.Println("Sunset:", data.Results.Sunset, data.TZID)
}
