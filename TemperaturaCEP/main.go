package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/temperatura", temperaturaHandler)
	log.Println("Servidor iniciado na porta: 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
