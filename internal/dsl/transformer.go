package dsl

import (
	"context"
	"errors"
	"github.com/tommed/ducto-dsl/model"
	transform2 "github.com/tommed/ducto-dsl/transform"
)

// Transformer applies DSL-defined transformations
type Transformer struct {
	reg *transform2.Registry
}

// New creates a new Transformer
func New() *Transformer {
	reg := transform2.NewDefaultRegistry()
	return &Transformer{reg: reg}
}

// Apply applies the given transformation definition
func (t *Transformer) Apply(ctx context.Context, input map[string]interface{}, prog *model.Program) (map[string]interface{}, error) {

	// Validate program before execution
	if err := ValidateProgram(t.reg, prog); err != nil {
		return nil, err
	}

	// Create a new context
	exec := transform2.NewExecutionContext(ctx, prog.OnError)

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
