//go:build integration
// +build integration

package dataaccess

import (
	"context"
	"testing"
	"time"

	"github.com/frhorschig/kant-search-backend/common/util"
	"github.com/frhorschig/kant-search-backend/dataaccess/esmodel"
	"github.com/stretchr/testify/assert"
)

func TestWorkRepo(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	repo := NewWorkRepo(dbClient)
	work := esmodel.Work{
		Code:         "GMS",
		Abbreviation: util.StrPtr("GMS"),
		Title:        "Grundlegung zur Metaphysik der Sitten",
		Year:         util.StrPtr("1785"),
		Ordinal:      1,
		Sections: []esmodel.Section{{
			Heading:    1,
			Paragraphs: []int32{2, 3, 4},
			Sections: []esmodel.Section{{
				Heading:    5,
				Paragraphs: []int32{6, 7},
			}},
		}},
	}

	// WHEN Insert
	err := repo.Insert(ctx, &work)
	// THEN
	assert.Nil(t, err)
	refreshWorks(t)

	// WHEN Get
	res, err := repo.Get(ctx, work.Code)
	// THEN
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, work.Code, res.Code)
	assert.Equal(t, work.Abbreviation, res.Abbreviation)
	assert.Equal(t, work.Title, res.Title)
	assert.Equal(t, work.Year, res.Year)
	assert.Equal(t, work.Ordinal, res.Ordinal)
	assert.Equal(t, work.Sections, res.Sections)

	// WHEN: Delete
	err = repo.Delete(ctx, work.Code)
	// THEN
	assert.Nil(t, err)
	refreshWorks(t)

	// WHEN: Get
	res, err = repo.Get(ctx, work.Code)
	// THEN
	assert.NotNil(t, err)
	assert.Nil(t, res)
	assert.Contains(t, err.Error(), "no work with code")
}

func refreshWorks(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	_, err := dbClient.Indices.Refresh().Index("works").Do(ctx)
	assert.Nil(t, err)
}
