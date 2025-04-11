package read

import (
	"github.com/frhorschig/kant-search-backend/core/read"
	"github.com/labstack/echo/v4"
)

type ReadHandler interface {
	Read(ctx echo.Context) error
}

type readHandlerImpl struct {
	readProcessor read.ReadProcessor
}

func NewReadHandler(readProcessor read.ReadProcessor) ReadHandler {
	return &readHandlerImpl{readProcessor: readProcessor}
}

func (rec *readHandlerImpl) Read(ctx echo.Context) error {
	// TODO implement me
	// workId, err := strconv.ParseInt(ctx.Param("workId"), 10, 32)
	// if err != nil {
	// 	log.Error().Err(err).Msgf("Error parsing work id: %v", err)
	// 	return errors.BadRequest(ctx, models.BAD_REQUEST_INVALID_WORK_SELECTION)
	// }

	// paragraphs, err := rec.paragraphRepo.SelectAll(ctx.Request().Context(), int32(workId))
	// if err != nil {
	// 	log.Error().Err(err).Msgf("Error reading paragraphs: %v", err)
	// 	return errors.InternalServerError(ctx)
	// }
	// if len(paragraphs) == 0 {
	// 	return errors.NotFound(ctx, models.NOT_FOUND_PARAGRAPHS)
	// }

	// apiParas := mapper.ParagraphsToApiModels(paragraphs)
	// return ctx.JSON(http.StatusOK, apiParas)
	return nil
}
