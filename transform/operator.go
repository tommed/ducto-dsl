package transform

type Operator interface {
	Name() string
	Validate(instr Instruction) error
	Apply(ctx *ExecutionContext, reg *Registry, input map[string]interface{}, instr Instruction) error
}
