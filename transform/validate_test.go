package transform

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidateProgram_BadInstruction(t *testing.T) {
	r := NewRegistry()
	r.Register(&SetOperator{})
	err := ValidateProgram(r, &Program{
		Version: 1,
		OnError: "fail",
		Instructions: []Instruction{
			{
				Op: "set",
			},
		},
	})
	assert.Error(t, err)
}

func TestValidateProgram_BadSubInstruction(t *testing.T) {
	r := NewRegistry()
	r.Register(&SetOperator{})
	r.Register(&MapOperator{})
	err := ValidateProgram(r, &Program{
		Version: 1,
		OnError: "fail",
		Instructions: []Instruction{
			{
				Op:  "map",
				Key: "a",
				Then: []Instruction{
					{
						Op: "set",
					},
				},
			},
		},
	})
	assert.Error(t, err)
}
