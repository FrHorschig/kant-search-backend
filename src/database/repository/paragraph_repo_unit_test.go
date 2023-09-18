//go:build unit
// +build unit

package repository

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestSelectParagraphDatabaseError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	repo := &paragraphRepoImpl{db: db}
	dbErr := fmt.Errorf("database error")

	// GIVEN
	mock.ExpectQuery(anyQuery).WillReturnError(dbErr)

	// WHEN
	para, err := repo.Select(context.Background(), 1, 1)

	// THEN
	assert.Equal(t, dbErr, err)
	assert.Empty(t, para)
}

func TestSelectParagraphNoRows(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	repo := &paragraphRepoImpl{db: db}

	// GIVEN
	mock.ExpectQuery(anyQuery).WillReturnError(sql.ErrNoRows)

	// WHEN
	para, err := repo.Select(context.Background(), 1, 1)

	// THEN
	assert.Nil(t, err)
	assert.Empty(t, para)
}

func TestSelectParagraphWrongRows(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	repo := &paragraphRepoImpl{db: db}

	// GIVEN
	mock.ExpectQuery(anyQuery).WillReturnRows(sqlmock.NewRows([]string{"abc", "def"}).AddRow(1, 1))

	// WHEN
	para, err := repo.Select(context.Background(), 1, 1)

	// THEN
	assert.NotNil(t, err)
	assert.Empty(t, para)
}

func TestSelectAllParagraphsDatabaseError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	repo := &paragraphRepoImpl{db: db}
	dbErr := fmt.Errorf("database error")

	// GIVEN
	mock.ExpectQuery(anyQuery).WillReturnError(dbErr)

	// WHEN
	paras, err := repo.SelectAll(context.Background(), 1)

	// THEN
	assert.Equal(t, dbErr, err)
	assert.Empty(t, paras)
}

func TestSelectAllParagraphsNoRows(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	repo := &paragraphRepoImpl{db: db}

	// GIVEN
	mock.ExpectQuery(anyQuery).WillReturnError(sql.ErrNoRows)

	// WHEN
	paras, err := repo.SelectAll(context.Background(), 1)

	// THEN
	assert.Nil(t, err)
	assert.Empty(t, paras)
}

func TestSelectAllParagraphsWrongRows(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	repo := &paragraphRepoImpl{db: db}

	// GIVEN
	mock.ExpectQuery(anyQuery).WillReturnRows(sqlmock.NewRows([]string{"abc", "def"}).AddRow(1, 1))

	// WHEN
	paras, err := repo.SelectAll(context.Background(), 1)

	// THEN
	assert.NotNil(t, err)
	assert.Empty(t, paras)
}
