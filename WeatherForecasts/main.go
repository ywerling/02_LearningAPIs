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
	tempVar string
	// Results struct {
	// Sunrise string `json:"sunrise"`
	// Sunset  string `json:"sunset"`
	// } `json:"results"`
	// Status string `json:"status"`
	// TZID   string `json:"tzid"`
}

// readUserInput prompts the user for latitude and longitude
func readUserInput() (string, string) {
	var lat, lng string
	fmt.Println("Enter the latitude and longitude in decimal format ('48.208 16.372' for Vienna):")
	_, err := fmt.Scanln(&lat, &lng)
	if err != nil {
		log.Fatalf("Invalid input: %v", err)
	}
	return lat, lng
}

// buildAPIURL constructs the API URL with query parameters
func buildAPIURL(lat, lng string) string {
	baseURL := "https://www.7timer.info/bin/astro.php"
	params := url.Values{}
	params.Add("lat", lat)
	params.Add("lng", lng)
	params.Add("ac", "0")
	params.Add("unit", "metric")
	params.Add("output", "json")
	params.Add("tzshift", "0")
	return fmt.Sprintf("%s?%s", baseURL, params.Encode())
}

// fetchWeatherData makes the HTTP request and parses the response into ApiResponse struct
func fetchWeatherData(apiURL string) (*ApiResponse, error) {
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

	return &data, nil

}

func printWeatherData(data *ApiResponse) {
	fmt.Println("Printing Weather Forecast for Location:")
}

func main() {
	lat, lng := readUserInput()
	apiURL := buildAPIURL(lat, lng)
	data, err := fetchWeatherData(apiURL)
	if err != nil {
		log.Fatalf("Error fetching data: %v", err)
	}

	printWeatherData(data)
}
