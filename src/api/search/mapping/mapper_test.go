//go:build unit
// +build unit

package mapping

import (
	"testing"

	"github.com/frhorschig/kant-search-api/generated/go/models"
	"github.com/frhorschig/kant-search-backend/common/model"
	"github.com/stretchr/testify/assert"
)

func TestCriteriaToCoreModel(t *testing.T) {
	criteria := models.SearchCriteria{
		WorkIds:      []int32{1, 2},
		SearchString: "search terms",
		Options:      models.SearchOptions{Scope: models.SearchScope("PARAGRAPH")},
	}

	result := CriteriaToCoreModel(criteria)

	assert.Equal(t, result.WorkIds, criteria.WorkIds)
	assert.Equal(t, result.SearchString, criteria.SearchString)
	assert.Equal(t, string(result.Options.Scope), string(criteria.Options.Scope))
}

func TestMatchesToApiModels(t *testing.T) {
	match1 := model.SearchResult{
		Snippet:     "snippet",
		Pages:       []int32{1, 2},
		SentenceId:  1,
		ParagraphId: 2,
		WorkId:      3,
	}
	match2 := model.SearchResult{
		Snippet:     "snippet2",
		Pages:       []int32{3, 4},
		SentenceId:  4,
		ParagraphId: 5,
		WorkId:      3,
	}
	match3 := model.SearchResult{
		Snippet:     "snippet3",
		Pages:       []int32{5, 6},
		SentenceId:  7,
		ParagraphId: 8,
		WorkId:      9,
	}
	matches := []model.SearchResult{
		match1, match2, match3,
	}

	results := MatchesToApiModels(matches)

	assert.Equal(t, len(results), 2)

	assert.Equal(t, results[0].WorkId, match1.WorkId)
	assert.Equal(t, results[0].Matches[0].Pages, match1.Pages)
	assert.Equal(t, results[0].Matches[0].Snippet, match1.Snippet)
	assert.Equal(t, results[0].Matches[0].SentenceId, match1.SentenceId)
	assert.Equal(t, results[0].Matches[0].ParagraphId, match1.ParagraphId)
	assert.Equal(t, results[0].Matches[1].Pages, match2.Pages)
	assert.Equal(t, results[0].Matches[1].Snippet, match2.Snippet)
	assert.Equal(t, results[0].Matches[1].SentenceId, match2.SentenceId)
	assert.Equal(t, results[0].Matches[1].ParagraphId, match2.ParagraphId)

	assert.Equal(t, results[1].WorkId, match3.WorkId)
	assert.Equal(t, results[1].Matches[0].Pages, match3.Pages)
	assert.Equal(t, results[1].Matches[0].Snippet, match3.Snippet)
	assert.Equal(t, results[1].Matches[0].SentenceId, match3.SentenceId)
	assert.Equal(t, results[1].Matches[0].ParagraphId, match3.ParagraphId)
}
