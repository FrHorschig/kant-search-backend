//go:build integration
// +build integration

package dataaccess

import (
	"context"
	"testing"
	"time"

	"github.com/frhorschig/kant-search-backend/common/util"
	"github.com/frhorschig/kant-search-backend/dataaccess/esmodel"
	"github.com/stretchr/testify/assert"
)

func TestContentRepo(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	repo := NewContentRepo(dbClient)

	workId := "work123"
	contents := []esmodel.Content{
		{
			Type:       esmodel.Footnote,
			Ref:        util.ToStrPtr("A121"),
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
			Ref:        util.ToStrPtr("A124"),
			FmtText:    "formatted text 3",
			SearchText: "search text 3",
			Pages:      []int32{4, 5},
			FnRefs:     []string{"fn3.4", "fn4.5"},
			WorkId:     workId,
		},
		{
			Type:       esmodel.Paragraph,
			Ref:        util.ToStrPtr("A124"),
			FmtText:    "formatted text 4",
			SearchText: "search text 4",
			Pages:      []int32{4, 5},
			FnRefs:     []string{"fn3.4", "fn4.5"},
			WorkId:     workId,
		},
		{
			Type:       esmodel.Summary,
			Ref:        util.ToStrPtr("A125"),
			FmtText:    "formatted text 5",
			SearchText: "search text 5",
			Pages:      []int32{4, 5},
			WorkId:     workId,
		},
	}

	// WHEN Insert
	err := repo.Insert(ctx, contents)
	// THEN
	assert.Nil(t, err)
	for _, c := range contents {
		assert.NotEmpty(t, c.Id)
	}
	refreshContents(t)

	// WHEN Get footnote
	fns, err := repo.GetFootnotesByWorkId(ctx, workId)
	// THEN
	assert.Nil(t, err)
	assert.Len(t, fns, 1)
	assert.Equal(t, contents[0].SearchText, fns[0].SearchText)
	// WHEN Get heading
	heads, err := repo.GetHeadingsByWorkId(ctx, workId)
	// THEN
	assert.Nil(t, err)
	assert.Len(t, heads, 1)
	assert.Equal(t, contents[1].SearchText, heads[0].SearchText)
	// WHEN Get paragraphs
	pars, err := repo.GetParagraphsByWorkId(ctx, workId)
	// THEN
	assert.Nil(t, err)
	assert.Len(t, heads, 1)
	assert.ElementsMatch(t,
		[]string{contents[2].SearchText, contents[3].SearchText},
		[]string{pars[0].SearchText, pars[1].SearchText},
	)
	// WHEN Get summary
	summ, err := repo.GetSummariesByWorkId(ctx, workId)
	// THEN
	assert.Nil(t, err)
	assert.Len(t, summ, 1)
	assert.Equal(t, contents[4].SearchText, summ[0].SearchText)

	// WHEN Delete
	err = repo.DeleteByWorkId(ctx, workId)
	// THEN
	assert.Nil(t, err)
	refreshContents(t)

	// WHEN Get footnote
	fns, err = repo.GetFootnotesByWorkId(ctx, workId)
	// THEN
	assert.Nil(t, err)
	assert.Len(t, fns, 0)
	// WHEN Get heading
	heads, err = repo.GetHeadingsByWorkId(ctx, workId)
	// THEN
	assert.Nil(t, err)
	assert.Len(t, heads, 0)
	// WHEN Get paragraphs
	pars, err = repo.GetParagraphsByWorkId(ctx, workId)
	// THEN
	assert.Nil(t, err)
	assert.Len(t, heads, 0)
	// WHEN Get summary
	summ, err = repo.GetSummariesByWorkId(ctx, workId)
	// THEN
	assert.Nil(t, err)
	assert.Len(t, summ, 0)
}

func refreshContents(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := dbClient.Indices.Refresh().Index("contents").Do(ctx)
	assert.Nil(t, err)
}
