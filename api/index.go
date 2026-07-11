package api

import (
	"net/http"

	"github.com/kenchi-huang/weather-aggregator/handler"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	handler.WeatherHandler(w, r)
}
