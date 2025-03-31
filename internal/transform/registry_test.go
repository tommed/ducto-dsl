package transform

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/tommed/dsl-transformer/internal/model"
	"testing"
)

func TestRegistry_Apply(t *testing.T) {
	type args struct {
		onError string
		input   map[string]interface{}
		instr   model.Instruction
	}
	tests := []struct {
		name        string
		args        args
		wantSuccess bool
		wantErrors  []error
	}{
		{
			name: "set op",
			args: args{
				onError: "fail",
				input:   map[string]interface{}{},
				instr: model.Instruction{
					Op:    "set",
					Key:   "a",
					Value: 123,
				},
			},
			wantSuccess: true,
			wantErrors:  make([]error, 0), // no errors
		},
		{
			name: "invalid op + fail",
			args: args{
				onError: "fail",
				input:   map[string]interface{}{},
				instr: model.Instruction{
					Op: "invalid_op",
				},
			},
			wantSuccess: false,
			wantErrors:  make([]error, 0), // no errors in list as using 'fail'
		},
		{
			name: "invalid apply + fail",
			args: args{
				onError: "fail",
				input:   map[string]interface{}{},
				instr: model.Instruction{
					Op: "set",
					// No Key or Value so invalid
				},
			},
			wantSuccess: false,
			wantErrors:  make([]error, 0), // no errors in list as using 'fail'
		},
		{
			name: "invalid op + error",
			args: args{
				onError: "error",
				input:   map[string]interface{}{},
				instr: model.Instruction{
					Op: "invalid_op",
				},
			},
			wantSuccess: true, // ignored and continued
			wantErrors:  []error{errors.New(`unknown op: "invalid_op"`)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exec := NewExecutionContext(context.Background(), tt.args.onError)
			r := NewRegistry()
			r.Register(&SetOperator{})
			assert.Equalf(t, tt.wantSuccess, r.Apply(exec, tt.args.input, tt.args.instr), "Apply(ctx, %v, %v)", tt.args.input, tt.args.instr)
			assert.Equal(t, tt.wantErrors, exec.Errors)
		})
	}
}
