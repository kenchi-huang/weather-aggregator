package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/kenchi-huang/weather-aggregator/handler"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, assuming in production mode")
	}
	http.HandleFunc("/", handler.WeatherHandler)
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		return
	}
}
