package repository

import (
	"context"
	"testing"

	"github.com/FrHorschig/kant-search-backend/database/model"
	"github.com/stretchr/testify/assert"
)

var workId = int32(1)
var para1 = model.Paragraph{
	Text:   "text1",
	Pages:  []int32{1},
	WorkId: workId,
}
var para2 = model.Paragraph{
	Text:   "test2",
	Pages:  []int32{2, 3},
	WorkId: workId,
}

func TestInsertParagraph(t *testing.T) {
	repo := &paragraphRepoImpl{db: testDb}
	ctx := context.Background()

	// WHEN
	id1, err1 := repo.Insert(ctx, para1)
	id2, err2 := repo.Insert(ctx, para2)

	// THEN
	assert.Nil(t, err1)
	assert.Nil(t, err2)
	assert.Greater(t, id1, int32(0))
	assert.Greater(t, id2, int32(0))

	testDb.Exec("DELETE FROM paragraphs")
}

func TestSelectParagraph(t *testing.T) {
	repo := &paragraphRepoImpl{db: testDb}
	ctx := context.Background()

	// GIVEN
	id1, _ := repo.Insert(ctx, para1)
	repo.Insert(ctx, para2)

	// WHEN
	para, err := repo.Select(ctx, workId, id1)

	// THEN
	assert.Nil(t, err)
	assert.Equal(t, id1, para.Id)
	assert.Equal(t, para1.Text, para.Text)
	assert.Equal(t, para1.Pages, para.Pages)
	assert.Equal(t, para1.WorkId, para.WorkId)

	testDb.Exec("DELETE FROM paragraphs")
}

func TestSelectAllParagraphs(t *testing.T) {
	repo := &paragraphRepoImpl{db: testDb}
	ctx := context.Background()

	// GIVEN
	id1, _ := repo.Insert(ctx, para1)
	id2, _ := repo.Insert(ctx, para2)

	// WHEN
	paras, err := repo.SelectAll(ctx, workId)

	// THEN
	assert.Nil(t, err)
	assert.Len(t, paras, 2)

	assert.Equal(t, id1, paras[0].Id)
	assert.Equal(t, para1.Text, paras[0].Text)
	assert.Equal(t, para1.Pages, paras[0].Pages)
	assert.Equal(t, para1.WorkId, paras[0].WorkId)

	assert.Equal(t, id2, paras[1].Id)
	assert.Equal(t, para2.Text, paras[1].Text)
	assert.Equal(t, para2.Pages, paras[1].Pages)
	assert.Equal(t, para2.WorkId, paras[1].WorkId)

	testDb.Exec("DELETE FROM paragraphs")
}
