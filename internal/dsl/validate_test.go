package dsl

import (
	"github.com/stretchr/testify/assert"
	"github.com/tommed/ducto-dsl/model"
	transform2 "github.com/tommed/ducto-dsl/transform"
	"testing"
)

func TestValidateProgram_BadInstruction(t *testing.T) {
	r := transform2.NewRegistry()
	r.Register(&transform2.SetOperator{})
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
	r := transform2.NewRegistry()
	r.Register(&transform2.SetOperator{})
	r.Register(&transform2.MapOperator{})
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
