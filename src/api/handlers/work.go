package handlers

import (
	"io"
	"net/http"
	"strconv"

	"github.com/frhorschig/kant-search-api/generated/go/models"
	"github.com/frhorschig/kant-search-backend/api/internal/errors"
	"github.com/frhorschig/kant-search-backend/api/internal/mapper"
	"github.com/frhorschig/kant-search-backend/core/upload"
	"github.com/frhorschig/kant-search-backend/database"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type WorkHandler interface {
	GetVolumes(ctx echo.Context) error
	GetWorks(ctx echo.Context) error
	PostWork(ctx echo.Context) error
}

type workHandlerImpl struct {
	volumeRepo      database.VolumeRepo
	workRepo        database.WorkRepo
	uploadProcessor upload.WorkUploadProcessor
}

func NewWorkHandler(volumeRepo database.VolumeRepo, workRepo database.WorkRepo, uploadProcessor upload.WorkUploadProcessor) WorkHandler {
	return &workHandlerImpl{
		volumeRepo:      volumeRepo,
		workRepo:        workRepo,
		uploadProcessor: uploadProcessor,
	}
}

func (rec *workHandlerImpl) GetVolumes(ctx echo.Context) error {
	volumes, err := rec.volumeRepo.SelectAll(ctx.Request().Context())
	if err != nil {
		log.Error().Err(err).Msgf("Error reading volumes: %v", err)
		return errors.InternalServerError(ctx)
	}

	if len(volumes) == 0 {
		return errors.NotFound(ctx, models.NOT_FOUND_VOLUMES)
	}

	apiVolumes := mapper.VolumesToApiModels(volumes)
	return ctx.JSON(http.StatusOK, apiVolumes)
}

func (rec *workHandlerImpl) GetWorks(ctx echo.Context) error {
	works, err := rec.workRepo.SelectAll(ctx.Request().Context())
	if err != nil {
		log.Error().Err(err).Msgf("Error reading works: %v", err)
		return errors.InternalServerError(ctx)
	}

	if len(works) == 0 {
		return errors.NotFound(ctx, models.NOT_FOUND_WORKS)
	}

	apiWorks := mapper.WorksToApiModels(works)
	return ctx.JSON(http.StatusOK, apiWorks)
}

func (rec *workHandlerImpl) PostWork(ctx echo.Context) error {
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
		return errors.CoreError(ctx, coreErr)
	}

	return ctx.NoContent(http.StatusCreated)
}
