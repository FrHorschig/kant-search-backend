package errors

import (
	"net/http"

	"github.com/FrHorschig/kant-search-api/models"
	"github.com/FrHorschig/kant-search-backend/core/errors"
	"github.com/labstack/echo/v4"
)

func NotFound(ctx echo.Context, err models.ErrorMessage) error {
	return ctx.JSON(http.StatusNotFound, models.HttpError{
		Code:    http.StatusNotFound,
		Message: err,
	})
}

func BadRequestFromCore(ctx echo.Context, err *errors.Error) error {
	return ctx.JSON(http.StatusBadRequest, models.HttpError{
		Code:    http.StatusBadRequest,
		Message: mapCoreEnum(err.Msg),
		Args:    err.Args,
	})
}

func BadRequest(ctx echo.Context, err models.ErrorMessage) error {
	return ctx.JSON(http.StatusBadRequest, models.HttpError{
		Code:    http.StatusBadRequest,
		Message: err,
	})
}

func InternalServerError(ctx echo.Context) error {
	return ctx.JSON(http.StatusInternalServerError, models.HttpError{
		Code:    http.StatusInternalServerError,
		Message: models.INTERNAL_SERVER_ERROR,
	})
}

func NotImplemented(ctx echo.Context, msg string) error {
	return ctx.JSON(http.StatusInternalServerError, models.HttpError{
		Code:    http.StatusNotImplemented,
		Message: models.NOT_IMPLEMENTED,
	})
}

func mapCoreEnum(err errors.ErrMsg) models.ErrorMessage {
	switch err {
	case errors.UNEXPECTED_TOKEN:
		return models.BAD_REQUEST_SYNTAX_UNEXPECTED_TOKEN
	default:
		return models.BAD_REQUEST_GENERIC
	}
}
