package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/FrHorschig/kant-search-backend/api/errors"
	"github.com/FrHorschig/kant-search-backend/api/mapper"
	"github.com/FrHorschig/kant-search-backend/database/repository"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type ParagraphHandler interface {
	GetParagraph(ctx echo.Context) error
	GetParagraphs(ctx echo.Context) error
}

type paragraphHandlerImpl struct {
	paragraphRepo repository.ParagraphRepo
}

func NewParagraphHandler(paragraphRepo repository.ParagraphRepo) ParagraphHandler {
	return &paragraphHandlerImpl{paragraphRepo: paragraphRepo}
}

func (rec *paragraphHandlerImpl) GetParagraph(ctx echo.Context) error {
	workId, err := strconv.ParseInt(ctx.Param("workId"), 10, 32)
	if err != nil {
		log.Error().Err(err).Msgf("Error parsing work id: %v", err)
		return errors.BadRequest(ctx, "Invalid work id")
	}
	paragraphId, err := strconv.ParseInt(ctx.Param("paragraphId"), 10, 32)
	if err != nil {
		log.Error().Err(err).Msgf("Error parsing paragraph id: %v", err)
		return errors.BadRequest(ctx, "Invalid paragraph id")
	}

	paragraph, err := rec.paragraphRepo.Select(ctx.Request().Context(), int32(workId), int32(paragraphId))
	// TODO frhorsch: return ptr so we don't have to check a database error here
	if err == sql.ErrNoRows {
		return errors.NotFound(ctx, fmt.Sprintf("Paragraph with id %d not found", paragraphId))
	}
	if err != nil {
		log.Error().Err(err).Msgf("Error reading paragraph: %v", err)
		return errors.InternalServerError(ctx)
	}
	return ctx.JSON(http.StatusOK, mapper.ParagraphToApiModel(paragraph))
}

func (rec *paragraphHandlerImpl) GetParagraphs(ctx echo.Context) error {
	workId, err := strconv.ParseInt(ctx.Param("workId"), 10, 32)
	if err != nil {
		log.Error().Err(err).Msgf("Error parsing work id: %v", err)
		return errors.BadRequest(ctx, "Invalid work id")
	}

	paragraphs, err := rec.paragraphRepo.SelectAll(ctx.Request().Context(), int32(workId))
	if err != nil {
		log.Error().Err(err).Msgf("Error reading paragraphs: %v", err)
		return errors.InternalServerError(ctx)
	}
	if len(paragraphs) == 0 {
		return errors.NotFound(ctx, "No paragraphs found")
	}

	apiParas := mapper.ParagraphsToApiModels(paragraphs)
	return ctx.JSON(http.StatusOK, apiParas)
}
