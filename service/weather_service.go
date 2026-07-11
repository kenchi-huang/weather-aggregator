package service

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/kenchi-huang/weather-aggregator/client/location"
	"github.com/kenchi-huang/weather-aggregator/client/openmeteo"
	"github.com/kenchi-huang/weather-aggregator/client/upstash"
	"github.com/kenchi-huang/weather-aggregator/client/wttr"
	"github.com/kenchi-huang/weather-aggregator/service/aggregator"
	"github.com/kenchi-huang/weather-aggregator/weather"
)

func GetWeatherForUser(lat string, lon string, remoteAddr string) (*weather.Weather, error) {
	latFloat, _ := strconv.ParseFloat(lat, 64)
	lonFloat, _ := strconv.ParseFloat(lon, 64)

	var locationData location.Location
	if lat == "" && lon == "" {
		ipAddress := remoteAddr
		if os.Getenv("DEV_MODE") == "true" {
			ipAddress = os.Getenv("DEV_IP")
		} else {
			ipAddress = strings.Split(ipAddress, ":")[0]
		}
		locData, err := location.GetCoordinatesFromIp(ipAddress)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		locationData = location.Location{
			locData.Country,
			locData.Region,
			locData.City,
			locData.Lat,
			locData.Lon,
		}
	} else {
		locData, err := location.GetLocationFromCoordinates(latFloat, lonFloat)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		locationData = location.Location{
			locData.Country,
			locData.Region,
			locData.City,
			locData.Lat,
			locData.Lon,
		}
	}

	cacheKey := fmt.Sprintf("%.2f,%.2f", locationData.Lat, locationData.Lon)
	cacheData, err := upstash.ReadCache(cacheKey)
	if err != nil {
		log.Println(err)
		// continue and try and get the user the weather data
	}

	if cacheData != nil {
		fmt.Println("Serving from cache! 🚀")
		return cacheData, nil
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

	aggregatedWeather := aggregator.BuildAggregatedWeather(conditions, locationData)

	err = upstash.WriteCache(cacheKey, &aggregatedWeather)
	if err != nil {
		return &aggregatedWeather, err
	}

	return &aggregatedWeather, nil
}
