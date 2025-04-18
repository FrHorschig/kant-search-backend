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
	}

	// WHEN Insert
	err := repo.Insert(ctx, &work)
	// THEN
	assert.Nil(t, err)
	assert.NotEmpty(t, work.Id)
	refreshWorks(t)

	// WHEN Get
	res, err := repo.Get(ctx, work.Id)
	// THEN
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, work.Code, res.Code)
	assert.Equal(t, work.Abbreviation, res.Abbreviation)
	assert.Equal(t, work.Title, res.Title)
	assert.Equal(t, work.Year, res.Year)

	// WHEN Update
	work.Sections = []esmodel.Section{{
		Heading:    "heading1Id",
		Paragraphs: []string{"par11Id", "par12Id", "par13Id"},
		Sections: []esmodel.Section{{
			Heading:    "heading2Id",
			Paragraphs: []string{"par21Id", "par22Id", "par23Id"},
		}},
	}}
	err = repo.Update(ctx, &work)
	// THEN
	assert.Nil(t, err)
	refreshWorks(t)

	// WHEN Get
	res, err = repo.Get(ctx, work.Id)
	// THEN
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, work.Sections, res.Sections)

	// WHEN: Delete
	err = repo.Delete(ctx, work.Id)
	// THEN
	assert.Nil(t, err)
	refreshWorks(t)

	// WHEN: Get
	res, err = repo.Get(ctx, work.Id)
	// THEN
	assert.NotNil(t, err)
	assert.Nil(t, res)
	assert.Contains(t, err.Error(), "not found")

	// WHEN: Delete nonexisting
	err = repo.Delete(ctx, work.Id)
	// THEN
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "unable to delete")

}

func refreshWorks(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	_, err := dbClient.Indices.Refresh().Index("works").Do(ctx)
	assert.Nil(t, err)
}
