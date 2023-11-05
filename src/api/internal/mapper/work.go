package mapper

import (
	"github.com/frhorschig/kant-search-api/src/go/models"
	"github.com/frhorschig/kant-search-backend/api/internal/util"
	"github.com/frhorschig/kant-search-backend/common/model"
)

func WorkUploadToCoreModel(work models.WorkUpload) model.WorkUpload {
	return model.WorkUpload{
		WorkId: work.WorkId,
		Text:   work.Text,
	}
}

func WorksToApiModels(works []model.Work) []models.Work {
	apiModels := make([]models.Work, 0)
	for _, work := range works {
		apiModels = append(apiModels, models.Work{
			Id:           work.Id,
			Title:        work.Title,
			Abbreviation: util.ToStrVal(work.Abbreviation),
			Ordinal:      work.Ordinal,
			Year:         util.ToStrVal(work.Year),
			VolumeId:     work.Volume,
		})
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
