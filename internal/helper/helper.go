package helper

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// Actor returns the acting user from the X-User header, defaulting to "system".
func Actor(c *gin.Context) string {
	if a := c.GetHeader("X-User"); a != "" {
		return a
	}
	return "system"
}

// Nullable returns nil for an empty string so the SP treats it as NULL / no filter.
func Nullable(s string) any {
	if strings.TrimSpace(s) == "" {
		return nil
	}
	return s
}

// PtrOrNil dereferences a *string for CallProcedure (which cannot take pointers),
// yielding nil (SQL NULL) when the field was omitted.
func PtrOrNil(p *string) any {
	if p == nil {
		return nil
	}
	return *p
}

// PtrInt64OrNil is PtrOrNil for *int64.
func PtrInt64OrNil(p *int64) any {
	if p == nil {
		return nil
	}
	return *p
}

// ParseInt64 parses s, returning def on failure.
func ParseInt64(s string, def int64) int64 {
	if v, err := strconv.ParseInt(strings.TrimSpace(s), 10, 64); err == nil {
		return v
	}
	return def
}

// ParseInt parses a positive int from s, returning def on failure.
func ParseInt(s string, def int) int {
	if v, err := strconv.Atoi(strings.TrimSpace(s)); err == nil && v > 0 {
		return v
	}
	return def
}
