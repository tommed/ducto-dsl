package dsl

import (
	"github.com/stretchr/testify/assert"
	"github.com/tommed/dsl-transformer/internal/model"
	"github.com/tommed/dsl-transformer/internal/transform"
	"testing"
)

func TestValidateProgram_BadInstruction(t *testing.T) {
	r := transform.NewRegistry()
	r.Register(&transform.SetOperator{})
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
	r := transform.NewRegistry()
	r.Register(&transform.SetOperator{})
	r.Register(&transform.MapOperator{})
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
