package transform

import (
	"errors"
	"github.com/tommed/dsl-transformer/internal/model"
)

type FailOperator struct {
}

func (f FailOperator) Name() string {
	return "fail"
}

func (f FailOperator) Apply(_ *ExecutionContext, _ *Registry, _ map[string]interface{}, instr model.Instruction) error {
	return errors.New(instr.Value.(string))
}
