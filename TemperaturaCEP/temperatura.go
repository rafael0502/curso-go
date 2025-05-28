package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type TemperaturaAPIResponse struct {
	Current struct {
		TempC float64 `json:"temp_c"`
	} `json:"current"`
}

func getTemperaturaByCity(city string) (float64, error) {
	escapedCity := url.QueryEscape(city)
	apiURL := fmt.Sprintf("http://api.weatherapi.com/v1/current.json?key=%s&q=%s", "f66eea1616bb484fad3183814252705", escapedCity)

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
