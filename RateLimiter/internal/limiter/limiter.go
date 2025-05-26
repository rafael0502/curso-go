package limiter

type RateLimiter interface {
	Permitir(key string, limit int, blockDuration int) (bool, error)
}
