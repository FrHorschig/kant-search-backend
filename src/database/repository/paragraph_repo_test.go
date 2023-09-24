//go:build integration
// +build integration

package repository

import (
	"context"
	"testing"

	"github.com/FrHorschig/kant-search-backend/database/model"
	"github.com/stretchr/testify/assert"
)

func TestInsertParagraph(t *testing.T) {
	repo := &paragraphRepoImpl{db: testDb}
	ctx := context.Background()

	// WHEN
	id1, err1 := repo.Insert(ctx, model.Paragraph{Text: "text", Pages: []int32{1, 2}, WorkId: 1})
	id2, err2 := repo.Insert(ctx, model.Paragraph{Text: "text", Pages: []int32{1, 2}, WorkId: 1})

	// THEN
	assert.Nil(t, err1)
	assert.Nil(t, err2)
	assert.Greater(t, id1, int32(0))
	assert.Greater(t, id2, int32(0))

	testDb.Exec("DELETE FROM paragraphs")
}

func TestSelectAllParagraphs(t *testing.T) {
	repo := &paragraphRepoImpl{db: testDb}
	ctx := context.Background()

	// GIVEN
	id1, _ := repo.insertParagraphs(1, "text1")
	id2, _ := repo.insertParagraphs(1, "text2")
	id3, _ := repo.insertParagraphs(2, "text3")

	// WHEN
	paras1, err1 := repo.SelectAll(ctx, 1)
	paras2, err2 := repo.SelectAll(ctx, 2)

	// THEN
	assert.Nil(t, err1)
	assert.Len(t, paras1, 2)
	assert.Equal(t, id1, paras1[0].Id)
	assert.Equal(t, "text1", paras1[0].Text)
	assert.Equal(t, int32(1), paras1[0].WorkId)
	assert.Equal(t, id2, paras1[1].Id)
	assert.Equal(t, "text2", paras1[1].Text)
	assert.Equal(t, int32(1), paras1[1].WorkId)

	assert.Nil(t, err2)
	assert.Len(t, paras2, 1)
	assert.Equal(t, id3, paras2[0].Id)
	assert.Equal(t, "text3", paras2[0].Text)
	assert.Equal(t, int32(2), paras2[0].WorkId)

	testDb.Exec("DELETE FROM paragraphs")
}

func TestSelectAllParagraphsNoResults(t *testing.T) {
	repo := &paragraphRepoImpl{db: testDb}
	ctx := context.Background()

	// WHEN
	paras, err := repo.SelectAll(ctx, 1)

	// THEN
	assert.Nil(t, err)
	assert.Len(t, paras, 0)

	testDb.Exec("DELETE FROM paragraphs")
}

func TestSearchParagraphsSingleMatch(t *testing.T) {
	sut := &paragraphRepoImpl{db: testDb}
	ctx := context.Background()

	criteria := model.SearchCriteria{
		WorkIds:     []int32{1},
		SearchTerms: []string{"Maxime"},
	}

	// GIVEN
	id, _ := sut.insertParagraphs(1, "Maxime")
	sut.insertParagraphs(1, "Wille")
	sut.insertParagraphs(2, "Maxime")

	// WHEN
	matches, err := sut.Search(ctx, criteria)

	// THEN
	assert.Nil(t, err)
	assert.NotNil(t, matches)
	assert.Len(t, matches, 1)

	assert.Equal(t, id, matches[0].ElementId)
	assert.Contains(t, matches[0].Snippet, "Maxime")
	assert.Contains(t, matches[0].Text, "Maxime")
	assert.Equal(t, int32(1), matches[0].WorkId)

	testDb.Exec("DELETE FROM paragraphs")
}

func TestSearchParagraphsIgnoreSpecialCharacters(t *testing.T) {
	sut := &paragraphRepoImpl{db: testDb}
	ctx := context.Background()

	criteria := model.SearchCriteria{
		WorkIds:     []int32{1},
		SearchTerms: []string{`\&`},
	}

	// GIVEN
	sut.insertParagraphs(1, `\&`)

	// WHEN
	matches, err := sut.Search(ctx, criteria)

	// THEN
	assert.Nil(t, err)
	assert.NotNil(t, matches)
	assert.Len(t, matches, 0)

	testDb.Exec("DELETE FROM paragraphs")
}

func TestSearchParagraphsMultiMatch(t *testing.T) {
	sut := &paragraphRepoImpl{db: testDb}
	ctx := context.Background()

	criteria := model.SearchCriteria{
		WorkIds:     []int32{1, 2},
		SearchTerms: []string{"Maxime"},
	}

	// GIVEN
	id1, _ := sut.insertParagraphs(1, "Maxime")
	sut.insertParagraphs(1, "Wille")
	id3, _ := sut.insertParagraphs(2, "Maxime")

	// WHEN
	matches, err := sut.Search(ctx, criteria)

	// THEN
	assert.Nil(t, err)
	assert.NotNil(t, matches)
	assert.Len(t, matches, 2)

	assert.Equal(t, id1, matches[0].ElementId)
	assert.Contains(t, matches[0].Snippet, "Maxime")
	assert.Contains(t, matches[0].Text, "Maxime")
	assert.Equal(t, int32(1), matches[0].WorkId)

	assert.Equal(t, id3, matches[1].ElementId)
	assert.Contains(t, matches[1].Snippet, "Maxime")
	assert.Contains(t, matches[1].Text, "Maxime")
	assert.Equal(t, int32(2), matches[1].WorkId)

	testDb.Exec("DELETE FROM paragraphs")
}

func TestSearchParagraphsNoMatch(t *testing.T) {
	sut := &paragraphRepoImpl{db: testDb}
	ctx := context.Background()

	criteria := model.SearchCriteria{
		WorkIds:     []int32{1},
		SearchTerms: []string{"Maxime"},
	}

	// WHEN
	matches, err := sut.Search(ctx, criteria)

	// THEN
	assert.Nil(t, err)
	assert.Len(t, matches, 0)
}
