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

	for _, condition := range conditions {
		totalTemp += condition.Temperature
		totalHumidity += condition.Humidity
		totalPrecipMM += condition.Precipitation
	}

	var aggregatedWeather = weather.Weather{
		Temperature:   totalTemp / float64(len(conditions)),
		Humidity:      totalHumidity / float64(len(conditions)),
		Precipitation: totalPrecipMM / float64(len(conditions)),
		Provider:      "Aggregated",
		City:          locationData.City,
	}
	tmpl, _ := template.ParseFiles("index.html")
	tmpl.Execute(w, aggregatedWeather)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, assuming in production mode")
	}
	http.HandleFunc("/", weatherHandler)
	http.ListenAndServe(":8080", nil)
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
	} else {
		defer res.Body.Close()
		jsonFile, _ := io.ReadAll(res.Body)
		json.Unmarshal(jsonFile, &data)
		return &data, nil
	}
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
	} else {
		defer res.Body.Close()
		jsonFile, _ := io.ReadAll(res.Body)
		json.Unmarshal(jsonFile, &data)
		return &data, nil
	}
}
