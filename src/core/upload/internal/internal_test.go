package internal

import (
	"testing"

	"github.com/frhorschig/kant-search-backend/common/util"
	"github.com/frhorschig/kant-search-backend/core/upload/internal/common/testutil"
	"github.com/frhorschig/kant-search-backend/core/upload/internal/metadatamapping/metadatamapping/metadata"
	"github.com/frhorschig/kant-search-backend/core/upload/internal/metadatamapping/metadatamapping/metadata/mocks"
	"github.com/frhorschig/kant-search-backend/dataaccess/model"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestXmlMapping(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	md := mocks.NewMockMetadata(ctrl)
	volMd := metadata.VolumeMetadata{
		VolumeNumber: 2,
		Title:        "vol 2 title",
		Works: []metadata.WorkMetadata{
			{Code: "C1"},
			{Code: "C2", Siglum: util.StrPtr("S2")},
		},
	}
	xmlMapper := NewXmlMapper(md)
	md.EXPECT().Read(volMd.VolumeNumber).Return(volMd, nil)

	const xml = `<?xml version="1.0" encoding="UTF-8" standalone="no"?>
<kant_abt1>
  <band nr="2">
    <titel>volume title</titel>

    <hauptteil>
      <hj>1234</hj>

      <seite nr="001"/>
      <h1>work 1</h1>
      <p><zeile nr="01"/>paragraph 1</p>

      <seite nr="002"/>
      <h2><zeile nr="01"/>heading 1.2</h2>
      <p><zeile nr="02"/>paragraph 2<fr seite="002" nr="1"/></p>
      <p><zeile nr="03"/>paragraph 3</p>

      <hj>5678</hj>
      <seite nr="003"/>
      <h1>work 2</h1>
      <p><zeile nr="01"/>paragraph 4</p>

      <seite nr="004"/>
      <h2><zeile nr="01"/>heading 2.2</h2>
      <p><zeile nr="02"/>paragraph 5</p>
      <p><zeile nr="03"/>paragraph 6<fr seite="004" nr="1"/></p>
    </hauptteil>

    <fussnoten>
      <fn seite="002" nr="1">
        <p>fn paragraph 2.1</p>
      </fn>
      
      <fn seite="004" nr="1">
        <p>fn paragraph 4.1</p>
      </fn>
    </fussnoten>

    <randtexte>
      <randtext seite="002" anfang="03">
        <p>summ paragraph 2.3</p>
      </randtext>

      <randtext seite="004" anfang="02">
        <p>summ paragraph 4.2</p>
      </randtext>
    </randtexte>
  <band>
</kant_abt1>`

	vol, contents, err := xmlMapper.MapXml(2, xml)

	if err.DomainError != nil {
		println("ERROR: ", err.DomainError.Error())
		return
	}
	if err.TechnicalError != nil {
		println("ERROR: ", err.TechnicalError.Error())
		return
	}
	assert.False(t, err.HasError)
	testutil.AssertDbVolume(t, model.Volume{
		VolumeNumber: 2,
		Title:        "vol 2 title",
		Works: []model.Work{
			{
				Code:       "C1",
				Title:      "Work 1",
				Year:       "1234",
				Ordinal:    1,
				Paragraphs: []int32{1},
				Sections: []model.Section{{
					Heading:    2,
					Paragraphs: []int32{3, 6},
					Sections:   []model.Section{},
				}},
			},
			{
				Code:       "C2",
				Siglum:     util.StrPtr("S2"),
				Title:      "Work 2",
				Year:       "5678",
				Ordinal:    2,
				Paragraphs: []int32{1},
				Sections: []model.Section{{
					Heading:    2,
					Paragraphs: []int32{4, 5},
					Sections:   []model.Section{},
				}},
			},
		},
	}, vol)
	testutil.AssertDbContents(t, []model.Content{
		{
			Type:         model.Paragraph,
			Ordinal:      1,
			FmtText:      "<ks-meta-page>1</ks-meta-page> <ks-meta-line>1</ks-meta-line> paragraph 1",
			SearchText:   "paragraph 1",
			Pages:        []int32{1},
			PageByIndex:  []model.IndexNumberPair{{I: 0, Num: 1}},
			LineByIndex:  []model.IndexNumberPair{{I: 31, Num: 1}},
			WordIndexMap: map[int32]int32{0: 62},
			WorkCode:     "C1",
		},
		{
			Type:         model.Heading,
			Ordinal:      2,
			FmtText:      "<ks-fmt-h1><ks-meta-page>2</ks-meta-page> <ks-meta-line>1</ks-meta-line> heading 1.2</ks-fmt-h1>",
			SearchText:   "heading 1.2",
			TocText:      util.StrPtr("Heading 1.2"),
			Pages:        []int32{2},
			PageByIndex:  []model.IndexNumberPair{{I: 11, Num: 2}},
			LineByIndex:  []model.IndexNumberPair{{I: 42, Num: 1}},
			WordIndexMap: map[int32]int32{0: 73},
			WorkCode:     "C1",
		},
		{
			Type:         model.Paragraph,
			Ordinal:      3,
			FmtText:      "<ks-meta-line>2</ks-meta-line> paragraph 2 <ks-meta-fnref>2.1</ks-meta-fnref>",
			SearchText:   "paragraph 2",
			Pages:        []int32{2},
			FnRefs:       []string{"2.1"},
			PageByIndex:  []model.IndexNumberPair{},
			LineByIndex:  []model.IndexNumberPair{{I: 0, Num: 2}},
			WordIndexMap: map[int32]int32{0: 31},
			WorkCode:     "C1",
		},
		{
			Type:         model.Paragraph,
			Ordinal:      6,
			FmtText:      "<ks-meta-line>3</ks-meta-line> paragraph 3",
			SearchText:   "paragraph 3",
			Pages:        []int32{2},
			SummaryRef:   util.StrPtr("2.3"),
			PageByIndex:  []model.IndexNumberPair{},
			LineByIndex:  []model.IndexNumberPair{{I: 0, Num: 3}},
			WordIndexMap: map[int32]int32{0: 31},
			WorkCode:     "C1",
		},
		{
			Type:         model.Footnote,
			Ordinal:      4,
			Ref:          util.StrPtr("2.1"),
			FmtText:      "fn paragraph 2.1",
			SearchText:   "fn paragraph 2.1",
			Pages:        []int32{2},
			PageByIndex:  []model.IndexNumberPair{},
			LineByIndex:  []model.IndexNumberPair{},
			WordIndexMap: map[int32]int32{0: 0, 3: 3},
			WorkCode:     "C1",
		},
		{
			Type:         model.Summary,
			Ordinal:      5,
			Ref:          util.StrPtr("2.3"),
			FmtText:      "summ paragraph 2.3",
			SearchText:   "summ paragraph 2.3",
			Pages:        []int32{2},
			PageByIndex:  []model.IndexNumberPair{},
			LineByIndex:  []model.IndexNumberPair{},
			WordIndexMap: map[int32]int32{0: 0, 5: 5},
			WorkCode:     "C1",
		},
		{
			Type:         model.Paragraph,
			Ordinal:      1,
			FmtText:      "<ks-meta-page>3</ks-meta-page> <ks-meta-line>1</ks-meta-line> paragraph 4",
			SearchText:   "paragraph 4",
			Pages:        []int32{3},
			PageByIndex:  []model.IndexNumberPair{{I: 0, Num: 3}},
			LineByIndex:  []model.IndexNumberPair{{I: 31, Num: 1}},
			WordIndexMap: map[int32]int32{0: 62},
			WorkCode:     "C2",
		},
		{
			Type:         model.Heading,
			Ordinal:      2,
			FmtText:      "<ks-fmt-h1><ks-meta-page>4</ks-meta-page> <ks-meta-line>1</ks-meta-line> heading 2.2</ks-fmt-h1>",
			SearchText:   "heading 2.2",
			TocText:      util.StrPtr("Heading 2.2"),
			Pages:        []int32{4},
			PageByIndex:  []model.IndexNumberPair{{I: 11, Num: 4}},
			LineByIndex:  []model.IndexNumberPair{{I: 42, Num: 1}},
			WordIndexMap: map[int32]int32{0: 73},
			WorkCode:     "C2",
		},
		{
			Type:         model.Paragraph,
			Ordinal:      4,
			FmtText:      "<ks-meta-line>2</ks-meta-line> paragraph 5",
			SearchText:   "paragraph 5",
			Pages:        []int32{4},
			SummaryRef:   util.StrPtr("4.2"),
			PageByIndex:  []model.IndexNumberPair{},
			LineByIndex:  []model.IndexNumberPair{{I: 0, Num: 2}},
			WordIndexMap: map[int32]int32{0: 31},
			WorkCode:     "C2",
		},
		{
			Type:         model.Paragraph,
			Ordinal:      5,
			FmtText:      "<ks-meta-line>3</ks-meta-line> paragraph 6 <ks-meta-fnref>4.1</ks-meta-fnref>",
			SearchText:   "paragraph 6",
			Pages:        []int32{4},
			FnRefs:       []string{"4.1"},
			PageByIndex:  []model.IndexNumberPair{},
			LineByIndex:  []model.IndexNumberPair{{I: 0, Num: 3}},
			WordIndexMap: map[int32]int32{0: 31},
			WorkCode:     "C2",
		},
		{
			Type:         model.Footnote,
			Ordinal:      6,
			Ref:          util.StrPtr("4.1"),
			FmtText:      "fn paragraph 4.1",
			SearchText:   "fn paragraph 4.1",
			Pages:        []int32{4},
			PageByIndex:  []model.IndexNumberPair{},
			LineByIndex:  []model.IndexNumberPair{},
			WordIndexMap: map[int32]int32{0: 0, 3: 3},
			WorkCode:     "C2",
		},
		{
			Type:         model.Summary,
			Ordinal:      3,
			Ref:          util.StrPtr("4.2"),
			FmtText:      "summ paragraph 4.2",
			SearchText:   "summ paragraph 4.2",
			Pages:        []int32{4},
			PageByIndex:  []model.IndexNumberPair{},
			LineByIndex:  []model.IndexNumberPair{},
			WordIndexMap: map[int32]int32{0: 0, 5: 5},
			WorkCode:     "C2",
		},
	}, contents)
}
