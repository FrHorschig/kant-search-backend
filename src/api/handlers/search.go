package handlers

import (
	"github.com/FrHorschig/kant-search-api/models"
	"github.com/FrHorschig/kant-search-backend/api/errors"
	"github.com/FrHorschig/kant-search-backend/api/mapper"
	"github.com/FrHorschig/kant-search-backend/core/search"
	"github.com/labstack/echo/v4"
)

type SearchHandler interface {
	SearchParagraphs(ctx echo.Context) error
}

type SearchHandlerImpl struct {
	paragraphSearcher search.ParagraphSearcher
}

func NewSearchHandler(paragraphSearcher search.ParagraphSearcher) SearchHandler {
	impl := SearchHandlerImpl{paragraphSearcher: paragraphSearcher}
	return &impl
}

func (rec *SearchHandlerImpl) SearchParagraphs(ctx echo.Context) error {
	criteria := new(models.SearchCriteria)
	if err := ctx.Bind(criteria); err != nil {
		return errors.BadRequest(ctx, err.Error())
	}

	paragraphs, err := rec.paragraphSearcher.Search(ctx.Request().Context(), criteria.WorkIds)
	if err != nil {
		return errors.InternalServerError(ctx)
	}

	apiParas := mapper.ParagraphsToApiModel(paragraphs)
	return ctx.JSON(200, apiParas)
}
