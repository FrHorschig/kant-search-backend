package model

type SearchResult struct {
	ElementId int32
	Snippet   string
	Text      string
	Pages     []int32
	WorkId    int32
}
