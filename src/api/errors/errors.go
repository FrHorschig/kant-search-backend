package errors

import (
	"net/http"

	"github.com/FrHorschig/kant-search-api/models"
	"github.com/labstack/echo/v4"
)

func NotFound(ctx echo.Context, msg string) error {
	return ctx.JSON(http.StatusNotFound, models.HttpError{
		Code:    http.StatusNotFound,
		Message: msg,
	})
}

func BadRequest(ctx echo.Context, msg string) error {
	return ctx.JSON(http.StatusBadRequest, models.HttpError{
		Code:    http.StatusBadRequest,
		Message: msg,
	})
}

func Conflict(ctx echo.Context, msg string) error {
	return ctx.JSON(http.StatusConflict, models.HttpError{
		Code:    http.StatusConflict,
		Message: msg,
	})
}

func InternalServerError(ctx echo.Context) error {
	return ctx.JSON(http.StatusInternalServerError, models.HttpError{
		Code:    http.StatusInternalServerError,
		Message: "Internal Server Error",
	})
}

func NotImplemented(ctx echo.Context, msg string) error {
	return ctx.JSON(http.StatusInternalServerError, models.HttpError{
		Code:    http.StatusNotImplemented,
		Message: msg,
	})
}
