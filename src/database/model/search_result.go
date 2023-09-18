package model

type SearchMatch struct {
	// TODO frhorsch: adjust struct and search_repo to match api spec
	Volume    int32
	WorkTitle string
	Snippet   string
	Pages     []int32
	WorkId    int32
	ElementId int32
}
