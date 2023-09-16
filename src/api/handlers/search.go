package handlers

import (
	"github.com/FrHorschig/kant-search-api/models"
	"github.com/FrHorschig/kant-search-backend/api/errors"
	"github.com/FrHorschig/kant-search-backend/api/mapper"
	"github.com/FrHorschig/kant-search-backend/database/repository"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type SearchHandler interface {
	SearchParagraphs(ctx echo.Context) error
}

type searchHandlerImpl struct {
	searchRepo repository.SearchRepo
}

func NewSearchHandler(searchRepo repository.SearchRepo) SearchHandler {
	return &searchHandlerImpl{searchRepo: searchRepo}
}

func (rec *searchHandlerImpl) SearchParagraphs(ctx echo.Context) error {
	criteria := new(models.SearchCriteria)
	if err := ctx.Bind(criteria); err != nil {
		log.Error().Err(err).Msgf("Error parsing search criteria: %v", err)
		return errors.BadRequest(ctx, err.Error())
	}

	c := mapper.CriteriaToCoreModel(*criteria)
	matches, err := rec.searchRepo.SearchParagraphs(ctx.Request().Context(), c)
	if err != nil {
		log.Error().Err(err).Msgf("Error searching for matches: %v", err)
		return errors.InternalServerError(ctx)
	}

	return ctx.JSON(200, mapper.MatchesToApiModels(matches))
}
