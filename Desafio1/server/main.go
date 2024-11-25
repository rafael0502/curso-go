package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type USDBRL struct {
	Code       string `json:"code"`
	Codein     string `json:"codein"`
	Name       string `json:"name"`
	High       string `json:"high"`
	Low        string `json:"low"`
	VarBid     string `json:"varBid"`
	PctChange  string `json:"pctChange"`
	Bid        string `json:"bid"`
	Ask        string `json:"ask"`
	Timestamp  string `json:"timestamp"`
	CreateDate string `json:"create_date"`
}

type Cotacao struct {
	Cot USDBRL `json:"USDBRL"`
}

func main() {
	http.HandleFunc("/cotacao", CambioDolarRealHandler)
	http.ListenAndServe(":8080", nil)
}

func CambioDolarRealHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 200*time.Millisecond)
	defer cancel()

	log.Println("Request iniciado")

	resp, err := http.Get("https://economia.awesomeapi.com.br/json/last/USD-BRL")
	if err != nil {
		log.Printf("Erro ao buscar cotação: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer resp.Body.Close()
	defer log.Println("Request finalizado")

	select {
	case <-time.After(200 * time.Millisecond):
		log.Println("Timeout ao acessar endpoint.")
		w.WriteHeader(http.StatusRequestTimeout)
		return
	case <-ctx.Done():
		log.Println("Request processado com sucesso!")
		body, _ := io.ReadAll(resp.Body)
		var cambio Cotacao
		err = json.Unmarshal(body, &cambio)

		if err != nil {
			log.Printf("Erro ao extrair informações da cotação: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Millisecond)
		defer cancel()

		err := GravarCotacao(cambio.Cot)

		if err != nil {
			log.Fatalf("%s\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		select {
		case <-time.After(10 * time.Millisecond):
			log.Println("Tempo excedido ao persistir os dados no BD.")
			w.WriteHeader(http.StatusRequestTimeout)
			return
		case <-ctx.Done():
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			data := map[string]string{"bid": cambio.Cot.Bid}

			json.NewEncoder(w).Encode(data)
		}
	}
}

func GravarCotacao(cot USDBRL) error {
	db, err := sql.Open("sqlite3", "./desafio1.db")

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	// Cria uma tabela
	sqlStmt := `CREATE TABLE IF NOT EXISTS cotacao ( 
	            id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, 
				code VARCHAR(3), 
				codein VARCHAR(3), 
				name VARCHAR(100), 
				high VARCHAR(20), 
				low VARCHAR(20), 
				varBid VARCHAR(20), 
				pctChange VARCHAR(10), 
				bid VARCHAR(20), 
				ask VARCHAR(20), 
				timestamp VARCHAR(20), 
				create_date VARCHAR(20) 
				);`
	_, err = db.Exec(sqlStmt)

	if err != nil {
		log.Fatalf("%q: %s\n", err, sqlStmt)
		return err
	}

	stmt, err := db.Prepare("INSERT INTO cotacao(code, codein, name, high, low, varBid, pctChange, bid, ask, timestamp, create_date) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")

	if err != nil {
		return err
	}

	defer stmt.Close()
	_, err = stmt.Exec(cot.Code, cot.Codein, cot.Name, cot.High, cot.Low, cot.VarBid, cot.PctChange, cot.Bid, cot.Ask, cot.Timestamp, cot.CreateDate)

	if err != nil {
		return err
	}

	return nil
}
