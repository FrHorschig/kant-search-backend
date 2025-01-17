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

type Randtext struct {
	Page string
	Line string
	Text string
}

type Footnote struct {
	// TODO frhorschig implement me
}
