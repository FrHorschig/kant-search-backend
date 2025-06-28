package model

type SearchTermNode struct {
	Left  *SearchTermNode
	Right *SearchTermNode
	Token *Token
}

type Token struct {
	IsAnd    bool
	IsOr     bool
	IsNot    bool
	IsWord   bool
	IsPhrase bool
	Text     string
}

type SearchOptions struct {
	IncludeHeadings  bool
	IncludeFootnotes bool
	IncludeSummaries bool
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
