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
	if instr.Key == "" {
		return fmt.Errorf("replace: missing 'key'")
	}
	if instr.Match == "" {
		return fmt.Errorf("replace: missing 'match'")
	}
	if instr.With == "" {
		return fmt.Errorf("replace: missing 'with' (replacement)")
	}
	if instr.To == "" {
		instr.To = instr.Key
	}
	return nil
}

func (o *ReplaceOperator) Apply(_ *ExecutionContext, _ *Registry, input map[string]interface{}, instr Instruction) error {
	val, ok := GetValueAtPath(input, instr.Key)
	if !ok {
		return fmt.Errorf("replace: 'key' path not found: %s", instr.Key)
	}
	src, ok := val.(string)
	if !ok {
		return fmt.Errorf("replace: value at 'key' is not a string")
	}

	match := instr.Match
	replacement := instr.With

	out := strings.ReplaceAll(src, match, replacement)
	return SetValueAtPath(input, instr.To, out)
}

// RegexReplaceOperator is for regex string replaces
type RegexReplaceOperator struct{}

func (o *RegexReplaceOperator) Name() string { return "regex_replace" }

func (o *RegexReplaceOperator) Validate(instr Instruction) error {
	if instr.Key == "" {
		return fmt.Errorf("regex_replace: missing 'key'")
	}
	if instr.Match == "" {
		return fmt.Errorf("regex_replace: missing 'match' (regex pattern)")
	}
	if instr.With == "" {
		return fmt.Errorf("regex_replace: missing 'with' (replacement string)")
	}
	if instr.To == "" {
		instr.To = instr.Key
	}
	_, err := regexp.Compile(instr.Match)
	if err != nil {
		return fmt.Errorf("regex_replace: invalid regex: %w", err)
	}
	return nil
}

func (o *RegexReplaceOperator) Apply(_ *ExecutionContext, _ *Registry, input map[string]interface{}, instr Instruction) error {
	val, ok := GetValueAtPath(input, instr.Key)
	if !ok {
		return fmt.Errorf("regex_replace: 'key' path not found: %s", instr.From)
	}
	src, ok := val.(string)
	if !ok {
		return fmt.Errorf("regex_replace: value at 'key' is not a string")
	}

	re, err := regexp.Compile(instr.Match)
	if err != nil {
		return fmt.Errorf("regex_replace: invalid regex: %w", err)
	}
	replacement := instr.With
	result := re.ReplaceAllString(src, replacement)
	return SetValueAtPath(input, instr.To, result)
}
