//go:build integration
// +build integration

package database

import (
	"context"
	"testing"

	"github.com/frhorschig/kant-search-backend/common/model"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestSentenceRepo_Insert(t *testing.T) {
	paragraphRepo := &paragraphRepoImpl{db: testDb}
	sut := NewSentenceRepo(testDb)
	ctx := context.Background()

	// GIVEN
	paragraphId, err := paragraphRepo.insertParagraphs(1, "text")
	sentences := []model.Sentence{
		{Text: "This is the first sentence", ParagraphId: paragraphId},
		{Text: "This is the second sentence", ParagraphId: paragraphId},
	}

	// WHEN
	ids, err := sut.Insert(ctx, sentences)

	// THEN
	assert.NoError(t, err)
	assert.Len(t, ids, len(sentences))
	contents := selectContents(t, ids)
	assert.Contains(t, contents, sentences[0].Text)
	assert.Contains(t, contents, sentences[1].Text)

	testDb.Exec("DELETE FROM sentences")
	testDb.Exec("DELETE FROM paragraphs")
}

func selectContents(t *testing.T, ids []int32) []string {
	rows, err := testDb.QueryContext(context.Background(), "SELECT content FROM sentences WHERE id = ANY($1)", pq.Array(ids))
	assert.NoError(t, err)
	defer rows.Close()
	var contents []string
	for rows.Next() {
		var content string
		err := rows.Scan(&content)
		assert.NoError(t, err)
		contents = append(contents, content)
	}
	return contents
}

func TestSearchSentencesSingleMatch(t *testing.T) {
	paraRepo := &paragraphRepoImpl{db: testDb}
	sut := &sentenceRepoImpl{db: testDb}
	ctx := context.Background()

	criteria := model.SearchCriteria{
		WorkIds:      []int32{1},
		SearchString: "Maxime",
	}

	// GIVEN
	pId1, _ := paraRepo.insertParagraphs(1, "Maxime Paragraph")
	pId2, _ := paraRepo.insertParagraphs(2, "Maxime")
	ids, _ := sut.insertSentences(pId1, []string{"Maxime", "Wille"})
	sut.insertSentences(pId2, []string{"Maxime"})

	// WHEN
	matches, err := sut.Search(ctx, criteria)

	// THEN
	assert.Nil(t, err)
	assert.NotNil(t, matches)
	assert.Len(t, matches, 1)

	assert.Equal(t, ids[0], matches[0].SentenceId)
	assert.Contains(t, matches[0].Snippet, "Maxime")
	assert.Contains(t, matches[0].Text, "Maxime")
	assert.NotContains(t, matches[0].Text, "Paragraph")
	assert.Equal(t, int32(1), matches[0].WorkId)

	testDb.Exec("DELETE FROM sentences")
	testDb.Exec("DELETE FROM paragraphs")
}

func TestSearchSentencesIgnoreSpecialCharacters(t *testing.T) {
	paraRepo := &paragraphRepoImpl{db: testDb}
	sut := &sentenceRepoImpl{db: testDb}
	ctx := context.Background()

	criteria := model.SearchCriteria{
		WorkIds:      []int32{1},
		SearchString: `\&`,
	}

	// GIVEN
	pId, _ := paraRepo.insertParagraphs(1, "Maxime Paragraph")
	sut.insertSentences(pId, []string{`\&`})

	// WHEN
	matches, err := sut.Search(ctx, criteria)

	// THEN
	assert.Nil(t, err)
	assert.NotNil(t, matches)
	assert.Len(t, matches, 0)

	testDb.Exec("DELETE FROM sentences")
	testDb.Exec("DELETE FROM paragraphs")
}

func TestSearchSentencesMultiMatch(t *testing.T) {
	paraRepo := &paragraphRepoImpl{db: testDb}
	sut := &sentenceRepoImpl{db: testDb}
	ctx := context.Background()

	criteria := model.SearchCriteria{
		WorkIds:      []int32{1, 2},
		SearchString: "Maxime",
	}

	// GIVEN
	pId1, _ := paraRepo.insertParagraphs(1, "Maxime Paragraph")
	pId2, _ := paraRepo.insertParagraphs(2, "Maxime Paragraph")
	ids1, _ := sut.insertSentences(pId1, []string{"Maxime", "Wille"})
	ids2, _ := sut.insertSentences(pId2, []string{"Maxime"})

	// WHEN
	matches, err := sut.Search(ctx, criteria)

	// THEN
	assert.Nil(t, err)
	assert.NotNil(t, matches)
	assert.Len(t, matches, 2)

	assert.Equal(t, ids1[0], matches[0].SentenceId)
	assert.Contains(t, matches[0].Snippet, "Maxime")
	assert.Contains(t, matches[0].Text, "Maxime")
	assert.NotContains(t, matches[0].Text, "Paragraph")
	assert.Equal(t, int32(1), matches[0].WorkId)

	assert.Equal(t, ids2[0], matches[1].SentenceId)
	assert.Contains(t, matches[1].Snippet, "Maxime")
	assert.Contains(t, matches[1].Text, "Maxime")
	assert.NotContains(t, matches[1].Text, "Paragraph")
	assert.Equal(t, int32(2), matches[1].WorkId)

	testDb.Exec("DELETE FROM sentences")
	testDb.Exec("DELETE FROM paragraphs")
}

func TestSearchSentencesNoMatch(t *testing.T) {
	sut := &sentenceRepoImpl{db: testDb}
	ctx := context.Background()

	criteria := model.SearchCriteria{
		WorkIds:      []int32{1},
		SearchString: "Maxime",
	}

	// WHEN
	matches, err := sut.Search(ctx, criteria)

	// THEN
	assert.Nil(t, err)
	assert.Len(t, matches, 0)
}

func TestDeleteSentences(t *testing.T) {
	paraRepo := &paragraphRepoImpl{db: testDb}
	sut := &sentenceRepoImpl{db: testDb}
	ctx := context.Background()

	// GIVEN
	pId1, _ := paraRepo.insertParagraphs(1, "text1")
	pId2, _ := paraRepo.insertParagraphs(2, "text2")
	sut.insertSentences(pId1, []string{"Maxime", "Wille"})
	sut.insertSentences(pId2, []string{"Vernunft"})

	// WHEN
	err := sut.DeleteByWorkId(ctx, 1)

	// THEN
	assert.Nil(t, err)
	var count int
	query := `SELECT COUNT(*) FROM sentences`
	err = testDb.QueryRowContext(ctx, query).Scan(&count)
	assert.Nil(t, err)
	assert.Equal(t, 1, count)
}

func TestDeleteSentencesOnEmptyTable(t *testing.T) {
	sut := &sentenceRepoImpl{db: testDb}
	ctx := context.Background()

	// WHEN
	err := sut.DeleteByWorkId(ctx, 1)

	// THEN
	assert.Nil(t, err)
}
