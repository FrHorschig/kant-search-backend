//go:build unit
// +build unit

package dataaccess

import (
	"testing"
)

func TestSelectAllVolumesDatabaseError(t *testing.T) {
	// repo := NewVolumeRepo(dbClient)
	// dbErr := fmt.Errorf("database error")

	// GIVEN
	// mock.ExpectQuery(anyQuery).WillReturnError(dbErr)

	// WHEN
	// volumes, err := repo.SelectAll(context.Background())

	// THEN
	// assert.Equal(t, dbErr, err)
	// assert.Empty(t, volumes)
}

func TestSelectAllVolumesNoRows(t *testing.T) {
	// repo := NewVolumeRepo(dbClient)

	// GIVEN
	// mock.ExpectQuery(anyQuery).WillReturnError(sql.ErrNoRows)

	// WHEN
	// volumes, err := repo.SelectAll(context.Background())

	// THEN
	// assert.Nil(t, err)
	// assert.Empty(t, volumes)
}

func TestSelectAllVolumesWrongRows(t *testing.T) {
	// repo := NewVolumeRepo(dbClient)

	// GIVEN
	// mock.ExpectQuery(anyQuery).WillReturnRows(sqlmock.NewRows([]string{"abc", "def"}).AddRow(1, "ghi"))

	// WHEN
	// volumes, err := repo.SelectAll(context.Background())

	// THEN
	// assert.NotNil(t, err)
	// assert.Empty(t, volumes)
}
