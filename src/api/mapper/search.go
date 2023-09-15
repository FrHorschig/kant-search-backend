package mapper

import (
	"github.com/FrHorschig/kant-search-api/models"
	"github.com/FrHorschig/kant-search-backend/database/model"
)

func CriteriaToCoreModel(criteria models.SearchCriteria) model.SearchCriteria {
	return model.SearchCriteria{
		SearchTerms: criteria.SearchTerms,
		WorkIds:     criteria.WorkIds,
		Scope:       mapSearchScope(criteria.Scope),
	}
}

func MatchesToApiModels(matches []model.SearchMatch) []models.SearchResult {
	resultByWorkId := make(map[int32][]models.Match)
	for _, match := range matches {
		apiMatch := models.Match{
			Snippet:   match.Snippet,
			Pages:     match.Pages,
			ElementId: match.ElementId,
		}

		arr, exists := resultByWorkId[match.WorkId]
		if !exists {
			arr = []models.Match{apiMatch}
		} else {
			arr = append(arr, apiMatch)
		}
		resultByWorkId[match.WorkId] = arr
	}

	var results []models.SearchResult
	for workId, apiMatches := range resultByWorkId {
		results = append(results, models.SearchResult{
			WorkId:  workId,
			Matches: apiMatches,
		})
	}
	return results
}

func mapSearchScope(scope models.SearchScope) model.SearchScope {
	switch scope {
	case models.PARAGRAPH:
		return model.PARAGRAPH
	case models.SENTENCE:
		return model.SENTENCE
	default:
		return model.PARAGRAPH
	}
}
