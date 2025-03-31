package transform

import (
	"context"
	"errors"
	"github.com/tommed/dsl-transformer/internal/model"
)

type Registry struct {
	ops map[string]Operator
}

func NewRegistry() *Registry {
	return &Registry{ops: map[string]Operator{}}
}

func (r *Registry) Register(op Operator) {
	r.ops[op.Name()] = op
}

func (r *Registry) Apply(ctx context.Context, input map[string]interface{}, instr model.Instruction) error {
	op, ok := r.ops[instr.Op]
	if !ok {
		return errors.New("unknown op: " + instr.Op)
	}
	return op.Apply(ctx, input, instr)
}
