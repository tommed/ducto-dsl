package dsl

import (
	"context"
	"github.com/tommed/dsl-transformer/internal/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransformer_Apply_NoOp(t *testing.T) {
	// Assemble
	tr := New()

	input := map[string]interface{}{"foo": "bar"}

	prog := &model.Program{
		Instructions: []model.Instruction{
			{Op: "set", Key: "greeting", Value: "hello world"},
		},
	}

	// Act
	out, err := tr.Apply(context.Background(), input, prog)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "bar", out["foo"])
	assert.Equal(t, "hello world", out["greeting"])
}

func TestTransformer_Apply_InvalidOp(t *testing.T) {
	// Assemble
	tr := New()
	input := map[string]interface{}{"foo": "bar"}

	prog := &model.Program{
		OnError: "fail",
		Instructions: []model.Instruction{
			{Op: "invalid-op", Key: "foo", Value: 2},
		},
	}

	// Act
	_, err := tr.Apply(context.Background(), input, prog)

	// Assert
	assert.Error(t, err)
}

func TestTransformer_Apply_ErrorsReturned(t *testing.T) {
	// Assemble
	tr := New()
	input := map[string]interface{}{}

	prog := &model.Program{
		OnError: "error",
		Instructions: []model.Instruction{
			{Op: ""},
		},
	}

	// Act
	output, err := tr.Apply(context.Background(), input, prog)

	// Assert
	assert.NoError(t, err) // should have been ignored due to OnError value
	errorList, ok := output["@dsl_errors"].([]error)
	assert.True(t, ok)
	assert.Len(t, errorList, 1)
}
