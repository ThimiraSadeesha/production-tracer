package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rashintha/logger"
	"github.com/thimira/production-tracer/api"
	"github.com/thimira/production-tracer/internal/config"
	"github.com/thimira/production-tracer/internal/db"
)

func init() {
	loc, err := time.LoadLocation("Asia/Colombo")
	if err != nil {
		logger.Errorf("❌ Failed to load location Asia/Colombo: %v", err)
	} else {
		time.Local = loc
		logger.Defaultln("🕒 Global timezone set to Asia/Colombo")
	}
}

func main() {
	logger.Defaultf("🔌 Connecting to database %s ", config.DB_NAME)
	dbc, err := db.NewGormClient(
		config.DB_HOST, config.DB_PORT, config.DB_USER, config.DB_PASS, config.DB_NAME,
	)
	if err != nil {
		logger.ErrorFatalf("❌ Failed to connect to database: %v", err)
	}
	logger.Defaultf("✅ Database connected: %s@%s:%s", config.DB_NAME, config.DB_HOST, config.DB_PORT)

	addr := "0.0.0.0:3000"
	go func() {
		logger.Defaultf("🚀 Starting %s API on %s", config.APP_NAME, addr)
		api.Run(addr)
	}()

	gracefulShutdown(dbc)
}

func gracefulShutdown(dbc *db.GormClient) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigChan
	logger.Defaultf("🧹 Caught signal: %s — shutting down gracefully...", sig.String())

	ctx, cancel := context.WithTimeout(context.Background(), config.SHUTDOWN_TIMEOUT)
	defer cancel()

	if dbc != nil && dbc.DB != nil {
		if sqlDB, err := dbc.DB.DB(); err == nil {
			done := make(chan struct{})
			go func() {
				if err := sqlDB.Close(); err != nil {
					logger.Errorf("❌ Failed to close database connections: %v", err)
				} else {
					logger.Defaultln("🔒 Database connections closed cleanly")
				}
				close(done)
			}()
			select {
			case <-done:
			case <-ctx.Done():
				logger.Errorln("⏱️ Timed out waiting for database connections to close")
			}
		}
	}

	logger.Defaultln("👋 Server stopped gracefully.")
	os.Exit(0)
}
