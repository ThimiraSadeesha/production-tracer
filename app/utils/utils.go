package utils

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type PaginationResult struct {
	NextCursor *int64 `json:"nextCursor,omitempty"`
	HasMore    bool   `json:"hasMore"`
}

func ProcessData[T any](data []T, index int) T {
	if len(data) == 0 {
		var zero T
		return zero
	}
	if index < 0 || index >= len(data) {
		return data[0]
	}
	return data[index]
}

func PaginationResponse(data []map[string]interface{}, limit int) PaginationResult {
	var nextCursor *int64
	if len(data) > 0 {
		lastRow := data[len(data)-1]
		var lastID int64
		if id, ok := lastRow["id"].(int64); ok {
			lastID = id
		} else if id, ok := lastRow["id"].(float64); ok {
			lastID = int64(id)
		} else {
			lastID = 0
		}
		if lastID != 0 {
			nextCursor = &lastID
		}
	}
	hasMore := len(data) == limit
	return PaginationResult{
		NextCursor: nextCursor,
		HasMore:    hasMore,
	}
}

func EmptyToNil(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}

// IsAbsentQueryParam reports whether a query value should be treated as omitted
// (empty, or the literal strings "null" / "undefined" from clients).
func IsAbsentQueryParam(s string) bool {
	s = strings.TrimSpace(s)
	return s == "" || strings.EqualFold(s, "null") || strings.EqualFold(s, "undefined")
}

// OptionalQueryString returns nil when the query param is absent or a sentinel like "null".
func OptionalQueryString(s string) *string {
	if IsAbsentQueryParam(s) {
		return nil
	}
	trimmed := strings.TrimSpace(s)
	return &trimmed
}

func ProcessJsonData[T any](row map[string]interface{}, keys ...string) {
	if row == nil {
		return
	}

	for _, key := range keys {
		val, exists := row[key]
		if !exists {
			var zero T
			row[key] = zero
			continue
		}
		switch v := val.(type) {
		case string:
			if v != "" {
				var parsed T
				if err := json.Unmarshal([]byte(v), &parsed); err == nil {
					row[key] = parsed
				} else {
					var zero T
					row[key] = zero
				}
			}
		case []byte:
			if len(v) == 0 {
				var zero T
				row[key] = zero
				break
			}
			var parsed T
			if err := json.Unmarshal(v, &parsed); err == nil {
				row[key] = parsed
			} else {
				var zero T
				row[key] = zero
			}
		case json.RawMessage:
			if len(v) == 0 {
				var zero T
				row[key] = zero
				break
			}
			var parsed T
			if err := json.Unmarshal(v, &parsed); err == nil {
				row[key] = parsed
			} else {
				var zero T
				row[key] = zero
			}
		case []interface{}:
			allMaps := len(v) > 0
			for _, item := range v {
				if _, ok := item.(map[string]interface{}); !ok {
					allMaps = false
					break
				}
			}
			if allMaps {
				row[key] = v
				break
			}
			parsedList := make([]interface{}, 0, len(v))
			for _, item := range v {
				if str, ok := item.(string); ok {
					var obj map[string]interface{}
					if err := json.Unmarshal([]byte(str), &obj); err == nil {
						parsedList = append(parsedList, obj)
					}
				} else if m, ok := item.(map[string]interface{}); ok {
					parsedList = append(parsedList, m)
				}
			}
			row[key] = parsedList

		default:
			var zero T
			row[key] = zero
		}
	}
}

// ProcedureIntFromMap reads an int from a MySQL procedure row; keys are often lowercased (e.g. operatorid vs operatorId).
func ProcedureIntFromMap(m map[string]interface{}, keyCandidates ...string) (int, error) {
	if m == nil {
		return 0, fmt.Errorf("nil result map")
	}
	for _, want := range keyCandidates {
		for k, v := range m {
			if strings.EqualFold(k, want) {
				return ToInt(v)
			}
		}
	}
	for _, want := range keyCandidates {
		wantNorm := strings.ToLower(strings.ReplaceAll(want, "_", ""))
		for k, v := range m {
			if strings.ToLower(strings.ReplaceAll(k, "_", "")) == wantNorm {
				return ToInt(v)
			}
		}
	}
	return 0, fmt.Errorf("procedure result missing keys %v (have: %v)", keyCandidates, mapKeys(m))
}

