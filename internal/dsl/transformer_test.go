package dsl

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransformer_Apply_NoOp(t *testing.T) {
	// Assemble
	tr := New()

	input := map[string]interface{}{"foo": "bar"}

	prog := &Program{
		Instructions: []Instruction{
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
