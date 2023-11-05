package model

type Work struct {
	Id           int32
	Code         string
	Abbreviation *string
	Ordinal      int32
	Year         *string
	Volume       int32
}

type WorkUpload struct {
	WorkId int32
	Text   string
}
