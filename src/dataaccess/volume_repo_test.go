//go:build integration
// +build integration

package dataaccess

import (
	"context"
	"testing"
	"time"

	"github.com/frhorschig/kant-search-backend/dataaccess/internal/esmodel"
	"github.com/stretchr/testify/assert"
)

func TestVolumeRepo(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	repo := NewVolumeRepo(dbClient)
	vol := esmodel.Volume{
		VolumeNumber: 1,
		Section:      2,
		Title:        "volume title",
		Works: []esmodel.WorkRef{{
			Id:    "work id",
			Code:  "code",
			Title: "work title",
		}},
	}

	// WHEN Insert
	err := repo.Insert(ctx, &vol)
	// THEN
	assert.Nil(t, err)
	assert.NotEmpty(t, vol.Id)
	refreshVolumes(t)

	// WHEN Insert duplicate
	err = repo.Insert(ctx, &vol)
	// THEN
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "already exists")
	assert.NotEmpty(t, vol.Id)

	// WHEN Get
	res, err := repo.Get(ctx, vol.VolumeNumber)
	// THEN
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, vol, *res)

	// WHEN Delete
	err = repo.Delete(ctx, vol.VolumeNumber)
	// THEN
	assert.Nil(t, err)
	refreshVolumes(t)
	res, err = repo.Get(ctx, vol.VolumeNumber)
	assert.Nil(t, err)
	assert.Nil(t, res)
}

func refreshVolumes(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	_, err := dbClient.Indices.Refresh().Index("volumes").Do(ctx)
	assert.Nil(t, err)
}
