package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type Cotacao struct {
	Bid string `json:"bid"`
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)

	if err != nil {
		log.Fatalf("Erro ao criar requisição: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Fatalf("Erro ao fazer requisição: %v", err)
	}

	defer resp.Body.Close()

	var quote Cotacao
	if err := json.NewDecoder(resp.Body).Decode(&quote); err != nil {
		log.Fatalf("Erro ao decodificar resposta: %v", err)
	}

	f, err := os.Create("cotacao.txt")

	if err != nil {
		log.Fatalf("Erro ao criar arquivo: %v\n", err)
	}

	defer f.Close()
	_, err = f.WriteString("Dólar: " + quote.Bid)

	fmt.Println("Dólar: " + quote.Bid)
}
