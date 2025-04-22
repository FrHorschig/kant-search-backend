//go:build integration
// +build integration

package dataaccess

import (
	"context"
	"testing"
	"time"

	"github.com/frhorschig/kant-search-backend/common/util"
	"github.com/frhorschig/kant-search-backend/dataaccess/esmodel"
	"github.com/frhorschig/kant-search-backend/dataaccess/model"
	"github.com/stretchr/testify/assert"
)

func TestContentInsertGetDelete(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	sut := NewContentRepo(dbClient)

	workId := "work123"
	contents := []esmodel.Content{
		{
			Type:       esmodel.Footnote,
			Ref:        util.StrPtr("A121"),
			FmtText:    "formatted text 1",
			SearchText: "search text 1",
			Pages:      []int32{1, 2, 3},
			WorkId:     workId,
		},
		{
			Type:       esmodel.Heading,
			FmtText:    "formatted text 2",
			SearchText: "search text 2",
			Pages:      []int32{1, 2, 3},
			FnRefs:     []string{"fn1.2", "fn2.3"},
			WorkId:     workId,
		},
		{
			Type:       esmodel.Paragraph,
			FmtText:    "formatted text 3",
			SearchText: "search text 3",
			Pages:      []int32{4, 5},
			FnRefs:     []string{"fn3.4", "fn4.5"},
			WorkId:     workId,
		},
		{
			Type:       esmodel.Paragraph,
			FmtText:    "formatted text 4",
			SearchText: "search text 4",
			Pages:      []int32{4, 5},
			FnRefs:     []string{"fn3.4", "fn4.5"},
			WorkId:     workId,
		},
		{
			Type:       esmodel.Summary,
			Ref:        util.StrPtr("A125"),
			FmtText:    "formatted text 5",
			SearchText: "search text 5",
			Pages:      []int32{4, 5},
			WorkId:     workId,
		},
	}

	// WHEN Insert
	err := sut.Insert(ctx, contents)
	// THEN
	assert.Nil(t, err)
	for _, c := range contents {
		assert.NotEmpty(t, c.Id)
	}
	refreshContents(t)

	// WHEN Get footnote
	fns, err := sut.GetFootnotesByWorkId(ctx, workId)
	// THEN
	assert.Nil(t, err)
	assert.Len(t, fns, 1)
	assert.Equal(t, contents[0].SearchText, fns[0].SearchText)
	// WHEN Get heading
	heads, err := sut.GetHeadingsByWorkId(ctx, workId)
	// THEN
	assert.Nil(t, err)
	assert.Len(t, heads, 1)
	assert.Equal(t, contents[1].SearchText, heads[0].SearchText)
	// WHEN Get paragraphs
	pars, err := sut.GetParagraphsByWorkId(ctx, workId)
	// THEN
	assert.Nil(t, err)
	assert.Len(t, heads, 1)
	assert.ElementsMatch(t,
		[]string{contents[2].SearchText, contents[3].SearchText},
		[]string{pars[0].SearchText, pars[1].SearchText},
	)
	// WHEN Get summary
	summ, err := sut.GetSummariesByWorkId(ctx, workId)
	// THEN
	assert.Nil(t, err)
	assert.Len(t, summ, 1)
	assert.Equal(t, contents[4].SearchText, summ[0].SearchText)

	// WHEN Delete
	err = sut.DeleteByWorkId(ctx, workId)
	// THEN
	assert.Nil(t, err)
	refreshContents(t)

	// WHEN Get footnote
	fns, err = sut.GetFootnotesByWorkId(ctx, workId)
	// THEN
	assert.Nil(t, err)
	assert.Len(t, fns, 0)
	// WHEN Get heading
	heads, err = sut.GetHeadingsByWorkId(ctx, workId)
	// THEN
	assert.Nil(t, err)
	assert.Len(t, heads, 0)
	// WHEN Get paragraphs
	pars, err = sut.GetParagraphsByWorkId(ctx, workId)
	// THEN
	assert.Nil(t, err)
	assert.Len(t, heads, 0)
	// WHEN Get summary
	summ, err = sut.GetSummariesByWorkId(ctx, workId)
	// THEN
	assert.Nil(t, err)
	assert.Len(t, summ, 0)
}

