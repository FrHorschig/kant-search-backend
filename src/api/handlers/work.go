package handlers

import (
	"net/http"

	"github.com/FrHorschig/kant-search-api/models"
	"github.com/FrHorschig/kant-search-backend/api/errors"
	"github.com/FrHorschig/kant-search-backend/api/mapper"
	"github.com/FrHorschig/kant-search-backend/core/read"
	processing "github.com/FrHorschig/kant-search-backend/core/upload"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type WorkHandler interface {
	PostWork(ctx echo.Context) error
	GetWork(ctx echo.Context) error
}

type WorkHandlerImpl struct {
	workProcessor processing.WorkProcessor
	workReader    read.WorkReader
}

func NewWorkHandler(workProcessor processing.WorkProcessor, workReader read.WorkReader) WorkHandler {
	impl := WorkHandlerImpl{
		workProcessor: workProcessor,
		workReader:    workReader,
	}
	return &impl
}

func (rec *WorkHandlerImpl) PostWork(ctx echo.Context) error {
	work := new(models.Work)
	if err := ctx.Bind(work); err != nil {
		log.Error().Err(err).Msg("Error reading request body")
		return errors.BadRequest(ctx, "Error reading request body")
	}

	context := ctx.Request().Context()
	coreModel := mapper.WorkToCoreModel(*work)
	err := rec.workProcessor.Process(context, coreModel)
	if err != nil {
		log.Error().Err(err).Msgf("Error processing work: %v", err)
		return errors.InternalServerError(ctx)
	}
	return ctx.NoContent(http.StatusOK)
}

func (rec *WorkHandlerImpl) GetWork(ctx echo.Context) error {
	works, err := rec.workReader.FindAll(ctx.Request().Context())
	if err != nil {
		log.Error().Err(err).Msgf("Error reading works: %v", err)
		return errors.InternalServerError(ctx)
	}

	apiWorks := mapper.WorkMetadataToApiModel(works)
	return ctx.JSON(http.StatusOK, apiWorks)
}
