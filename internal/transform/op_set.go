package transform

import (
	"context"
	"errors"
	"github.com/tommed/dsl-transformer/internal/model"
)

type SetOperator struct{}

func (s *SetOperator) Name() string { return "set" }

func (s *SetOperator) Apply(_ context.Context, input map[string]interface{}, instr model.Instruction) error {
	if instr.Key == "" {
		return errors.New("set op missing key")
	}
	input[instr.Key] = instr.Value
	return nil
}
