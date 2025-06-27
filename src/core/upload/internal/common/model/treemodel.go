package model

type TreeSection struct {
	Heading    TreeHeading
	Paragraphs []string
	Sections   []TreeSection
	Parent     *TreeSection
}

type TreeHeading struct {
	Level     TreeLevel
	TocTitle  string
	TextTitle string
	Year      string
}

type TreeLevel int32

const (
	HWork TreeLevel = iota + 0
	H1
	H2
	H3
	H4
	H5
	H6
	H7
	H8
)

type TreeSummary struct {
	Page int32
	Line int32
	Text string
}

type TreeFootnote struct {
	Page int32
	Nr   int32
	Text string
}
