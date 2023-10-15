package handlers

import (
	"net/http"

	"github.com/FrHorschig/kant-search-api/models"
	"github.com/FrHorschig/kant-search-backend/api/internal/errors"
	"github.com/FrHorschig/kant-search-backend/api/internal/mapper"
	processing "github.com/FrHorschig/kant-search-backend/core/upload"
	"github.com/FrHorschig/kant-search-backend/database/repository"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type WorkHandler interface {
	GetVolumes(ctx echo.Context) error
	GetWorks(ctx echo.Context) error
	PostWork(ctx echo.Context) error
}

type workHandlerImpl struct {
	volumeRepo    repository.VolumeRepo
	workRepo      repository.WorkRepo
	workProcessor processing.WorkUploadProcessor
}

func NewWorkHandler(volumeRepo repository.VolumeRepo, workRepo repository.WorkRepo, workProcessor processing.WorkUploadProcessor) WorkHandler {
	return &workHandlerImpl{
		volumeRepo:    volumeRepo,
		workRepo:      workRepo,
		workProcessor: workProcessor,
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
	work := new(models.WorkUpload)
	err := ctx.Bind(work)
	if err != nil {
		log.Error().Err(err).Msg("Error reading request body")
		return errors.BadRequest(ctx, models.BAD_REQUEST_GENERIC)
	}
	if work.WorkId < 1 {
		log.Error().Err(err).Msg("Empty work selection")
		return errors.BadRequest(ctx, models.BAD_REQUEST_EMPTY_WORKS_SELECTION)
	}
	if work.Text == "" {
		log.Error().Err(err).Msg("Empty text")
		return errors.BadRequest(ctx, models.BAD_REQUEST_EMPTY_WORK_TEXT)
	}

	coreModel := mapper.WorkUploadToCoreModel(*work)
	parseErr, err := rec.workProcessor.Process(ctx.Request().Context(), coreModel)
	if parseErr != nil {
		log.Error().Err(err).Msgf("Error parsing work text: %v", parseErr)
		return errors.BadRequestFromCore(ctx, parseErr)
	}
	if err != nil {
		log.Error().Err(err).Msgf("Error processing work: %v", err)
		return errors.InternalServerError(ctx)
	}

	return ctx.NoContent(http.StatusCreated)
}
