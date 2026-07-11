package aggregator

import (
	"testing"

	"github.com/kenchi-huang/weather-aggregator/client/location"
	"github.com/kenchi-huang/weather-aggregator/weather"
)

func TestBuildAggregatedWeather(t *testing.T) {
	// Arrange: Create dummy data from two different "providers"
	mockConditions := []weather.Weather{
		{
			Provider: "Provider A",
			Current: weather.CurrentWeather{
				Temperature:   20.0,
				MinTemp:       15.0,
				MaxTemp:       25.0,
				ConditionCode: 1, // Let's say 1 is Sunny
			},
			Daily: []weather.ForecastDay{
				{
					Date:          "2026-07-12",
					MinTemp:       10.0,
					MaxTemp:       20.0,
					ConditionCode: 1,
				},
			},
		},
		{
			Provider: "Provider B",
			Current: weather.CurrentWeather{
				Temperature:   30.0,
				MinTemp:       25.0,
				MaxTemp:       35.0,
				ConditionCode: 3, // Let's say 3 is Overcast
			},
			Daily: []weather.ForecastDay{
				{
					Date:          "2026-07-12",
					MinTemp:       20.0,
					MaxTemp:       30.0,
					ConditionCode: 3,
				},
			},
		},
	}

	mockLocation := location.Location{
		City: "TestCity",
	}

	// Act: Run our pure mathematical aggregator
	result := BuildAggregatedWeather(mockConditions, mockLocation)

	// Assert: Check the mathematical averages
	expectedTemp := 25.0
	if result.Current.Temperature != expectedTemp {
		t.Errorf("Expected current temperature to be %f, but got %f", expectedTemp, result.Current.Temperature)
	}

	if result.Current.MinTemp != 20.0 {
		t.Errorf("Expected min temperature to be 20.0, but got %f", result.Current.MinTemp)
	}

	if result.Current.MaxTemp != 30.0 {
		t.Errorf("Expected max temperature to be 30.0, but got %f", result.Current.MaxTemp)
	}

	// Assert: Check the max condition code (Should pick 3 over 1)
	if result.Current.ConditionCode != 3 {
		t.Errorf("Expected condition code to be 3 (max), but got %d", result.Current.ConditionCode)
	}

	// Assert: Check that the location data propagated correctly
	if result.City != "TestCity" {
		t.Errorf("Expected city to be TestCity, but got %s", result.City)
	}
}
