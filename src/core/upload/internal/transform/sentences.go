package transform

import (
	"github.com/FrHorschig/kant-search-backend/common/model"
	"github.com/FrHorschig/kant-search-backend/core/errors"
	"github.com/FrHorschig/kant-search-backend/core/upload/internal/pyutil"
)

func FindSentences(paragraphs []model.Paragraph, pyUtil pyutil.PythonUtil) ([]model.Sentence, *errors.Error) {
	sentencesByParagraphId, err := pyUtil.SplitIntoSentences(paragraphs)
	if err != nil {
		return nil, &errors.Error{
			Msg:    errors.GO_ERR,
			Params: []string{err.Error()},
		}
	}
	return createSentenceModels(sentencesByParagraphId), nil
}

func createSentenceModels(sentencesByParagraphId map[int32][]string) []model.Sentence {
	sentenceModels := make([]model.Sentence, 0)
	for pId, sentences := range sentencesByParagraphId {
		for _, s := range sentences {
			sentenceModels = append(sentenceModels, model.Sentence{
				ParagraphId: pId,
				Text:        s,
			})
		}
	}
	return sentenceModels
}
