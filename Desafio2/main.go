package main

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type ViaCEP struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

type BrasilAPI struct {
	Cep          string `json:"cep"`
	State        string `json:"state"`
	City         string `json:"city"`
	Neighborhood string `json:"neighborhood"`
	Street       string `json:"street"`
	Service      string `json:"service"`
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/{cep}", GetCepHandler)
	http.ListenAndServe(":8000", r)
}

func GetCepHandler(w http.ResponseWriter, r *http.Request) {
	cep := chi.URLParam(r, "cep")

	if cep == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	c1 := make(chan *ViaCEP)
	c2 := make(chan *BrasilAPI)

	go func() {
		c, err := BuscaViaCEP(cep)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		c1 <- c
	}()

	go func() {
		c, err := BuscaBrasilAPI(cep)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		c2 <- c
	}()

	select {
	case c := <-c1: //ViaCEP
		println("Endereço obtido pela API Via CEP\n CEP: " + c.Cep + "\n Estado: " + c.Uf + "\n Cidade: " + c.Localidade + "\n Bairro: " + c.Bairro + "\n Rua: " + c.Logradouro)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(c)
	case c := <-c2: //Brasil API
		println("Endereço obtido pela API Brasil API\n CEP: " + c.Cep + "\n Estado: " + c.State + "\n Cidade: " + c.City + "\n Bairro: " + c.Neighborhood + "\n Rua: " + c.Street)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(c)
	case <-time.After(time.Second * 1):
		println("timeout")
		w.WriteHeader(http.StatusRequestTimeout)
	}
}

func BuscaViaCEP(cep string) (*ViaCEP, error) {
	resp, err := http.Get("https://viacep.com.br/ws/" + cep + "/json")

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	var c ViaCEP

	err = json.Unmarshal(body, &c)

	if err != nil {
		return nil, err
	}

	return &c, nil
}

func BuscaBrasilAPI(cep string) (*BrasilAPI, error) {
	resp, err := http.Get("https://brasilapi.com.br/api/cep/v1/" + cep)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	var c BrasilAPI

	err = json.Unmarshal(body, &c)

	if err != nil {
		return nil, err
	}

	return &c, nil
}
