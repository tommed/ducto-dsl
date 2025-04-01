package transform

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConditions(t *testing.T) {
	input := map[string]interface{}{
		"level": "warn",
		"count": 5,
		"tags":  []string{"a", "b"},
	}

	tests := []struct {
		name      string
		condition map[string]interface{}
		want      bool
	}{
		{
			name:      "exists true",
			condition: map[string]interface{}{"exists": "level"},
			want:      true,
		},
		{
			name:      "exists false",
			condition: map[string]interface{}{"exists": "missing"},
			want:      false,
		},
		{
			name: "equals string match",
			condition: map[string]interface{}{
				"equals": map[string]interface{}{
					"key":   "level",
					"value": "warn",
				},
			},
			want: true,
		},
		{
			name: "equals string no match",
			condition: map[string]interface{}{
				"equals": map[string]interface{}{
					"key":   "level",
					"value": "error",
				},
			},
			want: false,
		},
		{
			name: "or true",
			condition: map[string]interface{}{
				"or": []interface{}{
					map[string]interface{}{"equals": map[string]interface{}{"key": "level", "value": "error"}},
					map[string]interface{}{"equals": map[string]interface{}{"key": "level", "value": "warn"}},
				},
			},
			want: true,
		},
		{
			name: "or false",
			condition: map[string]interface{}{
				"or": []interface{}{
					map[string]interface{}{"equals": map[string]interface{}{"key": "level", "value": "error"}},
				},
			},
			want: false,
		},
		{
			name: "and true",
			condition: map[string]interface{}{
				"and": []interface{}{
					map[string]interface{}{"exists": "level"},
					map[string]interface{}{"exists": "count"},
				},
			},
			want: true,
		},
		{
			name: "and false",
			condition: map[string]interface{}{
				"and": []interface{}{
					map[string]interface{}{"exists": "level"},
					map[string]interface{}{"exists": "missing"},
				},
			},
			want: false,
		},
		{
			name: "invalid condition",
			condition: map[string]interface{}{
				"invalid": []interface{}{},
			},
			want: false,
		},
		{
			name: "exists not a string",
			condition: map[string]interface{}{
				"exists": []string{"invalid", "array", "of", "strings"},
			},
			want: false,
		},
		{
			name: "equals not a map",
			condition: map[string]interface{}{
				"equals": []string{"invalid", "array", "of", "strings"},
			},
			want: false,
		},
		{
			name: "equals without key",
			condition: map[string]interface{}{
				"equals": map[string]interface{}{"needs-key": "will-fail"},
			},
			want: false,
		},
		{
			name: "equals without value",
			condition: map[string]interface{}{
				"equals": map[string]interface{}{"key": "a", "no-value": "will-fail"},
			},
			want: false,
		},
		{
			name: "no matching key in input",
			condition: map[string]interface{}{
				"equals": map[string]interface{}{"key": "not-a-real-field", "value": "123"},
			},
			want: false,
		},
		{
			name: "or not an array",
			condition: map[string]interface{}{
				"or": "bad-value",
			},
			want: false,
		},
		{
			name: "and not an array",
			condition: map[string]interface{}{
				"and": "bad-value",
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := evaluateCondition(input, tt.condition)
			assert.Equal(t, tt.want, got)
		})
	}
}
