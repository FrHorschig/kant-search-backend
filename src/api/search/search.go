package search

import (
	"strings"

	"github.com/frhorschig/kant-search-api/src/go/models"
	"github.com/frhorschig/kant-search-backend/api/search/internal/errors"
	"github.com/frhorschig/kant-search-backend/api/search/internal/mapping"
	"github.com/frhorschig/kant-search-backend/api/search/internal/validation"
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

	c := mapping.CriteriaToCoreModel(criteria)
	if len(c.WorkIds) == 0 {
		log.Error().Err(err).Msgf("empty work selection: %v", err)
		return errors.BadRequest(ctx, models.BAD_REQUEST_EMPTY_WORKS_SELECTION)
	}
	if len(strings.TrimSpace(c.SearchString)) == 0 {
		log.Error().Err(err).Msgf("empty search terms: %v", err)
		return errors.BadRequest(ctx, models.BAD_REQUEST_EMPTY_SEARCH_TERMS)
	}

	searchString, e := validation.CheckSyntax(c.SearchString)
	if e != nil {
		log.Error().Msgf("validation error in search string: %s", e.Msg)
		return errors.ValidationErrorToApiError(ctx, e)
	}
	c.SearchString = searchString

	matches, err := rec.searchProcessor.Search(ctx.Request().Context(), c)
	if err != nil {
		log.Error().Err(err).Msgf("error while searching for matches: %v", err)
		return errors.InternalServerError(ctx)
	}

	return ctx.JSON(200, mapping.MatchesToApiModels(matches))
}
