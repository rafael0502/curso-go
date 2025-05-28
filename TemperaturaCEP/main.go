package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	http.HandleFunc("/temperatura", temperaturaHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Servidor iniciado na porta:", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
