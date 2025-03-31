package dsl

import (
	"context"
	"errors"
	"github.com/tommed/dsl-transformer/internal/model"
	"github.com/tommed/dsl-transformer/internal/transform"
)

// Transformer applies DSL-defined transformations
type Transformer struct {
	reg *transform.Registry
}

// New creates a new Transformer
func New() *Transformer {
	reg := transform.NewRegistry()
	reg.Register(&transform.SetOperator{})
	reg.Register(&transform.CopyOperator{})
	reg.Register(&transform.DeleteOperator{})
	reg.Register(&transform.MapOperator{})
	reg.Register(&transform.FailOperator{})
	reg.Register(&transform.NoOperation{})
	return &Transformer{reg: reg}
}

// Apply applies the given transformation definition
func (t *Transformer) Apply(ctx context.Context, input map[string]interface{}, prog *model.Program) (map[string]interface{}, error) {

	// Create our output, start with the input values
	output := make(map[string]interface{})
	for k, v := range input {
		output[k] = v
	}

	// Create a new context
	exec := transform.NewExecutionContext(ctx, prog.OnError)

	// Apply instructions
	for _, instr := range prog.Instructions {
		if !t.reg.Apply(exec, t.reg, output, instr) {
			return nil, errors.New("execution halted due to an error")
		}
	}

	// HandleError errors
	if exec.OnError == "error" && len(exec.Errors) > 0 {
		output["@dsl_errors"] = exec.Errors
	}

	return output, nil
}
