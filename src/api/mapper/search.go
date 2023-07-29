package mapper

import (
	"github.com/FrHorschig/kant-search-api/models"
	"github.com/FrHorschig/kant-search-backend/core/model"
)

func CriteriaToCoreModel(criteria models.SearchCriteria) model.SearchCriteria {
	return model.SearchCriteria{
		SearchTerms: criteria.SearchTerms,
		WorkIds:     criteria.WorkIds,
	}
}

func MatchesToApiModel(matches []model.SearchMatch) []models.SearchResult {
	var results []models.SearchResult
	var currentResult *models.SearchResult

	for _, coreMatch := range matches {
		if isNewWork(currentResult, coreMatch) {
			if currentResult != nil {
				results = append(results, *currentResult)
			}
			currentResult = &models.SearchResult{
				Volume:    coreMatch.Volume,
				WorkTitle: coreMatch.WorkTitle,
				Matches:   []models.Match{},
			}
		}

		apiMatch := &models.Match{
			Snippet:   coreMatch.Snippet,
			Pages:     coreMatch.Pages,
			WorkId:    coreMatch.WorkId,
			ElementId: coreMatch.ElementId,
		}
		currentResult.Matches = append(currentResult.Matches, *apiMatch)
	}

	if currentResult != nil {
		results = append(results, *currentResult)
	}
	return results
}

func isNewWork(currentResult *models.SearchResult, match model.SearchMatch) bool {
	return currentResult == nil || currentResult.WorkTitle != match.WorkTitle
}
