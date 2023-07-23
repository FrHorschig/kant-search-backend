package model

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
	WorkId      int32
}
