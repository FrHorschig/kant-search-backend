package esmodel

import "github.com/elastic/go-elasticsearch/v8/typedapi/types"

// structs for volume-works tree data
type Volume struct {
	VolumeNumber int32     `json:"volumeNumber"`
	Section      int32     `json:"section"`
	Title        string    `json:"title"`
	Works        []WorkRef `json:"works"`
}

type WorkRef struct {
	Id    string `json:"id"`
	Code  string `json:"code"`
	Title string `json:"title"`
}

var VolumeMapping = &types.TypeMapping{
	Properties: map[string]types.Property{
		"volumeNumber": types.NewIntegerNumberProperty(),
		"section":      types.NewIntegerNumberProperty(),
		"title":        types.NewKeywordProperty(),
		"works": &types.NestedProperty{
			Properties: map[string]types.Property{
				"id":    types.NewKeywordProperty(),
				"code":  types.NewKeywordProperty(),
				"title": types.NewKeywordProperty(),
			},
		},
	},
}

// structs for works tree data without content
type Work struct {
	Id           string    `json:"id"`
	Code         string    `json:"code"`
	Abbreviation *string   `json:"abbreviation"`
	Title        string    `json:"title"`
	Year         *string   `json:"year"`
	Sections     []Section `json:"sections"`
}

type Section struct {
	Heading    string    `json:"heading"`
	Paragraphs []string  `json:"paragraphs"`
	Sections   []Section `json:"sections"`
}

var WorkMapping = &types.TypeMapping{
	Properties: map[string]types.Property{
		"id":           types.NewKeywordProperty(),
		"code":         types.NewKeywordProperty(),
		"abbreviation": types.NewKeywordProperty(),
		"title":        types.NewKeywordProperty(),
		"year":         types.NewKeywordProperty(),
		"sections": &types.NestedProperty{
			Properties: map[string]types.Property{
				"heading":    types.NewKeywordProperty(),
				"paragraphs": types.NewKeywordProperty(),
				"sections": &types.NestedProperty{
					Properties: map[string]types.Property{
						"heading":    types.NewKeywordProperty(),
						"paragraphs": types.NewKeywordProperty(),
						"sections": &types.NestedProperty{
							Properties: map[string]types.Property{
								"heading":    types.NewKeywordProperty(),
								"paragraphs": types.NewKeywordProperty(),
								"sections": &types.NestedProperty{
									Properties: map[string]types.Property{
										"heading":    types.NewKeywordProperty(),
										"paragraphs": types.NewKeywordProperty(),
										"sections": &types.NestedProperty{
											Properties: map[string]types.Property{
												"heading":    types.NewKeywordProperty(),
												"paragraphs": types.NewKeywordProperty(),
												"sections": &types.NestedProperty{
													Properties: map[string]types.Property{
														"heading":    types.NewKeywordProperty(),
														"paragraphs": types.NewKeywordProperty(),
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	},
}

// structs for actual content (headings, paragraphs, footnotes, summaries), stored in a linear structure to make searching and fetching it simple
type Type string

const (
	Heading   Type = "heading"
	Paragraph Type = "paragraph"
	Footnote  Type = "footnote"
	Summary   Type = "summary"
)

type Content struct {
	Type       Type     `json:"type"`
	Id         string   `json:"id"`
	Ref        *string  `json:"ref"`
	FmtText    string   `json:"fmtText"`
	SearchText string   `json:"searchText"`
	Pages      []int32  `json:"pages"`
	FnRefs     []string `json:"fnRefs"`
	SummaryRef *string  `json:"summaryRef"`
	WorkId     string   `json:"workId"`
}

var ContentMapping = &types.TypeMapping{
	Properties: map[string]types.Property{
		"type":       types.NewKeywordProperty(),
		"id":         types.NewKeywordProperty(),
		"ref":        types.NewKeywordProperty(),
		"fmtText":    types.NewKeywordProperty(),
		"searchText": types.NewTextProperty(),
		"pages":      types.NewKeywordProperty(),
		"fnRefs":     types.NewKeywordProperty(),
		"summaryRef": types.NewKeywordProperty(),
		"workId":     types.NewKeywordProperty(),
	},
}
