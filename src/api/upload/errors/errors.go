package errors

import (
	"github.com/labstack/echo/v4"
)

type HttpError struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
}

func JsonError(ctx echo.Context, code int, msg string) error {
	return ctx.JSON(code, HttpError{
		Code:    int32(code),
		Message: msg,
	})
}
