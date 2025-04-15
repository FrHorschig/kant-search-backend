package treemodel

type Section struct {
	Heading    Heading
	Paragraphs []string
	Sections   []Section
	Parent     *Section
}

type Heading struct {
	Level     Level
	TocTitle  string
	TextTitle string
	Year      string
}

type Level int32

const (
	HWork Level = iota + 0
	H1
	H2
	H3
	H4
	H5
	H6
	H7
	H8
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
