package transform

import (
	"context"
	"github.com/tommed/dsl-transformer/internal/model"
)

type Operator interface {
	Name() string
	Apply(ctx context.Context, input map[string]interface{}, instr model.Instruction) error
}
