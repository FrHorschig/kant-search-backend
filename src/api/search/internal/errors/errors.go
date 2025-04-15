package errors

import (
	"fmt"
	"net/http"

	"github.com/frhorschig/kant-search-api/src/go/models"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type ErrMsg string

const (
	UnexpectedToken         ErrMsg = "UNEXPECTED_TOKEN"
	WrongStartingChar       ErrMsg = "WRONG_STARTING_CHAR"
	WrongEndingChar         ErrMsg = "WRONG_ENDING_CHAR"
	UnexpectedEndOfInput    ErrMsg = "UNEXPECTED_END_OF_INPUT"
	MissingCloseParenthesis ErrMsg = "MISSING_CLOSING_PARENTHESIS"
	UnterminatedDoubleQuote ErrMsg = "UNTERMINATED_DOUBLE_QUOTE"
)

type ValidationError struct {
	Msg    ErrMsg
	Params []string
}

func ValidationErrorToApiError(ctx echo.Context, err *ValidationError) error {
	msg, e := mapValidationEnum(err.Msg)
	if e != nil {
		log.Error().Err(e).Msgf("error mapping validation error: %v", err)
		return InternalServerError(ctx)
	}
	return ctx.JSON(http.StatusBadRequest, models.HttpError{
		Code:    http.StatusBadRequest,
		Message: msg,
		Params:  err.Params,
	})
}

func BadRequest(ctx echo.Context, msg models.ErrorMessage) error {
	return ctx.JSON(http.StatusBadRequest, models.HttpError{
		Code:    http.StatusBadRequest,
		Message: msg,
	})
}

func InternalServerError(ctx echo.Context) error {
	return ctx.JSON(http.StatusInternalServerError, models.HttpError{
		Code:    http.StatusInternalServerError,
		Message: "",
	})
}

func mapValidationEnum(err ErrMsg) (models.ErrorMessage, error) {
	switch err {
	case UnexpectedToken:
		return models.BAD_REQUEST_VALIDATION_UNEXPECTED_TOKEN, nil
	case WrongStartingChar:
		return models.BAD_REQUEST_VALIDATION_WRONG_STARTING_CHAR, nil
	case WrongEndingChar:
		return models.BAD_REQUEST_VALIDATION_WRONG_ENDING_CHAR, nil
	case UnexpectedEndOfInput:
		return models.BAD_REQUEST_VALIDATION_UNEXPECTED_END_OF_INPUT, nil
	case MissingCloseParenthesis:
		return models.BAD_REQUEST_VALIDATION_MISSING_CLOSING_PARENTHESIS, nil
	case UnterminatedDoubleQuote:
		return models.BAD_REQUEST_VALIDATION_UNTERMINATED_DOUBLE_QUOTE, nil
	}
	return "", fmt.Errorf("unknown enum \"%s\"", err)
}
