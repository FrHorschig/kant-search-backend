package mapper

import (
	"github.com/FrHorschig/kant-search-api/models"
	"github.com/FrHorschig/kant-search-backend/core/model"
)

func WorkUploadToCoreModel(work models.WorkUpload) model.Work {
	return model.Work{
		Title:        work.Title,
		Abbreviation: work.Abbreviation,
		Text:         work.Text,
		Volume:       work.Volume,
		Ordinal:      work.Ordinal,
		Year:         work.Year,
	}
}

func WorkToApiModel(works []model.Work) []models.Work {
	apiModels := make([]models.Work, 0)
	for _, work := range works {
		apiModels = append(apiModels, models.Work{
			Id:           work.Id,
			Title:        work.Title,
			Abbreviation: work.Abbreviation,
			Volume:       work.Volume,
			Year:         work.Year,
		})
	}
	return apiModels
}
