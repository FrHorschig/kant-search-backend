package read

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/frhorschig/kant-search-api/src/go/models"
	"github.com/frhorschig/kant-search-backend/api/read/internal/errors"
	"github.com/frhorschig/kant-search-backend/api/read/internal/mapping"
	"github.com/frhorschig/kant-search-backend/core/read"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type ReadHandler interface {
	ReadVolumes(ctx echo.Context) error
	ReadFootnotes(ctx echo.Context) error
	ReadHeadings(ctx echo.Context) error
	ReadParagraphs(ctx echo.Context) error
	ReadSummaries(ctx echo.Context) error
}

type readHandlerImpl struct {
	readProcessor read.ReadProcessor
}

func NewReadHandler(readProcessor read.ReadProcessor) ReadHandler {
	return &readHandlerImpl{readProcessor: readProcessor}
}

func (rec *readHandlerImpl) ReadVolumes(ctx echo.Context) error {
	volumes, err := rec.readProcessor.ProcessVolumes(ctx.Request().Context())
	if err != nil {
		log.Error().Err(err).Msgf("error reading volumes: %v", err)
		return errors.InternalServerError(ctx)
	}

	if len(volumes) == 0 {
		return errors.NotFound(ctx)
	}
	apiVolumes := mapping.VolumesToApiModels(volumes)
	return ctx.JSON(http.StatusOK, apiVolumes)
}

func (rec *readHandlerImpl) ReadFootnotes(ctx echo.Context) error {
	workCode := ctx.Param("workCode")
	if workCode == "" {
		msg := "empty work code"
		log.Error().Msg(msg)
		return errors.BadRequest(ctx, models.BAD_REQUEST_GENERIC, msg)
	}

	ordsParam := ctx.QueryParam("ordinals")
	ordinals, err := findOrdinals(ordsParam)
	if err != nil {
		msg := fmt.Sprintf("invalid ordinal values: %v", ordsParam)
		log.Error().Err(err).Msg(msg)
		return errors.BadRequest(ctx, models.BAD_REQUEST_GENERIC, msg)
	}
	footnotes, err := rec.readProcessor.ProcessFootnotes(ctx.Request().Context(), workCode, ordinals)
	if err != nil {
		log.Error().Err(err).Msgf("error reading footnotes: %v", err)
		return errors.InternalServerError(ctx)
	}

	if len(footnotes) == 0 {
		return errors.NotFound(ctx)
	}
	apiFootnotes := mapping.FootnotesToApiModels(footnotes)
	return ctx.JSON(http.StatusOK, apiFootnotes)
}

func (rec *readHandlerImpl) ReadHeadings(ctx echo.Context) error {
	workCode := ctx.Param("workCode")
	if workCode == "" {
		msg := "empty work code"
		log.Error().Msg(msg)
		return errors.BadRequest(ctx, models.BAD_REQUEST_GENERIC, msg)
	}

	ordsParam := ctx.QueryParam("ordinals")
	ordinals, err := findOrdinals(ordsParam)
	if err != nil {
		msg := fmt.Sprintf("invalid ordinal values: %v", ordsParam)
		log.Error().Err(err).Msg(msg)
		return errors.BadRequest(ctx, models.BAD_REQUEST_GENERIC, msg)
	}
	headings, err := rec.readProcessor.ProcessHeadings(ctx.Request().Context(), workCode, ordinals)
	if err != nil {
		log.Error().Err(err).Msgf("error reading headings: %v", err)
		return errors.InternalServerError(ctx)
	}

	if len(headings) == 0 {
		return errors.NotFound(ctx)
	}
	apiHeadings := mapping.HeadingsToApiModels(headings)
	return ctx.JSON(http.StatusOK, apiHeadings)
}

func (rec *readHandlerImpl) ReadParagraphs(ctx echo.Context) error {
	workCode := ctx.Param("workCode")
	if workCode == "" {
		msg := "empty work code"
		log.Error().Msg(msg)
		return errors.BadRequest(ctx, models.BAD_REQUEST_GENERIC, msg)
	}

	ordsParam := ctx.QueryParam("ordinals")
	ordinals, err := findOrdinals(ordsParam)
	if err != nil {
		msg := fmt.Sprintf("invalid ordinal values: %v", ordsParam)
		log.Error().Err(err).Msg(msg)
		return errors.BadRequest(ctx, models.BAD_REQUEST_GENERIC, msg)
	}
	paragraphs, err := rec.readProcessor.ProcessParagraphs(ctx.Request().Context(), workCode, ordinals)
	if err != nil {
		log.Error().Err(err).Msgf("error reading paragraphs: %v", err)
		return errors.InternalServerError(ctx)
	}

	if len(paragraphs) == 0 {
		return errors.NotFound(ctx)
	}
	apiParagraphs := mapping.ParagraphsToApiModels(paragraphs)
	return ctx.JSON(http.StatusOK, apiParagraphs)
}

func (rec *readHandlerImpl) ReadSummaries(ctx echo.Context) error {
	workCode := ctx.Param("workCode")
	if workCode == "" {
		msg := "empty work code"
		log.Error().Msg(msg)
		return errors.BadRequest(ctx, models.BAD_REQUEST_GENERIC, msg)
	}

	ordsParam := ctx.QueryParam("ordinals")
	ordinals, err := findOrdinals(ordsParam)
	if err != nil {
		msg := fmt.Sprintf("invalid ordinal values: %v", ordsParam)
		log.Error().Err(err).Msg(msg)
		return errors.BadRequest(ctx, models.BAD_REQUEST_GENERIC, msg)
	}
	summaries, err := rec.readProcessor.ProcessSummaries(ctx.Request().Context(), workCode, ordinals)
	if err != nil {
		log.Error().Err(err).Msgf("error reading summaries: %v", err)
		return errors.InternalServerError(ctx)
	}

	if len(summaries) == 0 {
		return errors.NotFound(ctx)
	}
	apiSummaries := mapping.SummariesToApiModels(summaries)
	return ctx.JSON(http.StatusOK, apiSummaries)
}

func findOrdinals(ordsParam string) ([]int32, error) {
	ords := []int32{}
	parts := strings.Split(ordsParam, ",")
	for _, part := range parts {
		ordStr := strings.TrimSpace(part)
		if ordStr == "" {
			continue
		}
		ord, err := strconv.ParseInt(ordStr, 10, 32)
		if err != nil {
			return nil, err
		}
		ords = append(ords, int32(ord))
	}
	return ords, nil
}
