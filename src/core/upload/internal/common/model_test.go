//go:build unit
// +build unit

package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewOpen(t *testing.T) {
	result := NewOpen()
	assert.Equal(t, OPEN, result.Type)
	assert.Equal(t, "{", result.Text)
}

func TestNewClose(t *testing.T) {
	result := NewClose()
	assert.Equal(t, CLOSE, result.Type)
	assert.Equal(t, "}", result.Text)
}

func TestNewSeparator(t *testing.T) {
	result := NewSeparator()
	assert.Equal(t, SEPARATOR, result.Type)
	assert.Equal(t, "|", result.Text)
}

func TestNewType(t *testing.T) {
	text := "type"
	result := NewClass(text)
	assert.Equal(t, CLASS, result.Type)
	assert.Equal(t, text, result.Text)
}

func TestNewLocation(t *testing.T) {
	text := "123.456"
	result := NewParam(text)
	assert.Equal(t, PARAM, result.Type)
	assert.Equal(t, text, result.Text)
}

func TestNewText(t *testing.T) {
	text := "hello world"
	result := NewText(text)
	assert.Equal(t, TEXT, result.Type)
	assert.Equal(t, text, result.Text)
}
