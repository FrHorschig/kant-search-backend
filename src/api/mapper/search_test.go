//go:build unit
// +build unit

package mapper

import (
	"testing"

	"github.com/FrHorschig/kant-search-api/models"
	"github.com/FrHorschig/kant-search-backend/database/model"
	"github.com/stretchr/testify/assert"
)

func TestCriteriaToCoreModel(t *testing.T) {
	criteria := models.SearchCriteria{
		WorkIds:       []int32{1, 2},
		SearchTerms:   []string{"search", "terms"},
		ExcludedTerms: []string{"excluded", "terms"},
		OptionalTerms: []string{"optional", "terms"},
		Scope:         models.SearchScope("PARAGRAPH"),
	}

	result := CriteriaToCoreModel(criteria)

	assert.Equal(t, result.WorkIds, criteria.WorkIds)
	assert.Equal(t, result.SearchTerms, criteria.SearchTerms)
	assert.Equal(t, result.ExcludedTerms, criteria.ExcludedTerms)
	assert.Equal(t, result.OptionalTerms, criteria.OptionalTerms)
	assert.Equal(t, string(result.Scope), string(criteria.Scope))
}

func TestCriteriaToCoreModelWithEmptyStrings(t *testing.T) {
	criteria := models.SearchCriteria{
		SearchTerms:   []string{"search", "terms"},
		ExcludedTerms: []string{"", "  "},
		OptionalTerms: []string{"\t", "\n"},
	}

	result := CriteriaToCoreModel(criteria)

	assert.Len(t, result.SearchTerms, 2)
	assert.Len(t, result.ExcludedTerms, 0)
	assert.Len(t, result.OptionalTerms, 0)
}

func TestMatchesToApiModels(t *testing.T) {
	match1 := model.SearchResult{
		ElementId: 1,
		Snippet:   "snippet",
		Pages:     []int32{1, 2},
		WorkId:    1,
	}
	match2 := model.SearchResult{
		ElementId: 2,
		Snippet:   "snippet2",
		Pages:     []int32{3, 4},
		WorkId:    1,
	}
	match3 := model.SearchResult{
		ElementId: 3,
		Snippet:   "snippet3",
		Pages:     []int32{5, 6},
		WorkId:    2,
	}
	matches := []model.SearchResult{
		match1, match2, match3,
	}

	results := MatchesToApiModels(matches)

	assert.Equal(t, len(results), 2)
	assert.Equal(t, results[0].WorkId, match1.WorkId)
	assert.Equal(t, results[0].Matches[0].ElementId, match1.ElementId)
	assert.Equal(t, results[0].Matches[0].Pages, match1.Pages)
	assert.Equal(t, results[0].Matches[0].Snippet, match1.Snippet)
	assert.Equal(t, results[0].Matches[1].ElementId, match2.ElementId)
	assert.Equal(t, results[0].Matches[1].Pages, match2.Pages)
	assert.Equal(t, results[0].Matches[1].Snippet, match2.Snippet)
	assert.Equal(t, results[1].WorkId, match3.WorkId)
	assert.Equal(t, results[1].Matches[0].ElementId, match3.ElementId)
	assert.Equal(t, results[1].Matches[0].Pages, match3.Pages)
	assert.Equal(t, results[1].Matches[0].Snippet, match3.Snippet)
}
