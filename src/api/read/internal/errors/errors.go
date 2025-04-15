package errors

import (
	"net/http"

	"github.com/frhorschig/kant-search-api/src/go/models"
	"github.com/labstack/echo/v4"
)

func BadRequest(ctx echo.Context, msg models.ErrorMessage, params ...string) error {
	return ctx.JSON(http.StatusBadRequest, models.HttpError{
		Code:    http.StatusBadRequest,
		Message: msg,
		Params:  params,
	})
}

func NotFound(ctx echo.Context) error {
	return ctx.JSON(http.StatusNotFound, models.HttpError{
		Code:    http.StatusNotFound,
		Message: "",
	})
}

func InternalServerError(ctx echo.Context) error {
	return ctx.JSON(http.StatusInternalServerError, models.HttpError{
		Code:    http.StatusInternalServerError,
		Message: "",
	})
}
