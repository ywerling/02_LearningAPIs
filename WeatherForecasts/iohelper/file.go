package iohelper

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"

	"github.com/ywerling/02_LearningAPIs/weather"
)

// WriteToCSV writes weather forecast data to a CSV file
func WriteToCSV(data *weather.ApiResponse, lat, lng string) {
	file, err := os.Create("forecasts.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	headers := [][]string{
		{"Weather forecasts data for"},
		{"Latitude", lat},
		{"Longitude", lng},
		{"Initial Timestamp", data.Init},
		{"Product", data.Product},
		{"Timepoint", "Temperature 2m", "Cloud Cover", "Lifted Index", "Seeing", "Transparency", "Humidity", "Precipitation", "Wind 10m"},
	}

	for _, row := range headers {
		writer.Write(row)
	}

	for _, entry := range data.DataSeries {
		writer.Write([]string{
			strconv.Itoa(entry.Timepoint),
			entry.TemperatureString(),
			entry.CloudCoverString(),
			entry.LiftedIndexString(),
			entry.SeeingString(),
			entry.TransparencyString(),
			entry.HumidityString(),
			entry.PrecType,
			entry.WindSpeedString(),
		})
	}
}
