package upload

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/frhorschig/kant-search-backend/api/upload/errors"
	"github.com/frhorschig/kant-search-backend/core/upload"
	"github.com/frhorschig/kant-search-backend/core/upload/model/abt1"
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
		return errors.BadRequest(ctx, msg)
	}
	body, err := io.ReadAll(ctx.Request().Body)
	if err != nil {
		msg := fmt.Sprintf("Error reading request body: %v", err.Error())
		log.Error().Err(err).Msg(msg)
		return errors.BadRequest(ctx, msg)
	}

	switch {
	case volNum >= 1 && volNum <= 9:
		var vol abt1.Kantabt1
		err := xml.Unmarshal(body, &vol)
		if err != nil {
			msg := fmt.Sprintf("Error unmarshaling request body: %v", err.Error())
			log.Error().Err(err).Msg(msg)
			return errors.BadRequest(ctx, msg)
		}
		if err := rec.volumeProcessor.ProcessAbt1(ctx.Request().Context(), int32(volNum), vol); err != nil {
			msg := fmt.Sprintf("Error processing XML data for volume %d: %v", volNum, err.Error())
			log.Error().Err(err).Msg(msg)
			return errors.BadRequest(ctx, msg)
		}
	default:
		msg := "Uploading volumes greater than 9 is not yet implemented"
		log.Error().Msg(msg)
		return errors.BadRequest(ctx, msg)
	}

	return ctx.NoContent(http.StatusCreated)
}
