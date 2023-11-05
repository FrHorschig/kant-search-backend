package transform

import (
	"sort"

	"github.com/frhorschig/kant-search-backend/common/model"
	"github.com/frhorschig/kant-search-backend/core/errors"
	"github.com/frhorschig/kant-search-backend/core/upload/internal/pyutil"
)

type ByParagraphId []model.Sentence

func (a ByParagraphId) Len() int           { return len(a) }
func (a ByParagraphId) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByParagraphId) Less(i, j int) bool { return a[i].ParagraphId < a[j].ParagraphId }

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
	sort.Sort(ByParagraphId(sentenceModels))
	return sentenceModels
}