func mapKeys(m map[string]interface{}) []string {
	ks := make([]string, 0, len(m))
	for k := range m {
		ks = append(ks, k)
	}
	return ks
}

func ToInt(val any) (int, error) {
	switch v := val.(type) {
	case float64:
		return int(v), nil
	case int64:
		return int(v), nil
	case uint64:
		return int(v), nil
	case []uint8:
		n, err := strconv.Atoi(string(v))
		if err != nil {
			return 0, fmt.Errorf("failed to convert to int: %w", err)
		}
		return n, nil
	case nil:
		return 0, fmt.Errorf("value is nil")
	default:
		return 0, fmt.Errorf("unexpected type %T for value %v", val, val)
	}
}

func ToInt64(val any) (int64, error) {
	switch v := val.(type) {
	case float64:
		return int64(v), nil
	case int64:
		return v, nil
	case uint64:
		return int64(v), nil
	case []uint8:
		n, err := strconv.ParseInt(string(v), 10, 64)
		if err != nil {
			return 0, fmt.Errorf("failed to convert to int64: %w", err)
		}
		return n, nil
	case nil:
		return 0, fmt.Errorf("value is nil")
	default:
		return 0, fmt.Errorf("unexpected type %T for value %v", val, val)
	}
}
func GetPatchValue(updates map[string]interface{}, existing map[string]interface{}, updateKey string, existingKeys ...string) interface{} {
	if val, ok := updates[updateKey]; ok {
		return val
	}
	for _, key := range existingKeys {
		if val, ok := existing[key]; ok {
			return val
		}
	}
	val := existing[updateKey]
	if (updateKey == "updatedBy" || updateKey == "createdBy") && (val == nil || val == "") {
		return "system"
	}
	return val
}

func ParseNestedDateOperations(row map[string]interface{}, key string) {
	if row == nil {
		return
	}
	val, exists := row[key]
	if !exists {
		row[key] = map[string]interface{}{}
		return
	}

	switch v := val.(type) {
	case map[string]interface{}:
		out := make(map[string]interface{}, len(v))
		for date, inner := range v {
			out[date] = parseOperationSliceValue(inner)
		}
		row[key] = out
	case string:
		if v == "" {
			row[key] = map[string]interface{}{}
		} else {
			row[key] = decodeDateKeyedOperationsJSON([]byte(v))
		}
	case []byte:
		if len(v) == 0 {
			row[key] = map[string]interface{}{}
		} else {
			row[key] = decodeDateKeyedOperationsJSON(v)
		}
	default:
		row[key] = map[string]interface{}{}
	}
}

// decodeDateKeyedOperationsJSON parses {"YYYY-MM-DD": <operations>, ...}.
// Each value may be a JSON array (MySQL JSON_OBJECTAGG of JSON_ARRAYAGG) or a JSON string
// holding a serialized array (legacy double-encoded shape).
func decodeDateKeyedOperationsJSON(raw []byte) map[string]interface{} {
	var outer map[string]json.RawMessage
	if err := json.Unmarshal(raw, &outer); err != nil {
		return map[string]interface{}{}
	}
	result := make(map[string]interface{}, len(outer))
	for date, msg := range outer {
		result[date] = parseOperationSliceFromRawMessage(msg)
	}
	return result
}

func parseOperationSliceFromRawMessage(msg json.RawMessage) []interface{} {
	if len(msg) == 0 {
		return []interface{}{}
	}
	var direct []interface{}
	if err := json.Unmarshal(msg, &direct); err == nil {
		return direct
	}
	var innerJSON string
	if err := json.Unmarshal(msg, &innerJSON); err != nil {
		return []interface{}{}
	}
	if err := json.Unmarshal([]byte(innerJSON), &direct); err != nil {
		return []interface{}{}
	}
	return direct
}

