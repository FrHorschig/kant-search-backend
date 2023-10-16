//go:build integration
// +build integration

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

	krvB := works[30]
	assert.Equal(t, "Kritik der reinen Vernunft 2. Aufl.", krvB.Title)
	assert.Equal(t, "B", *krvB.Abbreviation)
	assert.Equal(t, int32(3), krvB.Volume)
	assert.Equal(t, int32(0), krvB.Ordinal)
	assert.Equal(t, "1787", *krvB.Year)

	work7_1 := works[40]
	assert.Equal(t, "Anthropologie in pragmatischer Hinsicht", work7_1.Title)
	assert.Equal(t, "Anth", *work7_1.Abbreviation)
	assert.Equal(t, int32(7), work7_1.Volume)
	assert.Equal(t, int32(1), work7_1.Ordinal)
	assert.Equal(t, "1798", *work7_1.Year)
}