package errors

import (
	"fmt"
	"net/http"

	"github.com/frhorschig/kant-search-api/generated/go/models"
	"github.com/frhorschig/kant-search-backend/core/search/errors"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

func SyntaxErrorToApiError(ctx echo.Context, err *errors.SyntaxError) error {
	msg, e := mapSyntaxEnum(err.Msg)
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

func mapSyntaxEnum(err errors.ErrMsg) (models.ErrorMessage, error) {
	switch err {
	case errors.UnexpectedToken:
		return models.BAD_REQUEST_SYNTAX_UNEXPECTED_TOKEN, nil
	case errors.WrongStartingChar:
		return models.BAD_REQUEST_SYNTAX_WRONG_STARTING_CHAR, nil
	case errors.WrongEndingChar:
		return models.BAD_REQUEST_SYNTAX_WRONG_ENDING_CHAR, nil
	case errors.UnexpectedEndOfInput:
		return models.BAD_REQUEST_SYNTAX_UNEXPECTED_END_OF_INPUT, nil
	case errors.MissingCloseParenthesis:
		return models.BAD_REQUEST_SYNTAX_MISSING_CLOSING_PARENTHESIS, nil
	case errors.UnterminatedDoubleQuote:
		return models.BAD_REQUEST_SYNTAX_UNTERMINATED_DOUBLE_QUOTE, nil
	}
	return "", fmt.Errorf("unknown enum \"%s\"", err)
}
