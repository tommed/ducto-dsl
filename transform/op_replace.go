package transform

import (
	"fmt"
	"regexp"
	"strings"
)

// ReplaceOperator is for simple string replace all
type ReplaceOperator struct{}

func (o *ReplaceOperator) Name() string { return "replace" }

func (o *ReplaceOperator) Validate(instr Instruction) error {
	if instr.From == "" {
		return fmt.Errorf("replace: missing 'from'")
	}
	if instr.To == "" {
		return fmt.Errorf("replace: missing 'to'")
	}
	if instr.Key == "" {
		return fmt.Errorf("replace: missing 'key' (match)")
	}
	if instr.Value == nil {
		return fmt.Errorf("replace: missing 'value' (replacement)")
	}
	return nil
}

func (o *ReplaceOperator) Apply(_ *ExecutionContext, _ *Registry, input map[string]interface{}, instr Instruction) error {
	val, ok := GetValueAtPath(input, instr.From)
	if !ok {
		return fmt.Errorf("replace: 'from' path not found: %s", instr.From)
	}
	src, ok := val.(string)
	if !ok {
		return fmt.Errorf("replace: value at 'from' is not a string")
	}

	match := instr.Key
	replacement, ok := instr.Value.(string)
	if !ok {
		return fmt.Errorf("replace: 'value' must be a string")
	}

	out := strings.ReplaceAll(src, match, replacement)
	return SetValueAtPath(input, instr.To, out)
}

// RegexReplaceOperator is for regex string replaces
type RegexReplaceOperator struct{}

func (o *RegexReplaceOperator) Name() string { return "regex_replace" }

func (o *RegexReplaceOperator) Validate(instr Instruction) error {
	if instr.From == "" {
		return fmt.Errorf("regex_replace: missing 'from'")
	}
	if instr.To == "" {
		return fmt.Errorf("regex_replace: missing 'to'")
	}
	if instr.Key == "" {
		return fmt.Errorf("regex_replace: missing 'key' (regex pattern)")
	}
	if instr.Value == nil {
		return fmt.Errorf("regex_replace: missing 'value' (replacement string)")
	}
	_, err := regexp.Compile(instr.Key)
	if err != nil {
		return fmt.Errorf("regex_replace: invalid regex: %w", err)
	}
	return nil
}

func (o *RegexReplaceOperator) Apply(_ *ExecutionContext, _ *Registry, input map[string]interface{}, instr Instruction) error {
	val, ok := GetValueAtPath(input, instr.From)
	if !ok {
		return fmt.Errorf("regex_replace: 'from' path not found: %s", instr.From)
	}
	src, ok := val.(string)
	if !ok {
		return fmt.Errorf("regex_replace: value at 'from' is not a string")
	}

	re, err := regexp.Compile(instr.Key)
	if err != nil {
		return fmt.Errorf("regex_replace: invalid regex: %w", err)
	}
	replacement, ok := instr.Value.(string)
	if !ok {
		return fmt.Errorf("regex_replace: 'value' must be a string")
	}

	result := re.ReplaceAllString(src, replacement)
	return SetValueAtPath(input, instr.To, result)
}
