package openmeteo

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/kenchi-huang/weather-aggregator/weather"
)

type Client struct{}

func (c Client) GetWeather(lat float64, lon float64) (*weather.Weather, error) {
	res, err := http.Get(fmt.Sprintf("https://api.open-meteo.com/v1/forecast?latitude=%f&longitude=%f&current=temperature_2m,relative_humidity_2m,precipitation", lat, lon))
	var data ForecastResponse
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	jsonFile, _ := io.ReadAll(res.Body)
	err = json.Unmarshal(jsonFile, &data)
	if err != nil {
		return nil, err
	}

	return &weather.Weather{
		Temperature:   float64(*data.Current.Temperature2m),
		City:          "",
		Provider:      "OpenMeteo",
		Humidity:      float64(*data.Current.RelativeHumidity2m),
		Precipitation: float64(*data.Current.Precipitation),
	}, nil
}
