package handlers

import (
	"io"
	"net/http"

	"github.com/FrHorschig/kant-search-backend/util/textprocessing"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type AddWorkHandler interface {
	PostWork(ctx echo.Context) error
}

type AddWorkHandlerImpl struct {
}

func NewAddWorkHandler() AddWorkHandler {
	handlers := AddWorkHandlerImpl{}
	return &handlers
}

func (handler *AddWorkHandlerImpl) PostWork(ctx echo.Context) error {
	body, err := io.ReadAll(ctx.Request().Body)
	if err != nil {
		log.Error().Err(err).Msg("Error reading request body")
		return ctx.JSON(http.StatusBadRequest, "Error reading request body")
	}

	paragraphByNumber := textprocessing.SplitIntoParagraphs(string(body))
	sentencesByNumber := make(map[int]string)
	sentenceNumbersByParagraphNumber := make(map[int][]int)
	sentenceNumber := 0
	for n, p := range paragraphByNumber {
		sentences, err := textprocessing.SplitIntoSentences(p)
		if err != nil {
			log.Error().Err(err).Msgf("Unable to split paragraph into sentences: %v", err)
			return ctx.JSON(http.StatusBadRequest, "Unable to split paragraph into sentences")
		}
		for _, s := range sentences {
			sentencesByNumber[sentenceNumber] = s
			sentenceNumbersByParagraphNumber[n] = append(sentenceNumbersByParagraphNumber[n], sentenceNumber)
			sentenceNumber++
		}
	}

	// TODO save sentences to database
	log.Info().Msgf("Number of paragraphs: %v", len(paragraphByNumber))
	log.Info().Msgf("Number of sentences: %v", len(sentencesByNumber))

	return ctx.JSON(http.StatusOK, "Hello World")
}
