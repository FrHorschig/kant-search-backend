package model

type Volume struct {
	VolumeNumber int32
	Section      int32
	Title        string
}

type Work struct {
	Code       string
	Siglum     *string
	Title      string
	Year       string
	Paragraphs []Paragraph
	Sections   []Section
	Footnotes  []Footnote
	Summaries  []Summary
}

type Section struct {
	Heading    Heading
	Paragraphs []Paragraph
	Sections   []Section
}

type Heading struct {
	Ordinal int32
	Text    string
	TocText string
	Pages   []int32
	FnRefs  []string
}

type Paragraph struct {
	Ordinal    int32
	Text       string
	Pages      []int32
	FnRefs     []string
	SummaryRef *string
}

type Footnote struct {
	Ordinal int32
	Ref     string
	Text    string
	Pages   []int32
}

type Summary struct {
	Ordinal int32
	Ref     string
	Text    string
	Pages   []int32
	FnRefs  []string
}
