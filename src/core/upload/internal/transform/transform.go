package transform

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/frhorschig/kant-search-backend/common/errors"
	"github.com/frhorschig/kant-search-backend/common/model"
	c "github.com/frhorschig/kant-search-backend/core/upload/internal/common"
)

func Transform(
	workId int32,
	exprs []c.Expression,
) ([]model.Paragraph, *errors.Error) {
	err := validateStartEnd(exprs)
	if err != nil {
		return nil, err
	}
	return buildParagraphs(workId, exprs)
}

func validateStartEnd(exprs []c.Expression) *errors.Error {
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
	exprs []c.Expression,
) ([]model.Paragraph, *errors.Error) {
	pars := make([]model.Paragraph, 0)
	var pn int32
	for i, e := range exprs {
		if e.Metadata.Class == "p" {
			pn = findPageNumber(e)
			if exprs[i+1].Content != nil {
				exprs[i+1].Content = &[]string{fmt.Sprintf("{p%d} %s", pn, *exprs[i+1].Content)}[0]
			} else {
				e.Metadata.Class = "paragraph"
				e.Content = &[]string{fmt.Sprintf("{p%d} <i>Inhalt nicht verfügbar</i>.", pn)}[0]
				par, err := createParagraph(workId, pn, e)
				if err != nil {
					return nil, err
				}
				pars = append(pars, par)
			}
		} else {
			par, err := createParagraph(workId, pn, e)
			if err != nil {
				return nil, err
			}
			pars = append(pars, par)
		}
	}
	return pars, nil
}

func findPageNumber(e c.Expression) int32 {
	// TODO frhorsch: here we "just know" that param is a number, fix when improving EBNF spec
	pn, _ := strconv.Atoi(*e.Metadata.Param)
	return int32(pn)
}

func createParagraph(
	workId int32,
	pn int32,
	e c.Expression,
) (model.Paragraph, *errors.Error) {
	par := model.Paragraph{
		// TODO frhorsch: here we "just know" that content exists, fix when improving EBNF spec
		Text:   strings.TrimSpace(*e.Content),
		Pages:  []int32{pn},
		WorkId: workId,
	}
	switch e.Metadata.Class {
	case "paragraph":
		// nothing to do
	case "heading":
		hl, _ := strconv.Atoi(*e.Metadata.Param)
		par.HeadingLevel = &[]int32{int32(hl)}[0]
	case "footnote":
		par.FootnoteName = e.Metadata.Param
	default:
		return model.Paragraph{}, &errors.Error{
			Msg:    errors.UNKNOWN_EXPRESSION_CLASS,
			Params: []string{e.Metadata.Class},
		}
	}
	return par, nil
}
