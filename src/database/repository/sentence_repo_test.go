//go:build integration
// +build integration

package repository

import (
	"context"
	"testing"

	"github.com/FrHorschig/kant-search-backend/database/model"
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
