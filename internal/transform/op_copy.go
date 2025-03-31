package transform

import (
	"errors"

	"github.com/tommed/dsl-transformer/internal/model"
)

type CopyOperator struct{}

func (c *CopyOperator) Name() string { return "copy" }

func (c *CopyOperator) Apply(ctx *ExecutionContext, input map[string]interface{}, instr model.Instruction) error {
	if instr.From == "" {
		return errors.New("copy op missing or invalid from")
	}
	if instr.To == "" {
		return errors.New("copy op missing or invalid to")
	}
	input[instr.To] = input[instr.From]
	return nil
}
