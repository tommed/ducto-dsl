package transform

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/tommed/ducto-dsl/internal/model"
	"testing"
)

func TestMapOperator_Apply(t *testing.T) {
	type args struct {
		instr model.Instruction
	}
	tests := []struct {
		name         string
		args         args
		wantErr      error
		wantContains string
	}{
		{
			name: "no key",
			args: args{
				instr: model.Instruction{
					Op:  "map",
					Key: "",
				},
			},
			wantErr: errors.New("map operator requires 'key' field"),
		},
		{
			name: "no then",
			args: args{
				instr: model.Instruction{
					Op:   "map",
					Key:  "a",
					Then: []model.Instruction{},
				},
			},
			wantErr: errors.New("map operator requires at least one instruction in 'then'"),
		},
		{
			name: "key not found",
			args: args{
				instr: model.Instruction{
					Op:   "map",
					Key:  "d",
					Then: []model.Instruction{{Op: "set", Key: "e", Value: "e"}},
				},
			},
			wantErr: errors.New("map operator: key \"d\" not found in input"),
		},
		{
			name: "key not array",
			args: args{
				instr: model.Instruction{
					Op:   "map",
					Key:  "b",
					Then: []model.Instruction{{Op: "set", Key: "e", Value: "e"}},
				},
			},
			wantErr: errors.New("map operator: input[\"b\"] is not an array"),
		},
		{
			name: "key not object",
			args: args{
				instr: model.Instruction{
					Op:   "map",
					Key:  "f",
					Then: []model.Instruction{{Op: "set", Key: "e", Value: "e"}},
				},
			},
			wantErr: errors.New("map operator: input[\"f\"] is not an array"),
		},
		{
			name: "key not object",
			args: args{
				instr: model.Instruction{
					Op:   "map",
					Key:  "f",
					Then: []model.Instruction{{Op: "fail", Value: "fail_on_purpose"}},
				},
			},
			wantErr: errors.New("map operator: input[\"f\"] is not an array"),
		},
		{
			name: "sub failure",
			args: args{
				instr: model.Instruction{
					Op:   "map",
					Key:  "a",
					Then: []model.Instruction{{Op: "fail", Value: "fail_on_purpose"}},
				},
			},
			wantErr: errors.New("map operator: sub-instruction failed at index 0"),
		},
		{
			name: "wrong array type",
			args: args{
				instr: model.Instruction{
					Op:   "map",
					Key:  "g",
					Then: []model.Instruction{{Op: "noop"}},
				},
			},
			wantErr: errors.New("map operator: array item at index 0 is not an object"),
		},
	}
	for _, tt := range tests {

		// Assemble
		input := map[string]interface{}{"a": []interface{}{
			map[string]interface{}{
				"foo": 1,
			},
			map[string]interface{}{
				"bar": 1,
			},
		},
			"b": "bad-prop",
			"f": []map[string]interface{}{}, // not []interface{} which is needed
			"g": []interface{}{
				"test",
				1,
			},
		}

		t.Run(tt.name, func(t *testing.T) {
			exec := NewExecutionContext(context.Background(), "fail")
			mapOp := &MapOperator{}
			r := NewRegistry()
			r.Register(mapOp)
			r.Register(&SetOperator{})
			r.Register(&FailOperator{})
			err := mapOp.Validate(tt.args.instr)
			if err == nil {
				err = mapOp.Apply(exec, r, input, tt.args.instr)
			}
			assert.Equal(t, tt.wantErr, err)
			if tt.wantErr == nil && err != nil {
				assert.NoError(t, err)
			}
			if err == nil {
				data, _ := json.Marshal(input)
				assert.Contains(t, string(data), tt.wantContains)
			}
		})
	}
}
