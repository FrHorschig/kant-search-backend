package handlers

import (
	"net/http"
	"strconv"

	"github.com/FrHorschig/kant-search-backend/database/repository"
	"github.com/FrHorschig/kant-search-backend/util/errors"
	"github.com/FrHorschig/kant-search-backend/util/mapper"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type TextHandler interface {
	GetTextById(ctx echo.Context) error
}

type TextHandlerImpl struct {
	textRepo repository.TextRepo
}

func NewTextHandler(textRepo repository.TextRepo) TextHandler {
	handlers := TextHandlerImpl{
		textRepo: textRepo,
	}
	return &handlers
}

func (handler *TextHandlerImpl) GetTextById(ctx echo.Context) error {
	id := ctx.Param("id")
	id_int, err := strconv.Atoi(id)
	if err != nil {
		log.Error().Err(err).Msg("Error parsing id")
		return errors.BadRequest(ctx, "Error parsing id")
	}
	text, err := handler.textRepo.Select(ctx.Request().Context(), int32(id_int))
	if err != nil {
		log.Error().Err(err).Msg("Error selecting text")
		return errors.InternalServerError(ctx)
	}

	mapped := mapper.MapTextFromDb(text)
	return ctx.JSON(http.StatusOK, mapped)
}
