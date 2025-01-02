package upload

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/frhorschig/kant-search-backend/core/upload"
	"github.com/frhorschig/kant-search-backend/core/upload/errors"
	"github.com/frhorschig/kant-search-backend/core/upload/model/abt1"
	"github.com/frhorschig/kant-search-backend/core/upload/model/abt2"
	"github.com/frhorschig/kant-search-backend/core/upload/model/abt31"
	"github.com/frhorschig/kant-search-backend/core/upload/model/abt32"
	"github.com/frhorschig/kant-search-backend/core/upload/model/vol14"
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
		var vol abt1.Band
		err = xml.Unmarshal(body, &vol)
		if err != nil {
			msg := fmt.Sprintf("Error reading request body: %v", err.Error())
			log.Error().Err(err).Msg(msg)
			return errors.BadRequest(ctx, msg)
		}
		err = rec.volumeProcessor.ProcessAbt1(ctx.Request().Context(), int32(volNum), vol)
		if err != nil {
			msg := fmt.Sprintf("Error processing xml data: %v", err.Error())
			log.Error().Err(err).Msg(msg)
			return errors.BadRequest(ctx, msg)
		}

	case volNum >= 10 && volNum <= 13:
		var vol abt2.Band
		err = xml.Unmarshal(body, &vol)
		if err != nil {
			msg := fmt.Sprintf("Error reading request body: %v", err.Error())
			log.Error().Err(err).Msg(msg)
			return errors.BadRequest(ctx, msg)
		}
		err = rec.volumeProcessor.ProcessAbt2(ctx.Request().Context(), int32(volNum), vol)
		if err != nil {
			msg := fmt.Sprintf("Error processing xml data: %v", err.Error())
			log.Error().Err(err).Msg(msg)
			return errors.BadRequest(ctx, msg)
		}

	case volNum == 14:
		var vol vol14.Band
		err = xml.Unmarshal(body, &vol)
		if err != nil {
			msg := fmt.Sprintf("Error reading request body: %v", err.Error())
			log.Error().Err(err).Msg(msg)
			return errors.BadRequest(ctx, msg)
		}
		err = rec.volumeProcessor.ProcessVol14(ctx.Request().Context(), int32(volNum), vol)
		if err != nil {
			msg := fmt.Sprintf("Error processing xml data: %v", err.Error())
			log.Error().Err(err).Msg(msg)
			return errors.BadRequest(ctx, msg)
		}

	case volNum >= 15 && volNum <= 19:
		var vol abt31.Band
		err = xml.Unmarshal(body, &vol)
		if err != nil {
			msg := fmt.Sprintf("Error reading request body: %v", err.Error())
			log.Error().Err(err).Msg(msg)
			return errors.BadRequest(ctx, msg)
		}
		err = rec.volumeProcessor.ProcessAbt31(ctx.Request().Context(), int32(volNum), vol)
		if err != nil {
			msg := fmt.Sprintf("Error processing xml data: %v", err.Error())
			log.Error().Err(err).Msg(msg)
			return errors.BadRequest(ctx, msg)
		}

	case volNum >= 20 && volNum <= 23:
		var vol abt32.Band
		err = xml.Unmarshal(body, &vol)
		if err != nil {
			msg := fmt.Sprintf("Error reading request body: %v", err.Error())
			log.Error().Err(err).Msg(msg)
			return errors.BadRequest(ctx, msg)
		}
		err = rec.volumeProcessor.ProcessAbt32(ctx.Request().Context(), int32(volNum), vol)
		if err != nil {
			msg := fmt.Sprintf("Error processing xml data: %v", err.Error())
			log.Error().Err(err).Msg(msg)
			return errors.BadRequest(ctx, msg)
		}

	default:
		msg := "The volume number must be a number from 1 to 23"
		log.Error().Msg(msg)
		return errors.BadRequest(ctx, msg)
	}

	return ctx.NoContent(http.StatusCreated)
}
