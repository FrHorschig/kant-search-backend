package search

import (
	"strings"

	"github.com/frhorschig/kant-search-api/src/go/models"
	"github.com/frhorschig/kant-search-backend/api/search/internal/errors"
	"github.com/frhorschig/kant-search-backend/api/search/internal/mapping"
	"github.com/frhorschig/kant-search-backend/core/search"
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
	if err != nil {
		log.Error().Err(err).Msgf("error parsing search criteria: %v", err)
		return errors.BadRequest(ctx, models.BAD_REQUEST_INVALID_SEARCH_CRITERIA)
	}

	searchTerms, options := mapping.CriteriaToCoreModel(criteria)
	if len(strings.TrimSpace(searchTerms)) == 0 {
		log.Error().Err(err).Msg("empty search terms")
		return errors.BadRequest(ctx, models.BAD_REQUEST_EMPTY_SEARCH_TERMS)
	}
	if len(options.WorkCodes) == 0 {
		log.Error().Err(err).Msg("empty work selection")
		return errors.BadRequest(ctx, models.BAD_REQUEST_EMPTY_WORKS_SELECTION)
	}

	matches, searchErr := rec.searchProcessor.Search(ctx.Request().Context(), searchTerms, options)
	if searchErr.HasError {
		if searchErr.SyntaxError != nil {
			e := searchErr.SyntaxError
			log.Error().Msgf("syntax error in search string: %s", e.Msg)
			return errors.SyntaxErrorToApiError(ctx, e)
		} else {
			log.Error().Err(err).Msgf("error while searching for matches: %v", err)
			return errors.InternalServerError(ctx)
		}
	}

	return ctx.JSON(200, mapping.HitsToApiModels(matches))
}
