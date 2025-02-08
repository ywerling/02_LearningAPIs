package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

type ApiResponse struct {
	Results struct {
		Sunrise string `json:"sunrise"`
		Sunset  string `json:"sunset"`
	} `json:"results"`
	Status string `json:"status"`
	TZID   string `json:"tzid"`
}

func main() {

	// Read latitude and longitude from the command line
	fmt.Println("Enter the latitude and longitude:")
	fmt.Println("Example Bratislava 48.14816 17.10674:")
	fmt.Println("Example Vienna: 48.210033 16.363449")
	fmt.Println("Enter values:\n")
	var lat, lng string
	n, err := fmt.Scanln(&lat, &lng)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("number of items read: %d\n", n)
	fmt.Printf("read line: %s %s\n", lat, lng)

	// Prepare http connection

	client := &http.Client{}

	// Define Sunrise Sunset URL
	suntimesURL := "https://api.sunrise-sunset.org/json"

	// Create URL with query parameters latitude and longitude
	// for Vienna 48.210033 16.363449
	params := url.Values{}
	// params.Add("lat", "48.210033")
	// params.Add("lng", "16.363449")
	params.Add("lat", lat)
	params.Add("lng", lng)
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
	var data ApiResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}
	//
	// Print relevant parameters response
	if data.Status == "OK" {
		fmt.Println("Sunrise:", data.Results.Sunrise, data.TZID)
		fmt.Println("Sunset:", data.Results.Sunset, data.TZID)
	} else {
		fmt.Println("Times of sunrise and sunset not available!")
	}

}
