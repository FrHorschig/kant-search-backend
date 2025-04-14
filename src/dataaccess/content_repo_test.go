//go:build integration
// +build integration

package dataaccess

import (
	"context"
	"testing"
	"time"

	"github.com/frhorschig/kant-search-backend/common/util"
	"github.com/frhorschig/kant-search-backend/dataaccess/internal/esmodel"
	"github.com/stretchr/testify/assert"
)

func TestContentRepo(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	repo := NewContentRepo(dbClient)

	workId := "work123"
	contentList := []esmodel.Content{
		{
			Type:       esmodel.Paragraph,
			Ref:        util.ToStrPtr("A123"),
			FmtText:    "formatted text 1",
			SearchText: "search text 1",
			Pages:      []int32{1, 2, 3},
			FnRefs:     []string{"fn1", "fn2"},
			SummaryRef: util.ToStrPtr("summ1"),
			WorkId:     workId,
		},
		{
			Type:       esmodel.Paragraph,
			Ref:        util.ToStrPtr("A124"),
			FmtText:    "formatted text 2",
			SearchText: "search text 2",
			Pages:      []int32{4, 5},
			FnRefs:     []string{},
			SummaryRef: nil,
			WorkId:     workId,
		},
	}

	// WHEN Insert
	err := repo.Insert(ctx, contentList)
	// THEN
	assert.Nil(t, err)
	for _, c := range contentList {
		assert.NotEmpty(t, c.Id)
	}
	refreshContents(t)

	// WHEN GetByWorkId
	res, err := repo.GetByWorkId(ctx, workId)
	// THEN
	assert.Nil(t, err)
	assert.Len(t, res, 2)
	assert.ElementsMatch(t,
		[]string{contentList[0].SearchText, contentList[1].SearchText},
		[]string{res[0].SearchText, res[1].SearchText},
	)

	// WHEN DeleteByWorkId
	err = repo.DeleteByWorkId(ctx, workId)
	// THEN
	assert.Nil(t, err)
	refreshContents(t)

	// WHEN GetByWorkId
	res, err = repo.GetByWorkId(ctx, workId)
	// THEN
	assert.Nil(t, err)
	assert.Len(t, res, 0)
}

func refreshContents(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := dbClient.Indices.Refresh().Index("contents").Do(ctx)
	assert.Nil(t, err)
}
