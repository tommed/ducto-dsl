package transform

import (
	"errors"
	"fmt"
	"strings"
)

type StringJoinOperator struct{}

func (o *StringJoinOperator) Name() string { return "string_join" }

func (o *StringJoinOperator) Validate(instr Instruction) error {
	if instr.From == "" {
		return errors.New("string_join operator requires a 'from' field")
	}
	if instr.To == "" {
		return errors.New("string_join operator requires a 'to' field")
	}
	_, ok := instr.Value.(string)
	if !ok {
		return errors.New("string_join operator requires joining value")
	}
	return nil
}

func (o *StringJoinOperator) Apply(_ *ExecutionContext, _ *Registry, input map[string]interface{}, instr Instruction) error {
	arrayRaw, ok := GetValueAtPath(input, instr.From)
	if !ok {
		return fmt.Errorf("string_join from '%s' did not yield an array", instr.From)
	}

	var arrayToJoin = make([]string, 0)
	switch array := arrayRaw.(type) {
	case []interface{}:
		for _, v := range array {
			arrayToJoin = append(arrayToJoin, fmt.Sprintf("%v", v))
		}
	case []string:
		arrayToJoin = array
	}
	return SetValueAtPath(input, instr.To, strings.Join(arrayToJoin, instr.Value.(string)))
}
