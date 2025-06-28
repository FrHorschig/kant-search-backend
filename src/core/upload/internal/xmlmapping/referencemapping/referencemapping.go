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

func addFnsToWorks(works []model.Work, footnotes []model.Footnote) errs.UploadError {
	prevMax := int32(1)
	for i := range works {
		var min int32 = prevMax
		var max int32 = 1
		findMinMaxPages(works[i].Paragraphs, works[i].Sections, &min, &max)
		if min < prevMax {
			return errs.New(fmt.Errorf("minimum page number %d of work '%s' is smaller than the maximum page number %d of the previous work", min, works[i].Title, prevMax), nil)
		}
		for j := range footnotes {
			pages := footnotes[j].Pages
			if pages[0] >= min && pages[0] <= max {
				works[i].Footnotes = append(works[i].Footnotes, footnotes[j])
			}
		}
		prevMax = max + 1
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
		findMinMaxParagraphHeadingPages(p.Pages, min, max)
	}
	for _, s := range sections {
		findMinMaxParagraphHeadingPages(s.Heading.Pages, min, max)
		findMinMaxPages(s.Paragraphs, s.Sections, min, max)
	}
}

func findMinMaxParagraphHeadingPages(pages []int32, min, max *int32) {
	if len(pages) > 0 {
		if pages[0] < *min {
			*min = pages[0]
		}
		if pages[len(pages)-1] > *max {
			*max = pages[len(pages)-1]
		}
	}
}

func insertSummaryRef(summary *model.Summary, paragraphs []model.Paragraph, sections []model.Section) errs.UploadError {
	page, line := findPageLine(summary.Ref)
	p := findSummaryParagraph(summary, paragraphs, sections)
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

func findSummaryParagraph(summary *model.Summary, paragraphs []model.Paragraph, sections []model.Section) *model.Paragraph {
	page, line := findPageLine(summary.Ref)
	for i := range paragraphs {
		p := &paragraphs[i]
		ok := isSummaryParagraph(p, page, line)
		if ok {
			return p
		}
	}
	for i := range sections {
		s := &sections[i]
		p := findSummaryParagraph(summary, s.Paragraphs, s.Sections)
		if p != nil {
			return p
		}
	}
	return nil
}

func findPageLine(name string) (int32, int32) {
	pageLine := strings.Split(name, ".")
	// ignore errs, because we know the format
	page, _ := strconv.ParseInt(pageLine[0], 10, 32)
	line, _ := strconv.ParseInt(pageLine[1], 10, 32)
	return int32(page), int32(line)
}

func isSummaryParagraph(p *model.Paragraph, page, line int32) bool {
	if !slices.Contains(p.Pages, page) {
		return false
	}
	pageIndex := strings.Index(p.Text, util.FmtPage(page))
	if pageIndex == -1 { // paragraph starts in the middle of the page
		pageIndex = 0
	}
	lineIndex := strings.Index(p.Text[pageIndex:], util.FmtLine(line))
	if lineIndex == -1 {
		return false
	}
	index := pageIndex + lineIndex + len(util.FmtLine(line))
	if !isSummaryAtStart(p.Text, index) {
		// TODO in this case the summary starts in the middle of a paragraph; this makes things complicated, so we simplify here and ignore it (we may improve this later)
		log.Debug().Msgf("found summary in the middle of a paragraph: %d.%d", page, line)
		return false
	}
	return true
}

func isSummaryAtStart(text string, startIndex int) bool {
	cleaned := util.RemoveTags(text[:startIndex])
	return cleaned == "" // text before summary is only formatting code, so the "real text" starts with the summary
}
