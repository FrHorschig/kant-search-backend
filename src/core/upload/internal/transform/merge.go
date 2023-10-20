package transform

import (
	"regexp"

	"github.com/FrHorschig/kant-search-backend/common/model"
)

type metadata struct {
	paragraph       *model.Paragraph
	isPageEndPara   bool
	isPageStartPara bool
	isIncomplete    bool
	FnRefs          []string
}

func MergeParagraphs(pars []model.Paragraph) []model.Paragraph {
	md, fnByName := mapToMetadata(pars)
	for i := len(md) - 1; i >= 1; i-- {
		curr := md[i]
		prev := md[i-1]
		if curr.isPageStartPara && prev.isPageEndPara && prev.isIncomplete {
			prev.paragraph.Text += curr.paragraph.Text
			prev.paragraph.Pages = append(prev.paragraph.Pages, curr.paragraph.Pages...)
			prev.FnRefs = append(prev.FnRefs, curr.FnRefs...)
			curr.paragraph = nil
		}
	}

	merged := make([]model.Paragraph, 0)
	for _, m := range md {
		if m.paragraph == nil {
			continue
		}
		merged = append(merged, *m.paragraph)
		for _, fnRef := range m.FnRefs {
			merged = append(merged, *fnByName[fnRef])
		}
	}
	return merged
}

func mapToMetadata(pars []model.Paragraph) ([]metadata, map[string]*model.Paragraph) {
	md := make([]metadata, 0)
	fnByName := make(map[string]*model.Paragraph)
	for i, p := range pars {
		if p.FootnoteName != nil {
			fnByName[*p.FootnoteName] = &p
			continue
		}
		if p.HeadingLevel != nil {
			md = append(md, metadata{paragraph: &p})
			continue
		}
		m := metadata{
			paragraph:    &p,
			isIncomplete: !isEndPunctuation(p.Text[len(p.Text)-1]),
			FnRefs:       findFnRefs(p.Text),
		}
		if i > 0 && p.Pages[0] != pars[i-1].Pages[0] {
			m.isPageStartPara = true
		}
		if i < len(pars)-1 && p.Pages[0] != pars[i+1].Pages[0] {
			m.isPageEndPara = true
		}
		md = append(md, m)
	}
	return md, fnByName
}

func isEndPunctuation(b byte) bool {
	return b == '.' || b == '!' || b == '?'
}

func findFnRefs(text string) []string {
	matches := regexp.MustCompile(`\{fn\d+\.\d+\}`).FindAllString(text, -1)
	return matches
}
