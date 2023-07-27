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

func ResultToApiModel(result model.ParagraphResults) models.ParagraphResults {
	return models.ParagraphResults{
		Paragraphs:  ParagraphsToApiModel(result.Paragraphs),
		SearchWords: result.MatchedWords,
	}
}
