package search

import (
	"context"
	"errors"
	"testing"

	"github.com/FrHorschig/kant-search-backend/database/model"
	"github.com/FrHorschig/kant-search-backend/database/repository/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestSearchHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	paragraphRepo := mocks.NewMockParagraphRepo(ctrl)
	sentenceRepo := mocks.NewMockSentenceRepo(ctrl)
	sut := &searchProcessorImpl{
		paragraphRepo: paragraphRepo,
		sentenceRepo:  sentenceRepo,
	}

	for scenario, fn := range map[string]func(t *testing.T, sut *searchProcessorImpl, paragraphRepo *mocks.MockParagraphRepo, sentenceRepo *mocks.MockSentenceRepo){
		"Search with paragraph scope": testSearchWithParagraphScope,
		"Search with sentence scope":  testSearchWithSentenceScope,
	} {
		t.Run(scenario, func(t *testing.T) {
			fn(t, sut, paragraphRepo, sentenceRepo)
		})
	}
}

func testSearchWithParagraphScope(t *testing.T, sut *searchProcessorImpl, paragraphRepo *mocks.MockParagraphRepo, sentenceRepo *mocks.MockSentenceRepo) {
	criteria := model.SearchCriteria{Scope: model.ParagraphScope}
	matches := []model.SearchResult{}
	err := errors.New("some error")
	// GIVEN
	paragraphRepo.EXPECT().Search(gomock.Any(), gomock.Any()).Return(matches, err)
	// WHEN
	result, errResult := sut.Search(context.Background(), criteria)
	// THEN
	assert.Equal(t, matches, result)
	assert.Equal(t, err, errResult)
}

func testSearchWithSentenceScope(t *testing.T, sut *searchProcessorImpl, paragraphRepo *mocks.MockParagraphRepo, sentenceRepo *mocks.MockSentenceRepo) {
	criteria := model.SearchCriteria{Scope: model.SentenceScope}
	matches := []model.SearchResult{}
	err := errors.New("some error")
	// GIVEN
	sentenceRepo.EXPECT().Search(gomock.Any(), gomock.Any()).Return(matches, err)
	// WHEN
	result, errResult := sut.Search(context.Background(), criteria)
	// THEN
	assert.Equal(t, matches, result)
	assert.Equal(t, err, errResult)
}
