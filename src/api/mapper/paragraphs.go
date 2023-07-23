package mapper

import (
	"github.com/FrHorschig/kant-search-api/models"
	"github.com/FrHorschig/kant-search-backend/core/model"
)

func ParagraphsToApiModel(paragraphs []model.Paragraph) []models.Paragraph {
	apiModels := make([]models.Paragraph, 0)
	for _, paragraph := range paragraphs {
		apiModels = append(apiModels, models.Paragraph{
			Id:     paragraph.Id,
			Text:   paragraph.Text,
			Pages:  paragraph.Pages,
			WorkId: paragraph.WorkId,
		})
	}
	return apiModels
}
