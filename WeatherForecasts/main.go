package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
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

func (ds DataSeries) CloudCoverString() string {
	return getCloudCover(ds.CloudCover)
}

func (ds DataSeries) SeeingString() string {
	return getSeeing(ds.Seeing)
}

func (ds DataSeries) TransparencyString() string {
	return getTransparency(ds.Transparency)
}

func (ds DataSeries) LiftedIndexString() string {
	return getLiftedIndex(ds.LiftedIndex)
}

func (ds DataSeries) WindSpeedString() string {
	return ds.Wind10m.Direction + " " + getWindSpeed(ds.Wind10m.Speed)
}

func (ds DataSeries) HumidityString() string {
	return getRelativeHumidity(ds.RH2m)
}

func (ds DataSeries) TemperatureString() string {
	return strconv.Itoa(ds.Temp2m) + " C"
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

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	// Read response body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}
	// Print the entire raw json response, uncomment when debugging
	// fmt.Println("Full Response (Raw Data):", string(body))

	// Parse JSON response
	var data ApiResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("parsing JSON: %w", err)
	}

	return &data, nil

}

// getMappedValue returns the mapped value for a given key using a provided map
func getMappedValue(key int, valueMap []string) string {
	if key >= 0 && key < len(valueMap) {
		return valueMap[key]
	}
	return "undefined"
}

var cloudCoverMap = []string{
	"undefined", "0-6 %", "6-19 %", "19-31 %", "31-44 %", "44-56 %", "56-69 %", "69-81 %", "81-94 %", "94-100 %",
}

func getCloudCover(cloudCover int) string {
	return getMappedValue(cloudCover, cloudCoverMap)
}

var seeingMap = []string{
	"undefined", "<0.5", "0.5-0.75", "0.75-1", "1-1.25", "1.25-1.5", "1.5-2", "2-2.5", ">2.5",
}

func getSeeing(seeing int) string {

	return getMappedValue(seeing, seeingMap)
}

var transparencyMap = []string{
	"undefined", "<0.3", "0.3-0.4", "0.4-0.5", "0.5-0.6", "0.6-0.7", "0.7-0.85", "0.85-1", ">1",
}

func getTransparency(transparency int) string {
	return getMappedValue(transparency, transparencyMap)
}

func getLiftedIndex(liftedIndex int) string {
	liftedIndexMap := map[int]string{
		-10: "below -7",
		-6:  "-7 to -5",
		-4:  "-5 to -3",
		-1:  "-3 to 0",
		2:   "0 to 4",
		6:   "4 to 8",
		10:  "8 to 11",
		15:  "over 11",
	}
	if val, ok := liftedIndexMap[liftedIndex]; ok {
		return val
	}
	return "undefined"
}

var windSpeedMap = []string{
	"undefined", "Below 0.3m/s (calm)", "0.3-3.4m/s (light)", "3.4-8.0m/s (moderate)", "8.0-10.8m/s (fresh)", "10.8-17.2m/s (strong)", "17.2-24.5m/s (gale)", "24.5-32.6m/s (storm)", "Over 32.6m/s (hurricane)",
}

func getWindSpeed(windSpeed int) string {
	return getMappedValue(windSpeed, windSpeedMap)
}

var relativeHumidityMap = map[int]string{
	-4: "0-5 %", -3: "5-10 %", -2: "10-15 %", -1: "15-20 %",
	0: "20-25 %", 1: "25-30 %", 2: "30-35 %", 3: "35-40 %",
	4: "40-45 %", 5: "45-50 %", 6: "50-55 %", 7: "55-60 %",
	8: "60-65 %", 9: "65-70 %", 10: "70-75 %", 11: "75-80 %",
	12: "80-85 %", 13: "85-90 %", 14: "90-95 %", 15: "95-99 %",
	16: "100 %",
}

func getRelativeHumidity(relHum int) string {
	if val, ok := relativeHumidityMap[relHum]; ok {
		return val
	}
	return "undefined"

}

func printWeatherData(data *ApiResponse) {
	fmt.Println("Printing Weather Forecast for Location:")
	fmt.Println("Initial timestamp:", data.Init)
	fmt.Println("Product:", data.Product)
	fmt.Println("Timepoint", data.DataSeries[0].Timepoint, "hours")
	// fmt.Println("Cloud cover", getCloudCover(data.DataSeries[0].CloudCover))
	fmt.Println("Cloud Cover:", data.DataSeries[0].CloudCoverString())
	// fmt.Println("Lifted Index", getLiftedIndex(data.DataSeries[0].LiftedIndex))
	fmt.Println("Lifted Index", data.DataSeries[0].LiftedIndexString())
	// fmt.Println("Temperature 2 meters", data.DataSeries[0].Temp2m, "°C")
	fmt.Println("Temperature 2 meters", data.DataSeries[0].TemperatureString())
	// fmt.Println("Seeing range", getSeeing(data.DataSeries[0].Seeing))
	fmt.Println("Seeing range:", data.DataSeries[0].SeeingString())
	// fmt.Println("Transparency range", getTransparency(data.DataSeries[0].Transparency))
	fmt.Println("Transparency range:", data.DataSeries[0].TransparencyString())
	// fmt.Println("Relative humidity 2 meters", getRelativeHumidity(data.DataSeries[0].RH2m))
	fmt.Println("Relative humidity 2 meters", data.DataSeries[0].HumidityString())
	fmt.Println("Precipitation type", data.DataSeries[0].PrecType)
	fmt.Println("Wind", data.DataSeries[0].WindSpeedString())

}

func printWeatherDataToCSV(data *ApiResponse, lat string, lng string) {
	// create a file
	file, err := os.Create("forecasts.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// initialize csv writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// write header
	rowData1 := []string{"Weather forecasts data for"}
	writer.Write(rowData1)
	rowData2 := []string{"latitude", lat}
	writer.Write(rowData2)
	rowData3 := []string{"longitude", lng}
	writer.Write(rowData3)
	rowData4 := []string{"Initial timestamp", data.Init}
	writer.Write(rowData4)
	rowData5 := []string{"Product", data.Product}
	writer.Write(rowData5)

	// write forecast series to file
	for i := 0; i < len(data.DataSeries); i++ {
		// rowData := []string{strconv.Itoa(data.DataSeries[i].Timepoint), getCloudCover(data.DataSeries[i].CloudCover)}
		rowData := []string{strconv.Itoa(data.DataSeries[i].Timepoint), data.DataSeries[i].CloudCoverString(), data.DataSeries[i].LiftedIndexString(), data.DataSeries[i].TemperatureString(), data.DataSeries[i].SeeingString(), data.DataSeries[i].TransparencyString(), data.DataSeries[i].HumidityString(), data.DataSeries[i].PrecType, data.DataSeries[i].WindSpeedString()}
		writer.Write(rowData)
	}

}

func main() {
	lat, lng := readUserInput()
	apiURL := buildAPIURL(lat, lng)
	data, err := fetchWeatherData(apiURL)
	if err != nil {
		log.Fatalf("Error fetching data: %v", err)
	}

	printWeatherData(data)

	printWeatherDataToCSV(data, lat, lng)
}
