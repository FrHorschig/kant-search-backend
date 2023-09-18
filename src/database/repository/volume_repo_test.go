package repository

import (
	"context"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestSelectAllVolumes(t *testing.T) {
	repo := &volumeRepoImpl{db: testDb}

	// WHEN
	volumes, err := repo.SelectAll(context.Background())

	// THEN
	assert.Nil(t, err)
	assert.Equal(t, len(volumes), 23)
	assert.Equal(t, volumes[0].Id, 1)
	assert.Equal(t, volumes[len(volumes)-1].Id, 23)
}
