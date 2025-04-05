package model

type Volume struct {
	Id      int32
	Section int32
	Title   string
	Works   []WorkRef
}

type WorkRef struct {
	Id    int32
	Code  string
	Title string
}

type Work struct {
	Id           int32
	Code         string
	Abbreviation *string
	Title        string
	Year         *string
	Sections     []Section
	Footnotes    []Footnote
	Summaries    []Summary
}

type Section struct {
	Id         int32
	Heading    Heading
	Paragraphs []Paragraph
	Sections   []Section
}

type Heading struct {
	Id      int32
	Text    string
	TocText string
	Pages   []int32
	FnRefs  []string
	WorkId  int32
}

type Paragraph struct {
	Id         int32
	Text       string
	Pages      []int32
	FnRefs     []string
	SummaryRef *string
	WorkId     int32
}

type Footnote struct {
	Id     int32
	Ref    string
	Text   string
	Pages  []int32
	WorkId int32
}

type Summary struct {
	Id     int32
	Ref    string
	Text   string
	Pages  []int32
	FnRefs []string
	WorkId int32
}
