package errors

import (
	"net/http"

	"github.com/frhorschig/kant-search-api/generated/go/models"
	"github.com/frhorschig/kant-search-backend/common/errors"
	"github.com/labstack/echo/v4"
)

func NotFound(ctx echo.Context, err models.ErrorMessage) error {
	return ctx.JSON(http.StatusNotFound, models.HttpError{
		Code:    http.StatusNotFound,
		Message: err,
	})
}

func CoreError(ctx echo.Context, err *errors.Error) error {
	code := http.StatusBadRequest
	if err.Msg == errors.GO_ERR {
		code = http.StatusInternalServerError
	}
	return ctx.JSON(code, models.HttpError{
		Code:    int32(code),
		Message: mapCoreEnum(err.Msg),
		Params:  err.Params,
	})
}

func UploadError(ctx echo.Context, err *errors.Error) error {
	code := http.StatusBadRequest
	if err.Msg == errors.GO_ERR {
		code = http.StatusInternalServerError
	}
	return ctx.JSON(code, mapUploadError(err.Msg))
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
	case errors.GO_ERR:
		return models.INTERNAL_SERVER_ERROR

	// search term parsing
	case errors.UNEXPECTED_TOKEN:
		return models.BAD_REQUEST_SYNTAX_UNEXPECTED_TOKEN
	case errors.WRONG_STARTING_CHAR:
		return models.BAD_REQUEST_SYNTAX_WRONG_STARTING_CHAR
	case errors.WRONG_ENDING_CHAR:
		return models.BAD_REQUEST_SYNTAX_WRONG_ENDING_CHAR
	case errors.UNEXPECTED_END_OF_INPUT:
		return models.BAD_REQUEST_SYNTAX_UNEXPECTED_END_OF_INPUT
	case errors.MISSING_CLOSING_PARENTHESIS:
		return models.BAD_REQUEST_SYNTAX_MISSING_CLOSING_PARENTHESIS
	case errors.UNTERMINATED_DOUBLE_QUOTE:
		return models.BAD_REQUEST_SYNTAX_UNTERMINATED_DOUBLE_QUOTE

	default:
		return models.BAD_REQUEST_GENERIC
	}
}

func mapUploadError(err errors.ErrMsg) string {
	switch err {
	case errors.UPLOAD_GO_ERR:
		return "INTERNAL_SERVER_ERROR"
	case errors.MISSING_EXPR_TYPE:
		return "BAD_REQUEST.MISSING_EXPR_TYPE"
	case errors.MISSING_CLOSING_BRACE:
		return "BAD_REQUEST.MISSING_CLOSING_BRACE"
	case errors.UPLOAD_WRONG_STARTING_CHAR:
		return "BAD_REQUEST.WRONG_STARTING_CHAR"
	case errors.WRONG_START_EXPRESSION:
		return "BAD_REQUEST.WRONG_START_EXPRESSION"
	case errors.WRONG_END_EXPRESSION:
		return "BAD_REQUEST.WRONG_END_EXPRESSION"
	case errors.UNKNOWN_EXPRESSION_CLASS:
		return "BAD_UNKNOWN.EXPRESSION_CLASS"
	default:
		return "BAD_REQUEST.GENERIC"
	}
}
