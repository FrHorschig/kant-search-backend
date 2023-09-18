package mapper

import (
	"sort"

	"github.com/FrHorschig/kant-search-api/models"
	"github.com/FrHorschig/kant-search-backend/database/model"
)

func CriteriaToCoreModel(criteria models.SearchCriteria) model.SearchCriteria {
	return model.SearchCriteria{
		SearchTerms: criteria.SearchTerms,
		WorkIds:     criteria.WorkIds,
	}
}

func MatchesToApiModels(matches []model.SearchResult) []models.SearchResult {
	resultByWorkId := make(map[int32][]models.Match)
	for _, match := range matches {
		apiMatch := models.Match{
			ElementId: match.ElementId,
			Snippet:   match.Snippet,
			Pages:     match.Pages,
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
	sort.Slice(results, func(i, j int) bool {
		return results[i].WorkId < results[j].WorkId
	})
	return results
}
