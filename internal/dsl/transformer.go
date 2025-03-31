package dsl

import (
	"context"
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
	return &Transformer{reg: reg}
}

// Apply applies the given transformation definition
func (t *Transformer) Apply(ctx context.Context, input map[string]interface{}, prog *model.Program) (map[string]interface{}, error) {

	// Create our output, start with the input values
	output := make(map[string]interface{})
	for k, v := range input {
		output[k] = v
	}

	for _, instr := range prog.Instructions {
		if err := t.reg.Apply(ctx, output, instr); err != nil {
			return nil, err
		}
	}

	return output, nil
}
