package handler

import (
	_ "embed"
	"html/template"
	"net/http"
	"os"

	"github.com/kenchi-huang/weather-aggregator/service"
	"github.com/kenchi-huang/weather-aggregator/weather"
)

//go:embed index.html
var htmlTemplate string

func WeatherHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/favicon.ico" {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	latStr, lonStr := r.URL.Query().Get("lat"), r.URL.Query().Get("lon")
	remoteAddr := r.RemoteAddr
	tmpl, _ := template.New("index").Parse(htmlTemplate)

	var aggregatedWeather *weather.Weather
	var err error
	if os.Getenv("RUN_WITH_LOCAL_CACHING") == "true" {
		aggregatedWeather, err = service.GetWeatherForUserWithLocalCaching(latStr, lonStr, remoteAddr)
	} else {
		aggregatedWeather, err = service.GetWeatherForUser(latStr, lonStr, remoteAddr)
	}

	if err != nil {
		return
	}

	tmplErr := tmpl.Execute(w, aggregatedWeather)
	if tmplErr != nil {
		return
	}
}
