package db

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/go-sql-driver/mysql"
	mysqldriver "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	dbTimezone  = "Asia%2FColombo"
	dbTZSession = "%27Asia%2FColombo%27"
)

var client *gorm.DB

type GormClient struct {
	DB *gorm.DB
}

func NewGormClient(host, port, username, password, dbName string) (*GormClient, error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=%s&time_zone=%s&multiStatements=true",
		username, password, host, port, dbName, dbTimezone, dbTZSession,
	)
	db, err := gorm.Open(mysqldriver.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB from gorm: %w", err)
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	client = db
	return &GormClient{DB: db}, nil
}

func CallProcedure[T any](name string, params ...any) (T, error) {
	var result T
	if client == nil {
		return result, fmt.Errorf("database not initialized")
	}
	query, args := buildCall(name, params)
	if err := client.Raw(query, args...).Scan(&result).Error; err != nil {
		return result, mapMySQLError(name, err)
	}

	// Dynamic check for database-level error messages
	if err := checkDBError(result); err != nil {
		return result, err
	}

	return result, nil
}

func checkDBError(res any) error {
	switch v := res.(type) {
	case map[string]interface{}:
		return inspectMap(v)
	case []map[string]interface{}:
		if len(v) > 0 {
			return inspectMap(v[0])
		}
	case *map[string]interface{}:
		if v != nil {
			return inspectMap(*v)
		}
	}
	return nil
}

func inspectMap(m map[string]interface{}) error {
	if raw, ok := m["error_message"]; ok {
		if msg := strings.TrimSpace(fmt.Sprint(raw)); msg != "" && msg != "<nil>" {
			return errors.New(msg)
		}
	}
	if raw, ok := m["Error"]; ok {
		if msg := strings.TrimSpace(fmt.Sprint(raw)); msg != "" && msg != "<nil>" {
			return errors.New(msg)
		}
	}
	return nil
}

func CallProcedureWithTimeout[T any](name string, params []any, met int, nrt int) (T, error) {
	var result T
	if client == nil {
		return result, fmt.Errorf("database not initialized")
	}

	sqlDB, err := client.DB()
	if err != nil {
		return result, fmt.Errorf("failed to get sql.DB: %w", err)
	}

	conn, err := sqlDB.Conn(context.Background())
	if err != nil {
		return result, fmt.Errorf("failed to get connection: %w", err)
	}
	defer conn.Close()

	var version string
	if err := conn.QueryRowContext(context.Background(), "SELECT VERSION()").Scan(&version); err != nil {
		return result, fmt.Errorf("failed to detect DB version: %w", err)
	}
	isMariaDB := strings.Contains(strings.ToLower(version), "mariadb")

	if nrt > 0 {
		if _, err := conn.ExecContext(context.Background(), fmt.Sprintf("SET SESSION NET_READ_TIMEOUT=%d", nrt)); err != nil {
			return result, fmt.Errorf("failed to set NET_READ_TIMEOUT: %w", err)
		}
	}
	if met > 0 {
		var timeoutSQL string
		if isMariaDB {
			timeoutSQL = fmt.Sprintf("SET SESSION max_statement_time=%d", (met+999)/1000)
		} else {
			timeoutSQL = fmt.Sprintf("SET SESSION MAX_EXECUTION_TIME=%d", met)
		}
		if _, err := conn.ExecContext(context.Background(), timeoutSQL); err != nil {
			return result, fmt.Errorf("failed to set execution timeout: %w", err)
		}
	}

	callSQL, args := buildCall(name, params)
	rows, err := conn.QueryContext(context.Background(), callSQL, args...)
	if err != nil {
		return result, mapMySQLError(name, err)
	}
	defer rows.Close()

	if err := client.ScanRows(rows, &result); err != nil {
		return result, mapMySQLError(name, err)
	}

	return result, nil
}

func buildCall(name string, params []any) (string, []any) {
	if len(params) == 0 {
		return fmt.Sprintf("CALL %s()", name), nil
	}
	placeholders := strings.Repeat("?,", len(params))
	placeholders = placeholders[:len(placeholders)-1]
	return fmt.Sprintf("CALL %s(%s)", name, placeholders), params
}

func mapMySQLError(procName string, err error) error {
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) {
		msg := mysqlErr.Message
		switch mysqlErr.Number {
		case 1062:
			return fmt.Errorf("conflict: duplicate entry in %s: %s", procName, msg)
		case 1451, 1452:
			return fmt.Errorf("bad request: foreign key constraint fails in %s: %s", procName, msg)
		case 1064:
			return fmt.Errorf("bad request: SQL syntax error in %s: %s", procName, msg)
		case 1045:
			return fmt.Errorf("forbidden: access denied to DB in %s: %s", procName, msg)
		case 1205, 1213:
			return fmt.Errorf("service unavailable: deadlock/lock wait timeout in %s: %s", procName, msg)
		case 1054:
			return fmt.Errorf("bad request: unknown column in %s: %s", procName, msg)
		case 1146:
			return fmt.Errorf("not found: table not found in %s: %s", procName, msg)
		default:
			return fmt.Errorf("internal error: MySQL [%d] in %s: %s", mysqlErr.Number, procName, msg)
		}
	}
	return fmt.Errorf("internal error: DB call failed for %s: %w", procName, err)
}
