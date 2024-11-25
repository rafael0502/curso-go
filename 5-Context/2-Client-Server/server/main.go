package main

import (
	"log"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log.Println("Request iniciado")
	defer log.Println("Request finalizado")
	select {
	case <-time.After(5 * time.Second):
		//imprime no command line stdout
		log.Println("Request processado com sucesso!")
		//imprime no browser
		w.Write([]byte("Request processado com sucesso"))
	case <-ctx.Done():
		//imprime no command line stdout
		log.Println("Request cancelada pelo cliente")
	}
}
