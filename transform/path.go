package transform

import (
	"fmt"
	"strings"
)

// GetValueAtPath retrieves a nested value using dot notation, e.g. "foo.bar.baz"
func GetValueAtPath(data map[string]interface{}, path string) (interface{}, bool) {
	parts := strings.Split(path, ".")
	var current interface{} = data

	for _, part := range parts {
		asMap, ok := current.(map[string]interface{})
		if !ok {
			return nil, false
		}
		current, ok = asMap[part]
		if !ok {
			return nil, false
		}
	}
	return current, true
}

// SetValueAtPath sets a value deeply within a map using dot notation
func SetValueAtPath(data map[string]interface{}, path string, value interface{}) error {
	parts := strings.Split(path, ".")
	if len(parts) == 0 {
		return fmt.Errorf("invalid path: empty")
	}

	current := data
	for _, part := range parts[:len(parts)-1] {
		next, exists := current[part]
		if !exists {
			newMap := make(map[string]interface{})
			current[part] = newMap
			current = newMap
			continue
		}
		asMap, ok := next.(map[string]interface{})
		if !ok {
			return fmt.Errorf("cannot descend into non-map value at %q", part)
		}
		current = asMap
	}
	current[parts[len(parts)-1]] = value
	return nil
}

// File: ducto-dsl/transform/path.go

// DeleteValueAtPath removes a key from a nested map structure via dot notation.
// If the path or intermediate keys don't exist, it exits silently.
func DeleteValueAtPath(data map[string]interface{}, path string) {
	parts := strings.Split(path, ".")
	if len(parts) == 0 {
		return
	}
	current := data

	for _, part := range parts[:len(parts)-1] {
		next, ok := current[part]
		if !ok {
			return // path does not exist
		}
		asMap, ok := next.(map[string]interface{})
		if !ok {
			return // not a map — can't proceed
		}
		current = asMap
	}

	delete(current, parts[len(parts)-1])
}
