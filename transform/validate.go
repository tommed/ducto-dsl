package transform

import (
	"fmt"
)

func ValidateProgram(r *Registry, prog *Program) error {
	if prog == nil {
		return fmt.Errorf("program is nil")
	}
	if prog.Version != 1 {
		return fmt.Errorf("program version %d is unsupported", prog.Version)
	}

	for i, instr := range prog.Instructions {
		op := r.Find(instr.Op)
		if op == nil {
			return fmt.Errorf("instruction #%d: unknown operator '%s'", i, instr.Op)
		}

		if err := op.Validate(instr); err != nil {
			return fmt.Errorf("instruction #%d (%s): %w", i, instr.Op, err)
		}

		// Validate nested instructions (e.g., map, if, filter)
		if len(instr.Then) > 0 {
			subProg := &Program{
				Version:      prog.Version,
				Instructions: instr.Then,
			}
			if err := ValidateProgram(r, subProg); err != nil {
				return fmt.Errorf("instruction #%d (%s): nested validation failed: %w", i, instr.Op, err)
			}
		}
	}
	return nil
}
