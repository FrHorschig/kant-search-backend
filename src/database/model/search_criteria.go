package model

type SearchScope string

const (
	ParagraphScope SearchScope = "PARAGRAPH"
	SentenceScope  SearchScope = "SENTENCE"
)

type SearchCriteria struct {
	WorkIds       []int32
	SearchTerms   string
	ExcludedTerms string
	OptionalTerms string
	Scope         SearchScope
}
