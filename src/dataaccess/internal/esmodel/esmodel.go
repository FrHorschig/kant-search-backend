package esmodel

// structs for volume-works tree data
type Volume struct {
	VolumeNumber int32     `json:"volumeNumber"`
	Section      int32     `json:"section"`
	Title        string    `json:"title"`
	Works        []WorkRef `json:"works"`
}

type WorkRef struct {
	Id    int32  `json:"id"`
	Code  string `json:"code"`
	Title string `json:"title"`
}

// structs for works tree data without content
type Work struct {
	Id           int32     `json:"id"`
	Code         string    `json:"code"`
	Abbreviation *string   `json:"abbreviation"`
	Title        string    `json:"title"`
	Year         *string   `json:"year"`
	Sections     []Section `json:"sections"`
	Footnotes    []int32   `json:"footnotes"`
	Summaries    []int32   `json:"summaries"`
}

type Section struct {
	Heading    int32     `json:"heading"`
	Paragraphs []int32   `json:"paragraphs"`
	Sections   []Section `json:"sections"`
}

// structs for actual content, stored in a linear structure to make searching and fetching it simple
type Heading struct {
	Id         int32    `json:"id"`
	FmtText    string   `json:"fmtText"`
	TocText    string   `json:"tocText"`
	SearchText string   `json:"searchText"`
	Pages      []int32  `json:"pages"`
	FnRefs     []string `json:"fnRefs"`
	WorkId     int32    `json:"workId"`
}

type Paragraph struct {
	Id         int32    `json:"id"`
	FmtText    string   `json:"fmtText"`
	SearchText string   `json:"searchText"`
	Pages      []int32  `json:"pages"`
	FnRefs     []string `json:"fnRefs"`
	SummaryRef *string  `json:"summaryRef"`
	WorkId     int32    `json:"workId"`
}

type Footnote struct {
	Id         int32   `json:"id"`
	Ref        string  `json:"ref"`
	FmtText    string  `json:"fmtText"`
	SearchText string  `json:"searchText"`
	Pages      []int32 `json:"pages"`
	WorkId     int32   `json:"workId"`
}

type Summary struct {
	Id         int32    `json:"id"`
	Ref        string   `json:"ref"`
	FmtText    string   `json:"fmtText"`
	SearchText string   `json:"searchText"`
	Pages      []int32  `json:"pages"`
	FnRefs     []string `json:"fnRefs"`
	WorkId     int32    `json:"workId"`
}
