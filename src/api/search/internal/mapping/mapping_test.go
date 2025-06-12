//go:build unit
// +build unit

package mapping

import (
	"reflect"
	"testing"

	"github.com/frhorschig/kant-search-api/src/go/models"
	"github.com/frhorschig/kant-search-backend/dataaccess/model"
	"github.com/stretchr/testify/assert"
)

func TestCriteriaToCoreModel(t *testing.T) {
	criteria := models.SearchCriteria{
		SearchTerms: "search terms",
		Options: models.SearchOptions{
			IncludeHeadings:  false,
			IncludeFootnotes: true,
			IncludeSummaries: false,
			Scope:            models.SearchScope("PARAGRAPH"),
			WithStemming:     true,
			WorkCodes:        []string{"id1", "id2"},
		},
	}

	ss, opts := CriteriaToCoreModel(&criteria)

	assert.Equal(t, ss, criteria.SearchTerms)
	assert.Equal(t, opts.WorkCodes, criteria.Options.WorkCodes)
	assert.Equal(t, string(opts.Scope), string(criteria.Options.Scope))
	assert.Equal(t, opts.IncludeHeadings, criteria.Options.IncludeHeadings)
	assert.Equal(t, opts.IncludeFootnotes, criteria.Options.IncludeFootnotes)
	assert.Equal(t, opts.IncludeSummaries, criteria.Options.IncludeSummaries)
}

func TestHitsToApiModels(t *testing.T) {
	wimInt := make(map[int32]int32)
	wimInt[3] = 5
	wimInt[9] = 47
	wimInt[36] = 184
	wimStr := make(map[string]int32)
	wimStr["3"] = 5
	wimStr["9"] = 47
	wimStr["36"] = 184

	tests := []struct {
		name     string
		input    []model.SearchResult
		expected []models.SearchResult
	}{
		{
			name:     "empty input returns empty output",
			input:    []model.SearchResult{},
			expected: []models.SearchResult{},
		},
		{
			name: "single result",
			input: []model.SearchResult{
				{
					WorkCode:     "w1",
					Snippets:     []string{"snippet1"},
					Pages:        []int32{1},
					Ordinal:      1,
					FmtText:      "fmtText",
					RawText:      "rawText",
					WordIndexMap: wimInt,
				},
			},
			expected: []models.SearchResult{
				{
					WorkCode: "w1",
					Hits: []models.Hit{{
						Snippets:     []string{"snippet1"},
						Pages:        []int32{1},
						Ordinal:      1,
						FmtText:      "fmtText",
						RawText:      "rawText",
						WordIndexMap: wimStr,
					}},
				},
			},
		},
		{
			name: "multiple results",
			input: []model.SearchResult{
				{
					WorkCode:     "w1",
					Snippets:     []string{"a"},
					Pages:        []int32{1},
					Ordinal:      1,
					FmtText:      "fmtText",
					RawText:      "rawText",
					WordIndexMap: wimInt,
				},
				{WorkCode: "w2",
					Snippets:     []string{"b"},
					Pages:        []int32{2},
					Ordinal:      2,
					FmtText:      "fmtText",
					RawText:      "rawText",
					WordIndexMap: wimInt,
				},
			},
			expected: []models.SearchResult{
				{
					WorkCode: "w1",
					Hits: []models.Hit{{
						Snippets:     []string{"a"},
						Pages:        []int32{1},
						Ordinal:      1,
						FmtText:      "fmtText",
						RawText:      "rawText",
						WordIndexMap: wimStr,
					}},
				},
				{
					WorkCode: "w2",
					Hits: []models.Hit{{
						Snippets:     []string{"b"},
						Pages:        []int32{2},
						Ordinal:      2,
						FmtText:      "fmtText",
						RawText:      "rawText",
						WordIndexMap: wimStr,
					}},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := HitsToApiModels(tt.input)
			if !equalSearchResults(actual, tt.expected) {
				t.Errorf("expected %+v, got %+v", tt.expected, actual)
			}
		})
	}
}

func equalSearchResults(a, b []models.SearchResult) bool {
	if len(a) != len(b) {
		return false
	}
	m1 := make(map[string][]models.Hit)
	m2 := make(map[string][]models.Hit)
	for _, r := range a {
		m1[r.WorkCode] = r.Hits
	}
	for _, r := range b {
		m2[r.WorkCode] = r.Hits
	}
	return reflect.DeepEqual(m1, m2)
}
