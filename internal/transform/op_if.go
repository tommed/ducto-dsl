package transform

import (
	"errors"
	"github.com/tommed/ducto-dsl/internal/model"
)

type IfOperator struct{}

func (o *IfOperator) Name() string { return "if" }

func (o *IfOperator) Validate(instr model.Instruction) error {
	if instr.Condition == nil {
		return errors.New("if operator missing 'condition'")
	}
	if len(instr.Then) == 0 {
		return errors.New("if operator missing 'then' instructions")
	}
	if err := validateConditions(instr.Condition); err != nil {
		return err
	}
	return nil
}

func (o *IfOperator) Apply(ctx *ExecutionContext, reg *Registry, input map[string]interface{}, instr model.Instruction) error {
	conditionResult := evaluateCondition(input, instr.Condition)
	if instr.Not {
		conditionResult = !conditionResult
	}
	if conditionResult {
		for _, sub := range instr.Then {
			op := reg.Find(sub.Op)
			if err := op.Apply(ctx, reg, input, sub); err != nil {
				return err
			}
		}
	}
	return nil
}
