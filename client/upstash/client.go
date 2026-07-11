package upstash

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/kenchi-huang/weather-aggregator/weather"
)

type CacheResponse struct {
	Result string `json:"result"`
}

func ReadCache(cacheKey string) (*weather.Weather, error) {
	url := fmt.Sprintf("%s/get/%s", os.Getenv("UPSTASH_REDIS_REST_URL"), cacheKey)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("UPSTASH_REDIS_REST_TOKEN")))

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	var data weather.Weather
	var cacheResponse CacheResponse
	defer response.Body.Close()
	jsonFile, _ := io.ReadAll(response.Body)
	err = json.Unmarshal(jsonFile, &cacheResponse)
	if err != nil {
		return nil, err
	} else if cacheResponse.Result == "" {
		return nil, nil
	}
	err = json.Unmarshal([]byte(cacheResponse.Result), &data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func WriteCache(cacheKey string, weather *weather.Weather) error {
	url := fmt.Sprintf(
		"%s/set/%s?ex=%s",
		os.Getenv("UPSTASH_REDIS_REST_URL"),
		cacheKey,
		os.Getenv("CACHE_EXPIRY_SECONDS"),
	)

	jsonData, err := json.Marshal(weather)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("UPSTASH_REDIS_REST_TOKEN")))
	if err != nil {
		return err
	}

	_, err = http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	return nil
}
