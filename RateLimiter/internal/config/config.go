package config

import (
	"os"
	"strconv"
	"strings"
)

func GetEnv(key string, defaultVal string) string {
	if val, exists := os.LookupEnv(key); exists {
		return val
	}
	return defaultVal
}

func GetInt(key string, defaultVal int) int {
	valStr := GetEnv(key, "")
	if val, err := strconv.Atoi(valStr); err == nil {
		return val
	}
	return defaultVal
}

func GetTokenLimits() map[string]int {
	limits := map[string]int{}
	for _, env := range os.Environ() {
		if strings.HasPrefix(env, "TOKEN_RATE_LIMIT_") {
			parts := strings.SplitN(env, "=", 2)
			token := strings.TrimPrefix(parts[0], "TOKEN_RATE_LIMIT_")
			limit, _ := strconv.Atoi(parts[1])
			limits[token] = limit
		}
	}
	return limits
}
