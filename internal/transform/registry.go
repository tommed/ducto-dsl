package transform

import (
	"fmt"
	"github.com/tommed/dsl-transformer/internal/model"
)

type Registry struct {
	ops map[string]Operator
}

func NewDefaultRegistry(optional ...Operator) *Registry {
	reg := NewRegistry()
	reg.Register(&SetOperator{})
	reg.Register(&CopyOperator{})
	reg.Register(&DeleteOperator{})
	reg.Register(&MapOperator{})
	reg.Register(&FailOperator{})
	reg.Register(&NoOperation{})
	reg.Register(&MergeOperator{})
	for _, op := range optional {
		reg.Register(op)
	}
	return reg
}

func NewRegistry() *Registry {
	return &Registry{ops: map[string]Operator{}}
}

func (r *Registry) Register(op Operator) {
	r.ops[op.Name()] = op
}

func (r *Registry) Apply(ctx *ExecutionContext, reg *Registry, input map[string]interface{}, instr model.Instruction) bool {
	op, ok := r.ops[instr.Op]
	if !ok {
		return ctx.HandleError(fmt.Errorf("unknown op: %q", instr.Op))
	}

	if err := op.Apply(ctx, reg, input, instr); err != nil {
		return ctx.HandleError(err)
	}

	return true
}

func (r *Registry) Find(op string) Operator {
	return r.ops[op]
}
