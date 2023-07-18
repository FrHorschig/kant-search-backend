package handlers

import (
	"io/ioutil"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type AddWorkHandler interface {
	PostWork(ctx echo.Context) error
}

type AddWorkHandlerImpl struct {
}

func NewAddWorkHandler() AddWorkHandler {
	handlers := AddWorkHandlerImpl{}
	return &handlers
}

func (handler *AddWorkHandlerImpl) PostWork(ctx echo.Context) error {
	body, err := ioutil.ReadAll(ctx.Request().Body)
	if err != nil {
		log.Error().Err(err).Msg("Error reading request body")
		return ctx.JSON(http.StatusBadRequest, "Error reading request body")
	}
	log.Info().Msg(string(body))
	return ctx.JSON(http.StatusOK, "Hello World")
}
