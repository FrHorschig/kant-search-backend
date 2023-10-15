//go:build unit
// +build unit

package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewOpen(t *testing.T) {
	result := newOpen()
	assert.Equal(t, OPEN, result.Type)
	assert.Equal(t, "{", result.Text)
}

func TestNewClose(t *testing.T) {
	result := newClose()
	assert.Equal(t, CLOSE, result.Type)
	assert.Equal(t, "}", result.Text)
}

func TestNewSeparator(t *testing.T) {
	result := newSeparator()
	assert.Equal(t, SEPARATOR, result.Type)
	assert.Equal(t, "|", result.Text)
}

func TestNewType(t *testing.T) {
	text := "type"
	result := newClass(text)
	assert.Equal(t, CLASS, result.Type)
	assert.Equal(t, text, result.Text)
}

func TestNewLocation(t *testing.T) {
	text := "123.456"
	result := newParam(text)
	assert.Equal(t, PARAM, result.Type)
	assert.Equal(t, text, result.Text)
}

func TestNewText(t *testing.T) {
	text := "hello world"
	result := newText(text)
	assert.Equal(t, TEXT, result.Type)
	assert.Equal(t, text, result.Text)
}
