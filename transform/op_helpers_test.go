package transform

type fakeOperator struct{}

func (f fakeOperator) Name() string {
	return "fake"
}

func (f fakeOperator) Validate(instr Instruction) error {
	return nil
}

func (f fakeOperator) Apply(_ *ExecutionContext, _ *Registry, _ map[string]interface{}, _ Instruction) error {
	return nil
}

type noNameOperator struct{}

func (n noNameOperator) Name() string {
	return ""
}

func (n noNameOperator) Validate(instr Instruction) error {
	return nil
}

func (n noNameOperator) Apply(_ *ExecutionContext, _ *Registry, _ map[string]interface{}, _ Instruction) error {
	return nil
}
