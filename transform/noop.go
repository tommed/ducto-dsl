package transform

// NoOperation is our nil implementation, it literally does nothing and never fails
type NoOperation struct{}

func (n NoOperation) Validate(instr Instruction) error {
	return nil
}

func (n NoOperation) Name() string {
	return "noop"
}

func (n NoOperation) Apply(_ *ExecutionContext, _ *Registry, _ map[string]interface{}, _ Instruction) error {
	return nil
}
