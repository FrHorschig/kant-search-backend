package handlers

import (
	"net/http"

	"github.com/FrHorschig/kant-search-api/models"
	"github.com/FrHorschig/kant-search-backend/api/errors"
	"github.com/FrHorschig/kant-search-backend/api/mapper"
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
	workProcessor processing.WorkUploadProcessor
	workRepo      repository.WorkRepo
}

func NewWorkHandler(workProcessor processing.WorkUploadProcessor, workRepo repository.WorkRepo) WorkHandler {
	impl := workHandlerImpl{
		workProcessor: workProcessor,
		workRepo:      workRepo,
	}
	return &impl
}

func (rec *workHandlerImpl) GetVolumes(ctx echo.Context) error {
	works, err := rec.workRepo.SelectAll(ctx.Request().Context())
	if err != nil {
		log.Error().Err(err).Msgf("Error reading works: %v", err)
		return errors.InternalServerError(ctx)
	}

	if len(works) == 0 {
		return errors.NotFound(ctx, "No works found")
	}

	apiWorks := mapper.WorksToApiModels(works)
	return ctx.JSON(http.StatusOK, apiWorks)
}

func (rec *workHandlerImpl) GetWorks(ctx echo.Context) error {
	works, err := rec.workRepo.SelectAll(ctx.Request().Context())
	if err != nil {
		log.Error().Err(err).Msgf("Error reading works: %v", err)
		return errors.InternalServerError(ctx)
	}

	if len(works) == 0 {
		return errors.NotFound(ctx, "No works found")
	}

	apiWorks := mapper.WorksToApiModels(works)
	return ctx.JSON(http.StatusOK, apiWorks)
}

func (rec *workHandlerImpl) PostWork(ctx echo.Context) error {
	work := new(models.WorkUpload)
	if err := ctx.Bind(work); err != nil {
		log.Error().Err(err).Msg("Error reading request body")
		return errors.BadRequest(ctx, "Error reading request body")
	}

	context := ctx.Request().Context()
	coreModel := mapper.WorkUploadToCoreModel(*work)
	err := rec.workProcessor.Process(context, coreModel)
	if err != nil {
		log.Error().Err(err).Msgf("Error processing work: %v", err)
		return errors.InternalServerError(ctx)
	}

	return ctx.NoContent(http.StatusCreated)
}
