//go:build unit
// +build unit

package mapping

import (
	"testing"

	"github.com/frhorschig/kant-search-api/src/go/models"
	"github.com/frhorschig/kant-search-backend/dataaccess/model"
	"github.com/stretchr/testify/assert"
)

func TestCriteriaToCoreModel(t *testing.T) {
	criteria := models.SearchCriteria{
		SearchString: "search terms",
		Options: models.SearchOptions{
			IncludeHeadings: false,
			Scope:           models.SearchScope("PARAGRAPH"),
			WorkIds:         []string{"id1", "id2"},
		},
	}

	ss, opts := CriteriaToCoreModel(&criteria)

	assert.Equal(t, ss, criteria.SearchString)
	assert.Equal(t, opts.IncludeHeadings, criteria.Options.IncludeHeadings)
	assert.Equal(t, string(opts.Scope), string(criteria.Options.Scope))
	assert.Equal(t, opts.WorkIds, criteria.Options.WorkIds)
}

func TestMatchesToApiModels(t *testing.T) {
	match1 := model.SearchResult{
		Snippet:   "snippet",
		Pages:     []int32{1, 2},
		ContentId: "content1Id",
		WorkId:    "work1Id",
	}
	match2 := model.SearchResult{
		Snippet:   "snippet2",
		Pages:     []int32{3, 4},
		ContentId: "content2Id",
		WorkId:    "work1Id",
	}
	match3 := model.SearchResult{
		Snippet:   "snippet3",
		Pages:     []int32{5, 6},
		ContentId: "content3Id",
		WorkId:    "work2Id",
	}
	matches := []model.SearchResult{
		match1, match2, match3,
	}

	results := MatchesToApiModels(matches)

	assert.Len(t, results, 2)

	assert.Equal(t, results[0].WorkId, match1.WorkId)
	assert.Equal(t, results[0].Matches[0].Pages, match1.Pages)
	assert.Equal(t, results[0].Matches[0].Snippet, match1.Snippet)
	assert.Equal(t, results[0].Matches[0].ContentId, match1.ContentId)
	assert.Equal(t, results[0].Matches[1].Pages, match2.Pages)
	assert.Equal(t, results[0].Matches[1].Snippet, match2.Snippet)
	assert.Equal(t, results[0].Matches[1].ContentId, match2.ContentId)

	assert.Equal(t, results[1].WorkId, match3.WorkId)
	assert.Equal(t, results[1].Matches[0].Pages, match3.Pages)
	assert.Equal(t, results[1].Matches[0].Snippet, match3.Snippet)
	assert.Equal(t, results[1].Matches[0].ContentId, match3.ContentId)
}
