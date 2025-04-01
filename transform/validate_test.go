package transform

import (
	"github.com/stretchr/testify/assert"
	"github.com/tommed/ducto-dsl/model"
	"testing"
)

func TestValidateProgram_BadInstruction(t *testing.T) {
	r := NewRegistry()
	r.Register(&SetOperator{})
	err := ValidateProgram(r, &model.Program{
		Version: 1,
		OnError: "fail",
		Instructions: []model.Instruction{
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
	err := ValidateProgram(r, &model.Program{
		Version: 1,
		OnError: "fail",
		Instructions: []model.Instruction{
			{
				Op:  "map",
				Key: "a",
				Then: []model.Instruction{
					{
						Op: "set",
					},
				},
			},
		},
	})
	assert.Error(t, err)
}
