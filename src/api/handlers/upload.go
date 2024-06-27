package handlers

import (
	"io"
	"net/http"
	"strconv"

	"github.com/frhorschig/kant-search-api/generated/go/models"
	"github.com/frhorschig/kant-search-backend/api/internal/errors"
	"github.com/frhorschig/kant-search-backend/core/upload"
	"github.com/frhorschig/kant-search-backend/database"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type UploadHandler interface {
	PostWork(ctx echo.Context) error
}

type uploadHandlerImpl struct {
	workRepo        database.WorkRepo
	uploadProcessor upload.WorkUploadProcessor
}

func NewUploadHandler(workRepo database.WorkRepo, uploadProcessor upload.WorkUploadProcessor) UploadHandler {
	return &uploadHandlerImpl{
		workRepo:        workRepo,
		uploadProcessor: uploadProcessor,
	}
}

func (rec *uploadHandlerImpl) PostWork(ctx echo.Context) error {
	workId, err := strconv.ParseInt(ctx.Param("workId"), 10, 32)
	if err != nil {
		log.Error().Err(err).Msgf("Error parsing work id: %v", err)
		return errors.BadRequest(ctx, models.BAD_REQUEST_INVALID_WORK_SELECTION)
	}

	body, err := io.ReadAll(ctx.Request().Body)
	if err != nil {
		log.Error().Err(err).Msg("Error reading request body")
		return errors.BadRequest(ctx, models.BAD_REQUEST_GENERIC)
	}
	text := string(body)
	if text == "" {
		log.Error().Err(err).Msg("Empty text")
		return errors.BadRequest(ctx, models.BAD_REQUEST_EMPTY_WORK_TEXT)
	}

	coreErr := rec.uploadProcessor.Process(ctx.Request().Context(), int32(workId), text)
	if coreErr != nil {
		log.Error().Str("Msg", string(coreErr.Msg)).Interface("Params", coreErr.Params).Msg("Error processing work")
		return errors.UploadError(ctx, coreErr)
	}

	return ctx.NoContent(http.StatusCreated)
}
