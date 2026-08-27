package config

import (
	"os"
	"time"

	"github.com/rashintha/logger"
	"github.com/thimira/production-tracer/internal/env"
)

// getEnv safely gets a variable from env.CONF first, then os.Getenv as fallback
func getEnv(key string) string {
	val := ""
	if v, ok := env.CONF[key]; ok && v != "" {
		val = v
	} else {
		val = os.Getenv(key)
	}
	if val == "" {
		logger.Warningln("⚠️ Missing environment variable:" + key)
	}
	return val
}

var (
	APP_ENV        = getEnv("NODE_ENV")
	DB_HOST        = getEnv("DB_HOST")
	DB_PORT        = getEnv("DB_PORT")
	DB_USER        = getEnv("DB_USER")
	DB_PASS        = getEnv("DB_PASS")
	DB_NAME        = getEnv("DB_NAME")
	JWT_SECRET     = getEnv("JWT_SECRET")
	JWT_EXPIRATION = getEnv("JWT_EXPIRATION")
	AES_SECRET_KEY = getEnv("AES_SECRET_KEY")
	AES_IV         = getEnv("AES_IV")
)

func GetJWTExpiration() time.Duration {
	expStr := getEnv("JWT_EXPIRATION")
	if expStr == "" {
		expStr = "24h"
	}
	d, err := time.ParseDuration(expStr)
	if err != nil {
		return 24 * time.Hour
	}
	return d
}
