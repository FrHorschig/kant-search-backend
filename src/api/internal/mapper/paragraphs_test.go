//go:build unit
// +build unit

package mapper

import (
	"testing"

	"github.com/FrHorschig/kant-search-backend/common/model"
	"github.com/stretchr/testify/assert"
)

func TestParagraphToApiModel(t *testing.T) {
	par := model.Paragraph{
		Id:           1,
		Text:         "text",
		Pages:        []int32{1, 2, 3},
		WorkId:       1,
		HeadingLevel: &[]int32{1}[0],
		FootnoteName: &[]string{"420.2"}[0],
	}

	result := ParagraphToApiModel(par)

	assert.Equal(t, result.Id, par.Id)
	assert.Equal(t, result.Text, par.Text)
	assert.Equal(t, result.Pages, par.Pages)
	assert.Equal(t, result.WorkId, par.WorkId)
	assert.Equal(t, result.HeadingLevel, par.HeadingLevel)
	assert.Equal(t, result.FootnoteName, par.FootnoteName)
}

func TestParagraphToApiModelWithNullPtrs(t *testing.T) {
	par := model.Paragraph{
		Id:     1,
		Text:   "text",
		Pages:  []int32{1, 2, 3},
		WorkId: 1,
	}

	result := ParagraphToApiModel(par)

	assert.Equal(t, result.Id, par.Id)
	assert.Equal(t, result.Text, par.Text)
	assert.Equal(t, result.Pages, par.Pages)
	assert.Equal(t, result.WorkId, par.WorkId)
	assert.Equal(t, result.HeadingLevel, int32(0))
	assert.Equal(t, result.FootnoteName, "")
}

func TestParagraphToApiModels(t *testing.T) {
	pars := []model.Paragraph{{
		Id:           1,
		Text:         "text",
		Pages:        []int32{1, 2, 3},
		WorkId:       1,
		HeadingLevel: &[]int32{1}[0],
		FootnoteName: &[]string{"420.2"}[0],
	}}

	result := ParagraphsToApiModels(pars)

	assert.Equal(t, len(result), len(pars))
	assert.Equal(t, result[0].Id, pars[0].Id)
	assert.Equal(t, result[0].Text, pars[0].Text)
	assert.Equal(t, result[0].Pages, pars[0].Pages)
	assert.Equal(t, result[0].WorkId, pars[0].WorkId)
	assert.Equal(t, result[0].HeadingLevel, pars[0].HeadingLevel)
	assert.Equal(t, result[0].FootnoteName, pars[0].FootnoteName)
}
