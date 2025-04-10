package transform

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSetValueAtPath_CreatesNestedStructure(t *testing.T) {
	data := map[string]interface{}{
		"existing": "keep me",
	}

	err := SetValueAtPath(data, "a.b.c", "nested value")

	assert.NoError(t, err)
	assert.Equal(t, "keep me", data["existing"])

	// Deep assertion
	a, ok := data["a"].(map[string]interface{})
	assert.True(t, ok)

	b, ok := a["b"].(map[string]interface{})
	assert.True(t, ok)

	assert.Equal(t, "nested value", b["c"])
}

func TestCoerceToMap(t *testing.T) {
	t.Run("map[string]interface{}", func(t *testing.T) {
		input := map[string]interface{}{"foo": 42}
		out, ok := CoerceToMap(input)
		assert.True(t, ok)
		assert.Equal(t, 42, out["foo"])
	})

	t.Run("map[string]string via reflection", func(t *testing.T) {
		input := map[string]string{"hello": "world"}
		out, ok := CoerceToMap(input)
		assert.True(t, ok)
		assert.Equal(t, "world", out["hello"])
	})

	t.Run("non-map input", func(t *testing.T) {
		input := "not a map"
		out, ok := CoerceToMap(input)
		assert.False(t, ok)
		assert.Nil(t, out)
	})

	t.Run("map with non-string keys", func(t *testing.T) {
		input := map[int]string{1: "one", 2: "two"}
		out, ok := CoerceToMap(input)
		assert.False(t, ok)
		assert.Nil(t, out)
	})
}
