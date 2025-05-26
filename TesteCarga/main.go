package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"
)

type result struct {
	status int
	err    error
}

func worker(url string, requests <-chan int, results chan<- result, wg *sync.WaitGroup) {
	defer wg.Done()
	for range requests {
		resp, err := http.Get(url)
		if err != nil {
			results <- result{status: 0, err: err}
			continue
		}
		results <- result{status: resp.StatusCode}
		resp.Body.Close()
	}
}

func main() {
	// CLI flags
	url := flag.String("url", "", "URL do serviço a ser testado (ex: http://localhost:8080)")
	totalRequests := flag.Int("requests", 100, "Número total de requests a serem feitos")
	concurrency := flag.Int("concurrency", 10, "Número de chamadas simultâneas")
	flag.Parse()

	// Validação de entrada
	if *url == "" || *totalRequests <= 0 || *concurrency <= 0 {
		fmt.Println("Uso: ./testecarga --url=http://exemplo.com --requests=100 --concurrency=10")
		flag.PrintDefaults()
		os.Exit(1)
	}

	startTime := time.Now()

	// Canais e WaitGroup
	requests := make(chan int, *totalRequests)
	results := make(chan result, *totalRequests)
	var wg sync.WaitGroup

	// Inicia workers
	for i := 0; i < *concurrency; i++ {
		wg.Add(1)
		go worker(*url, requests, results, &wg)
	}

	// Alimenta canal de requisições
	for i := 0; i < *totalRequests; i++ {
		requests <- i
	}
	close(requests)

	// Aguarda workers terminarem
	wg.Wait()
	close(results)

	// Processa resultados
	duration := time.Since(startTime)
	statusCount := make(map[int]int)
	var totalOK int
	var totalErr int

	for res := range results {
		if res.err != nil {
			totalErr++
			continue
		}
		if res.status == 200 {
			totalOK++
		}
		statusCount[res.status]++
	}

	// Geração do relatório
	fmt.Println("=== Relatório de Teste de Carga ===")
	fmt.Printf("URL: %s\n", *url)
	fmt.Printf("Total de requests: %d\n", *totalRequests)
	fmt.Printf("Concorrência: %d\n", *concurrency)
	fmt.Printf("Tempo total: %v\n", duration)
	fmt.Printf("Requests bem sucedidos (200): %d\n", totalOK)
	fmt.Printf("Erros de conexão: %d\n", totalErr)
	fmt.Println("Distribuição de códigos de status:")
	for code, count := range statusCount {
		fmt.Printf("  %d: %d\n", code, count)
	}
}
