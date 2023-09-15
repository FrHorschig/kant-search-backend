package mapper

import (
	"github.com/FrHorschig/kant-search-api/models"
	"github.com/FrHorschig/kant-search-backend/database/model"
)

func WorkUploadToCoreModel(work models.WorkUpload) model.WorkUpload {
	return model.WorkUpload{
		WorkId: work.WorkId,
		Text:   work.Text,
	}
}

func WorksToApiModels(works []model.Work) []models.Work {
	apiModels := make([]models.Work, 0)
	for i, work := range works {
		apiModels = append(apiModels, models.Work{
			Id:       work.Id,
			Title:    work.Title,
			Ordinal:  work.Ordinal,
			VolumeId: work.Volume,
		})
		if work.Abbreviation != nil {
			apiModels[i].Abbreviation = *work.Abbreviation
		}
		if work.Year != nil {
			apiModels[i].Year = *work.Year
		}
	}
	return apiModels
}

func VolumesToApiModels(volumes []model.Volume) []models.Volume {
	apiModels := make([]models.Volume, 0)
	for _, work := range volumes {
		apiModels = append(apiModels, models.Volume{
			Id:      work.Id,
			Title:   work.Title,
			Section: work.Section,
		})
	}
	return apiModels
}
