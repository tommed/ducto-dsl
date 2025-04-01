package transform

import (
	"fmt"
	"github.com/tommed/ducto-dsl/model"
)

type SetOperator struct{}

func (s *SetOperator) Name() string { return "set" }

func (s *SetOperator) Validate(instr model.Instruction) error {
	if instr.Key == "" {
		return fmt.Errorf("set operator missing 'key'")
	}
	return nil
}

func (s *SetOperator) Apply(_ *ExecutionContext, _ *Registry, input map[string]interface{}, instr model.Instruction) error {
	input[instr.Key] = instr.Value
	return nil
}
