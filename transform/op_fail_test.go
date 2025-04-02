package transform

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFailOperator(t *testing.T) {
	op := &FailOperator{}
	err := op.Validate(Instruction{Op: "fail", Value: "Failed on purpose"})
	assert.NoError(t, err)
}
