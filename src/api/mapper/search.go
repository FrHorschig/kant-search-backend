package mapper

import (
	"sort"
	"strings"

	"github.com/FrHorschig/kant-search-api/models"
	"github.com/FrHorschig/kant-search-backend/database/model"
)

func CriteriaToCoreModel(criteria models.SearchCriteria) model.SearchCriteria {
	return model.SearchCriteria{
		WorkIds:       criteria.WorkIds,
		SearchTerms:   removeEmptyStrings(criteria.SearchTerms),
		ExcludedTerms: removeEmptyStrings(criteria.ExcludedTerms),
		OptionalTerms: removeEmptyStrings(criteria.OptionalTerms),
		Scope:         model.SearchScope(criteria.Scope),
	}
}

func MatchesToApiModels(matches []model.SearchResult) []models.SearchResult {
	resultByWorkId := make(map[int32][]models.Match)
	for _, match := range matches {
		apiMatch := models.Match{
			Snippet:     match.Snippet,
			Text:        match.Text,
			Pages:       match.Pages,
			ParagraphId: match.ParagraphId,
			SentenceId:  match.SentenceId,
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

func removeEmptyStrings(arr []string) []string {
	var result []string
	for _, str := range arr {
		if len(strings.TrimSpace(str)) > 0 {
			result = append(result, str)
		}
	}
	return result
}
