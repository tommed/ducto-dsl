package transform

import (
	"errors"
)

type FailOperator struct {
}

func (f FailOperator) Validate(instr Instruction) error {
	return nil
}

func (f FailOperator) Name() string {
	return "fail"
}

func (f FailOperator) Apply(_ *ExecutionContext, _ *Registry, _ map[string]interface{}, instr Instruction) error {
	return errors.New(instr.Value.(string))
}
