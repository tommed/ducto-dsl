package transform

import (
	"errors"
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
	if _, ok := input[instr.Key]; !ok {
		input[instr.Key] = instr.Value
	}
	return nil
}
