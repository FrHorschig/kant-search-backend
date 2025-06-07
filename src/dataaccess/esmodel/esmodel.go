package esmodel

import (
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/frhorschig/kant-search-backend/common/util"
)

// structs for volume-works metadata
type Volume struct {
	VolumeNumber int32  `json:"volumeNumber"`
	Section      int32  `json:"section"`
	Title        string `json:"title"`
	Works        []Work `json:"works"`
}

type Work struct {
	Code         string    `json:"code"`
	Abbreviation *string   `json:"abbreviation"`
	Title        string    `json:"title"`
	Year         *string   `json:"year"`
	Ordinal      int32     `json:"ordinal"`
	Paragraphs   []int32   `json:"paragraphs"`
	Sections     []Section `json:"sections"`
}

type Section struct {
	Heading    int32     `json:"heading"`
	Paragraphs []int32   `json:"paragraphs"`
	Sections   []Section `json:"sections"`
}

var VolumeMapping = &types.TypeMapping{
	Properties: map[string]types.Property{
		"volumeNumber": types.NewIntegerNumberProperty(),
		"section":      &types.IntegerNumberProperty{Index: util.FalsePtr()},
		"title":        &types.TextProperty{Index: util.FalsePtr()},
		"works": &types.TypeMapping{
			Properties: map[string]types.Property{
				"code":       types.NewKeywordProperty(),
				"paragraphs": &types.IntegerNumberProperty{Index: util.FalsePtr()},
				"sections": &types.NestedProperty{
					Properties: map[string]types.Property{
						"heading":    &types.IntegerNumberProperty{Index: util.FalsePtr()},
						"paragraphs": &types.IntegerNumberProperty{Index: util.FalsePtr()},
						"sections": &types.NestedProperty{
							Properties: map[string]types.Property{
								"heading":    &types.IntegerNumberProperty{Index: util.FalsePtr()},
								"paragraphs": &types.IntegerNumberProperty{Index: util.FalsePtr()},
								"sections": &types.NestedProperty{
									Properties: map[string]types.Property{
										"heading":    &types.IntegerNumberProperty{Index: util.FalsePtr()},
										"paragraphs": &types.IntegerNumberProperty{Index: util.FalsePtr()},
										"sections": &types.NestedProperty{
											Properties: map[string]types.Property{
												"heading":    &types.IntegerNumberProperty{Index: util.FalsePtr()},
												"paragraphs": &types.IntegerNumberProperty{Index: util.FalsePtr()},
												"sections": &types.NestedProperty{
													Properties: map[string]types.Property{
														"heading":    &types.IntegerNumberProperty{Index: util.FalsePtr()},
														"paragraphs": &types.IntegerNumberProperty{Index: util.FalsePtr()},
														"sections": &types.NestedProperty{
															Properties: map[string]types.Property{
																"heading":    &types.IntegerNumberProperty{Index: util.FalsePtr()},
																"paragraphs": &types.IntegerNumberProperty{Index: util.FalsePtr()},
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
	Ordinal    int32    `json:"ordinal"` // used for sorting search results
	Ref        *string  `json:"ref"`     // only for footnotes and summaries
	FmtText    string   `json:"fmtText"`
	TocText    *string  `json:"tocText"`
	SearchText string   `json:"searchText"`
	Pages      []int32  `json:"pages"`
	FnRefs     []string `json:"fnRefs"`     // not for footnotes
	SummaryRef *string  `json:"summaryRef"` // only for paragraphs
	WorkCode   string   `json:"workCode"`
}

var ContentMapping = &types.TypeMapping{
	Properties: map[string]types.Property{
		"type":       types.NewKeywordProperty(),
		"ordinal":    types.NewIntegerNumberProperty(),
		"ref":        &types.TextProperty{Index: util.FalsePtr()},
		"fmtText":    &types.TextProperty{Index: util.FalsePtr()},
		"tocText":    &types.TextProperty{Index: util.FalsePtr()},
		"searchText": types.NewTextProperty(),
		"pages":      &types.TextProperty{Index: util.FalsePtr()},
		"fnRefs":     &types.TextProperty{Index: util.FalsePtr()},
		"summaryRef": &types.TextProperty{Index: util.FalsePtr()},
		"workCode":   types.NewKeywordProperty(),
	},
}
