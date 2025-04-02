package transform

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCoalesceOperator_Validate(t *testing.T) {
	op := &CoalesceOperator{}

	tests := []struct {
		name    string
		instr   Instruction
		wantErr bool
	}{
		{
			name:    "missing key",
			instr:   Instruction{Op: "coalesce", Value: "default"},
			wantErr: true,
		},
		{
			name:    "missing value",
			instr:   Instruction{Op: "coalesce", Key: "foo"},
			wantErr: true,
		},
		{
			name:    "valid coalesce",
			instr:   Instruction{Op: "coalesce", Key: "foo", Value: "default"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := op.Validate(tt.instr)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCoalesceOperator_Apply(t *testing.T) {
	op := &CoalesceOperator{}

	tests := []struct {
		name     string
		input    map[string]interface{}
		instr    Instruction
		expected map[string]interface{}
	}{
		{
			name:     "value missing - default applied",
			input:    map[string]interface{}{},
			instr:    Instruction{Op: "coalesce", Key: "foo", Value: "bar"},
			expected: map[string]interface{}{"foo": "bar"},
		},
		{
			name:     "value exists - default ignored",
			input:    map[string]interface{}{"foo": "baz"},
			instr:    Instruction{Op: "coalesce", Key: "foo", Value: "bar"},
			expected: map[string]interface{}{"foo": "baz"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := NewExecutionContext(context.Background(), "fail")
			require.NoError(t, op.Validate(tt.instr))
			err := op.Apply(ctx, nil, tt.input, tt.instr)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, tt.input)
		})
	}
}
