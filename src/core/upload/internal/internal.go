package internal

//go:generate mockgen -source=$GOFILE -destination=mocks/text_mapper_mock.go -package=mocks

import (
	"github.com/FrHorschig/kant-search-backend/common/model"
	"github.com/FrHorschig/kant-search-backend/core/errors"
	c "github.com/FrHorschig/kant-search-backend/core/upload/internal/common"
	"github.com/FrHorschig/kant-search-backend/core/upload/internal/parse"
	"github.com/FrHorschig/kant-search-backend/core/upload/internal/pyutil"
	"github.com/FrHorschig/kant-search-backend/core/upload/internal/tokenize"
	"github.com/FrHorschig/kant-search-backend/core/upload/internal/transform"
)

type TextMapper interface {
	Parse(tokens []c.Token) ([]c.Expression, *errors.Error)
	Tokenize(input string) ([]c.Token, *errors.Error)
	Transform(workId int32, exprs []c.Expression) ([]model.Paragraph, *errors.Error)
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

func (rec *textMapperImpl) Parse(tokens []c.Token) ([]c.Expression, *errors.Error) {
	return parse.Parse(tokens)
}

func (rec *textMapperImpl) Tokenize(input string) ([]c.Token, *errors.Error) {
	return tokenize.Tokenize(input)
}

func (rec *textMapperImpl) Transform(workId int32, exprs []c.Expression) ([]model.Paragraph, *errors.Error) {
	return transform.Transform(workId, exprs, rec.pyUtil)
}

func (rec *textMapperImpl) FindSentences(paragraphs []model.Paragraph) ([]model.Sentence, *errors.Error) {
	return transform.FindSentences(paragraphs, rec.pyUtil)
}
