package transform

import (
	"errors"
	"github.com/tommed/ducto-dsl/model"
)

type CopyOperator struct{}

func (c *CopyOperator) Validate(instr model.Instruction) error {
	if instr.From == "" {
		return errors.New("copy op missing or invalid from")
	}
	if instr.To == "" {
		return errors.New("copy op missing or invalid to")
	}
	return nil
}

func (c *CopyOperator) Name() string { return "copy" }

func (c *CopyOperator) Apply(_ *ExecutionContext, _ *Registry, input map[string]interface{}, instr model.Instruction) error {
	input[instr.To] = input[instr.From]
	return nil
}
