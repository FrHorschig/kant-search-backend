//go:build unit
// +build unit

package mapper

import (
	"testing"

	"github.com/FrHorschig/kant-search-backend/database/model"
	"github.com/stretchr/testify/assert"
)

func TestParagraphToApiModel(t *testing.T) {
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
}

func TestParagraphToApiModels(t *testing.T) {
	pars := []model.Paragraph{{
		Id:     1,
		Text:   "text",
		Pages:  []int32{1, 2, 3},
		WorkId: 1,
	}}

	result := ParagraphsToApiModels(pars)

	assert.Equal(t, len(result), len(pars))
	assert.Equal(t, result[0].Id, pars[0].Id)
	assert.Equal(t, result[0].Text, pars[0].Text)
	assert.Equal(t, result[0].Pages, pars[0].Pages)
	assert.Equal(t, result[0].WorkId, pars[0].WorkId)
}
