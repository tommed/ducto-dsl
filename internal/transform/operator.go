package transform

import (
	"github.com/tommed/dsl-transformer/internal/model"
)

type Operator interface {
	Name() string
	Validate(instr model.Instruction) error
	Apply(ctx *ExecutionContext, reg *Registry, input map[string]interface{}, instr model.Instruction) error
}
