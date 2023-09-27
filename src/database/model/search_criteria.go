package model

type SearchScope string

const (
	ParagraphScope SearchScope = "PARAGRAPH"
	SentenceScope  SearchScope = "SENTENCE"
)

type SearchOptions struct {
	Scope SearchScope
}
type SearchCriteria struct {
	WorkIds      []int32
	SearchString string
	Options      SearchOptions
}
