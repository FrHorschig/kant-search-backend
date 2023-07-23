package mapper

import (
	"github.com/FrHorschig/kant-search-api/models"
	"github.com/FrHorschig/kant-search-backend/core/model"
)

func WorkToCoreModel(work models.Work) model.Work {
	return model.Work{
		Title:        work.Title,
		Abbreviation: work.Abbreviation,
		Text:         work.Text,
		Volume:       work.Volume,
		Year:         work.Year,
	}
}

func WorkMetadataToApiModel(works []model.WorkMetadata) []models.WorkMetadata {
	apiModels := make([]models.WorkMetadata, 0)
	for _, work := range works {
		apiModels = append(apiModels, models.WorkMetadata{
			Id:           work.Id,
			Title:        work.Title,
			Abbreviation: work.Abbreviation,
			Volume:       work.Volume,
			Year:         work.Year,
		})
	}
	return apiModels
}
