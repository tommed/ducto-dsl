package transform

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tommed/dsl-transformer/internal/model"
)

func TestOperators(t *testing.T) {
	cases := []struct {
		name           string
		op             Operator
		otherOpsNeeded []Operator
		input          map[string]interface{}
		instr          model.Instruction
		expect         map[string]interface{}
		expectedName   string
		wantErr        bool
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
		{
			name:         "noop",
			op:           &NoOperation{},
			input:        map[string]interface{}{"a": 123},
			instr:        model.Instruction{},
			expectedName: "noop",
			expect:       map[string]interface{}{"a": 123},
		},
		{
			name:           "map_basic",
			op:             &MapOperator{},
			otherOpsNeeded: []Operator{&SetOperator{}},
			input: map[string]interface{}{"a": []interface{}{
				map[string]interface{}{
					"foo": 1,
				},
				map[string]interface{}{
					"bar": 1,
				},
			}},
			instr: model.Instruction{
				Op:  "map",
				Key: "a",
				Then: []model.Instruction{
					{
						Op:    "set",
						Key:   "status",
						Value: "ok",
					},
				},
			},
			expectedName: "map",
			expect: map[string]interface{}{"a": []interface{}{
				map[string]interface{}{
					"foo":    1,
					"status": "ok",
				},
				map[string]interface{}{
					"bar":    1,
					"status": "ok",
				},
			}},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			input := make(map[string]interface{})
			for k, v := range tt.input {
				input[k] = v
			}

			exec := NewExecutionContext(context.Background(), "fail")
			r := NewRegistry()
			r.Register(tt.op)
			for _, op := range tt.otherOpsNeeded {
				r.Register(op)
			}

			assert.Equal(t, tt.expectedName, tt.op.Name())
			err := tt.op.Apply(exec, r, input, tt.instr)
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
			exec := NewExecutionContext(context.Background(), "fail")
			r := NewRegistry() // only needed for nested ops like `if`, `map` etc.
			r.Register(tt.op)

			err := tt.op.Apply(exec, r, map[string]interface{}{}, tt.instruction)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.wantErr)
		})
	}
}
