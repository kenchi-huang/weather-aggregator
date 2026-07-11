package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/kenchi-huang/weather-aggregator/service"
)

func weatherHandler(w http.ResponseWriter, r *http.Request) {
	latStr, lonStr := r.URL.Query().Get("lat"), r.URL.Query().Get("lon")
	remoteAddr := r.RemoteAddr
	tmpl, _ := template.ParseFiles("index.html")

	aggregatedWeather, err := service.GetWeatherForUser(latStr, lonStr, remoteAddr)
	if err != nil {
		return
	}

	tmplErr := tmpl.Execute(w, aggregatedWeather)
	if tmplErr != nil {
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
