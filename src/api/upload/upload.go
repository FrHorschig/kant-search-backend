package upload

import (
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/beevik/etree"
	"github.com/frhorschig/kant-search-backend/api/upload/errors"
	"github.com/frhorschig/kant-search-backend/core/upload"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type UploadHandler interface {
	PostVolume(ctx echo.Context) error
}

type uploadHandlerImpl struct {
	volumeProcessor upload.VolumeUploadProcessor
}

func NewUploadHandler(volumeProcessor upload.VolumeUploadProcessor) UploadHandler {
	return &uploadHandlerImpl{
		volumeProcessor: volumeProcessor,
	}
}

func (rec *uploadHandlerImpl) PostVolume(ctx echo.Context) error {
	volNum, err := strconv.ParseInt(ctx.Param("volumeNumber"), 10, 32)
	if err != nil {
		msg := fmt.Sprintf("Error parsing volume number: %v", err.Error())
		log.Error().Err(err).Msg(msg)
		return errors.JsonError(ctx, http.StatusBadRequest, msg)
	}
	body, err := io.ReadAll(ctx.Request().Body)
	if err != nil {
		msg := fmt.Sprintf("Error reading request body: %v", err.Error())
		log.Error().Err(err).Msg(msg)
		return errors.JsonError(ctx, http.StatusBadRequest, msg)
	}

	if volNum < 1 {
		msg := "The volume number must be between 1 and 9"
		log.Error().Msg(msg)
		return errors.JsonError(ctx, http.StatusBadRequest, msg)
	} else if volNum > 9 {
		msg := "Uploading volumes greater than 9 is not yet implemented"
		log.Error().Msg(msg)
		return errors.JsonError(ctx, http.StatusNotImplemented, msg)
	}

	doc := etree.NewDocument()
	if err := doc.ReadFromBytes(body); err != nil {
		msg := fmt.Sprintf("Error unmarshaling request body: %v", err.Error())
		log.Error().Err(err).Msg(msg)
		return errors.JsonError(ctx, http.StatusBadRequest, msg)
	}

	if err := rec.volumeProcessor.Process(ctx.Request().Context(), int32(volNum), doc); err != nil {
		msg := fmt.Sprintf("Error processing XML data for volume %d: %v", volNum, err.Error())
		log.Error().Err(err).Msg(msg)
		return errors.JsonError(ctx, http.StatusBadRequest, msg)
	}

	return ctx.NoContent(http.StatusCreated)
}
