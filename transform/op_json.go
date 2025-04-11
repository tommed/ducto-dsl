package transform

import (
	"encoding/json"
	"fmt"
)

type ToJSONOperator struct{}

func (o *ToJSONOperator) Name() string { return "to_json" }

func (o *ToJSONOperator) Validate(instr Instruction) error {
	if instr.From == "" {
		return fmt.Errorf("to_json: missing 'from'")
	}
	if instr.To == "" {
		return fmt.Errorf("to_json: missing 'to'")
	}
	return nil
}

func (o *ToJSONOperator) Apply(_ *ExecutionContext, _ *Registry, input map[string]interface{}, instr Instruction) error {
	val, ok := GetValueAtPath(input, instr.From)
	if !ok {
		return fmt.Errorf("to_json: 'from' path not found: %s", instr.From)
	}
	data, err := json.Marshal(val)
	if err != nil {
		return fmt.Errorf("to_json: failed to marshal value: %w", err)
	}
	return SetValueAtPath(input, instr.To, string(data))
}

// ---

type FromJSONOperator struct{}

func (o *FromJSONOperator) Name() string { return "from_json" }

func (o *FromJSONOperator) Validate(instr Instruction) error {
	if instr.From == "" {
		return fmt.Errorf("from_json: missing 'from'")
	}
	if instr.To == "" {
		return fmt.Errorf("from_json: missing 'to'")
	}
	return nil
}

func (o *FromJSONOperator) Apply(_ *ExecutionContext, _ *Registry, input map[string]interface{}, instr Instruction) error {
	val, ok := GetValueAtPath(input, instr.From)
	if !ok {
		return fmt.Errorf("from_json: 'from' path not found: %s", instr.From)
	}
	str, ok := val.(string)
	if !ok {
		return fmt.Errorf("from_json: value at 'from' is not a string")
	}
	var out interface{}
	if err := json.Unmarshal([]byte(str), &out); err != nil {
		return fmt.Errorf("from_json: failed to unmarshal JSON: %w", err)
	}
	return SetValueAtPath(input, instr.To, out)
}
