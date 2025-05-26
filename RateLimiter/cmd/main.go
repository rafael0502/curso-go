package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/rafael0502/curso-go/RateLimiter/internal/limiter"
	"github.com/rafael0502/curso-go/RateLimiter/internal/middleware"
	"github.com/redis/go-redis/v9"
)

func main() {
	_ = godotenv.Load()

	redisClient := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	limiter := limiter.NewRedisLimiter(redisClient)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Estou funcionando!"))
	})

	handler := middleware.RateLimitMiddleware(limiter)(mux)

	log.Println("Servidor executando na porta 8081")
	if err := http.ListenAndServe(":8081", handler); err != nil {
		log.Fatalf("Falha ao iniciar o servidor: %v", err)
	}
}
