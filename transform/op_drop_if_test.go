package transform

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDropIfOperator(t *testing.T) {
	prog := &Program{
		Version: 1,
		Instructions: []Instruction{
			{Op: "drop_if", Key: "_flags.drop"},
			{Op: "set", Key: "should_not_run", Value: true},
		},
	}

	txf := New()

	input := map[string]interface{}{
		"_flags": map[string]interface{}{
			"drop": true,
		},
	}

	output, err := txf.Apply(context.Background(), input, prog)
	assert.Nil(t, output)
	assert.Nil(t, err)
}
