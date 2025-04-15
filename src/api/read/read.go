package read

import (
	"net/http"

	"github.com/frhorschig/kant-search-backend/api/internal/errors"
	"github.com/frhorschig/kant-search-backend/api/internal/mapping"
	"github.com/frhorschig/kant-search-backend/core/read"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type ReadHandler interface {
	ReadVolumes(ctx echo.Context) error
	ReadWork(ctx echo.Context) error
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

	apiVolumes := mapping.VolumesToApiModels(volumes)
	return ctx.JSON(http.StatusOK, apiVolumes)
}

func (rec *readHandlerImpl) ReadWork(ctx echo.Context) error {
	workId := ctx.Param("workId")
	if workId == "" {
		msg := "empty work ID"
		log.Error().Msgf(msg)
		return errors.BadRequest(ctx, msg)
	}

	work, err := rec.readProcessor.ProcessWork(ctx.Request().Context(), workId)
	if err != nil {
		log.Error().Err(err).Msgf("error reading work: %v", err)
		return errors.InternalServerError(ctx)
	}

	if work == nil {
		return errors.NotFound(ctx)
	}
	apiWork := mapping.WorkToApiModels(work)
	return ctx.JSON(http.StatusOK, apiWork)
}

func (rec *readHandlerImpl) ReadFootnotes(ctx echo.Context) error {
	workId := ctx.Param("workId")
	if workId == "" {
		msg := "empty work ID"
		log.Error().Msgf(msg)
		return errors.BadRequest(ctx, msg)
	}

	footnotes, err := rec.readProcessor.ProcessFootnotes(ctx.Request().Context(), workId)
	if err != nil {
		log.Error().Err(err).Msgf("error reading footnotes: %v", err)
		return errors.InternalServerError(ctx)
	}

	apiFootnotes := mapping.FootnotesToApiModels(footnotes)
	return ctx.JSON(http.StatusOK, apiFootnotes)
}

func (rec *readHandlerImpl) ReadHeadings(ctx echo.Context) error {
	workId := ctx.Param("workId")
	if workId == "" {
		msg := "empty work ID"
		log.Error().Msgf(msg)
		return errors.BadRequest(ctx, msg)
	}

	headings, err := rec.readProcessor.ProcessHeadings(ctx.Request().Context(), workId)
	if err != nil {
		log.Error().Err(err).Msgf("error reading headings: %v", err)
		return errors.InternalServerError(ctx)
	}

	apiHeadings := mapping.HeadingsToApiModels(headings)
	return ctx.JSON(http.StatusOK, apiHeadings)
}

func (rec *readHandlerImpl) ReadParagraphs(ctx echo.Context) error {
	workId := ctx.Param("workId")
	if workId == "" {
		msg := "empty work ID"
		log.Error().Msgf(msg)
		return errors.BadRequest(ctx, msg)
	}

	paragraphs, err := rec.readProcessor.ProcessParagraphs(ctx.Request().Context(), workId)
	if err != nil {
		log.Error().Err(err).Msgf("error reading paragraphs: %v", err)
		return errors.InternalServerError(ctx)
	}

	apiParagraphs := mapping.ParagraphsToApiModels(paragraphs)
	return ctx.JSON(http.StatusOK, apiParagraphs)
}

func (rec *readHandlerImpl) ReadSummaries(ctx echo.Context) error {
	workId := ctx.Param("workId")
	if workId == "" {
		msg := "empty work ID"
		log.Error().Msgf(msg)
		return errors.BadRequest(ctx, msg)
	}

	summaries, err := rec.readProcessor.ProcessSummaries(ctx.Request().Context(), workId)
	if err != nil {
		log.Error().Err(err).Msgf("error reading summaries: %v", err)
		return errors.InternalServerError(ctx)
	}

	apiSummaries := mapping.SummariesToApiModels(summaries)
	return ctx.JSON(http.StatusOK, apiSummaries)
}
