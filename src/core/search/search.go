package search

//go:generate mockgen -source=$GOFILE -destination=mocks/search_processor_mock.go -package=mocks

import (
	"context"

	"github.com/FrHorschig/kant-search-backend/database/model"
	"github.com/FrHorschig/kant-search-backend/database/repository"
	"github.com/FrHorschig/kant-search-backend/database/util"
)

type SearchProcessor interface {
	Search(ctx context.Context, criteria model.SearchCriteria) ([]model.SearchResult, error)
}

type searchProcessorImpl struct {
	paragraphRepo repository.ParagraphRepo
	sentenceRepo  repository.SentenceRepo
}

func NewSearchProcessor(paragraphRepo repository.ParagraphRepo, sentenceRepo repository.SentenceRepo) SearchProcessor {
	impl := searchProcessorImpl{paragraphRepo: paragraphRepo, sentenceRepo: sentenceRepo}
	return &impl
}

func (rec *searchProcessorImpl) Search(ctx context.Context, criteria model.SearchCriteria) ([]model.SearchResult, error) {
	criteria.SearchString = util.EscapeSpecialChars(criteria.SearchString)
	if criteria.Options.Scope == model.SentenceScope {
		return rec.sentenceRepo.Search(ctx, criteria)
	} else {
		return rec.paragraphRepo.Search(ctx, criteria)
	}
}
