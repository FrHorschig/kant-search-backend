package handlers

import (
	"github.com/FrHorschig/kant-search-api/models"
	"github.com/FrHorschig/kant-search-backend/api/errors"
	"github.com/FrHorschig/kant-search-backend/api/mapper"
	"github.com/FrHorschig/kant-search-backend/core/search"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type SearchHandler interface {
	SearchParagraphs(ctx echo.Context) error
}

type searchHandlerImpl struct {
	searcher search.Searcher
}

func NewSearchHandler(searcher search.Searcher) SearchHandler {
	impl := searchHandlerImpl{searcher: searcher}
	return &impl
}

func (rec *searchHandlerImpl) SearchParagraphs(ctx echo.Context) error {
	criteria := new(models.SearchCriteria)
	if err := ctx.Bind(criteria); err != nil {
		log.Error().Err(err).Msgf("Error parsing search criteria: %v", err)
		return errors.BadRequest(ctx, err.Error())
	}

	c := mapper.CriteriaToCoreModel(*criteria)
	matches, err := rec.searcher.SearchParagraphs(ctx.Request().Context(), c)
	if err != nil {
		log.Error().Err(err).Msgf("Error searching for matches: %v", err)
		return errors.InternalServerError(ctx)
	}

	return ctx.JSON(200, mapper.MatchesToApiModel(matches))
}
