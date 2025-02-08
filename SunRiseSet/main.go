package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type ApiResponse struct {
	Sunrise string `json:"sunrise"`
	Sunset  string `json:"sunset"`
}

func main() {
	client := &http.Client{}

	// Define Sunrise Sunset URL
	suntimesURL := "https://api.sunrise-sunset.org/json"

	// Create URL with query parameters latitude and longitude
	// for Vienna 48.210033 16.363449
	params := url.Values{}
	params.Add("lat", "48.210033")
	params.Add("lng", "16.363449")
	fullURL := fmt.Sprintf("%s?%s", suntimesURL, params.Encode())

	// Create a GET request
	request, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		fmt.Println("Error creating GET request:", err)
		return
	}

	// Set headers
	request.Header.Set("Content-Type", "application/json")

	// Execute the request
	response, err := client.Do(request)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer response.Body.Close()

	// Read the response
	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	// Print the full string response
	fmt.Println("Full Response (Raw Data):", string(body))

	// Parse JSON into the struct with Unmarshall
	// var data ApiResponse
	// err = json.Unmarshal(body, &data)
	// if err != nil {
	// fmt.Println("Error parsing JSON:", err)
	// return
	// }
	//
	// Print relevant parameters response
	// fmt.Println("Sunrise (UTC):", data.Sunrise)
	// fmt.Println("Sunset (UTC):", data.Sunset)

	// Use json.Decoder to decode the JSON response
	var data ApiResponse
	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(&data)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}

	// Print parsed response
	fmt.Println("Sunrise (UTC):", data.Sunrise)

}
