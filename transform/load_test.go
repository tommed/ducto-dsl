package transform

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/tommed/ducto-dsl/model"
	"testing"
)

func TestLoadProgram(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    *model.Program
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			args: args{
				path: "../examples/simplest.json",
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return err == nil
			},
			want: &model.Program{
				Version: 1,
				OnError: "fail",
				Instructions: []model.Instruction{
					{
						Op:    "set",
						Key:   "greeting",
						Value: "hello world",
					},
				},
			},
		},
		{
			name: "invalid path",
			args: args{
				path: "../examples/invalid-path.json",
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return err.Error() == "failed to read program: open ../examples/invalid-path.json: no such file or directory"
			},
		},
		{
			name: "invalid json",
			args: args{
				path: "../test/data/invalid-json.json",
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return err.Error() == "failed to parse program: unexpected end of JSON input"
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LoadProgram(tt.args.path)
			if !tt.wantErr(t, err, fmt.Sprintf("LoadProgram(%v)", tt.args.path)) {
				return
			}
			assert.Equalf(t, tt.want, got, "LoadProgram(%v)", tt.args.path)
		})
	}
}
