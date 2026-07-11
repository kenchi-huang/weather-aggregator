package aggregator

import (
	"fmt"
	"strings"

	"github.com/kenchi-huang/weather-aggregator/client/location"
	"github.com/kenchi-huang/weather-aggregator/weather"
)

func BuildAggregatedWeather(conditions []weather.Weather, locationData location.Location) weather.Weather {
	var totalTemp float64
	var totalHumidity float64
	var totalPrecipMM float64
	var totalMinTemp float64
	var totalMaxTemp float64
	var maxConditionCode int

	var aggregatedDaily []weather.ForecastDay
	for d := range 7 {
		var totalDailyMinTemp float64
		var totalDailyMaxTemp float64
		var maxDailyConditionCode = 0
		timeBuckets := make(map[string][]float64)
		conditionBuckets := make(map[string]int)
		var providersForDay int

		for _, condition := range conditions {
			if d >= len(condition.Daily) {
				continue
			}
			providersForDay++

			daily := condition.Daily[d]
			totalDailyMinTemp += daily.MinTemp
			totalDailyMaxTemp += daily.MaxTemp
			maxDailyConditionCode = max(maxDailyConditionCode, daily.ConditionCode)

			for _, hour := range daily.Hourly {
				timeBuckets[hour.Time] = append(timeBuckets[hour.Time], hour.Temperature)
				if hour.ConditionCode > conditionBuckets[hour.Time] {
					conditionBuckets[hour.Time] = hour.ConditionCode
				}
			}
		}

		var aggregatedHours []weather.ForecastHour
		for h := range 24 {
			formattedTimeString := fmt.Sprintf("%02d:00", h)
			tmps, exists := timeBuckets[formattedTimeString]

			if !exists {
				continue
			}

			var totalHourlyTemp float64
			for _, data := range tmps {
				totalHourlyTemp += data
			}

			hourlyConditionCode := conditionBuckets[formattedTimeString]

			aggregatedHours = append(aggregatedHours, weather.ForecastHour{
				Time:          formattedTimeString,
				Temperature:   totalHourlyTemp / float64(len(tmps)),
				ConditionCode: hourlyConditionCode,
			})
		}

		aggregatedDaily = append(aggregatedDaily, weather.ForecastDay{
			Date:          conditions[0].Daily[d].Date,
			MinTemp:       totalDailyMinTemp / float64(providersForDay),
			MaxTemp:       totalDailyMaxTemp / float64(providersForDay),
			ConditionCode: maxDailyConditionCode,
			Hourly:        aggregatedHours,
		})
	}

	var names []string

	for _, condition := range conditions {
		totalTemp += condition.Current.Temperature
		totalHumidity += condition.Current.Humidity
		totalPrecipMM += condition.Current.Precipitation
		totalMinTemp += condition.Current.MinTemp
		totalMaxTemp += condition.Current.MaxTemp
		maxConditionCode = max(condition.Current.ConditionCode, maxConditionCode)
		names = append(names, condition.Provider)
	}

	providerString := "Aggregated from " + strings.Join(names, ", ")

	return weather.Weather{
		Current: weather.CurrentWeather{
			Temperature:   totalTemp / float64(len(conditions)),
			Humidity:      totalHumidity / float64(len(conditions)),
			Precipitation: totalPrecipMM / float64(len(conditions)),
			MinTemp:       totalMinTemp / float64(len(conditions)),
			MaxTemp:       totalMaxTemp / float64(len(conditions)),
			ConditionCode: maxConditionCode,
		},
		Daily:    aggregatedDaily,
		Provider: providerString,
		City:     locationData.City,
	}
}
