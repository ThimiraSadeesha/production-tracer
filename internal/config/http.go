package config

import (
	"os"
	"time"

	"github.com/thimira/production-tracer/internal/env"
)

func optEnv(key, def string) string {
	if v, ok := env.CONF[key]; ok && v != "" {
		return v
	}
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func optDuration(key string, def time.Duration) time.Duration {
	if s := optEnv(key, ""); s != "" {
		if d, err := time.ParseDuration(s); err == nil {
			return d
		}
	}
	return def
}

// HTTP server settings (with production-safe defaults).
var (
	APP_NAME         = optEnv("APP_NAME", "production-tracer")
	PORT             = optEnv("PORT", "3000")
	SHUTDOWN_TIMEOUT = optDuration("SHUTDOWN_TIMEOUT", 20*time.Second)
)

func IsProduction() bool { return APP_ENV == "production" }
