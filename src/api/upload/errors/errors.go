package errors

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type HttpError struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
}

func BadRequest(ctx echo.Context, msg string) error {
	return ctx.JSON(http.StatusBadRequest, HttpError{
		Code:    http.StatusBadRequest,
		Message: msg,
	})
}

func UploadError(ctx echo.Context, msg string) error {
	return ctx.JSON(http.StatusBadRequest, HttpError{
		Code:    http.StatusInternalServerError,
		Message: msg,
	})
}
