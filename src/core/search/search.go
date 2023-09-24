package search

import (
	"context"

	"github.com/FrHorschig/kant-search-backend/database/model"
	"github.com/FrHorschig/kant-search-backend/database/repository"
	"github.com/FrHorschig/kant-search-backend/database/util"
)

//go:generate mockgen -source=$GOFILE -destination=mocks/work_mock.go -package=mocks

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
	escapeSpecialChars(&criteria)
	if criteria.Scope == model.SentenceScope {
		return rec.sentenceRepo.Search(ctx, criteria)
	} else {
		return rec.paragraphRepo.Search(ctx, criteria)
	}
}

func escapeSpecialChars(c *model.SearchCriteria) {
	for i := range c.SearchTerms {
		c.SearchTerms[i] = util.EscapeSpecialChars(c.ExcludedTerms[i])
	}
	for i := range c.ExcludedTerms {
		c.SearchTerms[i] = util.EscapeSpecialChars(c.ExcludedTerms[i])
	}
	for i := range c.OptionalTerms {
		c.SearchTerms[i] = util.EscapeSpecialChars(c.ExcludedTerms[i])
	}
}
