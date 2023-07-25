package mapper

import (
	"github.com/FrHorschig/kant-search-api/models"
	"github.com/FrHorschig/kant-search-backend/core/model"
)

func CriteriaToCoreModel(criteria models.SearchCriteria) model.SearchCriteria {
	return model.SearchCriteria{
		SearchWords: criteria.SearchWords,
		WorkIds:     criteria.WorkIds,
	}
}

func ResultToApiModel(result model.SearchResult) models.SearchResult {
	return models.SearchResult{
		Paragraphs: ParagraphsToApiModel(result.Paragraphs),
	}
}
