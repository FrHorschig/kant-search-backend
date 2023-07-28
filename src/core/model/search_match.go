package model

type SearchMatch struct {
	Volume    int32
	WorkTitle string
	Snippet   string
	Pages     []int32
	MatchId   *interface{}
}
