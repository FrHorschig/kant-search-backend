package handlers

import (
	"net/http"
	"strconv"

	"github.com/frhorschig/kant-search-api/models"
	"github.com/frhorschig/kant-search-backend/api/internal/errors"
	"github.com/frhorschig/kant-search-backend/api/internal/mapper"
	repository "github.com/frhorschig/kant-search-backend/database"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type ParagraphHandler interface {
	GetParagraphs(ctx echo.Context) error
}

type paragraphHandlerImpl struct {
	paragraphRepo repository.ParagraphRepo
}

func NewParagraphHandler(paragraphRepo repository.ParagraphRepo) ParagraphHandler {
	return &paragraphHandlerImpl{paragraphRepo: paragraphRepo}
}

func (rec *paragraphHandlerImpl) GetParagraphs(ctx echo.Context) error {
	workId, err := strconv.ParseInt(ctx.Param("workId"), 10, 32)
	if err != nil {
		log.Error().Err(err).Msgf("Error parsing work id: %v", err)
		return errors.BadRequest(ctx, models.BAD_REQUEST_INVALID_WORK_SELECTION)
	}

	paragraphs, err := rec.paragraphRepo.SelectAll(ctx.Request().Context(), int32(workId))
	if err != nil {
		log.Error().Err(err).Msgf("Error reading paragraphs: %v", err)
		return errors.InternalServerError(ctx)
	}
	if len(paragraphs) == 0 {
		return errors.NotFound(ctx, models.NOT_FOUND_PARAGRAPHS)
	}

	apiParas := mapper.ParagraphsToApiModels(paragraphs)
	return ctx.JSON(http.StatusOK, apiParas)
}
