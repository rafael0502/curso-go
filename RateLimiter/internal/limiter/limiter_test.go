package limiter_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/rafael0502/curso-go/RateLimiter/internal/limiter"
	"github.com/rafael0502/curso-go/RateLimiter/internal/middleware"
	"github.com/redis/go-redis/v9"
)

func setupTestLimiter() *limiter.RedisLimiter {
	os.Setenv("RATE_LIMIT", "2")
	os.Setenv("BLOCK_DURATION", "10") // segundos
	os.Setenv("TOKEN_RATE_LIMIT_testtoken", "3")
	os.Setenv("TOKEN_BLOCK_DURATION_testtoken", "10")

	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // certifique-se de que o Redis esteja rodando
		Password: "",
		DB:       1, // use outro DB para testes
	})

	// Limpa antes de cada teste
	redisClient.FlushDB(context.Background())

	return limiter.NewRedisLimiter(redisClient)
}

func performRequest(t *testing.T, handler http.Handler, token string, ip string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = ip + ":1234"
	if token != "" {
		req.Header.Set("API_KEY", token)
	}
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	return w
}

func TestRateLimitByIP(t *testing.T) {
	lim := setupTestLimiter()

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})
	handler := middleware.RateLimitMiddleware(lim)(mux)

	ip := "192.168.0.15"
	// 1ª requisição - deve passar
	if resp := performRequest(t, handler, "", ip); resp.Code != 200 {
		t.Errorf("esperado 200, recebido %d", resp.Code)
	}

	// 2ª requisição - deve passar
	if resp := performRequest(t, handler, "", ip); resp.Code != 200 {
		t.Errorf("esperado 200, recebido %d", resp.Code)
	}

	// 3ª requisição - deve ser bloqueada
	if resp := performRequest(t, handler, "", ip); resp.Code != 429 {
		t.Errorf("esperado 429, recebido %d", resp.Code)
	}
}

func TestRateLimitByToken(t *testing.T) {
	lim := setupTestLimiter()

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})
	handler := middleware.RateLimitMiddleware(lim)(mux)

	token := "testtoken"

	// Dentro do limite de 3
	for i := 0; i < 3; i++ {
		if resp := performRequest(t, handler, token, "1.1.1.1"); resp.Code != 200 {
			t.Errorf("esperado 200, recebido %d", resp.Code)
		}
	}

	// Estourou o limite
	if resp := performRequest(t, handler, token, "1.1.1.1"); resp.Code != 429 {
		t.Errorf("esperado 429, recebido %d", resp.Code)
	}
}

func TestTokenOverridesIPLimit(t *testing.T) {
	lim := setupTestLimiter()

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})
	handler := middleware.RateLimitMiddleware(lim)(mux)

	ip := "10.10.10.10"
	token := "testtoken" // tem limite maior (3)

	// Mesmo IP, mas usando token → deve respeitar limite do token
	for i := 0; i < 3; i++ {
		if resp := performRequest(t, handler, token, ip); resp.Code != 200 {
			t.Errorf("esperado 200, recebido %d", resp.Code)
		}
	}

	// Estourou o limite do token
	if resp := performRequest(t, handler, token, ip); resp.Code != 429 {
		t.Errorf("esperado 429, recebido %d", resp.Code)
	}
}

func TestUnblockAfterDuration(t *testing.T) {
	lim := setupTestLimiter()

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})
	handler := middleware.RateLimitMiddleware(lim)(mux)

	ip := "2.2.2.2"

	// Estoura o limite
	performRequest(t, handler, "", ip)
	performRequest(t, handler, "", ip)
	performRequest(t, handler, "", ip) // 429

	// Espera mais que o tempo de bloqueio
	time.Sleep(11 * time.Second)

	// Deve ser aceito novamente
	if resp := performRequest(t, handler, "", ip); resp.Code != 200 {
		t.Errorf("esperado 200 após bloqueio expirar, recebido %d", resp.Code)
	}
}
