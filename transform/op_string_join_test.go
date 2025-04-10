package transform

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStringJoinOperator_Validate(t *testing.T) {
	op := &StringJoinOperator{}

	tests := []struct {
		name    string
		instr   Instruction
		wantErr string
	}{
		{
			name:    "missing from",
			instr:   Instruction{To: "out", Value: ","},
			wantErr: "string_join operator requires a 'from' field",
		},
		{
			name:    "missing to",
			instr:   Instruction{From: "in", Value: ","},
			wantErr: "string_join operator requires a 'to' field",
		},
		{
			name:    "invalid value type",
			instr:   Instruction{From: "in", To: "out", Value: 123},
			wantErr: "string_join operator requires joining value",
		},
		{
			name:  "valid instruction",
			instr: Instruction{From: "in", To: "out", Value: ","},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := op.Validate(tt.instr)
			if tt.wantErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestStringJoinOperator_Apply(t *testing.T) {
	op := &StringJoinOperator{}

	tests := []struct {
		name        string
		input       map[string]interface{}
		instr       Instruction
		expected    map[string]interface{}
		expectError bool
	}{
		{
			name: "join []interface{}",
			input: map[string]interface{}{
				"arr": []interface{}{1, "two", true},
			},
			instr: Instruction{From: "arr", To: "joined", Value: "-"},
			expected: map[string]interface{}{
				"arr":    []interface{}{1, "two", true},
				"joined": "1-two-true",
			},
		},
		{
			name: "join []string",
			input: map[string]interface{}{
				"arr": []string{"a", "b", "c"},
			},
			instr: Instruction{From: "arr", To: "result", Value: ":"},
			expected: map[string]interface{}{
				"arr":    []string{"a", "b", "c"},
				"result": "a:b:c",
			},
		},
		{
			name:        "invalid from path",
			input:       map[string]interface{}{},
			instr:       Instruction{From: "missing", To: "out", Value: ","},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := op.Apply(nil, nil, tt.input, tt.instr)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, tt.input)
			}
		})
	}
}
