package mapping

import (
	"github.com/frhorschig/kant-search-api/src/go/models"
	"github.com/frhorschig/kant-search-backend/dataaccess/model"
)

func CriteriaToCoreModel(in *models.SearchCriteria) (string, model.SearchOptions) {
	return in.SearchString, model.SearchOptions{
		WorkIds:         in.Options.WorkIds,
		Scope:           model.SearchScope(in.Options.Scope),
		IncludeHeadings: in.Options.IncludeHeadings,
	}
}

func MatchesToApiModels(matches []model.SearchResult) []models.SearchResult {
	resultByWorkId := make(map[string][]models.Match)
	for _, match := range matches {
		apiMatch := models.Match{
			Snippet:   match.Snippet,
			Text:      match.Text,
			Pages:     match.Pages,
			ContentId: match.ContentId,
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
