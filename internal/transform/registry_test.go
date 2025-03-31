package transform

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/tommed/dsl-transformer/internal/model"
	"testing"
)

func TestRegistry_Register_DuplicatePanics(t *testing.T) {
	reg := NewRegistry()
	dummy := &NoOperation{}

	// Register once (should not panic)
	reg.Register(dummy)

	// Register again (should panic)
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic due to duplicate registration, got none")
		} else {
			t.Logf("panic caught as expected: %v", r)
		}
	}()

	reg.Register(dummy) // <- this triggers panic
}

func TestRegistry_Register_NoNameOperatorFails(t *testing.T) {
	reg := NewRegistry()
	badOp := &noNameOperator{}

	// Register again (should panic)
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic due to duplicate registration, got none")
		} else {
			t.Logf("panic caught as expected: %v", r)
		}
	}()

	reg.Register(badOp) // <- this triggers panic
}

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
		wantErrors  []string
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
			wantErrors:  []string{}, // no errors
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
			wantErrors:  []string{}, // no errors in list as using 'fail'
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
			wantErrors:  []string{}, // no errors in list as using 'fail'
		},
		{
			name: "invalid op + error",
			args: args{
				onError: "capture",
				input:   map[string]interface{}{},
				instr: model.Instruction{
					Op: "invalid_op",
				},
			},
			wantSuccess: false, // because this is validated not run on `apply`
			wantErrors:  []string{},
		},
		{
			name: "invalid at runtime (Apply)",
			args: args{
				onError: "fail",
				input:   map[string]interface{}{},
				instr: model.Instruction{
					Op:   "map",
					Key:  "a", // doesn't exist on input
					Then: []model.Instruction{{Op: "noop"}},
				},
			},
			wantSuccess: false, // because this is validated not run on `apply`
			wantErrors:  []string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exec := NewExecutionContext(context.Background(), tt.args.onError)
			r := NewDefaultRegistry(&fakeOperator{})
			err := r.Find("set").Validate(tt.args.instr)
			wantSuccess := false
			if err == nil {
				wantSuccess = r.Apply(exec, r, tt.args.input, tt.args.instr)
			}
			assert.Equalf(t, tt.wantSuccess, wantSuccess, "Apply(ctx, %v, %v)", tt.args.input, tt.args.instr)
			assert.Equal(t, tt.wantErrors, exec.Errors)
		})
	}
}
