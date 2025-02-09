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
	Product    string       `json:"product"`
	Init       string       `json:"init"`
	DataSeries []DataSeries `json:"dataseries"`
}

// Example of datapoint
// {
// "timepoint" : 72,
// "cloudcover" : 9,
// "seeing" : 3,
// "transparency" : 4,
// "lifted_index" : 2,
// "rh2m" : 11,
// "wind10m" : {
// "direction" : "SE",
// "speed" : 2
// },
// "temp2m" : 9,
// "prec_type" : "none"
// }
//

// DataSeries represents each entry in the dataseries array of weather forecasts
type DataSeries struct {
	Timepoint    int     `json:"timepoint"`
	CloudCover   int     `json:"cloudcover"`
	Seeing       int     `json:"seeing"`
	Transparency int     `json:"transparency"`
	LiftedIndex  int     `json:"lifted_index"`
	RH2m         int     `json:"rh2m"`
	Wind10m      Wind10m `json:"wind10m"`
	Temp2m       int     `json:"temp2m"`
	PrecType     string  `json:"prec_type"`
}

// Wind10m represents the wind information at 10 meters above ground.
type Wind10m struct {
	Direction string `json:"direction"`
	Speed     int    `json:"speed"`
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
	fmt.Println("Full Response (Raw Data):", string(body))

	// Parse JSON response
	var data ApiResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("parsing JSON: %w", err)
	}

	return &data, nil

}

func getCloudCover(cloudCover int) string {

	switch cloudCover {
	case 1:
		return "0-6 %"
	case 2:
		return "6-19 %"
	case 3:
		return "19-31 %"
	case 4:
		return "31-44 %"
	case 5:
		return "44-56 %"
	case 6:
		return "56-69 %"
	case 7:
		return "69-81 %"
	case 8:
		return "81-94 %"
	case 9:
		return "94-100 %"
	default:
		return "undefined"

	}
	return "N/A"
}

func printWeatherData(data *ApiResponse) {
	fmt.Println("Printing Weather Forecast for Location:")
	fmt.Println("Initial timestamp:", data.Init)
	fmt.Println("Product:", data.Product)
	fmt.Println("Timepoint", data.DataSeries[0].Timepoint, "hours")
	fmt.Println("Cloud cover", getCloudCover(data.DataSeries[0].CloudCover))

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
