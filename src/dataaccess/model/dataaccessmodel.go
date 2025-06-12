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
	HighlightText string
	FmtText       string
	Pages         []int32
	PageByIndex   []IndexNumberPair
	LineByIndex   []IndexNumberPair
	Ordinal       int32
	WorkCode      string
	WordIndexMap  map[int32]int32
}

type IndexNumberPair struct {
	I   int32
	Num int32
}