func parseOperationSliceValue(inner interface{}) []interface{} {
	inner = unwrapJSONValue(inner)
	switch x := inner.(type) {
	case []interface{}:
		return x
	case string:
		var ops []interface{}
		if err := json.Unmarshal([]byte(x), &ops); err == nil {
			return ops
		}
	case json.RawMessage:
		return parseOperationSliceFromRawMessage(x)
	}
	return []interface{}{}
}

func ParseOperationProcessTimeline(row map[string]interface{}) {
	if row == nil {
		return
	}
	val, exists := row["timeline"]
	if !exists {
		row["timeline"] = map[string]interface{}{}
		return
	}
	row["timeline"] = parseOperationProcessTimelineValue(val)
}

func parseOperationProcessTimelineValue(val interface{}) map[string]interface{} {
	timeline := asStringKeyedMap(unwrapJSONValue(val))
	if timeline == nil {
		return map[string]interface{}{}
	}

	out := make(map[string]interface{}, len(timeline))
	for date, machinesVal := range timeline {
		machines := asStringKeyedMap(unwrapJSONValue(machinesVal))
		if machines == nil {
			out[date] = map[string]interface{}{}
			continue
		}

		machinesOut := make(map[string]interface{}, len(machines))
		for code, machineVal := range machines {
			machine := asStringKeyedMap(unwrapJSONValue(machineVal))
			if machine == nil {
				continue
			}
			if ops, ok := machine["operations"]; ok {
				machine["operations"] = normalizeTimelineOperations(parseOperationSliceValue(ops))
			} else {
				machine["operations"] = []interface{}{}
			}
			machinesOut[code] = machine
		}
		out[date] = machinesOut
	}
	return out
}

// normalizeTimelineOperations unwraps JSON-encoded fields on each timeline operation.
func normalizeTimelineOperations(ops []interface{}) []interface{} {
	for i, op := range ops {
		opMap, ok := op.(map[string]interface{})
		if !ok {
			continue
		}
		if pp, exists := opMap["pp"]; exists {
			opMap["pp"] = unwrapJSONValue(pp)
		}
		ops[i] = opMap
	}
	return ops
}

func unwrapJSONValue(val interface{}) interface{} {
	for {
		switch v := val.(type) {
		case string:
			if strings.TrimSpace(v) == "" {
				return val
			}
			var next interface{}
			if err := json.Unmarshal([]byte(v), &next); err != nil {
				return val
			}
			val = next
		case []byte:
			if len(v) == 0 {
				return val
			}
			var next interface{}
			if err := json.Unmarshal(v, &next); err != nil {
				return val
			}
			val = next
		case json.RawMessage:
			if len(v) == 0 {
				return val
			}
			var next interface{}
			if err := json.Unmarshal(v, &next); err != nil {
				return val
			}
			val = next
		default:
			return val
		}
	}
}

func asStringKeyedMap(val interface{}) map[string]interface{} {
	switch v := val.(type) {
	case map[string]interface{}:
		return v
	case map[interface{}]interface{}:
		out := make(map[string]interface{}, len(v))
		for k, item := range v {
			out[fmt.Sprint(k)] = item
		}
		return out
	default:
		return nil
	}
}

func FormatDate(input string) string {
	if input == "" {
		return ""
	}

	// First, try parsing ISO 8601 / RFC3339 format
	if t, err := time.Parse(time.RFC3339, input); err == nil {
		return t.Format("2006-01-02") // YYYY-MM-DD
	}

	// Then, try your other custom layouts
	layouts := []string{
		"02.01.2006", // 15.12.2025
		"02/01/2006", // 15/12/2025
		"2006-01-02", // 2025-12-15
		"02-01-2006", // 15-12-2025
		"2006.01.02", // 2025.12.15
	}

	for _, layout := range layouts {
		if t, err := time.Parse(layout, input); err == nil {
			return t.Format("2006-01-02") // output YYYY-MM-DD
		}
	}

	return ""
}
