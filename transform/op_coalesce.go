package transform

import (
	"errors"
	"reflect"
)

type CoalesceOperator struct{}

func (o *CoalesceOperator) Name() string { return "coalesce" }

func (o *CoalesceOperator) Validate(instr Instruction) error {
	if instr.Key == "" {
		return errors.New("coalesce operator requires a key")
	}
	if instr.Value == nil {
		return errors.New("coalesce operator requires a default value")
	}
	return nil
}

func (o *CoalesceOperator) Apply(_ *ExecutionContext, _ *Registry, input map[string]interface{}, instr Instruction) error {
	val, exists := GetValueAtPath(input, instr.Key)
	if !exists || isZero(val) {
		return SetValueAtPath(input, instr.Key, instr.Value)
	}
	return nil
}

// isZero returns true for nil, empty strings, zero numbers, empty slices/maps
func isZero(v interface{}) bool {
	if v == nil {
		return true
	}
	switch val := v.(type) {
	case string:
		return val == ""
	case float64:
		return val == 0
	case int, int32, int64:
		return reflect.ValueOf(v).Int() == 0
	case uint, uint32, uint64:
		return reflect.ValueOf(v).Uint() == 0
	case bool:
		return false // don't treat false as "empty"
	case []interface{}:
		return len(val) == 0
	case map[string]interface{}:
		return len(val) == 0
	default:
		rv := reflect.ValueOf(v)
		return rv.Kind() == reflect.Ptr && rv.IsNil()
	}
}
