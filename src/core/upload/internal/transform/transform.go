package transform

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/FrHorschig/kant-search-backend/common/model"
	"github.com/FrHorschig/kant-search-backend/core/errors"
	"github.com/FrHorschig/kant-search-backend/core/upload/internal/parse"
	"github.com/FrHorschig/kant-search-backend/core/upload/internal/pyutils"
)

func Transform(workId int32, exprs []parse.Expression) ([]model.Paragraph, *errors.Error) {
	err := validateStartEnd(exprs)
	if err != nil {
		return nil, err
	}
	pars, boundaryIndices, err := buildParagraphs(workId, exprs)
	if err != nil {
		return nil, err
	}
	return mergePartialParagraphs(pars, boundaryIndices)
}

func validateStartEnd(exprs []parse.Expression) *errors.Error {
	firstClass := exprs[0].Metadata.Class
	if firstClass != "p" {
		return &errors.Error{
			Msg:    errors.WRONG_START_EXPRESSION,
			Params: []string{string(firstClass)},
		}
	}

	lastClass := exprs[len(exprs)-1].Metadata.Class
	if lastClass == "p" || lastClass == "l" {
		return &errors.Error{
			Msg:    errors.WRONG_END_EXPRESSION,
			Params: []string{string(lastClass)},
		}
	}
	return nil
}

func buildParagraphs(
	workId int32,
	exprs []parse.Expression,
) ([]*model.Paragraph, [][2]int, *errors.Error) {
	pars := make([]*model.Paragraph, 0)
	boundaryIndices := make([][2]int, 0)
	var pn int32
	for i, e := range exprs {
		if e.Metadata.Class == "p" {
			if i > 0 {
				// we know from validation that this is not the last expression
				if isParagraph(exprs[i-1]) && isParagraph(exprs[i+1]) {
					boundaryIndices = append(boundaryIndices, [2]int{len(pars) - 1, len(pars)})
				}
			}
			pn = findPageNumber(e)
		} else {
			par, err := createParagraph(workId, pn, e)
			if err != nil {
				return nil, nil, err
			}
			pars = append(pars, &par)
		}
	}
	return pars, boundaryIndices, nil
}

func isParagraph(e parse.Expression) bool {
	return e.Metadata.Class == "paragraph"
}

func findPageNumber(e parse.Expression) int32 {
	// TODO frhorsch: here we "just know" that param is a number, fix when improving EBNF spec
	pn, _ := strconv.Atoi(*e.Metadata.Param)
	return int32(pn)
}

func createParagraph(
	workId int32,
	pn int32,
	e parse.Expression,
) (model.Paragraph, *errors.Error) {
	par := model.Paragraph{
		// TODO frhorsch: here we "just know" that content exists, fix when improving EBNF spec
		Text:   fmt.Sprintf("{p%d} %s", pn, *e.Content),
		Pages:  []int32{pn},
		WorkId: workId,
	}
	hl, _ := strconv.Atoi(*e.Metadata.Param)
	switch e.Metadata.Class {
	case "paragraph":
		// nothing to do
	case "heading":
		par.HeadingLevel = int32(hl)
	case "footnote":
		par.FootnoteName = *e.Metadata.Param
	default:
		return model.Paragraph{}, &errors.Error{
			Msg:    errors.UNKNOWN_EXPRESSION_CLASS,
			Params: []string{e.Metadata.Class},
		}
	}
	return par, nil
}

// We merge paragraphs that end with incomplete sentences by merging them with the next paragraph, splitting them into sentences and checking if sentence splitting point is the paragraph boundary.
// Improvement: maybe get a list of pages that start with a new paragraph?
func mergePartialParagraphs(
	pars []*model.Paragraph,
	boundaryIndices [][2]int,
) ([]model.Paragraph, *errors.Error) {
	merged := make([]model.Paragraph, len(boundaryIndices))
	for i, b := range boundaryIndices {
		merged[i] = model.Paragraph{
			Id:   int32(b[1]),
			Text: pars[b[0]].Text + pars[b[1]].Text,
		}
	}

	sentencesByPageStartIndex, err := pyutils.SplitIntoSentences(merged)
	if err != nil {
		return nil, &errors.Error{
			Msg:    errors.GO_ERR,
			Params: []string{err.Error()},
		}
	}

	return mergeParagraphs(pars, sentencesByPageStartIndex), nil
}

func mergeParagraphs(pars []*model.Paragraph, sentencesByPageStartIndex map[int32][]string) []model.Paragraph {
	finalMerged := make([]model.Paragraph, 0)
	for i, p := range pars {
		sentences, isPageStart := sentencesByPageStartIndex[int32(i)]
		if !isPageStart || startsWithCompleteSentence(sentences, p) {
			finalMerged = append(finalMerged, *p)
		} else {
			finalMerged[len(finalMerged)-1].Text += p.Text
			finalMerged[len(finalMerged)-1].Pages = append(finalMerged[len(finalMerged)-1].Pages, p.Pages...)
		}
	}
	return finalMerged
}

func startsWithCompleteSentence(sentences []string, p *model.Paragraph) bool {
	for _, s := range sentences {
		if strings.HasPrefix(p.Text, s) {
			return true
		}
	}
	return false
}
