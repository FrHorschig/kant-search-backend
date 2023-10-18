//go:build unit
// +build unit

package database

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestSelectAllVolumesDatabaseError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	repo := NewVolumeRepo(db)
	dbErr := fmt.Errorf("database error")

	// GIVEN
	mock.ExpectQuery(anyQuery).WillReturnError(dbErr)

	// WHEN
	volumes, err := repo.SelectAll(context.Background())

	// THEN
	assert.Equal(t, dbErr, err)
	assert.Empty(t, volumes)
}

func TestSelectAllVolumesNoRows(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	repo := &volumeRepoImpl{db: db}

	// GIVEN
	mock.ExpectQuery(anyQuery).WillReturnError(sql.ErrNoRows)

	// WHEN
	volumes, err := repo.SelectAll(context.Background())

	// THEN
	assert.Nil(t, err)
	assert.Empty(t, volumes)
}

func TestSelectAllVolumesWrongRows(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	repo := &volumeRepoImpl{db: db}

	// GIVEN
	mock.ExpectQuery(anyQuery).WillReturnRows(sqlmock.NewRows([]string{"abc", "def"}).AddRow(1, 1))

	// WHEN
	volumes, err := repo.SelectAll(context.Background())

	// THEN
	assert.NotNil(t, err)
	assert.Empty(t, volumes)
}
