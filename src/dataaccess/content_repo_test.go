//go:build integration
// +build integration

package dataaccess

import (
	"context"
	"testing"
	"time"

	"github.com/frhorschig/kant-search-backend/common/util"
	"github.com/frhorschig/kant-search-backend/dataaccess/model"
	"github.com/stretchr/testify/assert"
)

func TestContentInsertGetDelete(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	sut := NewContentRepo(dbClient)

	workCode := "work123"
	contents := []model.Content{
		{
			Type:       model.Footnote,
			Ordinal:    1,
			Ref:        util.StrPtr("A121"),
			FmtText:    "formatted text 1",
			SearchText: "search text 1",
			Pages:      []int32{1, 2, 3},
			WorkCode:   workCode,
		},
		{
			Type:       model.Heading,
			Ordinal:    2,
			FmtText:    "formatted text 2",
			SearchText: "search text 2",
			Pages:      []int32{1, 2, 3},
			FnRefs:     []string{"fn1.2", "fn2.3"},
			WorkCode:   workCode,
		},
		{
			Type:       model.Paragraph,
			Ordinal:    3,
			FmtText:    "formatted text 3",
			SearchText: "search text 3",
			Pages:      []int32{4, 5},
			FnRefs:     []string{"fn3.4", "fn4.5"},
			WorkCode:   workCode,
		},
		{
			Type:       model.Paragraph,
			Ordinal:    4,
			FmtText:    "formatted text 4",
			SearchText: "search text 4",
			Pages:      []int32{4, 5},
			FnRefs:     []string{"fn3.4", "fn4.5"},
			WorkCode:   workCode,
		},
		{
			Type:       model.Summary,
			Ordinal:    5,
			Ref:        util.StrPtr("A125"),
			FmtText:    "formatted text 5",
			SearchText: "search text 5",
			Pages:      []int32{4, 5},
			WorkCode:   workCode,
		},
	}

	// WHEN Insert
	err := sut.Insert(ctx, contents)
	// THEN
	assert.Nil(t, err)
	refreshContents(t)

	// WHEN Get footnote
	fns, err := sut.GetFootnotesByWork(ctx, workCode, []int32{})
	// THEN
	assert.Nil(t, err)
	assert.Len(t, fns, 1)
	assert.Equal(t, contents[0].SearchText, fns[0].SearchText)
	// WHEN Get heading
	heads, err := sut.GetHeadingsByWork(ctx, workCode, []int32{})
	// THEN
	assert.Nil(t, err)
	assert.Len(t, heads, 1)
	assert.Equal(t, contents[1].SearchText, heads[0].SearchText)
	// WHEN Get paragraphs
	pars, err := sut.GetParagraphsByWork(ctx, workCode, []int32{})
	// THEN
	assert.Nil(t, err)
	assert.Len(t, heads, 1)
	assert.ElementsMatch(t,
		[]string{contents[2].SearchText, contents[3].SearchText},
		[]string{pars[0].SearchText, pars[1].SearchText},
	)
	// WHEN Get single paragraph
	pars, err = sut.GetParagraphsByWork(ctx, workCode, []int32{4})
	// THEN
	assert.Nil(t, err)
	assert.Len(t, pars, 1)
	assert.Equal(t, contents[3].SearchText, pars[0].SearchText)
	// WHEN Get summary
	summ, err := sut.GetSummariesByWork(ctx, workCode, []int32{})
	// THEN
	assert.Nil(t, err)
	assert.Len(t, summ, 1)
	assert.Equal(t, contents[4].SearchText, summ[0].SearchText)

	// WHEN Delete
	err = sut.DeleteByWork(ctx, workCode)
	// THEN
	assert.Nil(t, err)
	refreshContents(t)

	// WHEN Get footnote
	fns, err = sut.GetFootnotesByWork(ctx, workCode, []int32{})
	// THEN
	assert.Nil(t, err)
	assert.Len(t, fns, 0)
	// WHEN Get heading
	heads, err = sut.GetHeadingsByWork(ctx, workCode, []int32{})
	// THEN
	assert.Nil(t, err)
	assert.Len(t, heads, 0)
	// WHEN Get paragraphs
	pars, err = sut.GetParagraphsByWork(ctx, workCode, []int32{})
	// THEN
	assert.Nil(t, err)
	assert.Len(t, heads, 0)
	// WHEN Get summary
	summ, err = sut.GetSummariesByWork(ctx, workCode, []int32{})
	// THEN
	assert.Nil(t, err)
	assert.Len(t, summ, 0)
}

