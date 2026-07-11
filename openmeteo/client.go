package openmeteo

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/kenchi-huang/weather-aggregator/weather"
)

type Client struct{}

const HOURS_PER_DAY_IN_FORECAST = 24

func (c Client) GetWeather(lat float64, lon float64) (*weather.Weather, error) {
	u, _ := url.Parse("https://api.open-meteo.com/v1/forecast")
	q := u.Query()

	q.Add("latitude", fmt.Sprintf("%f", lat))
	q.Add("longitude", fmt.Sprintf("%f", lon))
	q.Add("current", "temperature_2m,relative_humidity_2m,precipitation,weather_code")
	q.Add("daily", "weather_code,temperature_2m_max,temperature_2m_min")
	q.Add("hourly", "temperature_2m,weather_code")

	u.RawQuery = q.Encode()
	finalUrl := u.String()

	res, err := http.Get(finalUrl)
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

	var dailyWeather []weather.ForecastDay

	for i := 0; i < len(*data.Daily.Time); i++ {
		var hourlyWeather []weather.ForecastHour

		for j := i * HOURS_PER_DAY_IN_FORECAST; j < ((i * HOURS_PER_DAY_IN_FORECAST) + HOURS_PER_DAY_IN_FORECAST); j++ {
			rawTime := (*data.Hourly.Time)[j]
			cleanTime := rawTime[len(rawTime)-5:]
			hourlyWeather = append(hourlyWeather, weather.ForecastHour{
				Time:          cleanTime,
				Temperature:   float64((*data.Hourly.Temperature2m)[j]),
				ConditionCode: int((*data.Hourly.WeatherCode)[j]),
			})
		}

		dailyWeather = append(dailyWeather, weather.ForecastDay{
			Date:          (*data.Daily.Time)[i],
			MinTemp:       float64((*data.Daily.Temperature2mMin)[i]),
			MaxTemp:       float64((*data.Daily.Temperature2mMax)[i]),
			ConditionCode: int((*data.Daily.WeatherCode)[i]),
			Hourly:        hourlyWeather,
		})
	}

	return &weather.Weather{
		Current: weather.CurrentWeather{
			Precipitation: float64(*data.Current.Precipitation),
			Humidity:      float64(*data.Current.RelativeHumidity2m),
			Temperature:   float64(*data.Current.Temperature2m),
			ConditionCode: int(*data.Current.WeatherCode),
			MinTemp:       dailyWeather[0].MinTemp,
			MaxTemp:       dailyWeather[0].MaxTemp,
		},
		Daily:    dailyWeather,
		Provider: "OpenMeteo",
	}, nil
}
