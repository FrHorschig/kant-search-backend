package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/FrHorschig/kant-search-backend/database/model"
	"github.com/lib/pq"
)

func buildParams() (snippetParams string, textParams string) {
	snippetParams = `FragmentDelimiter="...<br>... ",
		MaxFragments=10,
		MaxWords=16,
		MinWords=6`
	textParams = `MaxWords=99999, MinWords=99998`
	return
}

func buildTerms(c model.SearchCriteria) string {
	var builder strings.Builder
	builder.WriteString(strings.Join(c.SearchTerms, " & "))
	if len(c.ExcludedTerms) > 0 {
		builder.WriteString(" & !")
		builder.WriteString(strings.Join(c.ExcludedTerms, " & !"))
	}
	if len(c.OptionalTerms) > 0 {
		builder.WriteString(" | ")
		builder.WriteString(strings.Join(c.OptionalTerms, " | "))
	}
	return builder.String()
}

func scanSearchMatchRow(rows *sql.Rows) ([]model.SearchResult, error) {
	matches := make([]model.SearchResult, 0)
	for rows.Next() {
		var match model.SearchResult
		err := rows.Scan(&match.ElementId, &match.Snippet, &match.Text, pq.Array(&match.Pages), &match.WorkId)
		if err != nil {
			return nil, fmt.Errorf("search match row scan failed: %v", err)
		}
		matches = append(matches, match)
	}
	return matches, nil
}
