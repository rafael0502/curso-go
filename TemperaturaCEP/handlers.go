package main

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"
)

type WeatherResponse struct {
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

func temperaturaHandler(w http.ResponseWriter, r *http.Request) {
	cep := r.URL.Query().Get("cep")

	if !isValidCEP(cep) {
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}

	location, err := getCityFromCEP(cep)
	if err != nil {
		http.Error(w, "can not find zipcode", http.StatusNotFound)
		return
	}

	//log.Printf("Local encontrado: %s", location)

	tempC, err := getTemperaturaByCity(location)
	if err != nil {
		http.Error(w, "Temperatura não disponível para o CEP informado.", http.StatusInternalServerError)
		log.Printf("Erro: %s", err)
		return
	}

	resp := WeatherResponse{
		TempC: tempC,
		TempF: tempC*1.8 + 32,
		TempK: tempC + 273,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func isValidCEP(cep string) bool {
	match, _ := regexp.MatchString(`^\d{8}$`, cep)
	return match
}
