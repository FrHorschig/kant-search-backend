//go:build integration
// +build integration

package pyutil

import (
	"os"
	"testing"

	"github.com/frhorschig/kant-search-backend/common/model"
	"github.com/stretchr/testify/assert"
)

func TestSplitIntoSentences(t *testing.T) {
	pyUtil := NewPythonUtil()
	paragraphs := []model.Paragraph{
		{
			Id:   1,
			Text: "Das ist ein erster Satz. Das ist ein zweiter Satz! Ist das ein dritter Satz?",
		},
		{
			Id:   2,
			Text: "Das ist ein erster Satz. Das ist ein zweiter Satz ohne Punkt",
		},
		{
			Id:   3,
			Text: "Das ist ein erster Satz usw. der weiter geht. Das ist ein zweiter Satz.",
		},
		{
			Id:   4,
			Text: "Das ist ein erster Satz ... der weiter geht. Das ist ein zweiter Satz.",
		},
	}
	expectedResult := map[int32][]string{
		1: {"Das ist ein erster Satz.", "Das ist ein zweiter Satz!", "Ist das ein dritter Satz?"},
		2: {"Das ist ein erster Satz.", "Das ist ein zweiter Satz ohne Punkt"},
		3: {"Das ist ein erster Satz usw. der weiter geht.", "Das ist ein zweiter Satz."},
		4: {"Das ist ein erster Satz ... der weiter geht.", "Das ist ein zweiter Satz."},
	}
	// GIVEN
	os.Setenv("KSGO_PYTHON_BIN_PATH", "../../../../../src_py/.venv/bin/python3")
	os.Setenv("KSGO_PYTHON_SCRIPT_PATH", "../../../../../src_py/split_text.py")
	// WHEN
	result, err := pyUtil.SplitIntoSentences(paragraphs)
	// THEN
	assert.Nil(t, err)
	assert.Equal(t, expectedResult, result)
}

func TestSplitIntoSentencesWrongBinPath(t *testing.T) {
	pyUtil := NewPythonUtil()
	os.Setenv("KSGO_PYTHON_BIN_PATH", "../")
	os.Setenv("KSGO_PYTHON_SCRIPT_PATH", "../../../../../src_py/split_text.py")
	paragraphs := []model.Paragraph{
		{
			Id:   1,
			Text: "Das ist ein erster Satz. Das ist ein zweiter Satz.",
		},
	}
	// WHEN
	result, err := pyUtil.SplitIntoSentences(paragraphs)
	// THEN
	assert.NotNil(t, err)
	assert.Nil(t, result)
}

func TestSplitIntoSentencesWrongScriptPath(t *testing.T) {
	pyUtil := NewPythonUtil()
	os.Setenv("KSGO_PYTHON_BIN_PATH", "../../../../../src_py/.venv/bin/python3")
	os.Setenv("KSGO_PYTHON_SCRIPT_PATH", "../")
	paragraphs := []model.Paragraph{
		{
			Id:   1,
			Text: "Das ist ein erster Satz. Das ist ein zweiter Satz.",
		},
	}
	// WHEN
	result, err := pyUtil.SplitIntoSentences(paragraphs)
	// THEN
	assert.NotNil(t, err)
	assert.Nil(t, result)
}
