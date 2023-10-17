//go:build unit
// +build unit

package transform

import (
	"fmt"
	"testing"

	"github.com/FrHorschig/kant-search-backend/common/model"
	"github.com/FrHorschig/kant-search-backend/core/errors"
	"github.com/FrHorschig/kant-search-backend/core/upload/internal/pyutil/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestFindSentencesSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	pyUtil := mocks.NewMockPythonUtil(ctrl)

	paragraphs := []model.Paragraph{
		{Id: 1, Text: "Das ist ein erster Satz. Das ist ein zweiter Satz!"},
		{Id: 2, Text: "Das ist ein erster Satz. Das ist ein zweiter Satz ohne Punkt"},
		{Id: 3, Text: "der weiter geht. Das ist ein dritter Satz."},
		{Id: 4, Text: "Das ist ein Satz mit Abkürzung z.B. und weiter gehts."},
	}
	expected := []model.Sentence{
		{ParagraphId: 1, Text: "Das ist ein erster Satz."},
		{ParagraphId: 1, Text: "Das ist ein zweiter Satz!"},
		{ParagraphId: 2, Text: "Das ist ein erster Satz."},
		{ParagraphId: 2, Text: "Das ist ein zweiter Satz ohne Punkt"},
		{ParagraphId: 3, Text: "der weiter geht."},
		{ParagraphId: 3, Text: "Das ist ein dritter Satz."},
		{ParagraphId: 4, Text: "Das ist ein Satz mit Abkürzung z.B. und weiter gehts."},
	}

	// GIVEN
	pyUtil.EXPECT().SplitIntoSentences(paragraphs).Return(map[int32][]string{
		1: {"Das ist ein erster Satz.", "Das ist ein zweiter Satz!"},
		2: {"Das ist ein erster Satz.", "Das ist ein zweiter Satz ohne Punkt"},
		3: {"der weiter geht.", "Das ist ein dritter Satz."},
		4: {"Das ist ein Satz mit Abkürzung z.B. und weiter gehts."},
	}, nil)

	// WHEN
	result, err := FindSentences(paragraphs, pyUtil)

	// THEN
	assert.Nil(t, err)
	assert.Equal(t, expected, result)
}

func TestFindSentencesPyUtilError(t *testing.T) {
	ctrl := gomock.NewController(t)
	pyUtil := mocks.NewMockPythonUtil(ctrl)
	paragraphs := []model.Paragraph{
		{Id: 1, Text: "Das ist ein erster Satz. Das ist ein zweiter Satz!"},
	}
	expected := &errors.Error{
		Msg:    errors.GO_ERR,
		Params: []string{"error"},
	}

	// GIVEN
	pyUtil.EXPECT().SplitIntoSentences(paragraphs).Return(nil, fmt.Errorf("error"))

	// WHEN
	result, err := FindSentences(paragraphs, pyUtil)

	// THEN
	assert.Equal(t, expected, err)
	assert.Nil(t, result)
}
