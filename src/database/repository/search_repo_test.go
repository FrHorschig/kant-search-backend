//go:build integration
// +build integration

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
		WorkIds:     []int32{1},
		SearchTerms: []string{"Maxime"},
	}

	// GIVEN
	id, _ := paraRepo.insertParagraphs(1, "Maxime")
	paraRepo.insertParagraphs(1, "Wille")
	paraRepo.insertParagraphs(2, "Maxime")

	// WHEN
	matches, err := repo.SearchParagraphs(ctx, criteria)

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
	paraRepo := &paragraphRepoImpl{db: testDb}
	repo := &searchRepoImpl{db: testDb}
	ctx := context.Background()

	criteria := model.SearchCriteria{
		WorkIds:     []int32{1},
		SearchTerms: []string{"&"},
	}

	// GIVEN
	paraRepo.insertParagraphs(1, "&")

	// WHEN
	matches, err := repo.SearchParagraphs(ctx, criteria)

	// THEN
	assert.Nil(t, err)
	assert.NotNil(t, matches)
	assert.Len(t, matches, 0)

	testDb.Exec("DELETE FROM paragraphs")
}

func TestSearchParagraphsMultiMatch(t *testing.T) {
	paraRepo := &paragraphRepoImpl{db: testDb}
	repo := &searchRepoImpl{db: testDb}
	ctx := context.Background()

	criteria := model.SearchCriteria{
		WorkIds:     []int32{1, 2},
		SearchTerms: []string{"Maxime"},
	}

	// GIVEN
	id1, _ := paraRepo.insertParagraphs(1, "Maxime")
	paraRepo.insertParagraphs(1, "Wille")
	id3, _ := paraRepo.insertParagraphs(2, "Maxime")

	// WHEN
	matches, err := repo.SearchParagraphs(ctx, criteria)

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
	repo := &searchRepoImpl{db: testDb}
	ctx := context.Background()

	criteria := model.SearchCriteria{
		WorkIds:     []int32{1},
		SearchTerms: []string{"Maxime"},
	}

	// WHEN
	matches, err := repo.SearchParagraphs(ctx, criteria)

	// THEN
	assert.Nil(t, err)
	assert.Len(t, matches, 0)

	testDb.Exec("DELETE FROM paragraphs")
}
