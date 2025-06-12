package mapping

import (
	"fmt"

	"github.com/frhorschig/kant-search-api/src/go/models"
	"github.com/frhorschig/kant-search-backend/dataaccess/model"
)

func CriteriaToCoreModel(in *models.SearchCriteria) (string, model.SearchOptions) {
	return in.SearchTerms, model.SearchOptions{
		IncludeHeadings:  in.Options.IncludeHeadings,
		IncludeFootnotes: in.Options.IncludeFootnotes,
		IncludeSummaries: in.Options.IncludeSummaries,
		Scope:            model.SearchScope(in.Options.Scope),
		WithStemming:     in.Options.WithStemming,
		WorkCodes:        in.Options.WorkCodes,
	}
}

func HitsToApiModels(hits []model.SearchResult) []models.SearchResult {
	resultByWorkCode := make(map[string][]models.Hit)
	for _, hit := range hits {
		wim := make(map[string]int32)
		for k, v := range hit.WordIndexMap {
			wim[fmt.Sprint(k)] = v
		}

		apiHit := models.Hit{
			HighlightText: hit.HighlightText,
			FmtText:       hit.FmtText,
			PageByIndex:   mapIndexByNumberPairs(hit.PageByIndex),
			LineByIndex:   mapIndexByNumberPairs(hit.LineByIndex),
			Ordinal:       hit.Ordinal,
			WordIndexMap:  wim,
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

func mapIndexByNumberPairs(in []model.IndexNumberPair) []models.IndexNumberPair {
	result := []models.IndexNumberPair{}
	for _, pair := range in {
		result = append(result,
			models.IndexNumberPair{I: pair.I, Num: pair.Num},
		)
	}
	return result
}
