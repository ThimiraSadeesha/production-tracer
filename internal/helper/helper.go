package helper

import (
	"bytes"
	"encoding/json"
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

// ResolveActor prefers a value from the JSON body, then X-User, then "system".
func ResolveActor(c *gin.Context, fromBody string) string {
	if s := strings.TrimSpace(fromBody); s != "" {
		return s
	}
	return Actor(c)
}

// ActorFromMap reads createdBy/updatedBy from a JSON object.
func ActorFromMap(m map[string]interface{}) string {
	if m == nil {
		return ""
	}
	for _, k := range []string{"updatedBy", "createdBy", "updated_by", "created_by"} {
		if s, ok := m[k].(string); ok {
			if s = strings.TrimSpace(s); s != "" {
				return s
			}
		}
	}
	return ""
}

// ActorFromJSON reads createdBy/updatedBy from a JSON array or object.
func ActorFromJSON(raw json.RawMessage) string {
	trimmed := bytes.TrimSpace(raw)
	if len(trimmed) == 0 {
		return ""
	}
	if trimmed[0] == '[' {
		var arr []map[string]interface{}
		if err := json.Unmarshal(trimmed, &arr); err != nil {
			return ""
		}
		for _, item := range arr {
			if s := ActorFromMap(item); s != "" {
				return s
			}
		}
		return ""
	}
	var obj map[string]interface{}
	if err := json.Unmarshal(trimmed, &obj); err != nil {
		return ""
	}
	if s := ActorFromMap(obj); s != "" {
		return s
	}
	if data, ok := obj["data"]; ok {
		b, err := json.Marshal(data)
		if err == nil {
			return ActorFromJSON(b)
		}
	}
	return ""
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
