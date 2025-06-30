//go:build unit
// +build unit

package mapping

import (
	"reflect"
	"testing"

	"github.com/frhorschig/kant-search-api/generated/go/models"
	"github.com/frhorschig/kant-search-backend/dataaccess/model"
	"github.com/stretchr/testify/assert"
)

func TestCriteriaToCoreModel(t *testing.T) {
	criteria := models.SearchCriteria{
		SearchTerms: "search terms",
		Options: models.SearchOptions{
			IncludeHeadings:   false,
			IncludeFootnotes:  true,
			IncludeParagraphs: false,
			WithStemming:      true,
			WorkCodes:         []string{"id1", "id2"},
		},
	}

	ss, opts := CriteriaToCoreModel(&criteria)

	assert.Equal(t, ss, criteria.SearchTerms)
	assert.Equal(t, opts.WorkCodes, criteria.Options.WorkCodes)
	assert.Equal(t, opts.IncludeHeadings, criteria.Options.IncludeHeadings)
	assert.Equal(t, opts.IncludeFootnotes, criteria.Options.IncludeFootnotes)
	assert.Equal(t, opts.IncludeParagraphs, criteria.Options.IncludeParagraphs)
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
					HighlightText: "highlightText",
					FmtText:       "fmtText",
					Ordinal:       1,
					Pages:         []int32{1},
					PageByIndex:   []model.IndexNumberPair{{I: 1, Num: 2}, {I: 4, Num: 183}},
					LineByIndex:   []model.IndexNumberPair{{I: 32, Num: 54}},
					WordIndexMap:  wimInt,
					WorkCode:      "w1",
				},
			},
			expected: []models.SearchResult{
				{
					WorkCode: "w1",
					Hits: []models.Hit{{
						HighlightText: "highlightText",
						FmtText:       "fmtText",
						Pages:         []int32{1},
						PageByIndex:   []models.IndexNumberPair{{I: 1, Num: 2}, {I: 4, Num: 183}},
						LineByIndex:   []models.IndexNumberPair{{I: 32, Num: 54}},
						Ordinal:       1,
						WordIndexMap:  wimStr,
					}},
				},
			},
		},
		{
			name: "multiple results",
			input: []model.SearchResult{
				{
					HighlightText: "highlightText",
					FmtText:       "fmtText",
					Ordinal:       1,
					Pages:         []int32{1},
					PageByIndex:   []model.IndexNumberPair{{I: 1, Num: 2}, {I: 4, Num: 183}},
					LineByIndex:   []model.IndexNumberPair{{I: 32, Num: 54}},
					WordIndexMap:  wimInt,
					WorkCode:      "w1",
				},
				{
					HighlightText: "highlightText2",
					FmtText:       "fmtText2",
					Ordinal:       2,
					Pages:         []int32{2},
					PageByIndex:   []model.IndexNumberPair{{I: 12, Num: 37}},
					LineByIndex:   []model.IndexNumberPair{{I: 8, Num: 2481}},
					WordIndexMap:  wimInt,
					WorkCode:      "w2",
				},
			},
			expected: []models.SearchResult{
				{
					WorkCode: "w1",
					Hits: []models.Hit{{
						HighlightText: "highlightText",
						FmtText:       "fmtText",
						Ordinal:       1,
						Pages:         []int32{1},
						PageByIndex:   []models.IndexNumberPair{{I: 1, Num: 2}, {I: 4, Num: 183}},
						LineByIndex:   []models.IndexNumberPair{{I: 32, Num: 54}},
						WordIndexMap:  wimStr,
					}},
				},
				{
					WorkCode: "w2",
					Hits: []models.Hit{{
						HighlightText: "highlightText2",
						FmtText:       "fmtText2",
						Ordinal:       2,
						Pages:         []int32{2},
						PageByIndex:   []models.IndexNumberPair{{I: 12, Num: 37}},
						LineByIndex:   []models.IndexNumberPair{{I: 8, Num: 2481}},
						WordIndexMap:  wimStr,
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
