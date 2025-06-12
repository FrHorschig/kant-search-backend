package model

// TODO move these models to core/search/internal and create new ones with only the necessary fields; also check if phrase search works
type AstNode struct {
	Left  *AstNode
	Right *AstNode
	Token *Token
}

type Token struct {
	IsAnd    bool
	IsOr     bool
	IsNot    bool
	IsOpen   bool
	IsClose  bool
	IsWord   bool
	IsPhrase bool
	Text     string
}

type SearchScope string

const (
	ParagraphScope SearchScope = "PARAGRAPH"
)

type SearchOptions struct {
	IncludeHeadings  bool
	IncludeFootnotes bool
	IncludeSummaries bool
	Scope            SearchScope
	WithStemming     bool
	WorkCodes        []string
}

type SearchResult struct {
	Snippets []string
	Pages    []int32
	Ordinal  int32
	WorkCode string
	FmtText  string
	RawText  string
}
