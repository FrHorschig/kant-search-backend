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
	Search(ctx echo.Context) error
}

type searchHandlerImpl struct {
	searchProcessor search.SearchProcessor
}

func NewSearchHandler(searchProcessor search.SearchProcessor) SearchHandler {
	return &searchHandlerImpl{searchProcessor: searchProcessor}
}

func (rec *searchHandlerImpl) Search(ctx echo.Context) error {
	criteria := new(models.SearchCriteria)
	err := ctx.Bind(criteria)
	if err != nil || len(criteria.SearchTerms) == 0 || len(criteria.SearchTerms[0]) == 0 || len(criteria.WorkIds) == 0 {
		log.Error().Err(err).Msgf("Error parsing search criteria: %v", err)
		return errors.BadRequest(ctx, "Error parsing search criteria")
	}

	c := mapper.CriteriaToCoreModel(*criteria)
	matches, err := rec.searchProcessor.Search(ctx.Request().Context(), c)
	if err != nil {
		log.Error().Err(err).Msgf("Error searching for matches: %v", err)
		return errors.InternalServerError(ctx)
	}
	if len(matches) == 0 {
		return errors.NotFound(ctx, "No matches found")
	}

	return ctx.JSON(200, mapper.MatchesToApiModels(matches))
}
