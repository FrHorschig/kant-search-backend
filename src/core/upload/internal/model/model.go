package model

type Volume struct {
	VolumeNumber int32
	Section      int32
	Title        string
}

type Work struct {
	Code         string
	Abbreviation *string
	Title        string
	Year         string
	Paragraphs   []Paragraph
	Sections     []Section
	Footnotes    []Footnote
	Summaries    []Summary
}

type Section struct {
	Heading    Heading
	Paragraphs []Paragraph
	Sections   []Section
}

type Heading struct {
	Text    string
	TocText string
	Pages   []int32
	FnRefs  []string
}

type Paragraph struct {
	Text       string
	Pages      []int32
	FnRefs     []string
	SummaryRef *string
}

type Footnote struct {
	Ref   string
	Text  string
	Pages []int32
}

type Summary struct {
	Ref    string
	Text   string
	Pages  []int32
	FnRefs []string
}
