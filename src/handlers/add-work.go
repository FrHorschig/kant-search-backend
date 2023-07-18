package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
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
	return ctx.JSON(http.StatusOK, "Hello World")
}
