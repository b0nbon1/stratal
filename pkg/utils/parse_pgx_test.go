package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseUUID_Valid(t *testing.T) {
	validUUIDStr := "550e8400-e29b-41d4-a716-446655440000"

	parsedUUID, err := ParseUUID(validUUIDStr)
	require.NoError(t, err)
	assert.True(t, parsedUUID.Valid)
}

func TestParseUUID_Invalid(t *testing.T) {
	invalidUUIDStr := "invalid-uuid-string"

	parsedUUID, err := ParseUUID(invalidUUIDStr)
	require.Error(t, err)
	assert.False(t, parsedUUID.Valid)
}

func TestParseText(t *testing.T) {
	text := "Hello, World!"
	parsedText := ParseText(text)

	assert.True(t, parsedText.Valid)
	assert.Equal(t, text, parsedText.String)
}

func TestParseText_Empty(t *testing.T) {
	text := ""
	parsedText := ParseText(text)

	assert.True(t, parsedText.Valid)
	assert.Equal(t, text, parsedText.String)
}
