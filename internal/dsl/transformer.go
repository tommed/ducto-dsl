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
	reg := transform.NewDefaultRegistry()
	return &Transformer{reg: reg}
}

// Apply applies the given transformation definition
func (t *Transformer) Apply(ctx context.Context, input map[string]interface{}, prog *model.Program) (map[string]interface{}, error) {

	// Validate program before execution
	if err := ValidateProgram(t.reg, prog); err != nil {
		return nil, err
	}

	// Create a new context
	exec := transform.NewExecutionContext(ctx, prog.OnError)

	// Create our output, start with the input values
	output := make(map[string]interface{})
	for k, v := range input {
		output[k] = v
	}

	// Apply instructions
	for _, instr := range prog.Instructions {
		if !t.reg.Apply(exec, t.reg, output, instr) {
			return nil, errors.New("execution halted due to an error")
		}
	}

	// HandleError errors
	if exec.OnError == "capture" && len(exec.Errors) > 0 {
		output["@dsl_errors"] = exec.Errors
	}

	return output, nil
}
