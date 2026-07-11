package location

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Location struct {
	Country string
	Region  string
	City    string
	Lat     float64
	Lon     float64
}

type DataFromIp struct {
	Country string  `json:"country"`
	Region  string  `json:"region"`
	City    string  `json:"city"`
	Lat     float64 `json:"lat"`
	Lon     float64 `json:"lon"`
}

func GetCoordinatesFromIp(ip string) (*DataFromIp, error) {
	res, err := http.Get(fmt.Sprintf("http://ip-api.com/json/%s", ip))
	var data DataFromIp
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

type DataFromCoordinates struct {
	Country string  `json:"countryName"`
	Region  string  `json:"principalSubdivision"`
	City    string  `json:"city"`
	Lat     float64 `json:"latitude"`
	Lon     float64 `json:"longitude"`
}

func GetLocationFromCoordinates(lat float64, lon float64) (*DataFromCoordinates, error) {
	res, err := http.Get(fmt.Sprintf("https://api.bigdatacloud.net/data/reverse-geocode-client?latitude=%f&longitude=%f&localityLanguage=en", lat, lon))
	var data DataFromCoordinates
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
