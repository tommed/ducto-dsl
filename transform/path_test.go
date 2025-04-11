package transform

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCoerceToArray(t *testing.T) {
	cases := []struct {
		name      string
		in        interface{}
		expectLen int
		ok        bool
	}{
		{"nil", nil, 0, false},
		{"[]interface{}", []interface{}{1, 2, 3}, 3, true},
		{"[]map[string]interface{}", []map[string]interface{}{{"x": 1}, {"y": 2}}, 2, true},
		{"[]string via reflect", []string{"a", "b"}, 2, true},
		{"int (invalid)", 42, 0, false},
		{"string (invalid)", "not an array", 0, false},
		{"map (invalid)", map[string]interface{}{"a": 1}, 0, false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			out, ok := CoerceToArray(tc.in)
			assert.Equal(t, tc.ok, ok)
			if ok {
				assert.Len(t, out, tc.expectLen)
			}
		})
	}
}

func TestCoerceToArray_FastPaths(t *testing.T) {
	t.Run("[]interface{} is returned directly", func(t *testing.T) {
		input := []interface{}{1, 2, 3}
		out, ok := CoerceToArray(input)
		assert.True(t, ok)
		assert.Len(t, out, 3)
	})

	t.Run("[]map[string]interface{} is converted to []interface{}", func(t *testing.T) {
		input := []map[string]interface{}{
			{"a": 1},
			{"b": 2},
		}
		out, ok := CoerceToArray(input)
		assert.True(t, ok)
		assert.Len(t, out, 2)
	})
}

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
