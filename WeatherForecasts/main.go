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

// func getCloudCover(cloudCover int) string {
//
// switch cloudCover {
// case 1:
// return "0-6 %"
// case 2:
// return "6-19 %"
// case 3:
// return "19-31 %"
// case 4:
// return "31-44 %"
// case 5:
// return "44-56 %"
// case 6:
// return "56-69 %"
// case 7:
// return "69-81 %"
// case 8:
// return "81-94 %"
// case 9:
// return "94-100 %"
// default:
// return "undefined"
//
// }
//
// }
//

func getCloudCover(cloudCover int) string {
	cloudCoverMap := []string{
		"undefined", // Index 0 (not used, as valid values start from 1)
		"0-6 %",     // Index 1
		"6-19 %",    // Index 2
		"19-31 %",   // Index 3
		"31-44 %",   // Index 4
		"44-56 %",   // Index 5
		"56-69 %",   // Index 6
		"69-81 %",   // Index 7
		"81-94 %",   // Index 8
		"94-100 %",  // Index 9
	}

	if cloudCover >= 1 && cloudCover <= 9 {
		return cloudCoverMap[cloudCover]
	}
	return "undefined"
}

func getLiftedIndex(liftedIndex int) string {
	switch liftedIndex {
	case -10:
		return "below -7"
	case -6:
		return "-7 to -5"
	case -4:
		return "-5 to -3"
	case -1:
		return "-3 to 0"
	case 2:
		return "0 to 4"
	case 6:
		return "4 to 8"
	case 10:
		return "8 to 11"
	case 15:
		return "over 11"
	default:
		return "undefined"

	}
}

func getWindSpeed(windSpeed int) string {
	windSpeedMap := []string{
		"undefined",                // Index 0 (not used, as valid values start from 1)
		"Below 0.3m/s (calm)",      // Index 1
		"0.3-3.4m/s (light)",       // Index 2
		"3.4-8.0m/s (moderate)",    // Index 3
		"8.0-10.8m/s (fresh)",      // Index 4
		"10.8-17.2m/s (strong)",    // Index 5
		"17.2-24.5m/s (gale)",      // Index 6
		"24.5-32.6m/s (storm)",     // Index 7
		"Over 32.6m/s (hurricane)", // Index 8
	}

	if windSpeed >= 1 && windSpeed <= 8 {
		return windSpeedMap[windSpeed]
	}
	return "undefined"
}

func getRelativeHumidity(relHum int) string {
	switch relHum {
	case -4:
		return "0-5 %"
	case -3:
		return "5-10 %"
	case -2:
		return "10-15 %"
	case -1:
		return "15-20 %"
	case 0:
		return "20-25 %"
	case 1:
		return "25-30 %"
	case 2:
		return "30-35 %"
	case 3:
		return "35-40 %"
	case 4:
		return "40-45 %"
	case 5:
		return "45-50 %"
	case 6:
		return "50-55 %"
	case 7:
		return "55-60 %"
	case 8:
		return "60-65 %"
	case 9:
		return "65-70 %"
	case 10:
		return "70-75 %"
	case 11:
		return "75-80 %"
	case 12:
		return "80-85 %"
	case 13:
		return "85-90 %"
	case 14:
		return "90-95 %"
	case 15:
		return "95-99 %"
	case 16:
		return "100 %"
	default:
		return "undefined"

	}
}

func printWeatherData(data *ApiResponse) {
	fmt.Println("Printing Weather Forecast for Location:")
	fmt.Println("Initial timestamp:", data.Init)
	fmt.Println("Product:", data.Product)
	fmt.Println("Timepoint", data.DataSeries[0].Timepoint, "hours")
	fmt.Println("Cloud cover", getCloudCover(data.DataSeries[0].CloudCover))
	fmt.Println("Lifted Index", getLiftedIndex(data.DataSeries[0].LiftedIndex))
	fmt.Println("Temperature 2 meters", data.DataSeries[0].Temp2m, "°C")
	fmt.Println("Relative humidity 2 meters", getRelativeHumidity(data.DataSeries[0].RH2m))
	fmt.Println("Precipitation type", data.DataSeries[0].PrecType)
	fmt.Println("Wind", data.DataSeries[0].Wind10m.Direction, getWindSpeed(data.DataSeries[0].Wind10m.Speed))

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