func TestSearch(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	sut := NewContentRepo(dbClient)

	workCode := "work123"
	workCode2 := "456work"
	testdata := []struct {
		name        string
		dbInput     []model.Content
		searchTerms *model.SearchTermNode
		options     model.SearchOptions
		hitCount    int
	}{
		{
			name: "test complex query",
			dbInput: []model.Content{
				{Type: model.Paragraph, SearchText: "dog night bird", WorkCode: workCode},
				{Type: model.Paragraph, SearchText: "cat night bird", WorkCode: workCode},
				{Type: model.Paragraph, SearchText: "dog mice night bird", WorkCode: workCode},
				{Type: model.Paragraph, SearchText: "dog mouse night bird", WorkCode: workCode},
				{Type: model.Paragraph, SearchText: "dog knight bird", WorkCode: workCode},
				{Type: model.Paragraph, SearchText: "cat night burd", WorkCode: workCode},
				{Type: model.Paragraph, SearchText: "dog night bird 2", WorkCode: workCode2},
				{Type: model.Paragraph, SearchText: "cat night bird 2", WorkCode: workCode2},
			},
			searchTerms: &model.SearchTermNode{ // (dog | cat) & !mouse & "night bird"
				Token: newAnd(),
				Left: &model.SearchTermNode{
					Token: newAnd(),
					Left: &model.SearchTermNode{
						Token: newOr(),
						Left:  &model.SearchTermNode{Token: newWord("dog")},
						Right: &model.SearchTermNode{Token: newWord("cat")},
					},
					Right: &model.SearchTermNode{
						Token: newNot(),
						Left:  &model.SearchTermNode{Token: newWord("mouse")},
					},
				},
				Right: &model.SearchTermNode{Token: newPhrase("night bird")},
			},
			options:  model.SearchOptions{WorkCodes: []string{workCode}},
			hitCount: 3,
		},
		{
			name: "test includeHeadings option",
			dbInput: []model.Content{
				{Type: model.Paragraph, SearchText: "paragraph text", WorkCode: workCode},
				{Type: model.Heading, SearchText: "heading text", WorkCode: workCode},
				{Type: model.Footnote, SearchText: "footnote text", WorkCode: workCode},
				{Type: model.Summary, SearchText: "summary text", WorkCode: workCode},
			},
			searchTerms: &model.SearchTermNode{Token: newWord("text")},
			options: model.SearchOptions{
				WorkCodes:       []string{workCode},
				IncludeHeadings: true,
			},
			hitCount: 2,
		},
		{
			name: "test includeFootnotes option",
			dbInput: []model.Content{
				{Type: model.Paragraph, SearchText: "paragraph text", WorkCode: workCode},
				{Type: model.Heading, SearchText: "heading text", WorkCode: workCode},
				{Type: model.Footnote, SearchText: "footnote text", WorkCode: workCode},
				{Type: model.Summary, SearchText: "summary text", WorkCode: workCode},
			},
			searchTerms: &model.SearchTermNode{Token: newWord("text")},
			options: model.SearchOptions{
				WorkCodes:        []string{workCode},
				IncludeFootnotes: true,
			},
			hitCount: 2,
		},
		{
			name: "test includeSummaries option",
			dbInput: []model.Content{
				{Type: model.Paragraph, SearchText: "paragraph text", WorkCode: workCode},
				{Type: model.Heading, SearchText: "heading text", WorkCode: workCode},
				{Type: model.Footnote, SearchText: "footnote text", WorkCode: workCode},
				{Type: model.Summary, SearchText: "summary text", WorkCode: workCode},
			},
			searchTerms: &model.SearchTermNode{Token: newWord("text")},
			options: model.SearchOptions{
				WorkCodes:        []string{workCode},
				IncludeSummaries: true,
			},
			hitCount: 2,
		},
	}

	for _, tc := range testdata {
		err := sut.Insert(ctx, tc.dbInput)
		if err != nil {
			t.Fatal("content insertion failure")
		}
		refreshContents(t)

		t.Run(tc.name, func(t *testing.T) {
			result, err := sut.Search(ctx, tc.searchTerms, tc.options)
			assert.Nil(t, err)
			assert.Len(t, result, tc.hitCount)
		})

		err = sut.DeleteByWork(ctx, workCode)
		if err != nil {
			t.Fatal("content deletion failure")
		}
		sut.DeleteByWork(ctx, workCode2)
		if err != nil {
			t.Fatal("content deletion failure")
		}
	}
}

func refreshContents(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := dbClient.Indices.Refresh().Index("contents").Do(ctx)
	assert.Nil(t, err)
}

func newAnd() *model.Token {
	return &model.Token{IsAnd: true, Text: "&"}
}
func newOr() *model.Token {
	return &model.Token{IsOr: true, Text: "|"}
}
func newNot() *model.Token {
	return &model.Token{IsNot: true, Text: "!"}
}
func newWord(text string) *model.Token {
	return &model.Token{IsWord: true, Text: text}
}
func newPhrase(text string) *model.Token {
	return &model.Token{IsPhrase: true, Text: text}
}
