package referencemapping

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/frhorschig/kant-search-backend/common/errs"
	"github.com/frhorschig/kant-search-backend/core/upload/internal/common/model"
	"github.com/frhorschig/kant-search-backend/core/upload/internal/common/util"
	"github.com/rs/zerolog/log"
)

func MapReferences(works []model.Work, footnotes []model.Footnote, summaries []model.Summary) errs.UploadError {
	err := addFnsToWorks(works, footnotes)
	if err.HasError {
		return err
	}
	leftovers := findLeftoverFns(works, footnotes)
	if len(leftovers) > 0 {
		return errs.New(fmt.Errorf("unable to assign the following footnotes ('page.number') to a work: %v", leftovers), nil)
	}

	if len(summaries) > 0 {
		err = addSummariesToWorks(works, summaries)
		if err.HasError {
			return err
		}
		err = insertSummaryRefs(works)
		if err.HasError {
			return err
		}
		leftovers = findLeftoverSumms(works, summaries)
		if len(leftovers) > 0 {
			log.Debug().Msgf("unable to assign the following footnotes ('page.number') to a work: %v", leftovers)
		}
	}
	return errs.Nil()
}

func addFnsToWorks(works []model.Work, fns []model.Footnote) errs.UploadError {
	prevMax := int32(1)
	for i := range works {
		var min int32 = prevMax
		var max int32 = 1
		findMinMaxPages(works[i].Paragraphs, works[i].Sections, &min, &max)
		if min < prevMax {
			return errs.New(fmt.Errorf("minimum page number %d of work '%s' is smaller than the maximum page number %d of the previous work", min, works[i].Title, prevMax), nil)
		}
		for j := range fns {
			pages := fns[j].Pages
			if pages[0] >= min && pages[len(pages)-1] <= max {
				works[i].Footnotes = append(works[i].Footnotes, fns[j])
			}
		}
		prevMax = max
	}
	return errs.Nil()
}

func findLeftoverFns(works []model.Work, footnotes []model.Footnote) []string {
	workFns := make(map[string]struct{})
	for _, w := range works {
		for _, fn := range w.Footnotes {
			workFns[fn.Ref] = struct{}{}
		}
	}
	result := []string{}
	for _, fn := range footnotes {
		if _, exists := workFns[fn.Ref]; !exists {
			result = append(result, fn.Ref)
		}
	}
	return result
}

func addSummariesToWorks(works []model.Work, summaries []model.Summary) errs.UploadError {
	prevMax := int32(1)
	for i := range works {
		var min int32 = prevMax + 1
		var max int32 = 1
		findMinMaxPages(works[i].Paragraphs, works[i].Sections, &min, &max)
		if min < prevMax {
			return errs.New(fmt.Errorf("minimum page number %d of work '%s' is smaller than the maximum page number %d of the previous work", min, works[i].Title, prevMax), nil)
		}
		for j := range summaries {
			pages := summaries[j].Pages
			if pages[0] >= min && pages[len(pages)-1] <= max {
				works[i].Summaries = append(works[i].Summaries, summaries[j])
			}
		}
		prevMax = max
	}
	return errs.Nil()
}

func insertSummaryRefs(works []model.Work) errs.UploadError {
	for i := range works {
		w := &works[i]
		for j := range w.Summaries {
			summary := &w.Summaries[j]
			err := insertSummaryRef(summary, w.Paragraphs, w.Sections)
			if err.HasError {
				return err
			}
		}
	}
	return errs.Nil()
}

func findLeftoverSumms(works []model.Work, summaries []model.Summary) []string {
	workSumms := make(map[string]struct{})
	for _, w := range works {
		for _, summ := range w.Summaries {
			workSumms[summ.Ref] = struct{}{}
		}
	}
	result := []string{}
	for _, summ := range summaries {
		if _, exists := workSumms[summ.Ref]; !exists {
			result = append(result, summ.Ref)
		}
	}
	return result
}

