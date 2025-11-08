package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContainsSubstring(t *testing.T) {
	tests := []struct {
		name     string
		str      string
		substr   string
		expected bool
	}{
		{
			name:     "Substring present",
			str:      "Hello, world!",
			substr:   "world",
			expected: true,
		},
		{
			name:     "Substring absent",
			str:      "Hello, world!",
			substr:   "gopher",
			expected: false,
		},
		{
			name:     "Empty substring",
			str:      "Hello, world!",
			substr:   "",
			expected: true,
		},
		{
			name:     "Empty string",
			str:      "",
			substr:   "gopher",
			expected: false,
		},
		{
			name:     "Both empty",
			str:      "",
			substr:   "",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ContainsSubstring(tt.str, tt.substr)
			assert.Equal(t, tt.expected, result)
		})
	}
}