package mapper

import (
	"testing"

	"github.com/FrHorschig/kant-search-api/models"
	"github.com/FrHorschig/kant-search-backend/database/model"
	"github.com/stretchr/testify/assert"
)

func TestCriteriaToCoreModel(t *testing.T) {
	criteria := models.SearchCriteria{
		SearchTerms: []string{"search", "terms"},
		WorkIds:     []int32{1, 2},
		Scope:       models.PARAGRAPH,
	}

	result := CriteriaToCoreModel(criteria)

	assert.Equal(t, result.SearchTerms, criteria.SearchTerms)
	assert.Equal(t, result.WorkIds, criteria.WorkIds)
	assert.Equal(t, result.Scope, model.PARAGRAPH)
}

func TestCriteriaToCoreModelSentenceScope(t *testing.T) {
	criteria := models.SearchCriteria{
		SearchTerms: []string{"search", "terms"},
		WorkIds:     []int32{1, 2},
		Scope:       models.SENTENCE,
	}

	result := CriteriaToCoreModel(criteria)

	assert.Equal(t, result.Scope, model.SENTENCE)
}

func TestCriteriaToCoreModelDefaultScope(t *testing.T) {
	criteria := models.SearchCriteria{
		SearchTerms: []string{"search", "terms"},
		WorkIds:     []int32{1, 2},
	}

	result := CriteriaToCoreModel(criteria)

	assert.Equal(t, result.Scope, model.PARAGRAPH)
}

func TestMatchesToApiModels(t *testing.T) {
	match1 := model.SearchResult{
		WorkId:    1,
		Snippet:   "snippet",
		Pages:     []int32{1, 2},
		ElementId: 1,
	}
	match2 := model.SearchResult{
		WorkId:    1,
		Snippet:   "snippet2",
		Pages:     []int32{3, 4},
		ElementId: 2,
	}
	match3 := model.SearchResult{
		WorkId:    2,
		Snippet:   "snippet3",
		Pages:     []int32{5, 6},
		ElementId: 3,
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
