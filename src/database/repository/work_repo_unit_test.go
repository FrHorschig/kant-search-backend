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

func TestSelectAllWorksDatabaseError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	repo := &workRepoImpl{db: db}
	dbErr := fmt.Errorf("database error")

	// GIVEN
	mock.ExpectQuery(anyQuery).WillReturnError(dbErr)

	// WHEN
	works, err := repo.SelectAll(context.Background())

	// THEN
	assert.Equal(t, dbErr, err)
	assert.Empty(t, works)
}

func TestSelectAllWorksNoRows(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	repo := &workRepoImpl{db: db}

	// GIVEN
	mock.ExpectQuery(anyQuery).WillReturnError(sql.ErrNoRows)

	// WHEN
	works, err := repo.SelectAll(context.Background())

	// THEN
	assert.Nil(t, err)
	assert.Empty(t, works)
}

func TestSelectAllWorksWrongRows(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	repo := &workRepoImpl{db: db}

	// GIVEN
	mock.ExpectQuery(anyQuery).WillReturnRows(sqlmock.NewRows([]string{"abc", "def"}).AddRow(1, 1))

	// WHEN
	works, err := repo.SelectAll(context.Background())

	// THEN
	assert.NotNil(t, err)
	assert.Empty(t, works)
}
