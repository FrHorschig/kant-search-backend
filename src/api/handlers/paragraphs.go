package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/FrHorschig/kant-search-backend/api/errors"
	"github.com/FrHorschig/kant-search-backend/api/mapper"
	"github.com/FrHorschig/kant-search-backend/core/read"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type ParagraphHandler interface {
	GetParagraph(ctx echo.Context) error
	GetParagraphs(ctx echo.Context) error
}

type paragraphHandlerImpl struct {
	paragraphReader read.ParagraphReader
}

func NewParagraphHandler(paragraphReader read.ParagraphReader) ParagraphHandler {
	handlers := paragraphHandlerImpl{
		paragraphReader: paragraphReader,
	}
	return &handlers
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

	paragraph, err := rec.paragraphReader.Find(ctx.Request().Context(), int32(workId), int32(paragraphId))
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
	workId, err := strconv.ParseInt(ctx.Param("id"), 10, 32)
	if err != nil {
		log.Error().Err(err).Msgf("Error parsing work id: %v", err)
		return errors.BadRequest(ctx, "Invalid work id")
	}
	start, end, err := findPages(ctx.QueryParam("pages"))
	if err != nil {
		log.Error().Err(err).Msgf("Error parsing page range: %v", err)
		return errors.BadRequest(ctx, "Invalid page range")
	}

	paragraphs, err := rec.paragraphReader.FindOfPages(ctx.Request().Context(), int32(workId), start, end)
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

func findPages(pageRange string) (start int32, end int32, err error) {
	parseError := fmt.Errorf("invalid page range: %s", pageRange)
	parts := strings.Split(pageRange, "-")
	if len(parts) != 2 {
		return -1, -1, parseError
	}
	s, err := strconv.ParseInt(parts[0], 10, 32)
	if err != nil {
		return -1, -1, parseError
	}
	e, err := strconv.ParseInt(parts[1], 10, 32)
	if err != nil {
		return -1, -1, parseError
	}
	return int32(s), int32(e), nil
}
