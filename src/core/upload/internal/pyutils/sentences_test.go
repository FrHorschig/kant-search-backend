//go:build integration
// +build integration

package pyutils

import (
	"os"
	"testing"

	"github.com/FrHorschig/kant-search-backend/common/model"
	"github.com/stretchr/testify/assert"
)

func TestSplitIntoSentences(t *testing.T) {
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
	os.Setenv("PYTHON_BIN_PATH", "../../../../../src_py")
	// WHEN
	result, err := SplitIntoSentences(paragraphs)
	// THEN
	assert.Nil(t, err)
	assert.Equal(t, expectedResult, result)
}

func TestSplitIntoSentencesDefaultPath(t *testing.T) {
	os.Setenv("PYTHON_BIN_PATH", "")
	paragraphs := []model.Paragraph{
		{
			Id:   1,
			Text: "Das ist ein erster Satz. Das ist ein zweiter Satz.",
		},
	}
	// WHEN
	result, err := SplitIntoSentences(paragraphs)
	// THEN
	assert.NotNil(t, err)
	assert.Nil(t, result)
}
