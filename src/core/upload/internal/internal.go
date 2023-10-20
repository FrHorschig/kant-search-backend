package internal

//go:generate mockgen -source=$GOFILE -destination=mocks/text_mapper_mock.go -package=mocks

import (
	"github.com/FrHorschig/kant-search-backend/common/model"
	"github.com/FrHorschig/kant-search-backend/core/errors"
	"github.com/FrHorschig/kant-search-backend/core/upload/internal/parse"
	"github.com/FrHorschig/kant-search-backend/core/upload/internal/pyutil"
	"github.com/FrHorschig/kant-search-backend/core/upload/internal/tokenize"
	"github.com/FrHorschig/kant-search-backend/core/upload/internal/transform"
)

type TextMapper interface {
	FindParagraphs(text string, workId int32) ([]model.Paragraph, *errors.Error)
	FindSentences(paragraphs []model.Paragraph) ([]model.Sentence, *errors.Error)
}

type textMapperImpl struct {
	pyUtil pyutil.PythonUtil
}

func NewTextMapper() TextMapper {
	impl := textMapperImpl{
		pyUtil: pyutil.NewPythonUtil(),
	}
	return &impl
}

func (rec *textMapperImpl) FindParagraphs(text string, workId int32) ([]model.Paragraph, *errors.Error) {
	tokens, err := tokenize.Tokenize(text)
	if err != nil {
		return nil, err
	}
	exprs, err := parse.Parse(tokens)
	if err != nil {
		return nil, err
	}
	pars, err := transform.Transform(workId, exprs)
	if err != nil {
		return nil, err
	}
	return transform.MergeParagraphs(pars), nil
}

func (rec *textMapperImpl) FindSentences(paragraphs []model.Paragraph) ([]model.Sentence, *errors.Error) {
	return transform.FindSentences(paragraphs, rec.pyUtil)
}
