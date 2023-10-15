package model

type SearchResult struct {
	Snippet     string
	Text        string
	Pages       []int32
	ParagraphId int32
	SentenceId  int32
	WorkId      int32
}
