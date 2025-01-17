package model

type Section struct {
	Heading    Heading
	Paragraphs []string
	Sections   []Section
	Parent     *Section
}

type Heading struct {
	TocTitle  string
	TextTitle string
	Level     Level
}

type Level int

const (
	H1 Level = iota + 1
	H2
	H3
	H4
	H5
	H6
	H7
	H8
	H9
)

type Summary struct {
	Page int32
	Line int32
	Text string
}

type Footnote struct {
	Page int32
	Nr   int32
	Text string
}
