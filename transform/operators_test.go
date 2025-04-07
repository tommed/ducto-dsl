package transform

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOperators(t *testing.T) {
	cases := []struct {
		name           string
		op             Operator
		otherOpsNeeded []Operator
		input          map[string]interface{}
		instr          Instruction
		expect         map[string]interface{}
		expectedName   string
		wantErr        bool
	}{
		{
			name:         "set_basic",
			op:           &SetOperator{},
			input:        map[string]interface{}{},
			instr:        Instruction{Key: "foo", Value: "bar"},
			expectedName: "set",
			expect:       map[string]interface{}{"foo": "bar"},
		},
		{
			name:         "delete_basic",
			op:           &DeleteOperator{},
			input:        map[string]interface{}{"foo": "bar"},
			instr:        Instruction{Key: "foo"},
			expectedName: "delete",
			expect:       map[string]interface{}{},
		},
		{
			name:         "delete_at_path",
			op:           &DeleteOperator{},
			input:        map[string]interface{}{"a": map[string]interface{}{"b": "c", "d": "e"}},
			instr:        Instruction{Key: "a.b"},
			expectedName: "delete",
			expect:       map[string]interface{}{"a": map[string]interface{}{"d": "e"}},
		},
		{
			name:         "delete_at_bad_path",
			op:           &DeleteOperator{},
			input:        map[string]interface{}{"a": map[string]interface{}{"d": "e"}},
			instr:        Instruction{Key: "a.b"},
			expectedName: "delete",
			expect:       map[string]interface{}{"a": map[string]interface{}{"d": "e"}},
		},
		{
			name:         "delete_at_bad_path_deep",
			op:           &DeleteOperator{},
			input:        map[string]interface{}{"a": map[string]interface{}{"d": "e"}},
			instr:        Instruction{Key: "a.d.c"},
			expectedName: "delete",
			expect:       map[string]interface{}{"a": map[string]interface{}{"d": "e"}},
		},
		{
			name:         "copy_basic",
			op:           &CopyOperator{},
			input:        map[string]interface{}{"a": 123},
			instr:        Instruction{From: "a", To: "b"},
			expectedName: "copy",
			expect:       map[string]interface{}{"a": 123, "b": 123},
		},
		{
			name:         "noop",
			op:           &NoOperation{},
			input:        map[string]interface{}{"a": 123},
			instr:        Instruction{},
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
			instr: Instruction{
				Op:  "map",
				Key: "a",
				Then: []Instruction{
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

			err := tt.op.Validate(tt.instr)
			if err == nil {
				err = tt.op.Apply(exec, r, input, tt.instr)
			}

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
		instruction Instruction
		wantErr     string
	}{
		{
			name:        "set_no_key",
			op:          &SetOperator{},
			instruction: Instruction{},
			wantErr:     "operator missing 'key'",
		},
		{
			name:        "delete_no_key",
			op:          &DeleteOperator{},
			instruction: Instruction{},
			wantErr:     "operator missing 'key'",
		},
		{
			name:        "copy_no_from",
			op:          &CopyOperator{},
			instruction: Instruction{},
			wantErr:     "op missing or invalid from",
		},
		{
			name:        "copy_no_to",
			op:          &CopyOperator{},
			instruction: Instruction{From: "a"},
			wantErr:     "op missing or invalid to",
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			exec := NewExecutionContext(context.Background(), "fail")
			r := NewRegistry() // only needed for nested ops like `if`, `map` etc.
			r.Register(tt.op)

			// Determine which error we're looking for:
			// Validation or Application errors
			err := tt.op.Validate(tt.instruction)
			if err == nil {
				err = tt.op.Apply(exec, r, map[string]interface{}{}, tt.instruction)
			}

			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.wantErr)
		})
	}
}
