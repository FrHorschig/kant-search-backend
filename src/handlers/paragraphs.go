package handlers

import (
	"github.com/FrHorschig/kant-search-backend/database/repository"
	"github.com/labstack/echo/v4"
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
	// TODO implement me
	return nil
}
