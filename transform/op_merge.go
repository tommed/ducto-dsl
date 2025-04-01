package transform

import (
	"fmt"
	"github.com/tommed/ducto-dsl/model"
)

type MergeOperator struct{}

func (o *MergeOperator) Name() string {
	return "merge"
}

func (o *MergeOperator) Validate(instr model.Instruction) error {
	if instr.Value == nil {
		return fmt.Errorf("merge operator missing 'value'")
	}
	if _, ok := instr.Value.(map[string]interface{}); !ok {
		return fmt.Errorf("merge operator: 'value' must be an object")
	}
	return nil
}

func (o *MergeOperator) Apply(_ *ExecutionContext, _ *Registry, input map[string]interface{}, instr model.Instruction) error {
	for k, v := range instr.Value.(map[string]interface{}) {
		// If `if_not_set` is true, only set if missing
		if instr.IfNotSet {
			if _, exists := input[k]; exists {
				continue
			}
		}
		input[k] = v
	}

	return nil
}
