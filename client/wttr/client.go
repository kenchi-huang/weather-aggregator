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
	tempFloat, _ := strconv.ParseFloat(*condition.TempC, 64)
	humFloat, _ := strconv.ParseFloat(*condition.Humidity, 64)
	precipFloat, _ := strconv.ParseFloat(*condition.PrecipMM, 64)

	var dailyWeather []weather.ForecastDay

	for _, daily := range *data.Weather {
		var hourlyWeather []weather.ForecastHour

		for _, hourly := range *daily.Hourly {
			temp, _ := strconv.ParseFloat(*hourly.TempC, 64)
			conditionCode, _ := strconv.Atoi(*hourly.WeatherCode)
			rawTimeInt, _ := strconv.Atoi(*hourly.Time)
			cleanTime := fmt.Sprintf("%02d:00", rawTimeInt/100)

			hourlyWeather = append(hourlyWeather, weather.ForecastHour{
				Time:          cleanTime,
				Temperature:   temp,
				ConditionCode: translateToWMO(conditionCode),
			})
		}

		minTemp, _ := strconv.ParseFloat(*daily.MintempC, 64)
		maxTemp, _ := strconv.ParseFloat(*daily.MaxtempC, 64)
		conditionCode, _ := strconv.Atoi(*(*daily.Hourly)[4].WeatherCode)

		dailyWeather = append(dailyWeather, weather.ForecastDay{
			Date:          *daily.Date,
			MinTemp:       minTemp,
			MaxTemp:       maxTemp,
			ConditionCode: translateToWMO(conditionCode),
			Hourly:        hourlyWeather,
		})
	}

	return &weather.Weather{
		Current: weather.CurrentWeather{
			Temperature:   tempFloat,
			Humidity:      humFloat,
			Precipitation: precipFloat,
			MinTemp:       dailyWeather[0].MinTemp,
			MaxTemp:       dailyWeather[0].MaxTemp,
			ConditionCode: translateToWMO(dailyWeather[0].ConditionCode),
		},
		Daily:    dailyWeather,
		Provider: "wttr",
	}, nil
}

func translateToWMO(wttrCode int) int {
	switch wttrCode {
	case 113:
		return 0 // Clear/Sunny
	case 116:
		return 2 // Partly cloudy
	case 119, 122:
		return 3 // Cloudy/Overcast
	case 143, 248, 260:
		return 45 // Fog
	case 176, 263, 266, 281, 284:
		return 51 // Drizzle
	case 293, 296, 299, 302, 305, 308:
		return 61 // Rain
	case 353, 356, 359:
		return 80 // Rain showers
	case 227, 230, 320, 323, 326, 329, 332, 335, 338, 362, 365, 368, 371:
		return 71 // Snow
	case 200, 386, 389, 392, 395:
		return 95 // Thunderstorm
	default:
		return 0 // Default to Sunny if unknown
	}
}
