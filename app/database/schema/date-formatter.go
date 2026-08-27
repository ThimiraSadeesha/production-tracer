package schema

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

type LocalDateTime time.Time

const localDateTimeMySQL = "2006-01-02 15:04:05"

func (t LocalDateTime) Time() time.Time { return time.Time(t) }

func (t LocalDateTime) GormDataType() string { return "datetime" }

func (t LocalDateTime) Value() (driver.Value, error) {
	tm := time.Time(t)
	if tm.IsZero() {
		return nil, nil
	}
	return tm, nil
}

func (t *LocalDateTime) Scan(src interface{}) error {
	if src == nil {
		*t = LocalDateTime(time.Time{})
		return nil
	}
	switch v := src.(type) {
	case time.Time:
		*t = LocalDateTime(v)
		return nil
	case []byte:
		return t.Scan(string(v))
	case string:
		if v == "" {
			*t = LocalDateTime(time.Time{})
			return nil
		}
		tm, err := time.ParseInLocation(localDateTimeMySQL, v, time.Local)
		if err != nil {
			return err
		}
		*t = LocalDateTime(tm)
		return nil
	default:
		return fmt.Errorf("cannot scan %T into LocalDateTime", src)
	}
}

func (t LocalDateTime) MarshalJSON() ([]byte, error) {
	tm := time.Time(t)
	if tm.IsZero() {
		return []byte("null"), nil
	}
	return json.Marshal(tm.In(time.Local).Format(localDateTimeMySQL))
}

func (t *LocalDateTime) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		*t = LocalDateTime(time.Time{})
		return nil
	}

	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	if s == "" {
		*t = LocalDateTime(time.Time{})
		return nil
	}

	// ── 1. Formats that carry timezone info (parse as-is, then normalize to local) ──
	zoneFormats := []string{
		time.RFC3339Nano,            // 2006-01-02T15:04:05.999999999Z07:00
		time.RFC3339,                // 2006-01-02T15:04:05Z07:00
		"2006-01-02T15:04:05-07:00", // explicit offset, no Z alias
		"2006-01-02 15:04:05-07:00",
		"2006-01-02T15:04:05 MST",
		"2006-01-02 15:04:05 MST",
		"Mon, 02 Jan 2006 15:04:05 MST",   // RFC1123
		"Mon, 02 Jan 2006 15:04:05 -0700", // RFC1123Z
		time.UnixDate,                     // Mon Jan _2 15:04:05 MST 2006
		time.RubyDate,                     // Mon Jan 02 15:04:05 -0700 2006
		time.ANSIC,                        // Mon Jan _2 15:04:05 2006
	}

	for _, layout := range zoneFormats {
		if tm, err := time.Parse(layout, s); err == nil {
			*t = LocalDateTime(tm.In(time.Local)) // normalize to local
			return nil
		}
	}

	// ── 2. Naive formats (no timezone — assume local) ──
	naiveFormats := []string{
		"2006-01-02 15:04:05.999999999", // with nanoseconds
		"2006-01-02 15:04:05.999",       // with milliseconds
		localDateTimeMySQL,              // 2006-01-02 15:04:05
		"2006-01-02 15:04",              // no seconds
		"2006-01-02 15",                 // hour only
		"2006-01-02T15:04:05.999999999", // T-sep with nanoseconds
		"2006-01-02T15:04:05.999",       // T-sep with milliseconds
		"2006-01-02T15:04:05",           // T-sep full
		"2006-01-02T15:04",              // datetime-local input (HTML)
		"2006-01-02T15",                 // T-sep hour only
		"2006-01-02",                    // date only → 00:00:00
		"02/01/2006 15:04:05",           // DD/MM/YYYY with time
		"02/01/2006 15:04",
		"02/01/2006",          // DD/MM/YYYY
		"01/02/2006 15:04:05", // MM/DD/YYYY with time
		"01/02/2006 15:04",
		"01/02/2006",          // MM/DD/YYYY
		"2006/01/02 15:04:05", // slash-separated YMD
		"2006/01/02 15:04",
		"2006/01/02",
		"Jan 2, 2006 15:04:05", // human readable
		"Jan 2, 2006 15:04",
		"Jan 2, 2006",
		"2 Jan 2006 15:04:05",
		"2 Jan 2006",
		"January 2, 2006 15:04:05",
		"January 2, 2006",
	}

	for _, layout := range naiveFormats {
		if tm, err := time.ParseInLocation(layout, s, time.Local); err == nil {
			*t = LocalDateTime(tm)
			return nil
		}
	}

	return fmt.Errorf("LocalDateTime: unsupported format %q", s)
}
