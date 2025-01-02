package search

//go:generate mockgen -source=$GOFILE -destination=mocks/search_processor_mock.go -package=mocks

import (
	"context"

	"github.com/frhorschig/kant-search-backend/common/model"
	"github.com/frhorschig/kant-search-backend/dataaccess"
)

type SearchProcessor interface {
	Search(ctx context.Context, criteria model.SearchCriteria) ([]model.SearchResult, error)
}

type searchProcessorImpl struct {
	paragraphRepo dataaccess.ParagraphRepo
	sentenceRepo  dataaccess.SentenceRepo
}

func NewSearchProcessor(paragraphRepo dataaccess.ParagraphRepo, sentenceRepo dataaccess.SentenceRepo) SearchProcessor {
	impl := searchProcessorImpl{paragraphRepo: paragraphRepo, sentenceRepo: sentenceRepo}
	return &impl
}

func (rec *searchProcessorImpl) Search(ctx context.Context, criteria model.SearchCriteria) ([]model.SearchResult, error) {
	if criteria.Options.Scope == model.SentenceScope {
		return rec.sentenceRepo.Search(ctx, criteria)
	} else {
		return rec.paragraphRepo.Search(ctx, criteria)
	}
}
