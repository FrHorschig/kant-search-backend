package repository

import (
	"strings"

	"github.com/FrHorschig/kant-search-backend/database/model"
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
		builder.WriteString(" & ! ")
		builder.WriteString(strings.Join(c.ExcludedTerms, " & ! "))
	}
	if len(c.OptionalTerms) > 0 {
		builder.WriteString(" & ( ")
		builder.WriteString(strings.Join(c.OptionalTerms, " | "))
		builder.WriteString(" )")
	}
	return builder.String()
}
