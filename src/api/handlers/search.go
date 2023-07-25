package handlers

import (
	"github.com/FrHorschig/kant-search-backend/core/search"
	"github.com/labstack/echo/v4"
)

type SearchHandler interface {
	SearchParagraphs(ctx echo.Context) error
}

type SearchHandlerImpl struct {
	paragraphSearcher search.ParagraphSearcher
}

func NewSearchHandler(paragraphSearcher search.ParagraphSearcher) SearchHandler {
	impl := SearchHandlerImpl{paragraphSearcher: paragraphSearcher}
	return &impl
}

func (rec *SearchHandlerImpl) SearchParagraphs(ctx echo.Context) error {
	// TODO implement me
	return nil
}
