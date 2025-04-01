package transform

import (
	"github.com/tommed/ducto-dsl/model"
)

type fakeOperator struct{}

func (f fakeOperator) Name() string {
	return "fake"
}

func (f fakeOperator) Validate(instr model.Instruction) error {
	return nil
}

func (f fakeOperator) Apply(_ *ExecutionContext, _ *Registry, _ map[string]interface{}, _ model.Instruction) error {
	return nil
}

type noNameOperator struct{}

func (n noNameOperator) Name() string {
	return ""
}

func (n noNameOperator) Validate(instr model.Instruction) error {
	return nil
}

func (n noNameOperator) Apply(_ *ExecutionContext, _ *Registry, _ map[string]interface{}, _ model.Instruction) error {
	return nil
}
