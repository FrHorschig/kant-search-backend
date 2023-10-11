package mapper

import (
	"github.com/FrHorschig/kant-search-api/models"
	"github.com/FrHorschig/kant-search-backend/database/model"
)

func ParagraphToApiModel(paragraph model.Paragraph) models.Paragraph {
	return models.Paragraph{
		Id:     paragraph.Id,
		Text:   paragraph.Text,
		Pages:  paragraph.Pages,
		WorkId: paragraph.WorkId,
	}
}

func ParagraphsToApiModels(paragraphs []model.Paragraph) []models.Paragraph {
	apiModels := make([]models.Paragraph, 0)
	for _, paragraph := range paragraphs {
		apiModels = append(apiModels, ParagraphToApiModel(paragraph))
	}
	return apiModels
}
