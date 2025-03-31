package transform

import (
	"errors"
	"github.com/tommed/dsl-transformer/internal/model"
)

type DeleteOperator struct{}

func (d *DeleteOperator) Validate(instr model.Instruction) error {
	if instr.Key == "" {
		return errors.New("delete operator missing 'key'")
	}
	return nil
}

func (d *DeleteOperator) Name() string { return "delete" }

func (d *DeleteOperator) Apply(_ *ExecutionContext, _ *Registry, input map[string]interface{}, instr model.Instruction) error {
	delete(input, instr.Key)
	return nil
}
