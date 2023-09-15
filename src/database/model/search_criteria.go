package model

type SearchScope string

// List of SearchScope
const (
	PARAGRAPH SearchScope = "paragraph"
	SENTENCE  SearchScope = "sentence"
)

type SearchCriteria struct {
	SearchTerms []string
	WorkIds     []int32
	Scope       SearchScope
}
