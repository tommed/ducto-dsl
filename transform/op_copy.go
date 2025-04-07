package transform

import (
	"errors"
)

type CopyOperator struct{}

func (c *CopyOperator) Validate(instr Instruction) error {
	if instr.From == "" {
		return errors.New("copy op missing or invalid from")
	}
	if instr.To == "" {
		return errors.New("copy op missing or invalid to")
	}
	return nil
}

func (c *CopyOperator) Name() string { return "copy" }

func (c *CopyOperator) Apply(_ *ExecutionContext, _ *Registry, input map[string]interface{}, instr Instruction) error {
	val, _ := GetValueAtPath(input, instr.From)
	return SetValueAtPath(input, instr.To, val)
}
