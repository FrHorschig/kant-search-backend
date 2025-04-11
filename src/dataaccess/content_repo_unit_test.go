//go:build unit
// +build unit

package dataaccess

import (
	"testing"

	_ "github.com/lib/pq"
)

func TestInsertParagraphsDatabaseError(t *testing.T) {
	// repo := NewContentRepo(dbClient)
	// dbErr := fmt.Errorf("database error")
	// paragraph := model.Paragraph{
	// 	Text:   "text",
	// 	Pages:  []int32{1, 2, 3},
	// 	WorkId: 1,
	// }

	// GIVEN
	// mock.ExpectQuery(anyQuery).WillReturnError(dbErr)

	// WHEN
	// id, err := repo.Insert(context.Background(), paragraph)

	// THEN
	// assert.Equal(t, int32(0), id)
	// assert.NotNil(t, err)
}

func TestSelectAllParagraphsDatabaseError(t *testing.T) {
	// repo := NewContentRepo(dbClient)
	// dbErr := fmt.Errorf("database error")

	// GIVEN
	// mock.ExpectQuery(anyQuery).WillReturnError(dbErr)

	// WHEN
	// paras, err := repo.SelectAll(context.Background(), 1)

	// THEN
	// assert.Equal(t, dbErr, err)
	// assert.Empty(t, paras)
}

func TestSelectAllParagraphsNoRows(t *testing.T) {
	// repo := NewContentRepo(dbClient)

	// GIVEN
	// mock.ExpectQuery(anyQuery).WillReturnError(sql.ErrNoRows)

	// WHEN
	// paras, err := repo.SelectAll(context.Background(), 1)

	// THEN
	// assert.Nil(t, err)
	// assert.Empty(t, paras)
}

func TestSelectAllParagraphsWrongRows(t *testing.T) {
	// repo := NewContentRepo(dbClient)

	// GIVEN
	// mock.ExpectQuery(anyQuery).WillReturnRows(sqlmock.NewRows([]string{"abc", "def"}).AddRow(1, 1))

	// WHEN
	// paras, err := repo.SelectAll(context.Background(), 1)

	// THEN
	// assert.NotNil(t, err)
	// assert.Empty(t, paras)
}

func TestSearchParagraphsDatabaseError(t *testing.T) {
	// repo := NewContentRepo(dbClient)
	// dbErr := fmt.Errorf("database error")
	// criteria := model.SearchCriteria{
	// 	WorkIds:      []int32{1},
	// 	SearchString: "Maxime",
	// }

	// GIVEN
	// mock.ExpectQuery(anyQuery).WillReturnError(dbErr)

	// WHEN
	// result, err := repo.Search(context.Background(), criteria)

	// THEN
	// assert.Nil(t, result)
	// assert.NotNil(t, err)
}

func TestSearchParagraphsNoRows(t *testing.T) {
	// repo := NewContentRepo(dbClient)

	// criteria := model.SearchCriteria{
	// 	WorkIds:      []int32{1},
	// 	SearchString: "Maxime",
	// }

	// GIVEN
	// mock.ExpectQuery(anyQuery).WillReturnError(sql.ErrNoRows)

	// WHEN
	// matches, err := repo.Search(context.Background(), criteria)

	// THEN
	// assert.Nil(t, err)
	// assert.Empty(t, matches)
}

func TestSearchParagraphsWrongRows(t *testing.T) {
	// repo := NewContentRepo(dbClient)

	// criteria := model.SearchCriteria{
	// 	WorkIds:      []int32{1},
	// 	SearchString: "Maxime",
	// }

	// GIVEN
	// mock.ExpectQuery(anyQuery).WillReturnRows(sqlmock.NewRows([]string{"abc", "def"}).AddRow(1, 1))

	// WHEN
	// matches, err := repo.Search(context.Background(), criteria)

	// THEN
	// assert.NotNil(t, err)
	// assert.Empty(t, matches)
}

func TestDeleteParagraphDatabaseError(t *testing.T) {
	// repo := NewContentRepo(dbClient)
	// dbErr := fmt.Errorf("database error")

	// GIVEN
	// mock.ExpectQuery(anyQuery).WillReturnError(dbErr)

	// WHEN
	// err = repo.DeleteByWorkId(context.Background(), 1)

	// THEN
	// assert.NotNil(t, err)
}
