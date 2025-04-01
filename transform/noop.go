package transform

import (
	"github.com/tommed/ducto-dsl/model"
)

// NoOperation is our nil implementation, it literally does nothing and never fails
type NoOperation struct{}

func (n NoOperation) Validate(instr model.Instruction) error {
	return nil
}

func (n NoOperation) Name() string {
	return "noop"
}

func (n NoOperation) Apply(_ *ExecutionContext, _ *Registry, _ map[string]interface{}, _ model.Instruction) error {
	return nil
}
