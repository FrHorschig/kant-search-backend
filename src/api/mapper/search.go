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

func ResultToApiModel(result []model.SearchMatch) []models.SearchResult {
	results := make([]models.SearchResult, len(result))
	// TODO implement me
	return results
}
