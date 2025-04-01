package transform

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tommed/ducto-dsl/internal/model"
)

func TestIfOperator_Validate(t *testing.T) {
	op := &IfOperator{}
	type args struct {
		instr model.Instruction
	}
	var tests = []struct {
		name string
		args args
		want error
	}{
		{
			name: "success",
			args: args{
				instr: model.Instruction{
					Op:   "if",
					Then: []model.Instruction{{Op: "noop"}},
					Condition: map[string]interface{}{
						"equals": map[string]interface{}{
							"key":   "a",
							"value": "hello",
						},
					},
				},
			},
		},
		{
			name: "no condition",
			args: args{
				instr: model.Instruction{
					Op:   "if",
					Then: []model.Instruction{{Op: "noop"}},
				},
			},
			want: errors.New("if operator missing 'condition'"),
		},
		{
			name: "no then",
			args: args{
				instr: model.Instruction{
					Op: "if",
					Condition: map[string]interface{}{
						"equals": map[string]interface{}{
							"key":   "a",
							"value": "hello",
						},
					},
				},
			},
			want: errors.New("if operator missing 'then' instructions"),
		},
		{
			name: "empty conditions",
			args: args{
				instr: model.Instruction{
					Op:        "if",
					Then:      []model.Instruction{{Op: "noop"}},
					Condition: map[string]interface{}{},
				},
			},
			want: errors.New("no conditions defined"),
		},
		{
			name: "multiple conditions",
			args: args{
				instr: model.Instruction{
					Op:   "if",
					Then: []model.Instruction{{Op: "noop"}},
					Condition: map[string]interface{}{
						"equals": "",
						"exists": "",
					},
				},
			},
			want: errors.New("only one condition type is allowed per condition block, got: map[equals: exists:]"),
		},
		{
			name: "invalid condition",
			args: args{
				instr: model.Instruction{
					Op:   "if",
					Then: []model.Instruction{{Op: "noop"}},
					Condition: map[string]interface{}{
						"invalid": "",
					},
				},
			},
			want: errors.New(`unknown condition "invalid"`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := op.Validate(tt.args.instr)
			assert.Equal(t, tt.want, err)
		})
	}
}

func TestIfOperator_Apply(t *testing.T) {
	op := &IfOperator{}
	reg := NewDefaultRegistry()

	tests := []struct {
		name     string
		input    map[string]interface{}
		instr    model.Instruction
		expected map[string]interface{}
		wantErr  error
	}{
		{
			name:  "condition matches - executes then block",
			input: map[string]interface{}{"foo": "bar"},
			instr: model.Instruction{
				Op: "if",
				Condition: map[string]interface{}{
					"exists": "foo",
				},
				Then: []model.Instruction{
					{
						Op:    "set",
						Key:   "ran",
						Value: true,
					},
				},
			},
			expected: map[string]interface{}{"foo": "bar", "ran": true},
		},
		{
			name:  "condition does not match - no then executed",
			input: map[string]interface{}{"baz": "qux"},
			instr: model.Instruction{
				Op: "if",
				Condition: map[string]interface{}{
					"exists": "foo",
				},
				Then: []model.Instruction{
					{
						Op:    "set",
						Key:   "ran",
						Value: true,
					},
				},
			},
			expected: map[string]interface{}{"baz": "qux"},
		},
		{
			name:  "sub failure",
			input: map[string]interface{}{"baz": "qux"},
			instr: model.Instruction{
				Op: "if",
				Condition: map[string]interface{}{
					"exists": "baz",
				},
				Then: []model.Instruction{{Op: "fail", Value: "intention failure"}},
			},
			wantErr: errors.New("intention failure"),
		},
		{
			name:  "condition matches but negated - should skip",
			input: map[string]interface{}{"foo": "bar"},
			instr: model.Instruction{
				Op: "if",
				Condition: map[string]interface{}{
					"exists": "foo",
				},
				Not: true,
				Then: []model.Instruction{
					{
						Op:    "set",
						Key:   "ran",
						Value: true,
					},
				},
			},
			expected: map[string]interface{}{"foo": "bar"},
		},
		{
			name:  "condition does not match but negated - should run",
			input: map[string]interface{}{"baz": "qux"},
			instr: model.Instruction{
				Op: "if",
				Condition: map[string]interface{}{
					"exists": "foo",
				},
				Not: true,
				Then: []model.Instruction{
					{
						Op:    "set",
						Key:   "ran",
						Value: true,
					},
				},
			},
			expected: map[string]interface{}{"baz": "qux", "ran": true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := NewExecutionContext(context.Background(), "fail")
			require.NoError(t, op.Validate(tt.instr))
			err := op.Apply(ctx, reg, tt.input, tt.instr)
			if tt.wantErr != nil {
				assert.EqualError(t, err, tt.wantErr.Error())
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, tt.input)
			}
		})
	}
}
