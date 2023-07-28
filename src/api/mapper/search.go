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

func MatchesToApiModel(matches []model.SearchMatch) []models.SearchResult {
	var results []models.SearchResult
	var currentResult *models.SearchResult

	for _, coreMatch := range matches {
		if isNewVolume(currentResult, coreMatch) {
			if currentResult != nil {
				results = append(results, *currentResult)
			}
			currentResult = &models.SearchResult{
				Volume:  coreMatch.Volume,
				Matches: []models.Match{},
			}
		}

		apiMatch := &models.Match{
			WorkTitle: coreMatch.WorkTitle,
			Snippet:   coreMatch.Snippet,
			Pages:     coreMatch.Pages,
			MatchId:   coreMatch.MatchId,
		}
		currentResult.Matches = append(currentResult.Matches, *apiMatch)
	}

	if currentResult != nil {
		results = append(results, *currentResult)
	}
	return results
}

func isNewVolume(currentResult *models.SearchResult, match model.SearchMatch) bool {
	return currentResult == nil || currentResult.Volume != match.Volume
}
