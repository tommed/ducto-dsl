package dsl

import (
	"context"
	"errors"
)

type Instruction struct {
	Op    string      `json:"op"`
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

type Program struct {
	Instructions []Instruction `json:"instructions"`
}

// Transformer applies DSL-defined transformations
type Transformer struct{}

// New creates a new Transformer
func New() *Transformer {
	return &Transformer{}
}

// Apply applies the given transformation definition
func (t *Transformer) Apply(ctx context.Context, input map[string]interface{}, prog *Program) (map[string]interface{}, error) {

	// Create our output, start with the input values
	output := make(map[string]interface{})
	for k, v := range input {
		output[k] = v
	}

	for _, instr := range prog.Instructions {
		switch instr.Op {
		case "set":
			output[instr.Key] = instr.Value
		default:
			return nil, errors.New("unsupported op: " + instr.Op)
		}
	}

	return output, nil
}
