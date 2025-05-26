package limiter

import (
	"fmt"
	"time"

	"context"

	"github.com/redis/go-redis/v9"
)

type RedisLimiter struct {
	Client *redis.Client
}

func NewRedisLimiter(client *redis.Client) *RedisLimiter {
	return &RedisLimiter{Client: client}
}

func (r *RedisLimiter) Permitir(key string, limit int, blockDuration int) (bool, error) {
	ctx := context.Background()

	blockedKey := fmt.Sprintf("blocked:%s", key)
	if blocked, _ := r.Client.Get(ctx, blockedKey).Result(); blocked == "1" {
		return false, nil
	}

	countKey := fmt.Sprintf("rate:%s", key)
	count, _ := r.Client.Incr(ctx, countKey).Result()

	if count == 1 {
		r.Client.Expire(ctx, countKey, time.Second)
	}

	if count > int64(limit) {
		r.Client.Set(ctx, blockedKey, "1", time.Duration(blockDuration)*time.Second)
		return false, nil
	}

	return true, nil
}
