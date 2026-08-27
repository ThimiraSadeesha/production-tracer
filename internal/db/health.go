package db

import (
	"context"
	"database/sql"
	"fmt"
)

// Ping verifies the database connection is alive. Used by health checks.
func Ping(ctx context.Context) error {
	if client == nil {
		return fmt.Errorf("database not initialized")
	}
	sqlDB, err := client.DB()
	if err != nil {
		return err
	}
	return sqlDB.PingContext(ctx)
}

// Stats returns the underlying connection-pool statistics.
func Stats() (sql.DBStats, error) {
	if client == nil {
		return sql.DBStats{}, fmt.Errorf("database not initialized")
	}
	sqlDB, err := client.DB()
	if err != nil {
		return sql.DBStats{}, err
	}
	return sqlDB.Stats(), nil
}
