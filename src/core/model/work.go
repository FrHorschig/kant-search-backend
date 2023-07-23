package model

type Work struct {
	Title        string
	Abbreviation string
	Text         string
	Volume       int32
	Year         int32
}

type WorkMetadata struct {
	Id           int32
	Title        string
	Abbreviation string
	Volume       int32
	Year         int32
}
