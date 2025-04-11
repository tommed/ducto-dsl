package transform

import (
	"fmt"
	"strings"
)

type FilterOperator struct{}

func (o *FilterOperator) Name() string { return "filter" }

func (o *FilterOperator) Validate(instr Instruction) error {
	if instr.From == "" {
		return fmt.Errorf("filter: missing 'from'")
	}
	if err := validateConditions(instr.Condition); err != nil {
		return fmt.Errorf("filter: invalid condition: %w", err)
	}
	if len(instr.Then) == 0 {
		return fmt.Errorf("filter: missing 'then' block")
	}
	if instr.As != "" && strings.Contains(instr.As, ".") {
		return fmt.Errorf("filter: 'as' must be a top-level field without dots")
	}
	return nil
}

func (o *FilterOperator) Apply(ctx *ExecutionContext, reg *Registry, input map[string]interface{}, instr Instruction) error {
	val, ok := GetValueAtPath(input, instr.From)
	if !ok {
		return fmt.Errorf("filter: 'from' path not found: %s", instr.From)
	}

	arr, ok := CoerceToArray(val)
	if !ok {
		return fmt.Errorf("filter: value at 'from' path is not an array")
	}

	var filtered []interface{}
	for _, item := range arr {
		itemMap, ok := CoerceToMap(item)
		if !ok {
			continue
		}
		if evaluateCondition(itemMap, instr.Condition) {
			filtered = append(filtered, item)
		}
	}

	if len(filtered) == 0 {
		return nil
	}

	targetKey := instr.As
	if targetKey == "" {
		targetKey = "_ctx"
	}

	if err := SetValueAtPath(input, targetKey, filtered); err != nil {
		return fmt.Errorf("filter: failed to set '%s': %w", targetKey, err)
	}

	for _, nested := range instr.Then {
		op := reg.Find(nested.Op)
		if err := op.Apply(ctx, reg, input, nested); err != nil {
			return err
		}
	}

	DeleteValueAtPath(input, targetKey)
	return nil
}
