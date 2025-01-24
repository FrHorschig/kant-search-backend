package model

type Volume struct {
	Section int32
	Title   string
	Works   []Work
}

type Work struct {
	Id           int32
	Code         string
	Title        string
	Abbreviation *string
	Year         *string
	TextData     Section
	Footnotes    []Footnote
	Volume       int32
}

type Section struct {
	Id         int32
	Heading    Heading
	Paragraphs []Paragraph
	Sections   []Section
	Parent     *Section
}

type Heading struct {
	Id        int32
	Level     Level
	TocTitle  string
	TextTitle string
	WorkId    int32
}

type Level int32

const (
	HWork Level = iota + 1
	H1
	H2
	H3
	H4
	H5
	H6
	H7
)

type Paragraph struct {
	Id           int32
	Text         string
	Pages        []int32
	FnReferences []int32
	Sentences    []Sentence
	WorkId       int32
}

type Sentence struct {
	Id          int32
	Text        string
	Pages       []int32
	ParagraphId int32
	WorkId      int32
}

type Footnote struct {
	Id     int32
	Name   string
	Text   string
	Pages  []int32
	WorkId int32
}
