package transform

import (
	"fmt"
	"github.com/tommed/ducto-dsl/model"
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
	reg.Register(&IfOperator{})
	reg.Register(&CoalesceOperator{})
	for _, op := range optional {
		reg.Register(op)
	}
	return reg
}

func NewRegistry() *Registry {
	return &Registry{ops: map[string]Operator{}}
}

func (r *Registry) Register(op Operator) {
	name := op.Name()
	if name == "" {
		panic("operator has no name")
	}
	if _, exists := r.ops[name]; exists {
		panic(fmt.Sprintf("operator with name '%s' is already registered", name))
	}
	r.ops[name] = op
}

func (r *Registry) Apply(ctx *ExecutionContext, reg *Registry, input map[string]interface{}, instr model.Instruction) bool {
	op := r.ops[instr.Op]
	if err := op.Apply(ctx, reg, input, instr); err != nil {
		return ctx.HandleError(err)
	}

	return true
}

func (r *Registry) Find(op string) Operator {
	return r.ops[op]
}
