package transform

import (
	"context"
	"errors"

	"github.com/tommed/dsl-transformer/internal/model"
)

type DeleteOperator struct{}

func (d *DeleteOperator) Name() string { return "delete" }

func (d *DeleteOperator) Apply(_ context.Context, input map[string]interface{}, instr model.Instruction) error {
	if instr.Key == "" {
		return errors.New("delete op missing key")
	}
	delete(input, instr.Key)
	return nil
}
