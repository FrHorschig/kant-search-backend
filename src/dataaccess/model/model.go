package model

import "github.com/frhorschig/kant-search-backend/common/model"

type Volume struct {
	Id      int32
	Section int32
	Title   string
	Works   []Work
}

type Work struct {
	Id           int32
	Code         string
	Abbreviation *string
	Year         *string
	Volume       int32
}

type Heading struct {
	Id     int32
	Level  model.Level
	Text   string
	WorkId int32
}

type Paragraph struct {
	Id     int32
	Text   string
	Pages  []int32
	WorkId int32
}

type Sentence struct {
	Id          int32
	Text        string
	ParagraphId int32
}

type Footnote struct {
	Id     int32
	Name   string
	Text   string
	Pages  []int32
	WorkId int32
}
