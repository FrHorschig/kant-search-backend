package model

type SearchScope string

const (
	ParagraphScope SearchScope = "PARAGRAPH"
)

type SearchOptions struct {
	IncludeHeadings bool
	Scope           SearchScope
}
type SearchCriteria struct {
	WorkIds      []string
	SearchString string
	Options      SearchOptions
}

type SearchResult struct {
	Snippet   string
	Text      string
	Pages     []int32
	ContentId string
	WorkId    string
}
