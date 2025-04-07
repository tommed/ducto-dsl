package transform

import (
	"errors"
)

type DeleteOperator struct{}

func (d *DeleteOperator) Validate(instr Instruction) error {
	if instr.Key == "" {
		return errors.New("delete operator missing 'key'")
	}
	return nil
}

func (d *DeleteOperator) Name() string { return "delete" }

func (d *DeleteOperator) Apply(_ *ExecutionContext, _ *Registry, input map[string]interface{}, instr Instruction) error {
	DeleteValueAtPath(input, instr.Key)
	return nil
}
