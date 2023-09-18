package handlers

import (
	"github.com/FrHorschig/kant-search-api/models"
	"github.com/FrHorschig/kant-search-backend/api/errors"
	"github.com/FrHorschig/kant-search-backend/api/mapper"
	"github.com/FrHorschig/kant-search-backend/database/model"
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
	err := ctx.Bind(criteria)
	if err != nil || len(criteria.SearchTerms) == 0 || len(criteria.WorkIds) == 0 {
		log.Error().Err(err).Msgf("Error parsing search criteria: %v", err)
		return errors.BadRequest(ctx, "Error parsing search criteria")
	}

	c := mapper.CriteriaToCoreModel(*criteria)
	if c.Scope == model.PARAGRAPH {
		matches, err := rec.searchRepo.SearchParagraphs(ctx.Request().Context(), c)
		if err != nil {
			log.Error().Err(err).Msgf("Error searching for matches: %v", err)
			return errors.InternalServerError(ctx)
		}
		if len(matches) == 0 {
			return errors.NotFound(ctx, "No matches found")
		}
		return ctx.JSON(200, mapper.MatchesToApiModels(matches))
	} else {
		return errors.NotImplemented(ctx, "Sentence scope is not implemented")
	}
}
