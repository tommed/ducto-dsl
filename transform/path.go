package transform

import (
	"fmt"
	"reflect"
	"strings"
)

// CoerceToMap attempts to cast or reflect a map[string]T into map[string]interface{}
func CoerceToMap(v interface{}) (map[string]interface{}, bool) {
	switch m := v.(type) {
	case map[string]interface{}:
		return m, true
	default:
		rv := reflect.ValueOf(v)
		if rv.Kind() != reflect.Map || rv.Type().Key().Kind() != reflect.String {
			return nil, false
		}

		coerced := make(map[string]interface{}, rv.Len())
		for _, key := range rv.MapKeys() {
			val := rv.MapIndex(key)
			coerced[key.String()] = val.Interface()
		}
		return coerced, true
	}
}

// CoerceToArray makes sure an array is always a []interface{}.
// This isn't a problem for JSON, but if the program is deserialized
// from YAML (as per ducto-orchestrator), it might do smart things like
// `[]map[string]interface{}` which wouldn't narrow cast as Go doesn't
// support this.
func CoerceToArray(val interface{}) ([]interface{}, bool) {
	if val == nil {
		return nil, false
	}

	// Use reflection first for general slices
	rv := reflect.ValueOf(val)
	if rv.Kind() == reflect.Slice {
		out := make([]interface{}, rv.Len())
		for i := 0; i < rv.Len(); i++ {
			out[i] = rv.Index(i).Interface()
		}
		return out, true
	}

	// Now check explicitly for known decoded types
	switch v := val.(type) {
	case []interface{}:
		return v, true
	case []map[string]interface{}:
		out := make([]interface{}, len(v))
		for i := range v {
			out[i] = v[i]
		}
		return out, true
	default:
		return nil, false
	}
}

func GetValueAtPath(data map[string]interface{}, path string) (interface{}, bool) {
	if path == "" {
		return data, true
	}

	parts := strings.Split(path, ".")
	var current interface{} = data

	for _, part := range parts {
		m, ok := CoerceToMap(current)
		if !ok {
			return nil, false
		}
		current, ok = m[part]
		if !ok {
			return nil, false
		}
	}

	return current, true
}

func SetValueAtPath(data map[string]interface{}, path string, value interface{}) error {
	if path == "" {
		return fmt.Errorf("cannot set empty path")
	}

	parts := strings.Split(path, ".")
	var current interface{} = data

	for i, part := range parts {
		isLast := i == len(parts)-1

		m, ok := CoerceToMap(current)
		if !ok {
			return fmt.Errorf("path segment %q is not a map", part)
		}

		if isLast {
			m[part] = value
			return nil
		}

		next, exists := m[part]
		if !exists {
			newMap := map[string]interface{}{}
			m[part] = newMap
			current = newMap
			continue
		}

		current = next
	}

	return nil
}

func DeleteValueAtPath(data map[string]interface{}, path string) {
	if path == "" {
		return
	}

	parts := strings.Split(path, ".")
	if len(parts) == 0 {
		return
	}

	var current interface{} = data

	for _, part := range parts[:len(parts)-1] {
		m, ok := CoerceToMap(current)
		if !ok {
			return
		}
		current = m[part]
	}

	m, ok := CoerceToMap(current)
	if !ok {
		return
	}

	delete(m, parts[len(parts)-1])
}
