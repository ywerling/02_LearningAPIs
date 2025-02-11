package weather

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

type ApiResponse struct {
	Product    string       `json:"product"`
	Init       string       `json:"init"`
	DataSeries []DataSeries `json:"dataseries"`
}

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

func (ds DataSeries) CloudCoverString() string   { return getCloudCover(ds.CloudCover) }
func (ds DataSeries) SeeingString() string       { return getSeeing(ds.Seeing) }
func (ds DataSeries) TransparencyString() string { return getTransparency(ds.Transparency) }
func (ds DataSeries) LiftedIndexString() string  { return getLiftedIndex(ds.LiftedIndex) }
func (ds DataSeries) WindSpeedString() string {
	return ds.Wind10m.Direction + " " + getWindSpeed(ds.Wind10m.Speed)
}
func (ds DataSeries) HumidityString() string    { return getRelativeHumidity(ds.RH2m) }
func (ds DataSeries) TemperatureString() string { return strconv.Itoa(ds.Temp2m) + " C" }

type Wind10m struct {
	Direction string `json:"direction"`
	Speed     int    `json:"speed"`
}

// FetchWeatherData makes an HTTP request and parses the response into ApiResponse
func FetchWeatherData(apiURL string) (*ApiResponse, error) {
	client := &http.Client{}

	request, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("creating GET request: %w", err)
	}
	request.Header.Set("Content-Type", "application/json")

	response, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("making request: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

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

// Functions that returned strings for the various parameters received via API
// Using the mapping definition provided in the API description

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
