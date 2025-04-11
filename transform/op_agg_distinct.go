package transform

import (
	"fmt"
	"reflect"
)

type AggDistinctOperator struct{}

func (o *AggDistinctOperator) Name() string { return "agg_distinct_value" }

func (o *AggDistinctOperator) Validate(instr Instruction) error {
	if instr.From == "" {
		return fmt.Errorf("agg_distinct_value: missing 'from'")
	}
	if instr.To == "" {
		return fmt.Errorf("agg_distinct_value: missing 'to'")
	}
	if instr.Key == "" {
		return fmt.Errorf("agg_distinct_value: missing 'key'")
	}
	return nil
}

func (o *AggDistinctOperator) Apply(_ *ExecutionContext, _ *Registry, input map[string]interface{}, instr Instruction) error {
	val, ok := GetValueAtPath(input, instr.From)
	if !ok {
		return fmt.Errorf("agg_distinct_value: 'from' path not found: %s", instr.From)
	}

	arr, ok := CoerceToArray(val)
	if !ok {
		return fmt.Errorf("agg_distinct_value: value at 'from' path is not an array")
	}

	seen := make(map[interface{}]struct{})
	var distinct []interface{}

	for _, item := range arr {
		itemMap, ok := CoerceToMap(item)
		if !ok {
			continue
		}
		fieldVal, exists := GetValueAtPath(itemMap, instr.Key)
		if !exists {
			continue
		}

		key := normaliseKey(fieldVal)
		if _, exists := seen[key]; !exists {
			seen[key] = struct{}{}
			distinct = append(distinct, fieldVal)
		}
	}

	return SetValueAtPath(input, instr.To, distinct)
}

// normaliseKey ensures that the map key is comparable
func normaliseKey(v interface{}) interface{} {
	if v == nil {
		return "<invalid>"
	}

	switch val := v.(type) {
	case string, bool, float64, int, int64:
		return val
	default:
		rv := reflect.ValueOf(v)
		if !rv.IsValid() {
			return "<invalid>"
		}
		return fmt.Sprintf("%v", v)
	}
}
