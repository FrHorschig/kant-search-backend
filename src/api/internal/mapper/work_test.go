//go:build unit
// +build unit

package mapper

import (
	"testing"

	"github.com/FrHorschig/kant-search-api/models"
	"github.com/FrHorschig/kant-search-backend/database/model"
	"github.com/stretchr/testify/assert"
)

func TestWorkUploadToCoreModel(t *testing.T) {
	upload := models.WorkUpload{
		WorkId: 1,
		Text:   "text",
	}

	result := WorkUploadToCoreModel(upload)

	assert.Equal(t, result.WorkId, upload.WorkId)
	assert.Equal(t, result.Text, upload.Text)
}

func TestWorkToApiModel(t *testing.T) {
	abbr := "abbr"
	year := "1785"
	works := []model.Work{{
		Id:           1,
		Title:        "title",
		Abbreviation: &abbr,
		Ordinal:      1,
		Year:         &year,
		Volume:       1,
	}}

	results := WorksToApiModels(works)

	assert.Equal(t, len(results), len(works))
	assert.Equal(t, results[0].Id, works[0].Id)
	assert.Equal(t, results[0].Title, works[0].Title)
	assert.Equal(t, results[0].Abbreviation, *works[0].Abbreviation)
	assert.Equal(t, results[0].Ordinal, works[0].Ordinal)
	assert.Equal(t, results[0].Year, *works[0].Year)
	assert.Equal(t, results[0].VolumeId, works[0].Volume)
}

func TestWorkToApiModelNilStrings(t *testing.T) {
	works := []model.Work{{
		Id:           1,
		Title:        "title",
		Abbreviation: nil,
		Ordinal:      1,
		Year:         nil,
		Volume:       1,
	}}

	results := WorksToApiModels(works)

	assert.Equal(t, len(results), len(works))
	assert.Empty(t, results[0].Abbreviation)
	assert.Empty(t, results[0].Year)
}

func TestVolumeToApiModel(t *testing.T) {
	volumes := []model.Volume{{
		Id:      1,
		Title:   "title",
		Section: 1,
	}}

	results := VolumesToApiModels(volumes)

	assert.Equal(t, len(results), len(volumes))
	assert.Equal(t, results[0].Id, volumes[0].Id)
	assert.Equal(t, results[0].Title, volumes[0].Title)
	assert.Equal(t, results[0].Section, volumes[0].Section)
}