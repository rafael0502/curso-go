package middleware

import (
	"net/http"
	"strings"

	"github.com/rafael0502/curso-go/RateLimiter/internal/config"
	"github.com/rafael0502/curso-go/RateLimiter/internal/limiter"
)

func RateLimitMiddleware(limiter limiter.RateLimiter) func(http.Handler) http.Handler {
	tokenLimits := config.GetTokenLimits()
	rateLimit := config.GetInt("RATE_LIMIT", 5)
	blockDuration := config.GetInt("BLOCK_DURATION", 300)

	//log.Printf("rateLimit: %d, blockDuration: %d, tokenLimits: %v", rateLimit, blockDuration, tokenLimits)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := r.RemoteAddr
			limit := rateLimit
			block := blockDuration

			token := r.Header.Get("API_KEY")
			//log.Printf("token: %s", token)

			if token != "" {
				key = "token:" + token
				if tokenLimit, ok := tokenLimits[token]; ok {
					limit = tokenLimit
					block = config.GetInt("TOKEN_BLOCK_DURATION_"+token, blockDuration)
				}
			} else {
				key = "ip:" + strings.Split(r.RemoteAddr, ":")[0]
			}

			permitido, err := limiter.Permitir(key, limit, block)
			if err != nil || !permitido {
				w.WriteHeader(http.StatusTooManyRequests)
				w.Write([]byte("Você alcançou o número máximo de requisições ou ações permitidas dentro de um determinado período de tempo"))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
