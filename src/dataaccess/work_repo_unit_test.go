//go:build unit
// +build unit

package dataaccess

import (
	"testing"

	_ "github.com/lib/pq"
)

func TestSelectAllWorksDatabaseError(t *testing.T) {
	// repo := NewWorkRepo(dbClient)
	// dbErr := fmt.Errorf("database error")

	// GIVEN
	// mock.ExpectQuery(anyQuery).WillReturnError(dbErr)

	// WHEN
	// works, err := repo.SelectAll(context.Background())

	// THEN
	// assert.Equal(t, dbErr, err)
	// assert.Empty(t, works)
}

func TestSelectAllWorksNoRows(t *testing.T) {
	// repo := NewWorkRepo(dbClient)

	// GIVEN
	// mock.ExpectQuery(anyQuery).WillReturnError(sql.ErrNoRows)

	// WHEN
	// works, err := repo.SelectAll(context.Background())

	// THEN
	// assert.Nil(t, err)
	// assert.Empty(t, works)
}

func TestSelectAllWorksWrongRows(t *testing.T) {
	// repo := NewWorkRepo(dbClient)

	// GIVEN
	// mock.ExpectQuery(anyQuery).WillReturnRows(sqlmock.NewRows([]string{"abc", "def"}).AddRow(1, 1))

	// WHEN
	// works, err := repo.SelectAll(context.Background())

	// THEN
	// assert.NotNil(t, err)
	// assert.Empty(t, works)
}
