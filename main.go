package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"github.com/kenchi-huang/weather-aggregator/openmeteo"
	"github.com/kenchi-huang/weather-aggregator/weather"
	"github.com/kenchi-huang/weather-aggregator/wttr"
)

func weatherHandler(w http.ResponseWriter, r *http.Request) {
	latStr := r.URL.Query().Get("lat")
	latFloat, _ := strconv.ParseFloat(latStr, 64)
	lonStr := r.URL.Query().Get("lon")
	lonFloat, _ := strconv.ParseFloat(lonStr, 64)

	fmt.Println("Latitude: " + latStr + ", Longitude: " + lonStr)

	var locationData location
	if latStr == "" && lonStr == "" {
		ipAddress := r.RemoteAddr
		if os.Getenv("DEV_MODE") == "true" {
			ipAddress = os.Getenv("DEV_IP")
		} else {
			ipAddress = strings.Split(ipAddress, ":")[0]
		}
		locData, err := getCoordinatesFromIp(ipAddress)
		if err != nil {
			log.Println(err)
			return
		}
		locationData = location{
			locData.Country,
			locData.Region,
			locData.City,
			locData.Lat,
			locData.Lon,
		}
	} else {
		locData, err := getLocationFromCoordinates(latFloat, lonFloat)
		if err != nil {
			log.Println(err)
			return
		}
		locationData = location{
			locData.Country,
			locData.Region,
			locData.City,
			locData.Lat,
			locData.Lon,
		}
	}

	providers := []weather.Provider{
		openmeteo.Client{},
		wttr.Client{},
	}

	var conditions []weather.Weather

	for _, provider := range providers {
		currWeather, err := provider.GetWeather(locationData.Lat, locationData.Lon)
		if err != nil {
			println(err.Error())
			continue
		}
		conditions = append(conditions, *currWeather)
	}

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

	var aggregatedWeather = weather.Weather{
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
	tmpl, _ := template.ParseFiles("index.html")
	err := tmpl.Execute(w, aggregatedWeather)
	if err != nil {
		return
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, assuming in production mode")
	}
	http.HandleFunc("/", weatherHandler)
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		return
	}
}

type location struct {
	Country string
	Region  string
	City    string
	Lat     float64
	Lon     float64
}

type locationDataFromIp struct {
	Country string  `json:"country"`
	Region  string  `json:"region"`
	City    string  `json:"city"`
	Lat     float64 `json:"lat"`
	Lon     float64 `json:"lon"`
}

func getCoordinatesFromIp(ip string) (*locationDataFromIp, error) {
	res, err := http.Get(fmt.Sprintf("http://ip-api.com/json/%s", ip))
	var data locationDataFromIp
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	jsonFile, _ := io.ReadAll(res.Body)
	err = json.Unmarshal(jsonFile, &data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

type locationDataFromCoordinates struct {
	Country string  `json:"countryName"`
	Region  string  `json:"principalSubdivision"`
	City    string  `json:"city"`
	Lat     float64 `json:"latitude"`
	Lon     float64 `json:"longitude"`
}

func getLocationFromCoordinates(lat float64, lon float64) (*locationDataFromCoordinates, error) {
	res, err := http.Get(fmt.Sprintf("https://api.bigdatacloud.net/data/reverse-geocode-client?latitude=%f&longitude=%f&localityLanguage=en", lat, lon))
	var data locationDataFromCoordinates
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	jsonFile, _ := io.ReadAll(res.Body)
	err = json.Unmarshal(jsonFile, &data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}
