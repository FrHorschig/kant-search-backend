package read

import (
	"net/http"

	"github.com/frhorschig/kant-search-api/generated/go/models"
	"github.com/frhorschig/kant-search-backend/api/common/errors"
	"github.com/frhorschig/kant-search-backend/api/read/mapper"
	database "github.com/frhorschig/kant-search-backend/dataaccess"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type WorkHandler interface {
	GetVolumes(ctx echo.Context) error
	GetWorks(ctx echo.Context) error
}

type workHandlerImpl struct {
	volumeRepo database.VolumeRepo
	workRepo   database.WorkRepo
}

func NewWorkHandler(volumeRepo database.VolumeRepo, workRepo database.WorkRepo) WorkHandler {
	return &workHandlerImpl{
		volumeRepo: volumeRepo,
		workRepo:   workRepo,
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
