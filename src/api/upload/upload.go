package upload

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

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
	volumeProcessor upload.UploadProcessor
}

func NewUploadHandler(volumeProcessor upload.UploadProcessor) UploadHandler {
	return &uploadHandlerImpl{
		volumeProcessor: volumeProcessor,
	}
}

func (rec *uploadHandlerImpl) PostVolume(ctx echo.Context) error {
	body, err := io.ReadAll(ctx.Request().Body)
	if err != nil {
		msg := fmt.Sprintf("error reading request body: %v", err.Error())
		log.Error().Msg(msg)
		return errors.JsonError(ctx, http.StatusBadRequest, msg)
	}
	doc := etree.NewDocument()
	if err := doc.ReadFromBytes(body); err != nil {
		msg := fmt.Sprintf("error unmarshaling request body: %v", err.Error())
		log.Error().Msg(msg)
		return errors.JsonError(ctx, http.StatusBadRequest, msg)
	}

	band := doc.FindElement("//band")
	if band == nil || band.SelectAttr("nr") == nil {
		msg := "missing element 'band' with attribute 'nr'"
		log.Error().Msg(msg)
		return errors.JsonError(ctx, http.StatusBadRequest, msg)
	}
	nrStr := strings.TrimLeft(band.SelectAttr("nr").Value, "0")
	if nrStr == "" {
		msg := "the volume number is 0, but it must be between 1 and 9"
		log.Error().Msg(msg)
		return errors.JsonError(ctx, http.StatusBadRequest, msg)
	}
	volNum, err := strconv.Atoi(nrStr)
	if err != nil {
		msg := fmt.Sprintf("attribute 'nr' of element 'band' can't be converted to a number: %v", err.Error())
		log.Error().Msg(msg)
		return errors.JsonError(ctx, http.StatusBadRequest, msg)
	}
	if volNum < 1 {
		msg := fmt.Sprintf("the volume number is %d, but it must be between 1 and 9", volNum)
		log.Error().Msg(msg)
		return errors.JsonError(ctx, http.StatusBadRequest, msg)
	} else if volNum > 9 {
		msg := "uploading volumes greater than 9 is not yet implemented"
		log.Error().Msg(msg)
		return errors.JsonError(ctx, http.StatusNotImplemented, msg)
	}

	if err := rec.volumeProcessor.Process(ctx.Request().Context(), doc); err != nil {
		msg := fmt.Sprintf("error processing XML data for volume %d", volNum)
		log.Error().Err(err).Msg(msg)
		return errors.JsonError(ctx, http.StatusInternalServerError, msg)
	}

	return ctx.NoContent(http.StatusCreated)
}
