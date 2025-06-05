//go:build integration
// +build integration

package dataaccess

import (
	"context"
	"testing"
	"time"

	"github.com/frhorschig/kant-search-backend/dataaccess/esmodel"
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
			Code:  "code",
			Title: "work title",
		}},
	}
	vol2 := esmodel.Volume{
		VolumeNumber: 2,
		Section:      3,
		Title:        "volume title 2",
		Works: []esmodel.WorkRef{{
			Code:  "code2",
			Title: "work title 2",
		}},
	}

	// WHEN Insert
	err := repo.Insert(ctx, &vol)
	// THEN
	assert.Nil(t, err)
	refreshVolumes(t)

	// WHEN Insert duplicate
	err = repo.Insert(ctx, &vol)
	// THEN
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "already exists")

	// WHEN Insert 2nd
	err = repo.Insert(ctx, &vol2)
	// THEN
	assert.Nil(t, err)
	refreshVolumes(t)

	// WHEN Get
	res, err := repo.GetAll(ctx)
	// THEN
	assert.Nil(t, err)
	assert.Len(t, res, 2)
	assert.ElementsMatch(t,
		[]string{vol.Title, vol2.Title},
		[]string{res[0].Title, res[1].Title},
	)

	// WHEN Delete
	err = repo.Delete(ctx, vol.VolumeNumber)
	// THEN
	assert.Nil(t, err)
	refreshVolumes(t)
	res, err = repo.GetAll(ctx)
	assert.Nil(t, err)
	assert.Len(t, res, 1)
	assert.Equal(t, vol2.Title, res[0].Title)

	// WHEN Get
	singleRes, err := repo.GetByVolumeNumber(ctx, vol2.VolumeNumber)
	// THEN
	assert.Nil(t, err)
	assert.NotNil(t, singleRes)
	assert.Equal(t, vol2.Title, singleRes.Title)

	// WHEN Delete 2ns
	err = repo.Delete(ctx, vol2.VolumeNumber)
	// THEN
	assert.Nil(t, err)
	refreshVolumes(t)
	res, err = repo.GetAll(ctx)
	assert.Nil(t, err)
	assert.Len(t, res, 0)

	// WHEN Get 2nd
	singleRes, err = repo.GetByVolumeNumber(ctx, vol2.VolumeNumber)
	// THEN
	assert.Nil(t, err)
	assert.Nil(t, singleRes)
}

// TODO test GetAll

func refreshVolumes(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	_, err := dbClient.Indices.Refresh().Index("volumes").Do(ctx)
	assert.Nil(t, err)
}
