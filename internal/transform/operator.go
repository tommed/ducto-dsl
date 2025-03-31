package transform

import (
	"github.com/tommed/dsl-transformer/internal/model"
)

type Operator interface {
	Name() string
	Apply(ctx *ExecutionContext, input map[string]interface{}, instr model.Instruction) error
}
