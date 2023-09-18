//go:build unit
// +build unit

package repository

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/FrHorschig/kant-search-backend/database/model"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestSearchParagraphsDatabaseError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	repo := &searchRepoImpl{db: db}
	dbErr := fmt.Errorf("database error")

	criteria := model.SearchCriteria{
		SearchTerms: []string{"Maxime"},
		WorkIds:     []int32{1},
		Scope:       model.PARAGRAPH,
	}

	// GIVEN
	mock.ExpectQuery(anyQuery).WillReturnError(dbErr)

	// WHEN
	matches, err := repo.SearchParagraphs(context.Background(), criteria)

	// THEN
	assert.Equal(t, dbErr, err)
	assert.Empty(t, matches)
}

func TestSearchParagraphsNoRows(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	repo := &searchRepoImpl{db: db}

	criteria := model.SearchCriteria{
		SearchTerms: []string{"Maxime"},
		WorkIds:     []int32{1},
		Scope:       model.PARAGRAPH,
	}

	// GIVEN
	mock.ExpectQuery(anyQuery).WillReturnError(sql.ErrNoRows)

	// WHEN
	matches, err := repo.SearchParagraphs(context.Background(), criteria)

	// THEN
	assert.Nil(t, err)
	assert.Empty(t, matches)
}

func TestSearchParagraphsWrongRows(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	repo := &searchRepoImpl{db: db}

	criteria := model.SearchCriteria{
		SearchTerms: []string{"Maxime"},
		WorkIds:     []int32{1},
		Scope:       model.PARAGRAPH,
	}

	// GIVEN
	mock.ExpectQuery(anyQuery).WillReturnRows(sqlmock.NewRows([]string{"abc", "def"}).AddRow(1, 1))

	// WHEN
	matches, err := repo.SearchParagraphs(context.Background(), criteria)

	// THEN
	assert.NotNil(t, err)
	assert.Empty(t, matches)
}
