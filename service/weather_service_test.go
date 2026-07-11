package service

import (
	"testing"

	"github.com/joho/godotenv"
)

// This is an INTEGRATION TEST.
// It will actually hit the real Open-Meteo, Wttr, and Upstash APIs!
func TestGetWeatherForUser_Integration(t *testing.T) {
	godotenv.Load("../.env")
	// Arrange: Use coordinates for London
	lat := "51.5072"
	lon := "-0.1276"
	remoteAddr := "127.0.0.1"

	// Act: Run the full service pipeline
	weatherData, err := GetWeatherForUser(lat, lon, remoteAddr)

	// Assert: Ensure no errors occurred during the API calls
	if err != nil {
		t.Fatalf("Expected no error, but got: %v", err)
	}

	// Assert: Ensure we got valid data back
	if weatherData == nil {
		t.Fatalf("Expected weather data, but got nil")
	}

	// Assert: Ensure the city resolved correctly
	if weatherData.City == "" {
		t.Errorf("Expected City to be populated, but it was empty")
	}

	// Assert: Ensure the mathematical aggregation ran
	if weatherData.Current.Temperature == 0 && weatherData.Current.Humidity == 0 {
		t.Errorf("Weather data appears empty, aggregation may have failed")
	}
}

func TestGetWeatherForUserWithLocalCaching_Integration(t *testing.T) {
	godotenv.Load("../.env")
	// Arrange: Use coordinates for Tokyo
	lat := "35.6762"
	lon := "139.6503"
	remoteAddr := "127.0.0.1"

	// Act: Run the local cache pipeline
	weatherData, err := GetWeatherForUserWithLocalCaching(lat, lon, remoteAddr)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, but got: %v", err)
	}

	if weatherData == nil {
		t.Fatalf("Expected weather data, but got nil")
	}

	if weatherData.City == "" {
		t.Errorf("Expected City to be populated, but it was empty")
	}
}
