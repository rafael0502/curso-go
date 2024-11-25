package main

import (
	"log"
	"net/http"
)

func recoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Recuperou o panico: %v\n", r)
				//debug.PrintStack()
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Daeeee povo!"))
	})

	mux.HandleFunc("/panic", func(w http.ResponseWriter, r *http.Request) {
		panic("deu ruim")
	})

	log.Println("Escutando na porta", ":3000")

	if err := http.ListenAndServe(":3000", recoverMiddleware(mux)); err != nil {
		log.Fatalf("Não pode escutar em %s: %v\n", ":3000", err)
	}

}
