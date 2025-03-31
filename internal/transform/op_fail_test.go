package transform

import (
	"github.com/stretchr/testify/assert"
	"github.com/tommed/dsl-transformer/internal/model"
	"testing"
)

func TestFailOperator(t *testing.T) {
	op := &FailOperator{}
	err := op.Validate(model.Instruction{Op: "fail", Value: "Failed on purpose"})
	assert.NoError(t, err)
}
