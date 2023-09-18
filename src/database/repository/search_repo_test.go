package repository

import (
	"context"
	"testing"

	"github.com/FrHorschig/kant-search-backend/database/model"
	"github.com/stretchr/testify/assert"
)

func TestSearchParagraphsSingleMatch(t *testing.T) {
	paraRepo := &paragraphRepoImpl{db: testDb}
	repo := &searchRepoImpl{db: testDb}
	ctx := context.Background()

	criteria := model.SearchCriteria{
		SearchTerms: []string{"Maxime"},
		WorkIds:     []int32{workId1},
		Scope:       model.PARAGRAPH,
	}

	// GIVEN
	id, _ := paraRepo.Insert(ctx, para1)
	paraRepo.Insert(ctx, para2)
	paraRepo.Insert(ctx, para3)

	// WHEN
	matches, err := repo.SearchParagraphs(ctx, criteria)

	// THEN
	assert.Nil(t, err)
	assert.NotNil(t, matches)
	assert.Len(t, matches, 1)

	assert.Equal(t, id, matches[0].ElementId)
	assert.Contains(t, matches[0].Snippet, "Maxime")
	assert.Equal(t, para1.Pages, matches[0].Pages)
	assert.Equal(t, para1.WorkId, matches[0].WorkId)

	testDb.Exec("DELETE FROM paragraphs")
}

func TestSearchParagraphsMultiMatch(t *testing.T) {
	paraRepo := &paragraphRepoImpl{db: testDb}
	repo := &searchRepoImpl{db: testDb}
	ctx := context.Background()

	criteria := model.SearchCriteria{
		SearchTerms: []string{"Kant"},
		WorkIds:     []int32{workId1, workId2},
		Scope:       model.PARAGRAPH,
	}

	// GIVEN
	id1, _ := paraRepo.Insert(ctx, para1)
	id2, _ := paraRepo.Insert(ctx, para2)
	id3, _ := paraRepo.Insert(ctx, para3)

	// WHEN
	matches, err := repo.SearchParagraphs(ctx, criteria)

	// THEN
	assert.Nil(t, err)
	assert.NotNil(t, matches)
	assert.Len(t, matches, 3)

	assert.Equal(t, id1, matches[0].ElementId)
	assert.Contains(t, matches[0].Snippet, "Kant")
	assert.Equal(t, para1.Pages, matches[0].Pages)
	assert.Equal(t, para1.WorkId, matches[0].WorkId)

	assert.Equal(t, id2, matches[1].ElementId)
	assert.Contains(t, matches[1].Snippet, "Kant")
	assert.Equal(t, para2.Pages, matches[1].Pages)
	assert.Equal(t, para2.WorkId, matches[1].WorkId)

	assert.Equal(t, id3, matches[2].ElementId)
	assert.Contains(t, matches[2].Snippet, "Kant")
	assert.Equal(t, para3.Pages, matches[2].Pages)
	assert.Equal(t, para3.WorkId, matches[2].WorkId)

	testDb.Exec("DELETE FROM paragraphs")
}

func TestSearchParagraphsNoMatch(t *testing.T) {
	repo := &searchRepoImpl{db: testDb}
	ctx := context.Background()

	criteria := model.SearchCriteria{
		SearchTerms: []string{"Maxime"},
		WorkIds:     []int32{workId1},
		Scope:       model.PARAGRAPH,
	}

	// WHEN
	matches, err := repo.SearchParagraphs(ctx, criteria)

	// THEN
	assert.Nil(t, err)
	assert.Len(t, matches, 0)

	testDb.Exec("DELETE FROM paragraphs")
}
