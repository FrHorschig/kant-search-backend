package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type DebugHandler interface {
	GetDebugInfo(ctx echo.Context) error
}

type DebugHandlerImpl struct {
}

func NewDebugHandler() DebugHandler {
	handlers := DebugHandlerImpl{}
	return &handlers
}

func (handler *DebugHandlerImpl) GetDebugInfo(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, "Hello World")
}
