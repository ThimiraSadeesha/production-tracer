// Command migrate applies the database schema using GORM AutoMigrate over the
// models registered in app/database/schema (schema.Migrations), then drops any
// columns that were removed from the models.
//
// Usage:
//
//	migrate            apply the schema, then drop removed columns
//	migrate status     list the models that would be migrated
//
// AutoMigrate is additive (creates tables, adds missing columns/indexes) but
// never drops columns. To fully sync — including removing a column when you
// delete its field from a model — this command follows AutoMigrate with a
// drop pass (see dropRemovedColumns).
package main

import (
	"fmt"
	"os"

	"github.com/rashintha/logger"
	"github.com/thimira/production-tracer/app/database/schema"
	"github.com/thimira/production-tracer/internal/config"
	"github.com/thimira/production-tracer/internal/db"
	"gorm.io/gorm"
)

func main() {
	cmd := "up"
	if args := os.Args[1:]; len(args) > 0 {
		cmd = args[0]
	}

	if len(schema.Migrations) == 0 {
		logger.Warningln("schema.Migrations is empty — register model structs in app/database/schema/migrator.go")
		return
	}

	switch cmd {
	case "status":
		logger.Defaultf("%d model(s) registered for migration:", len(schema.Migrations))
		for _, m := range schema.Migrations {
			logger.Infof("  - %T", m)
		}
		return

	case "up", "":
		dbc, err := db.NewGormClient(
			config.DB_HOST, config.DB_PORT, config.DB_USER, config.DB_PASS, config.DB_NAME,
		)
		if err != nil {
			logger.ErrorFatalf("database connection failed: %v", err)
		}
		defer func() {
			if sqlDB, err := dbc.DB.DB(); err == nil {
				_ = sqlDB.Close()
			}
		}()

		logger.Defaultf("running AutoMigrate for %d model(s)...", len(schema.Migrations))
		if err := dbc.DB.AutoMigrate(schema.Migrations...); err != nil {
			logger.ErrorFatalf("auto-migrate failed: %v", err)
		}
		logger.Defaultln("✅ AutoMigrate complete")

		// AutoMigrate never drops columns; remove ones no longer on the models.
		if err := dropRemovedColumns(dbc.DB, schema.Migrations); err != nil {
			logger.ErrorFatalf("column sync failed: %v", err)
		}
		logger.Defaultln("✅ migration complete — schema is up to date")

	default:
		logger.Defaultln("usage: migrate [up | status]")
		os.Exit(2)
	}
}

// dropRemovedColumns drops any column that exists in a table but is no longer
// defined on its model. This makes deleting a field from a schema struct and
// running migrate also remove the column from the database.
func dropRemovedColumns(gdb *gorm.DB, models []interface{}) error {
	migrator := gdb.Migrator()

	for _, model := range models {
		if !migrator.HasTable(model) {
			continue
		}

		// Columns the model still defines.
		stmt := &gorm.Statement{DB: gdb}
		if err := stmt.Parse(model); err != nil {
			return fmt.Errorf("parse %T: %w", model, err)
		}
		table := stmt.Schema.Table
		desired := make(map[string]struct{}, len(stmt.Schema.DBNames))
		for _, name := range stmt.Schema.DBNames {
			desired[name] = struct{}{}
		}

		// Columns the table currently has.
		cols, err := migrator.ColumnTypes(model)
		if err != nil {
			return fmt.Errorf("read columns of %s: %w", table, err)
		}
		for _, col := range cols {
			name := col.Name()
			if _, keep := desired[name]; keep {
				continue
			}
			logger.Warningf("dropping removed column %s.%s", table, name)
			if err := migrator.DropColumn(model, name); err != nil {
				logger.Errorf("could not drop %s.%s: %v", table, name, err)
			}
		}
	}
	return nil
}
