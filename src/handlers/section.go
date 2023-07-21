package handlers

import (
	"github.com/FrHorschig/kant-search-backend/database/repository"
	"github.com/labstack/echo/v4"
)

type SectionHandler interface {
	GetSection(ctx echo.Context) error
}

type SectionHandlerImpl struct {
	paragraphRepo repository.ParagraphRepo
}

func NewSectionHandler(paragraphRepo repository.ParagraphRepo) SectionHandler {
	handlers := SectionHandlerImpl{
		paragraphRepo: paragraphRepo,
	}
	return &handlers
}

func (handler *SectionHandlerImpl) GetSection(ctx echo.Context) error {
	// TODO implement me
	return nil
}