func TestSearch(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	sut := NewContentRepo(dbClient)

	workId := "work123"
	workId2 := "456work"
	testdata := []struct {
		name        string
		dbInput     []esmodel.Content
		searchTerms *model.AstNode
		options     model.SearchOptions
		hitCount    int
	}{
		{
			name: "test complex query",
			dbInput: []esmodel.Content{
				{Type: esmodel.Paragraph, SearchText: "dog night bird", WorkId: workId},
				{Type: esmodel.Paragraph, SearchText: "cat night bird", WorkId: workId},
				{Type: esmodel.Paragraph, SearchText: "dog mice night bird", WorkId: workId},
				{Type: esmodel.Paragraph, SearchText: "dog mouse night bird", WorkId: workId},
				{Type: esmodel.Paragraph, SearchText: "dog knight bird", WorkId: workId},
				{Type: esmodel.Paragraph, SearchText: "cat night burd", WorkId: workId},
				{Type: esmodel.Paragraph, SearchText: "dog night bird 2", WorkId: workId2},
				{Type: esmodel.Paragraph, SearchText: "cat night bird 2", WorkId: workId2},
			},
			searchTerms: &model.AstNode{ // (dog | cat) & !mouse & "night bird"
				Token: newAnd(),
				Left: &model.AstNode{
					Token: newAnd(),
					Left: &model.AstNode{
						Token: newOr(),
						Left:  &model.AstNode{Token: newWord("dog")},
						Right: &model.AstNode{Token: newWord("cat")},
					},
					Right: &model.AstNode{
						Token: newNot(),
						Left:  &model.AstNode{Token: newWord("mouse")},
					},
				},
				Right: &model.AstNode{Token: newPhrase("night bird")},
			},
			options:  model.SearchOptions{WorkIds: []string{workId}},
			hitCount: 3,
		},
		{
			name: "test includeHeadings option",
			dbInput: []esmodel.Content{
				{Type: esmodel.Paragraph, SearchText: "paragraph text", WorkId: workId},
				{Type: esmodel.Heading, SearchText: "heading text", WorkId: workId},
				{Type: esmodel.Footnote, SearchText: "footnote text", WorkId: workId},
				{Type: esmodel.Summary, SearchText: "summary text", WorkId: workId},
			},
			searchTerms: &model.AstNode{Token: newWord("text")},
			options: model.SearchOptions{
				WorkIds:         []string{workId},
				IncludeHeadings: true,
			},
			hitCount: 2,
		},
		{
			name: "test includeFootnotes option",
			dbInput: []esmodel.Content{
				{Type: esmodel.Paragraph, SearchText: "paragraph text", WorkId: workId},
				{Type: esmodel.Heading, SearchText: "heading text", WorkId: workId},
				{Type: esmodel.Footnote, SearchText: "footnote text", WorkId: workId},
				{Type: esmodel.Summary, SearchText: "summary text", WorkId: workId},
			},
			searchTerms: &model.AstNode{Token: newWord("text")},
			options: model.SearchOptions{
				WorkIds:          []string{workId},
				IncludeFootnotes: true,
			},
			hitCount: 2,
		},
		{
			name: "test includeSummaries option",
			dbInput: []esmodel.Content{
				{Type: esmodel.Paragraph, SearchText: "paragraph text", WorkId: workId},
				{Type: esmodel.Heading, SearchText: "heading text", WorkId: workId},
				{Type: esmodel.Footnote, SearchText: "footnote text", WorkId: workId},
				{Type: esmodel.Summary, SearchText: "summary text", WorkId: workId},
			},
			searchTerms: &model.AstNode{Token: newWord("text")},
			options: model.SearchOptions{
				WorkIds:          []string{workId},
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

		err = sut.DeleteByWorkId(ctx, workId)
		if err != nil {
			t.Fatal("content deletion failure")
		}
		sut.DeleteByWorkId(ctx, workId2)
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
