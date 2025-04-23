package mapping

import (
	"github.com/frhorschig/kant-search-api/src/go/models"
	"github.com/frhorschig/kant-search-backend/dataaccess/model"
)

func CriteriaToCoreModel(in *models.SearchCriteria) (string, model.SearchOptions) {
	return in.SearchString, model.SearchOptions{
		WorkIds:          in.Options.WorkIds,
		Scope:            model.SearchScope(in.Options.Scope),
		IncludeHeadings:  in.Options.IncludeHeadings,
		IncludeFootnotes: in.Options.IncludeFootnotes,
		IncludeSummaries: in.Options.IncludeSummaries,
	}
}

func MatchesToApiModels(hits []model.SearchResult) []models.SearchResult {
	resultByWorkId := make(map[string][]models.Hit)
	for _, hit := range hits {
		apiHit := models.Hit{
			Snippets:  hit.Snippets,
			Pages:     hit.Pages,
			ContentId: hit.ContentId,
		}

		arr, exists := resultByWorkId[hit.WorkId]
		if exists {
			arr = append(arr, apiHit)
		} else {
			arr = []models.Hit{apiHit}
		}
		resultByWorkId[hit.WorkId] = arr
	}

	var results []models.SearchResult
	for workId, apiHits := range resultByWorkId {
		results = append(results, models.SearchResult{
			WorkId: workId,
			Hits:   apiHits,
		})
	}
	return results
}
