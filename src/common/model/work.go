package model

type Work struct {
	Id           int32
	Title        string
	Abbreviation *string
	Ordinal      int32
	Year         *string
	Volume       int32
}

type WorkUpload struct {
	WorkId int32
	Text   string
}
