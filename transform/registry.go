package transform

import (
	"fmt"
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
	reg.Register(&DropIfOperator{})
	reg.Register(&StringJoinOperator{})
	reg.Register(&ArrayLengthOperator{})
	reg.Register(&AggSumOperator{})
	reg.Register(&AggDistinctOperator{})
	reg.Register(&FilterOperator{})
	reg.Register(&ToJSONOperator{})
	reg.Register(&FromJSONOperator{})
	reg.Register(&ReplaceOperator{})
	reg.Register(&RegexReplaceOperator{})
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

func (r *Registry) Apply(ctx *ExecutionContext, reg *Registry, input map[string]interface{}, instr Instruction) bool {
	op := r.ops[instr.Op]
	if err := op.Apply(ctx, reg, input, instr); err != nil {
		return ctx.HandleError(err)
	}
	return true
}

func (r *Registry) Find(op string) Operator {
	return r.ops[op]
}
