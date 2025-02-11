package main

import (
	"fmt"
	"log"
	"net/url"

	"github.com/ywerling/02_LearningAPIs/iohelper"
	"github.com/ywerling/02_LearningAPIs/weather"
)

func main() {
	lat, lng := readUserInput()
	apiURL := buildAPIURL(lat, lng)
	data, err := weather.FetchWeatherData(apiURL)
	if err != nil {
		log.Fatalf("Error fetching data: %v", err)
	}

	printWeatherData(data)
	iohelper.WriteToCSV(data, lat, lng)
}

// readUserInput prompts the user for latitude and longitude
func readUserInput() (string, string) {
	var lat, lng string
	fmt.Println("Enter the latitude and longitude (e.g., '48.208 16.372' for Vienna):")
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

func printWeatherData(data *weather.ApiResponse) {
	ds := data.DataSeries[0]
	fmt.Printf(`Weather Forecast:
---------------------------
Location Timestamp: %s
Product: %s
Timepoint: %d hours
Cloud Cover: %s
Lifted Index: %s
Temperature: %s
Seeing Range: %s
Transparency: %s
Relative Humidity: %s
Precipitation Type: %s
Wind: %s 
---------------------------
`, data.Init, data.Product, ds.Timepoint, ds.CloudCoverString(),
		ds.LiftedIndexString(), ds.TemperatureString(), ds.SeeingString(),
		ds.TransparencyString(), ds.HumidityString(),
		ds.PrecType, ds.WindSpeedString())

	fmt.Printf("Full 72h forecasts available in the forecsts.csv file")
}
