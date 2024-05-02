package mapper

import (
	"sort"

	"github.com/frhorschig/kant-search-api/generated/go/models"
	"github.com/frhorschig/kant-search-backend/common/model"
)

func CriteriaToCoreModel(criteria models.SearchCriteria) model.SearchCriteria {
	return model.SearchCriteria{
		WorkIds:      criteria.WorkIds,
		SearchString: criteria.SearchString,
		Options:      searchOptionsToCoreModel(criteria.Options),
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

func searchOptionsToCoreModel(options models.SearchOptions) model.SearchOptions {
	return model.SearchOptions{
		Scope: model.SearchScope(options.Scope),
	}
}
