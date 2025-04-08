package transform

import (
	"fmt"
)

type MergeOperator struct{}

func (o *MergeOperator) Name() string {
	return "merge"
}

func (o *MergeOperator) Validate(instr Instruction) error {
	if instr.Value == nil {
		return fmt.Errorf("merge operator missing 'value'")
	}
	if _, ok := instr.Value.(map[string]interface{}); !ok {
		return fmt.Errorf("merge operator: 'value' must be an object")
	}
	return nil
}

func (o *MergeOperator) Apply(_ *ExecutionContext, _ *Registry, input map[string]interface{}, instr Instruction) error {
	inputAsMap, ok := CoerceToMap(instr.Value)
	if !ok {
		return fmt.Errorf("merge operator: value must be an object")
	}
	for k, v := range inputAsMap {

		// If `if_not_set` is true, only set if missing
		if instr.IfNotSet {
			_, exists := GetValueAtPath(input, k)
			if exists {
				continue
			}
		}

		// Set value (don't need a path here as it's always the root object)
		input[k] = v
	}

	return nil
}
