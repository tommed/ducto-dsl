package transform

import (
	"fmt"
)

type AggSumOperator struct{}

func (o *AggSumOperator) Name() string { return "agg_sum" }

func (o *AggSumOperator) Validate(instr Instruction) error {
	if instr.From == "" {
		return fmt.Errorf("agg_sum: missing 'from'")
	}
	if instr.To == "" {
		return fmt.Errorf("agg_sum: missing 'to'")
	}
	if instr.Key == "" {
		return fmt.Errorf("agg_sum: missing 'key'")
	}
	return nil
}

func (o *AggSumOperator) Apply(_ *ExecutionContext, _ *Registry, input map[string]interface{}, instr Instruction) error {
	val, ok := GetValueAtPath(input, instr.From)
	if !ok {
		return fmt.Errorf("agg_sum: 'from' path not found: %s", instr.From)
	}

	arr, ok := CoerceToArray(val)
	if !ok {
		return fmt.Errorf("agg_sum: value at 'from' path is not an array")
	}

	var sum float64

	for i, item := range arr {
		itemMap, ok := CoerceToMap(item)
		if !ok {
			continue // skip invalid items
		}
		fieldVal, exists := GetValueAtPath(itemMap, instr.Key)
		if !exists {
			continue
		}

		switch v := fieldVal.(type) {
		case float64:
			sum += v
		case int:
			sum += float64(v)
		case int64:
			sum += float64(v)
		case float32:
			sum += float64(v)
		default:
			// silently skip
			_ = i // force usage of index to avoid compile warning if needed
		}
	}

	return SetValueAtPath(input, instr.To, sum)
}
