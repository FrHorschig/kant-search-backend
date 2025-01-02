//go:build integration
// +build integration

package dataaccess

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
	assert.Equal(t, 23, len(volumes))
	assert.Equal(t, int32(1), volumes[0].Id)
	assert.Equal(t, int32(23), volumes[len(volumes)-1].Id)
}
