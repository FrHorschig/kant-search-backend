package repository

import (
	"context"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestSelectAllWorks(t *testing.T) {
	repo := &workRepoImpl{db: testDb}

	// WHEN
	works, err := repo.SelectAll(context.Background())

	// THEN
	assert.Nil(t, err)
	assert.Greater(t, len(works), 0)

	krv := works[30]
	assert.Equal(t, "Kritik der reinen Vernunft A", krv.Title)
	assert.Equal(t, "A", *krv.Abbreviation)
	assert.Equal(t, int32(4), krv.Volume)
	assert.Equal(t, int32(0), krv.Ordinal)
	assert.Equal(t, "1781", *krv.Year)

	work8_1 := works[40]
	assert.Equal(t, "Anzeige des Lambert'schen Briefwechsels", work8_1.Title)
	assert.Nil(t, work8_1.Abbreviation)
	assert.Equal(t, int32(8), work8_1.Volume)
	assert.Equal(t, int32(0), work8_1.Ordinal)
	assert.Equal(t, "1782", *work8_1.Year)
}
