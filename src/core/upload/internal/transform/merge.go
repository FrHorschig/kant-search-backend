package transform

import (
	"regexp"

	"github.com/frhorschig/kant-search-backend/common/model"
)

type metadata struct {
	paragraph       *model.Paragraph
	isPageEndPara   bool
	isPageStartPara bool
	isIncomplete    bool
	fnRefs          []string
}

func MergeParagraphs(pars []model.Paragraph) []model.Paragraph {
	md, fnByName := mapToMetadata(pars)
	for i := len(md) - 1; i >= 1; i-- {
		h := i - 1
		if md[i].isPageStartPara && md[h].isPageEndPara && md[h].isIncomplete {
			md[h].paragraph.Text += " " + md[i].paragraph.Text
			md[h].paragraph.Pages = append(md[h].paragraph.Pages, md[i].paragraph.Pages...)
			md[h].fnRefs = append(md[h].fnRefs, md[i].fnRefs...)
			md[i].paragraph = nil
		}
	}

	merged := make([]model.Paragraph, 0)
	for _, m := range md {
		if m.paragraph == nil {
			continue
		}
		merged = append(merged, *m.paragraph)
		for _, fnRef := range m.fnRefs {
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
			fnByName[*p.FootnoteName] = &pars[i]
			continue
		}
		if p.HeadingLevel != nil {
			md = append(md, metadata{paragraph: &pars[i]})
			continue
		}
		m := metadata{
			paragraph:    &pars[i],
			isIncomplete: !isEndPunctuation(p.Text[len(p.Text)-1]),
			fnRefs:       findFnRefs(p.Text),
		}
		if i > 0 && p.Pages[0] != pars[i-1].Pages[0] || i == 0 {
			m.isPageStartPara = true
		}
		if i < len(pars)-1 && p.Pages[0] != pars[i+1].Pages[0] || i == len(pars)-1 {
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
	matches := regexp.MustCompile(`\{fn(\d+\.\d+)\}`).FindAllStringSubmatch(text, -1)
	fnRefs := make([]string, 0)
	for _, m := range matches {
		fnRefs = append(fnRefs, m[1])
	}
	return fnRefs
}
