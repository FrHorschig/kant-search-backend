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
		SearchString: "search terms",
		Options: models.SearchOptions{
			IncludeHeadings: false,
			Scope:           models.SearchScope("PARAGRAPH"),
			WorkIds:         []string{"id1", "id2"},
		},
	}

	ss, opts := CriteriaToCoreModel(&criteria)

	assert.Equal(t, ss, criteria.SearchString)
	assert.Equal(t, opts.WorkIds, criteria.Options.WorkIds)
	assert.Equal(t, string(opts.Scope), string(criteria.Options.Scope))
	assert.Equal(t, opts.IncludeHeadings, criteria.Options.IncludeHeadings)
	assert.Equal(t, opts.IncludeFootnotes, criteria.Options.IncludeFootnotes)
	assert.Equal(t, opts.IncludeSummaries, criteria.Options.IncludeSummaries)
}

func TestMatchesToApiModels(t *testing.T) {
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
					WorkId:    "w1",
					Snippets:  []string{"snippet1"},
					Pages:     []int32{1},
					ContentId: "c1",
				},
			},
			expected: []models.SearchResult{
				{
					WorkId: "w1",
					Hits:   []models.Hit{{Snippets: []string{"snippet1"}, Pages: []int32{1}, ContentId: "c1"}},
				},
			},
		},
		{
			name: "multiple results",
			input: []model.SearchResult{
				{WorkId: "w1", Snippets: []string{"a"}, Pages: []int32{1}, ContentId: "x"},
				{WorkId: "w2", Snippets: []string{"b"}, Pages: []int32{2}, ContentId: "y"},
			},
			expected: []models.SearchResult{
				{
					WorkId: "w1",
					Hits:   []models.Hit{{Snippets: []string{"a"}, Pages: []int32{1}, ContentId: "x"}},
				},
				{
					WorkId: "w2",
					Hits:   []models.Hit{{Snippets: []string{"b"}, Pages: []int32{2}, ContentId: "y"}},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := MatchesToApiModels(tt.input)
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
		m1[r.WorkId] = r.Hits
	}
	for _, r := range b {
		m2[r.WorkId] = r.Hits
	}
	return reflect.DeepEqual(m1, m2)
}
