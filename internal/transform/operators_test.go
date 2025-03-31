package transform

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tommed/dsl-transformer/internal/model"
)

func TestOperators(t *testing.T) {
	cases := []struct {
		name         string
		op           Operator
		input        map[string]interface{}
		instr        model.Instruction
		expect       map[string]interface{}
		expectedName string
		wantErr      bool
	}{
		{
			name:         "set_basic",
			op:           &SetOperator{},
			input:        map[string]interface{}{},
			instr:        model.Instruction{Key: "foo", Value: "bar"},
			expectedName: "set",
			expect:       map[string]interface{}{"foo": "bar"},
		},
		{
			name:         "delete_basic",
			op:           &DeleteOperator{},
			input:        map[string]interface{}{"foo": "bar"},
			instr:        model.Instruction{Key: "foo"},
			expectedName: "delete",
			expect:       map[string]interface{}{},
		},
		{
			name:         "copy_basic",
			op:           &CopyOperator{},
			input:        map[string]interface{}{"a": 123},
			instr:        model.Instruction{From: "a", To: "b"},
			expectedName: "copy",
			expect:       map[string]interface{}{"a": 123, "b": 123},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			input := make(map[string]interface{})
			for k, v := range tt.input {
				input[k] = v
			}
			assert.Equal(t, tt.expectedName, tt.op.Name())
			err := tt.op.Apply(context.Background(), input, tt.instr)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expect, input)
			}
		})
	}
}

func TestOperators_ValidationErrors(t *testing.T) {
	cases := []struct {
		name        string
		op          Operator
		instruction model.Instruction
		wantErr     string
	}{
		{
			name:        "set_no_key",
			op:          &SetOperator{},
			instruction: model.Instruction{},
			wantErr:     "missing key",
		},
		{
			name:        "delete_no_key",
			op:          &DeleteOperator{},
			instruction: model.Instruction{},
			wantErr:     "missing key",
		},
		{
			name:        "copy_no_from",
			op:          &CopyOperator{},
			instruction: model.Instruction{},
			wantErr:     "op missing or invalid from",
		},
		{
			name:        "copy_no_to",
			op:          &CopyOperator{},
			instruction: model.Instruction{From: "a"},
			wantErr:     "op missing or invalid to",
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.op.Apply(context.Background(), map[string]interface{}{}, tt.instruction)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.wantErr)
		})
	}
}
