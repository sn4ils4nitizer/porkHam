package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type WeatherResponse struct {
	Location struct {
		Name    string `json:"name"`
		Region  string `json:"region"`
		Country string `json:"country"`
	} `json:"location"`
	Current struct {
		TempC     float64 `json:"temp_c"`
		Condition struct {
			Text string `json:"text"`
			Code int    `json:"code"`
		} `json:"condition"`
		Pressure      float64 `json:"pressure_mb"`
		WindSpeed     float64 `json:"wind_kph"`
		WindDirection string  `json:"wind_dir"`
		Humidity      int     `json:"humidity"`
		Precipitation float64 `json:"precip_mm"`
		Gust          float64 `json:"gust_kph"`
	} `json:"current"`
}

// TODO func read api
func readAPIKey(filename string) (string, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(data)), nil
}

func GetWeather(location string) (WeatherResponse, error) {
	//url :=
	apiKey, err := readAPIKey("apikey.txt")
	if err != nil {
		log.Fatalf("Error reading API key: %v", err)
	}

	url := fmt.Sprintf("http://api.weatherapi.com/v1/current.json?key=%s&q=%s", apiKey, location)

	resp, err := http.Get(url)
	if err != nil {
		return WeatherResponse{}, fmt.Errorf("Error fetching weather: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return WeatherResponse{}, fmt.Errorf("Error reading response body: %v", err)
	}

	var weather WeatherResponse
	err = json.Unmarshal(body, &weather)
	if err != nil {
		return WeatherResponse{}, fmt.Errorf("Error parsing JSON: %v", err)
	}
	fmt.Println("Weather request successful.")
	return weather, nil
}

func WeatherHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	location := vars["location"]

	weather, err := GetWeather(location)
	if err != nil {
		http.Error(w, "Failed to fetch weather", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(weather)
}
