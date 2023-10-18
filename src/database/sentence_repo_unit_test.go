//go:build unit
// +build unit

package database

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/FrHorschig/kant-search-backend/common/model"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestInsertSentencesDatabaseError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	repo := &sentenceRepoImpl{db: db}
	dbErr := fmt.Errorf("database error")
	sentence := model.Sentence{
		Text:        "text",
		ParagraphId: 1,
	}

	// GIVEN
	mock.ExpectQuery(anyQuery).WillReturnError(dbErr)

	// WHEN
	ids, err := repo.Insert(context.Background(), []model.Sentence{sentence})

	// THEN
	assert.Nil(t, ids)
	assert.NotNil(t, err)
}

func TestInsertSentencesScanError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	repo := &sentenceRepoImpl{db: db}
	sentence := model.Sentence{
		Text:        "text",
		ParagraphId: 1,
	}

	// GIVEN
	mock.ExpectQuery(anyQuery).WillReturnRows(sqlmock.NewRows([]string{"abc", "def"}).AddRow("x", "y"))

	// WHEN
	ids, err := repo.Insert(context.Background(), []model.Sentence{sentence})

	// THEN
	assert.Nil(t, ids)
	assert.NotNil(t, err)
}

func TestSearchSentencesDatabaseError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	repo := &sentenceRepoImpl{db: db}
	dbErr := fmt.Errorf("database error")

	criteria := model.SearchCriteria{
		WorkIds:      []int32{1},
		SearchString: "Maxime",
	}

	// GIVEN
	mock.ExpectQuery(anyQuery).WillReturnError(dbErr)

	// WHEN
	matches, err := repo.Search(context.Background(), criteria)

	// THEN
	assert.Equal(t, dbErr, err)
	assert.Empty(t, matches)
}

func TestSearchSentencesNoRows(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	repo := &sentenceRepoImpl{db: db}

	criteria := model.SearchCriteria{
		WorkIds:      []int32{1},
		SearchString: "Maxime",
	}

	// GIVEN
	mock.ExpectQuery(anyQuery).WillReturnError(sql.ErrNoRows)

	// WHEN
	matches, err := repo.Search(context.Background(), criteria)

	// THEN
	assert.Nil(t, err)
	assert.Empty(t, matches)
}

func TestSearchSentencesWrongRows(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	repo := &sentenceRepoImpl{db: db}

	criteria := model.SearchCriteria{
		WorkIds:      []int32{1},
		SearchString: "Maxime",
	}

	// GIVEN
	mock.ExpectQuery(anyQuery).WillReturnRows(sqlmock.NewRows([]string{"abc", "def"}).AddRow(1, 1))

	// WHEN
	matches, err := repo.Search(context.Background(), criteria)

	// THEN
	assert.NotNil(t, err)
	assert.Empty(t, matches)
}

func TestDeleteSentencesDatabaseError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	repo := &sentenceRepoImpl{db: db}
	dbErr := fmt.Errorf("database error")

	// GIVEN
	mock.ExpectQuery(anyQuery).WillReturnError(dbErr)

	// WHEN
	err = repo.DeleteByWorkId(context.Background(), 1)

	// THEN
	assert.NotNil(t, err)
}
