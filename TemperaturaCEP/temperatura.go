package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

type TemperaturaAPIResponse struct {
	Current struct {
		TempC float64 `json:"temp_c"`
	} `json:"current"`
}

func getTemperaturaByCity(city string) (float64, error) {
	apiKey := os.Getenv("WEATHER_API_KEY")
	if apiKey == "" {
		return 0, fmt.Errorf("variável de ambiente WEATHER_API_KEY não está configurada")
	}

	escapedCity := url.QueryEscape(city)
	apiURL := fmt.Sprintf("https://api.weatherapi.com/v1/current.json?key=%s&q=%s", apiKey, escapedCity)

	resp, err := http.Get(apiURL)
	if err != nil {
		return 0, fmt.Errorf("falha no request da API weather: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		return 0, fmt.Errorf("weather API retornou o status %d: %s", resp.StatusCode, string(body))
	}

	var data TemperaturaAPIResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return 0, fmt.Errorf("falha no parse do JSON da API weather: %w", err)
	}

	return data.Current.TempC, nil
}
