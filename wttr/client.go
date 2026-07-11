package wttr

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/kenchi-huang/weather-aggregator/weather"
)

type Client struct{}

func (c Client) GetWeather(lat float64, lon float64) (*weather.Weather, error) {
	res, err := http.Get(fmt.Sprintf("https://wttr.in/%f,%f?format=j1", lat, lon))
	var data WttrResponse
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	jsonFile, _ := io.ReadAll(res.Body)
	err = json.Unmarshal(jsonFile, &data)
	if err != nil {
		return nil, err
	}

	condition := (*data.CurrentCondition)[0]
	tempStr := *condition.TempC
	tempFloat, _ := strconv.ParseFloat(tempStr, 64)
	humStr := *condition.Humidity
	humFloat, _ := strconv.ParseFloat(humStr, 64)
	precipStr := *condition.PrecipMM
	precipFloat, _ := strconv.ParseFloat(precipStr, 64)

	return &weather.Weather{
		Temperature:   tempFloat,
		City:          "",
		Provider:      "wttr",
		Humidity:      humFloat,
		Precipitation: precipFloat,
	}, nil
}
