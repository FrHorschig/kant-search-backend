package mapper

import (
	api "github.com/FrHorschig/kant-search-api/models"
	"github.com/FrHorschig/kant-search-backend/database/models"
)

func MapTextFromDb(text models.Text) api.Text {
	return api.Text{
		Id:   text.Id,
		Text: text.Text.String,
	}
}
