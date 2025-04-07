package transform

import (
	"fmt"
)

type SetOperator struct{}

func (s *SetOperator) Name() string { return "set" }

func (s *SetOperator) Validate(instr Instruction) error {
	if instr.Key == "" {
		return fmt.Errorf("set operator missing 'key'")
	}
	return nil
}

func (s *SetOperator) Apply(_ *ExecutionContext, _ *Registry, input map[string]interface{}, instr Instruction) error {
	return SetValueAtPath(input, instr.Key, instr.Value)
}
