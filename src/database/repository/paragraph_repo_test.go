package repository

import (
	"context"
	"testing"

	"github.com/FrHorschig/kant-search-backend/database/model"
	"github.com/stretchr/testify/assert"
)

var workId1 = int32(1)
var workId2 = int32(2)
var para1 = model.Paragraph{
	Text:   "Kant Wille Maxime",
	Pages:  []int32{1},
	WorkId: workId1,
}
var para2 = model.Paragraph{
	Text:   "Kant Kategorischer Imperativ",
	Pages:  []int32{2},
	WorkId: workId1,
}
var para3 = model.Paragraph{
	Text:   "Kant Vernunft Kategorie",
	Pages:  []int32{3},
	WorkId: workId2,
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
	para, err := repo.Select(ctx, workId1, id1)

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
	id3, _ := repo.Insert(ctx, para3)

	// WHEN
	paras1, err1 := repo.SelectAll(ctx, workId1)
	paras2, err2 := repo.SelectAll(ctx, workId2)

	// THEN
	assert.Nil(t, err1)
	assert.Len(t, paras1, 2)
	assert.Equal(t, id1, paras1[0].Id)
	assert.Equal(t, para1.Text, paras1[0].Text)
	assert.Equal(t, para1.Pages, paras1[0].Pages)
	assert.Equal(t, para1.WorkId, paras1[0].WorkId)
	assert.Equal(t, id2, paras1[1].Id)
	assert.Equal(t, para2.Text, paras1[1].Text)
	assert.Equal(t, para2.Pages, paras1[1].Pages)
	assert.Equal(t, para2.WorkId, paras1[1].WorkId)

	assert.Nil(t, err2)
	assert.Len(t, paras2, 1)
	assert.Equal(t, id3, paras2[0].Id)
	assert.Equal(t, para3.Text, paras2[0].Text)
	assert.Equal(t, para3.Pages, paras2[0].Pages)
	assert.Equal(t, para3.WorkId, paras2[0].WorkId)

	testDb.Exec("DELETE FROM paragraphs")
}
