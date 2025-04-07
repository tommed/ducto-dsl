// File: ducto-dsl/transform/op_if.go
package transform

import (
	"encoding/json"
	"errors"
	"fmt"
)

// validateConditions ensures the conditions map uses valid clauses
func validateConditions(condition map[string]interface{}) error {
	if len(condition) == 0 {
		return errors.New("no conditions defined")
	}
	if len(condition) > 1 {
		return fmt.Errorf("only one condition type is allowed per condition block, got: %v", condition)
	}
	for key := range condition {
		switch key {
		case "exists", "equals", "or", "and":
			break
		default:
			return fmt.Errorf("unknown condition %q", key)
		}
	}
	return nil
}

// evaluateCondition evaluates the conditions
func evaluateCondition(input map[string]interface{}, condition map[string]interface{}) bool {
	switch {
	case condition["exists"] != nil:
		return conditionExists(input, condition)
	case condition["equals"] != nil:
		return conditionEquals(input, condition)
	case condition["or"] != nil:
		return conditionOr(input, condition)
	case condition["and"] != nil:
		return conditionAnd(input, condition)
	default:
		return false
	}
}

func conditionExists(input map[string]interface{}, condition map[string]interface{}) bool {
	key, ok := condition["exists"].(string)
	if !ok {
		return false
	}
	_, exists := GetValueAtPath(input, key)
	return exists
}

func conditionEquals(input map[string]interface{}, condition map[string]interface{}) bool {
	data, ok := condition["equals"].(map[string]interface{})
	if !ok {
		return false
	}
	key, ok := data["key"].(string)
	if !ok {
		return false
	}
	expected := data["value"]

	actual, exists := GetValueAtPath(input, key)
	if !exists {
		return false
	}

	return jsonEqual(actual, expected)
}

func conditionOr(input map[string]interface{}, condition map[string]interface{}) bool {
	conds, ok := condition["or"].([]interface{})
	if !ok {
		return false
	}
	for _, c := range conds {
		if sub, ok := c.(map[string]interface{}); ok {
			if evaluateCondition(input, sub) {
				return true
			}
		}
	}
	return false
}

func conditionAnd(input map[string]interface{}, condition map[string]interface{}) bool {
	conds, ok := condition["and"].([]interface{})
	if !ok {
		return false
	}
	for _, c := range conds {
		if sub, ok := c.(map[string]interface{}); ok {
			if !evaluateCondition(input, sub) {
				return false
			}
		}
	}
	return true
}

func jsonEqual(a, b interface{}) bool {
	aj, _ := json.Marshal(a)
	bj, _ := json.Marshal(b)
	return string(aj) == string(bj)
}