func findMinMaxPages(paragraphs []model.Paragraph, sections []model.Section, min, max *int32) {
	for _, p := range paragraphs {
		if len(p.Pages) > 0 {
			if p.Pages[0] < *min {
				*min = p.Pages[0]
			}
			if p.Pages[len(p.Pages)-1] > *max {
				*max = p.Pages[len(p.Pages)-1]
			}
		}
	}
	for _, s := range sections {
		if len(s.Heading.Pages) > 0 {
			if s.Heading.Pages[0] < *min {
				*min = s.Heading.Pages[0]
			}
			if s.Heading.Pages[len(s.Heading.Pages)-1] > *max {
				*max = s.Heading.Pages[len(s.Heading.Pages)-1]
			}
		}
		findMinMaxPages(s.Paragraphs, s.Sections, min, max)
	}
}

func insertSummaryRef(summary *model.Summary, paragraphs []model.Paragraph, sections []model.Section) errs.UploadError {
	page, line := findPageLine(summary.Ref)
	p, err := findSummaryParagraph(summary, paragraphs, sections)
	if err.HasError {
		// in this case the summary starts in the middle of a paragraph, this is probably an error in the text, so we ignore the summary
		// TODO improve this behavior
		log.Debug().Msgf("found summary in the middle of a paragraph: %d.%d", page, line)
		return errs.Nil()
	}
	if p == nil {
		return errs.New(fmt.Errorf("could not find a paragraph for summary on page %d line %d", page, line), nil)
	}

	// duplicate page ref in the summary, so that summary and paragraph can be displayed independently from each other without loosing the page ref
	if line == 1 && !strings.Contains(summary.Text, util.FmtPage(page)) {
		summary.Text = util.FmtPage(page) + summary.Text
	}
	// line references should already by included in the summary text

	p.SummaryRef = &summary.Ref
	return errs.Nil()
}

func findSummaryParagraph(summary *model.Summary, paragraphs []model.Paragraph, sections []model.Section) (*model.Paragraph, errs.UploadError) {
	page, line := findPageLine(summary.Ref)
	for i := range paragraphs {
		p := &paragraphs[i]
		ok, err := isSummaryParagraph(p, page, line)
		if err.HasError {
			return nil, err
		}
		if ok {
			return p, errs.Nil()
		}
	}
	for i := range sections {
		s := &sections[i]
		p, err := findSummaryParagraph(summary, s.Paragraphs, s.Sections)
		if err.HasError {
			return nil, err
		}
		if p != nil {
			return p, errs.Nil()
		}
	}
	return nil, errs.Nil()
}

func findPageLine(name string) (int32, int32) {
	pageLine := strings.Split(name, ".")
	// ignore errs, because we know the format
	page, _ := strconv.ParseInt(pageLine[0], 10, 32)
	line, _ := strconv.ParseInt(pageLine[1], 10, 32)
	return int32(page), int32(line)
}

func isSummaryParagraph(p *model.Paragraph, page, line int32) (bool, errs.UploadError) {
	if !slices.Contains(p.Pages, page) {
		return false, errs.Nil()
	}
	pageIndex := strings.Index(p.Text, util.FmtPage(page))
	if pageIndex == -1 { // paragraph starts in the middle of the page
		pageIndex = 0
	}
	lineIndex := strings.Index(p.Text[pageIndex:], util.FmtLine(line))
	if lineIndex == -1 {
		return false, errs.Nil()
	}
	index := pageIndex + lineIndex + len(util.FmtLine(line))
	if !isSummaryAtStart(p.Text, index) {
		return false, errs.New(fmt.Errorf("summary on page %d line %d is not at the start of paragraph", page, line), nil)
	}
	return true, errs.Nil()
}

func isSummaryAtStart(text string, startIndex int) bool {
	cleaned := util.RemoveTags(text[:startIndex])
	return cleaned == "" // text before summary is only formatting code, so the "real text" starts with the summary
}
