package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type ViaCEPResponse struct {
	Localidade string `json:"localidade"`
}

func getCityFromCEP(cep string) (string, error) {
	resp, err := http.Get(fmt.Sprintf("http://viacep.com.br/ws/%s/json/", cep))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var data ViaCEPResponse
	if err := json.Unmarshal(body, &data); err != nil || data.Localidade == "" {
		return "", errors.New("cidade nao encontrada")
	}
	return data.Localidade, nil
}
