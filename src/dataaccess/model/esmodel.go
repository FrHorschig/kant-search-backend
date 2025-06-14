package model

import (
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/frhorschig/kant-search-backend/common/util"
)

// structs for volume-works metadata
type Volume struct {
	VolumeNumber int32  `json:"volumeNumber"`
	Title        string `json:"title"`
	Works        []Work `json:"works"`
}

type Work struct {
	Code       string    `json:"code"`
	Siglum     *string   `json:"siglum,omitempty"`
	Title      string    `json:"title"`
	Year       string    `json:"year"`
	Ordinal    int32     `json:"ordinal"`
	Paragraphs []int32   `json:"paragraphs"`
	Sections   []Section `json:"sections"`
}

type Section struct {
	Heading    int32     `json:"heading"`
	Paragraphs []int32   `json:"paragraphs"`
	Sections   []Section `json:"sections"`
}

var VolumeMapping = &types.TypeMapping{
	Properties: map[string]types.Property{
		"volumeNumber": types.NewIntegerNumberProperty(),
		"title":        &types.TextProperty{Index: util.FalsePtr()},
		// TODO do we really need this here, would ObjectProperty not be enough?
		"works": &types.TypeMapping{
			Properties: map[string]types.Property{
				"code":       types.NewKeywordProperty(),
				"siglum":     &types.TextProperty{Index: util.FalsePtr()},
				"title":      &types.TextProperty{Index: util.FalsePtr()},
				"year":       &types.TextProperty{Index: util.FalsePtr()},
				"ordinal":    &types.TextProperty{Index: util.FalsePtr()},
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

type Analyzer string

const (
	NoStemming     Analyzer = "noStemming"
	GermanStemming Analyzer = "germanStemming"
)

type Content struct {
	// text data
	FmtText    string  `json:"fmtText"`
	TocText    *string `json:"tocText"` // only for headings
	SearchText string  `json:"searchText"`

	// sort and filter fields
	Type     Type   `json:"type"`
	Ordinal  int32  `json:"ordinal"`
	WorkCode string `json:"workCode"`

	// metadata
	Pages        []int32           `json:"pages"`
	PageByIndex  []IndexNumberPair `json:"pageByIndex"`
	LineByIndex  []IndexNumberPair `json:"lineByIndex"`
	WordIndexMap map[int32]int32   `json:"wordIndexMap"`
	FnRefs       []string          `json:"fnRefs"`     // not for footnotes
	SummaryRef   *string           `json:"summaryRef"` // only for paragraphs
	Ref          *string           `json:"ref"`        // for fns and summaries
}

var ContentMapping = &types.TypeMapping{
	Properties: map[string]types.Property{
		"fmtText": &types.TextProperty{Index: util.FalsePtr()},
		"tocText": &types.TextProperty{Index: util.FalsePtr()},
		"searchText": types.TextProperty{
			Fields: map[string]types.Property{
				string(NoStemming): &types.TextProperty{
					Analyzer: util.StrPtr(string(NoStemming)),
				},
				string(GermanStemming): &types.TextProperty{
					Analyzer: util.StrPtr(string(GermanStemming)),
				},
			},
		},

		"ref":      &types.TextProperty{Index: util.FalsePtr()},
		"workCode": types.NewKeywordProperty(),
		"type":     types.NewKeywordProperty(),
		"ordinal":  types.NewIntegerNumberProperty(),

		"pages":        &types.TextProperty{Index: util.FalsePtr()},
		"pageByIndex":  &types.ObjectProperty{Enabled: util.FalsePtr()},
		"lineByIndex":  &types.ObjectProperty{Enabled: util.FalsePtr()},
		"wordIndexMap": &types.ObjectProperty{Enabled: util.FalsePtr()},
		"fnRefs":       &types.TextProperty{Index: util.FalsePtr()},
		"summaryRef":   &types.TextProperty{Index: util.FalsePtr()},
	},
}
