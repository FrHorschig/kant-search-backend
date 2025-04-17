package esmodel

import (
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/frhorschig/kant-search-backend/common/util"
)

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
		"section":      &types.IntegerNumberProperty{Index: util.FalsePtr()},
		"title":        &types.TextProperty{Index: util.FalsePtr()},
		"works": &types.NestedProperty{
			Properties: map[string]types.Property{
				"id":    types.NewKeywordProperty(),
				"code":  &types.TextProperty{Index: util.FalsePtr()},
				"title": &types.TextProperty{Index: util.FalsePtr()},
			},
		},
	},
}

// structs for works tree data without content
type Work struct {
	Id           string    `json:"id"`
	Ordinal      int32     `json:"ordinal"`
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
		"ordinal":      types.NewIntegerNumberProperty(),
		"code":         &types.TextProperty{Index: util.FalsePtr()},
		"abbreviation": &types.TextProperty{Index: util.FalsePtr()},
		"title":        &types.TextProperty{Index: util.FalsePtr()},
		"year":         &types.TextProperty{Index: util.FalsePtr()},
		"sections": &types.NestedProperty{
			Properties: map[string]types.Property{
				"heading":    &types.TextProperty{Index: util.FalsePtr()},
				"paragraphs": &types.TextProperty{Index: util.FalsePtr()},
				"sections": &types.NestedProperty{
					Properties: map[string]types.Property{
						"heading":    &types.TextProperty{Index: util.FalsePtr()},
						"paragraphs": &types.TextProperty{Index: util.FalsePtr()},
						"sections": &types.NestedProperty{
							Properties: map[string]types.Property{
								"heading":    &types.TextProperty{Index: util.FalsePtr()},
								"paragraphs": &types.TextProperty{Index: util.FalsePtr()},
								"sections": &types.NestedProperty{
									Properties: map[string]types.Property{
										"heading":    &types.TextProperty{Index: util.FalsePtr()},
										"paragraphs": &types.TextProperty{Index: util.FalsePtr()},
										"sections": &types.NestedProperty{
											Properties: map[string]types.Property{
												"heading":    &types.TextProperty{Index: util.FalsePtr()},
												"paragraphs": &types.TextProperty{Index: util.FalsePtr()},
												"sections": &types.NestedProperty{
													Properties: map[string]types.Property{
														"heading":    &types.TextProperty{Index: util.FalsePtr()},
														"paragraphs": &types.TextProperty{Index: util.FalsePtr()},
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

// structs for actual content (headings, paragraphs, footnotes, summaries), stored in a linear structure to simplify searching and fetching
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
	Ordinal    int32    `json:"ordinal"`
	Ref        *string  `json:"ref"`
	FmtText    string   `json:"fmtText"`
	TocText    *string  `json:"tocText"`
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
		"ordinal":    types.NewIntegerNumberProperty(),
		"ref":        &types.TextProperty{Index: util.FalsePtr()},
		"fmtText":    &types.TextProperty{Index: util.FalsePtr()},
		"tocText":    &types.TextProperty{Index: util.FalsePtr()},
		"searchText": types.NewTextProperty(),
		"pages":      &types.TextProperty{Index: util.FalsePtr()},
		"fnRefs":     &types.TextProperty{Index: util.FalsePtr()},
		"summaryRef": &types.TextProperty{Index: util.FalsePtr()},
		"workId":     types.NewKeywordProperty(),
	},
}
