package transform

import (
	"fmt"
)

type DropIfOperator struct{}

func (d *DropIfOperator) Name() string { return "drop_if" }

func (d *DropIfOperator) Validate(instr Instruction) error {
	if instr.Key == "" {
		return fmt.Errorf("drop_if requires a 'key'")
	}
	return nil
}

func (d *DropIfOperator) Apply(ctx *ExecutionContext, _ *Registry, input map[string]interface{}, instr Instruction) error {
	val, ok := GetValueAtPath(input, instr.Key)

	drop := ok && val == true
	if instr.Not {
		drop = !drop
	}

	if drop {
		ctx.Dropped = true
		return nil // I know... but this just looks more right to me!
	}

	return nil
}
