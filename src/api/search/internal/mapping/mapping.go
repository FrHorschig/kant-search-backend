package mapping

import (
	"github.com/frhorschig/kant-search-api/src/go/models"
	"github.com/frhorschig/kant-search-backend/dataaccess/model"
)

func CriteriaToCoreModel(in *models.SearchCriteria) (string, model.SearchOptions) {
	return in.SearchTerms, model.SearchOptions{
		IncludeHeadings:  in.Options.IncludeHeadings,
		IncludeFootnotes: in.Options.IncludeFootnotes,
		IncludeSummaries: in.Options.IncludeSummaries,
		Scope:            model.SearchScope(in.Options.Scope),
		WorkCodes:        in.Options.WorkCodes,
	}
}

func HitsToApiModels(hits []model.SearchResult) []models.SearchResult {
	resultByWorkCode := make(map[string][]models.Hit)
	for _, hit := range hits {
		apiHit := models.Hit{
			Snippets: hit.Snippets,
			Pages:    hit.Pages,
			Ordinal:  hit.Ordinal,
			Text:     hit.Text,
		}

		arr, exists := resultByWorkCode[hit.WorkCode]
		if exists {
			arr = append(arr, apiHit)
		} else {
			arr = []models.Hit{apiHit}
		}
		resultByWorkCode[hit.WorkCode] = arr
	}

	var results []models.SearchResult
	for workCode, apiHits := range resultByWorkCode {
		results = append(results, models.SearchResult{
			Hits:     apiHits,
			WorkCode: workCode,
		})
	}
	return results
}
