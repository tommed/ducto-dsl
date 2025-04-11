package transform

import (
	"fmt"
	"sort"
)

type AggAvgOperator struct{}

func (o *AggAvgOperator) Name() string { return "agg_avg" }

func (o *AggAvgOperator) Validate(instr Instruction) error {
	if instr.From == "" {
		return fmt.Errorf("agg_avg: missing 'from'")
	}
	if instr.To == "" {
		return fmt.Errorf("agg_avg: missing 'to'")
	}
	if instr.Key == "" {
		return fmt.Errorf("agg_avg: missing 'key'")
	}
	return nil
}

func (o *AggAvgOperator) Apply(_ *ExecutionContext, _ *Registry, input map[string]interface{}, instr Instruction) error {
	val, ok := GetValueAtPath(input, instr.From)
	if !ok {
		return fmt.Errorf("agg_avg: 'from' path not found: %s", instr.From)
	}

	arr, ok := CoerceToArray(val)
	if !ok {
		return fmt.Errorf("agg_avg: value at 'from' path is not an array")
	}

	var values []float64
	for _, item := range arr {
		itemMap, ok := CoerceToMap(item)
		if !ok {
			continue
		}
		fval, exists := GetValueAtPath(itemMap, instr.Key)
		if !exists {
			continue
		}
		switch v := fval.(type) {
		case float64:
			values = append(values, v)
		case int:
			values = append(values, float64(v))
		case int64:
			values = append(values, float64(v))
		case float32:
			values = append(values, float64(v))
		}
	}

	if len(values) == 0 {
		return SetValueAtPath(input, instr.To, nil)
	}

	switch instr.Variant {
	case "mode":
		return SetValueAtPath(input, instr.To, computeMode(values))
	case "median":
		return SetValueAtPath(input, instr.To, computeMedian(values))
	default:
		return SetValueAtPath(input, instr.To, computeMean(values))
	}
}

func computeMean(vals []float64) float64 {
	sum := 0.0
	for _, v := range vals {
		sum += v
	}
	return sum / float64(len(vals))
}

func computeMedian(vals []float64) float64 {
	sorted := append([]float64(nil), vals...)
	sort.Float64s(sorted)
	n := len(sorted)
	mid := n / 2
	if n%2 == 0 {
		return (sorted[mid-1] + sorted[mid]) / 2
	}
	return sorted[mid]
}

func computeMode(vals []float64) float64 {
	freq := make(map[float64]int)
	maxFreq := 0
	var mode float64
	for _, v := range vals {
		freq[v]++
		if freq[v] > maxFreq {
			maxFreq = freq[v]
			mode = v
		}
	}
	return mode
}
