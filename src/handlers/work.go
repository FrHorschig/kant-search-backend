package handlers

import (
	"database/sql"
	"net/http"
	"sort"

	"github.com/FrHorschig/kant-search-api/models"
	"github.com/FrHorschig/kant-search-backend/database/model"
	"github.com/FrHorschig/kant-search-backend/database/repository"
	"github.com/FrHorschig/kant-search-backend/util/errors"
	"github.com/FrHorschig/kant-search-backend/util/textprocessing"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type WorkHandler interface {
	PostWork(ctx echo.Context) error
	GetWork(ctx echo.Context) error
}

type WorkHandlerImpl struct {
	workRepo      repository.WorkRepo
	paragraphRepo repository.ParagraphRepo
	sentenceRepo  repository.SentenceRepo
}

func NewWorkHandler(workRepo repository.WorkRepo, paragraphRepo repository.ParagraphRepo, sentenceRepo repository.SentenceRepo) WorkHandler {
	handlers := WorkHandlerImpl{
		workRepo:      workRepo,
		paragraphRepo: paragraphRepo,
		sentenceRepo:  sentenceRepo,
	}
	return &handlers
}

func (handler *WorkHandlerImpl) PostWork(ctx echo.Context) error {
	work := new(models.Work)
	if err := ctx.Bind(work); err != nil {
		log.Error().Err(err).Msg("Error reading request body")
		return errors.BadRequest(ctx, "Error reading request body")
	}

	workId, err := handler.workRepo.Insert(ctx.Request().Context(), model.Work{Title: work.Title, Abbrev: work.Abbreviation, Volume: work.Volume, Year: work.Year})
	if err != nil {
		log.Error().Err(err).Msg("Error inserting work")
		return errors.InternalServerError(ctx)
	}

	pByNumber, sByPNumber, err := textprocessing.GetParagraphsAndSentences(work.Text)
	if err != nil {
		log.Error().Err(err).Msg("Error processing text")
		return ctx.JSON(http.StatusBadRequest, "Error processing text")
	}

	sortedKeys := make([]int, 0)
	for k := range pByNumber {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Ints(sortedKeys)
	for n := range sortedKeys {
		text := pByNumber[n]
		paragraphId, err := handler.insertParagraph(ctx, text, workId)
		if err != nil {
			log.Error().Err(err).Msg("Error inserting paragraph")
			return errors.InternalServerError(ctx)
		}
		err = handler.insertSentences(ctx, sByPNumber[n], paragraphId, workId)
		if err != nil {
			log.Error().Err(err).Msg("Error inserting sentences")
			return errors.InternalServerError(ctx)
		}
	}

	return ctx.JSON(http.StatusOK, nil)
}

func (handler *WorkHandlerImpl) GetWork(ctx echo.Context) error {
	works, err := handler.workRepo.SelectAll(ctx.Request().Context())
	result := make([]models.WorkMetadata, 0)
	if err != nil {
		if err == sql.ErrNoRows {
			return ctx.JSON(http.StatusOK, result)
		}
		log.Error().Err(err).Msg("Error selecting works")
		return errors.InternalServerError(ctx)
	}

	for _, w := range works {
		result = append(result, models.WorkMetadata{Id: w.Id, Title: w.Title, Abbreviation: w.Abbrev, Volume: w.Volume, Year: w.Year})
	}
	return ctx.JSON(http.StatusOK, result)
}

func (handler *WorkHandlerImpl) insertParagraph(ctx echo.Context, text string, workId int32) (int32, error) {
	pages, err := textprocessing.GetPages(text)
	if err != nil {
		return -1, err
	}
	id, err := handler.paragraphRepo.Insert(ctx.Request().Context(), model.Paragraph{Text: text, Pages: pages, WorkId: workId})
	if err != nil {
		log.Error().Err(err).Msg("Error inserting paragraph")
		return -1, err
	}
	return id, nil
}

func (handler *WorkHandlerImpl) insertSentences(ctx echo.Context, sentences []string, paragraphId int32, workId int32) error {
	sModels := make([]model.Sentence, 0)
	for _, s := range sentences {
		sModels = append(sModels, model.Sentence{Text: s, ParagraphId: paragraphId, WorkId: workId})
	}
	_, err := handler.sentenceRepo.Insert(ctx.Request().Context(), sModels)
	return err
}
