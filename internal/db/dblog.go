package db

import (
	"io"
	"log"
	"os"
	"time"

	gormlogger "gorm.io/gorm/logger"
)

// Route GORM's logger to BOTH the console and log/log.log, matching the app
// logger (github.com/rashintha/logger), which also writes to both. GORM's
// default logger writes only to stdout; wrapping stdout + the log file in an
// io.MultiWriter mirrors database logs (slow queries, errors, migration output)
// into the shared log file as well.
//
// NewGormClient opens with &gorm.Config{} (nil Logger), so GORM falls back to
// gormlogger.Default at connect time — reassigning it here, before any
// connection is opened, takes effect.
var _ = func() bool {
	var w io.Writer = os.Stdout
	_ = os.MkdirAll("log", 0o755)
	if f, err := os.OpenFile("log/log.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0o644); err == nil {
		w = io.MultiWriter(os.Stdout, f)
	}
	gormlogger.Default = gormlogger.New(
		log.New(w, "", log.LstdFlags),
		gormlogger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  gormlogger.Warn,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)
	return true
}()
