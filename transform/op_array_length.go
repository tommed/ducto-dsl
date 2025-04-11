package transform

import "fmt"

type ArrayLengthOperator struct{}

func (o *ArrayLengthOperator) Name() string { return "array_length" }

func (o *ArrayLengthOperator) Validate(instr Instruction) error {
	if instr.From == "" {
		return fmt.Errorf("array_length: missing 'from'")
	}
	if instr.To == "" {
		return fmt.Errorf("array_length: missing 'to'")
	}
	return nil
}

func (o *ArrayLengthOperator) Apply(_ *ExecutionContext, _ *Registry, input map[string]interface{}, instr Instruction) error {
	val, ok := GetValueAtPath(input, instr.From)
	if !ok {
		return fmt.Errorf("array_length: 'from' path not found: %s", instr.From)
	}

	arr, ok := CoerceToArray(val)
	if !ok {
		return fmt.Errorf("array_length: value at 'from' path is not an array")
	}

	return SetValueAtPath(input, instr.To, len(arr))
}
