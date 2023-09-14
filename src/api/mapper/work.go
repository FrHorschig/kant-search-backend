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

func WorkToApiModel(works []model.Work) []models.Work {
	apiModels := make([]models.Work, 0)
	for _, work := range works {
		apiModels = append(apiModels, models.Work{
			Id:           work.Id,
			Title:        work.Title,
			Abbreviation: *work.Abbreviation,
			Year:         *work.Year,
			VolumeId:     work.Volume,
		})
	}
	return apiModels
}

func VolumeToApiModel(volumes []model.Volume) []models.Volume {
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
