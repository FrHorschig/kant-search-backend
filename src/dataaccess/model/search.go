package model

// TODO move these models to core/search/internal and create new ones with only the necessary fields
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
	IncludeHeadings bool
	Scope           SearchScope
	WorkIds         []string
}

type SearchResult struct {
	Snippet   string
	Text      string
	Pages     []int32
	ContentId string
	WorkId    string
}
