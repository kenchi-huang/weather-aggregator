package service

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/kenchi-huang/weather-aggregator/client/location"
	"github.com/kenchi-huang/weather-aggregator/client/openmeteo"
	"github.com/kenchi-huang/weather-aggregator/client/wttr"
	"github.com/kenchi-huang/weather-aggregator/service/aggregator"
	"github.com/kenchi-huang/weather-aggregator/weather"
)

type CacheEntry struct {
	Weather *weather.Weather
	Expiry  time.Time
}

var cache = make(map[string]CacheEntry)
var cacheMutex sync.RWMutex

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
	cacheMutex.RLock()
	entry, exists := cache[cacheKey]
	cacheMutex.RUnlock()

	if exists && time.Now().Before(entry.Expiry) {
		fmt.Println("Serving from cache! 🚀")
		return entry.Weather, nil
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

	cacheMutex.Lock()
	cache[cacheKey] = CacheEntry{
		Weather: &aggregatedWeather,
		Expiry:  time.Now().Add(time.Hour * 1),
	}
	cacheMutex.Unlock()

	return &aggregatedWeather, nil
}
