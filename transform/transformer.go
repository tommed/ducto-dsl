package transform

import (
	"context"
	"errors"
)

// Transformer applies DSL-defined transformations
type Transformer struct {
	reg *Registry
}

// New creates a new Transformer
func New() *Transformer {
	reg := NewDefaultRegistry()
	return &Transformer{reg: reg}
}

// Apply applies the given transformation definition.
// NOTE: There is a scenario where both results are nil, meaning this input should be
// dropped/disregarded as requested by the policy author.
func (t *Transformer) Apply(ctx context.Context, input map[string]interface{}, prog *Program) (map[string]interface{}, error) {

	// Validate program before execution
	if err := ValidateProgram(t.reg, prog); err != nil {
		return nil, err
	}

	// Create a new context
	exec := NewExecutionContext(ctx, prog.OnError)
	debug, _ := ctx.Value(ContextKeyDebug).(bool)
	exec.Debug = debug

	// Create our output, start with the input values
	output := make(map[string]interface{})
	for k, v := range input {
		output[k] = v
	}

	// Apply instructions
	for _, instr := range prog.Instructions {
		ok := t.reg.Apply(exec, t.reg, output, instr)

		if exec.Dropped {
			return nil, nil
		}
		if !ok {
			return nil, errors.New("execution halted due to an error")
		}
	}

	// HandleError errors
	if exec.OnError == "capture" && len(exec.Errors) > 0 {
		output["@dsl_errors"] = exec.Errors
	}

	// Debug information
	if exec.Debug {
		output["@dsl_debug"] = map[string]interface{}{
			"applied_instructions": len(prog.Instructions),
			"errors":               len(exec.Errors),
			// could later include operator timings, traces etc.
		}
	}

	return output, nil
}
