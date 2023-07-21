package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/FrHorschig/kant-search-api/models"
	"github.com/FrHorschig/kant-search-backend/database/repository"
	"github.com/FrHorschig/kant-search-backend/util/errors"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type ParagraphHandler interface {
	GetParagraphs(ctx echo.Context) error
}

type ParagraphHandlerImpl struct {
	paragraphRepo repository.ParagraphRepo
}

func NewParagraphHandler(paragraphRepo repository.ParagraphRepo) ParagraphHandler {
	handlers := ParagraphHandlerImpl{
		paragraphRepo: paragraphRepo,
	}
	return &handlers
}

func (handler *ParagraphHandlerImpl) GetParagraphs(ctx echo.Context) error {
	workId, err := strconv.ParseInt(ctx.Param("id"), 10, 32)
	if err != nil {
		log.Error().Err(err).Msgf("Error parsing work id: %v", err)
		return errors.InternalServerError(ctx)
	}

	start, end, err := findPages(ctx.QueryParam("range"))
	if err != nil {
		log.Error().Err(err).Msgf("Error parsing page range: %v", err)
		return errors.InternalServerError(ctx)
	}

	paragraphs, err := handler.paragraphRepo.SelectRange(ctx.Request().Context(), int32(workId), start, end)
	results := make([]models.Paragraph, 0)
	if err != nil {
		if err == sql.ErrNoRows {
			return ctx.JSON(http.StatusOK, results)
		}
		log.Error().Err(err).Msgf("Error selecting paragraphs: %v", err)
		return errors.InternalServerError(ctx)
	}

	for _, paragraph := range paragraphs {
		results = append(results, models.Paragraph{
			Id:     paragraph.Id,
			Text:   paragraph.Text,
			Pages:  paragraph.Pages,
			WorkId: paragraph.WorkId,
		})
	}

	return ctx.JSON(http.StatusOK, results)
}

func findPages(pageRange string) (start int32, end int32, err error) {
	onError := fmt.Errorf("invalid page range: %s", pageRange)
	parts := strings.Split(pageRange, "-")
	if len(parts) != 2 {
		return -1, -1, onError
	}
	s, err := strconv.ParseInt(parts[0], 10, 32)
	if err != nil {
		return -1, -1, onError
	}
	e, err := strconv.ParseInt(parts[1], 10, 32)
	if err != nil {
		return -1, -1, onError
	}
	return int32(s), int32(e), nil
}
