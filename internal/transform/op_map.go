package transform

import (
	"fmt"
	"github.com/tommed/dsl-transformer/internal/model"
)

type MapOperator struct{}

func (o *MapOperator) Validate(instr model.Instruction) error {
	if instr.Key == "" {
		return fmt.Errorf("map operator requires 'key' field")
	}

	if len(instr.Then) == 0 {
		return fmt.Errorf("map operator requires at least one instruction in 'then'")
	}
	return nil
}

func (o *MapOperator) Name() string {
	return "map"
}

func (o *MapOperator) Apply(ctx *ExecutionContext, reg *Registry, input map[string]interface{}, instr model.Instruction) error {
	arrVal, ok := input[instr.Key]
	if !ok {
		return fmt.Errorf("map operator: key %q not found in input", instr.Key)
	}

	arr, ok := arrVal.([]interface{})
	if !ok {
		return fmt.Errorf("map operator: input[%q] is not an array", instr.Key)
	}

	for i, item := range arr {
		itemMap, ok := item.(map[string]interface{})
		if !ok {
			return fmt.Errorf("map operator: array item at index %d is not an object", i)
		}

		for _, subInstr := range instr.Then {
			if !reg.Apply(ctx, reg, itemMap, subInstr) {
				return fmt.Errorf("map operator: sub-instruction failed at index %d", i)
			}
		}
	}

	return nil
}
