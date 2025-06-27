package metadataextraction

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/frhorschig/kant-search-backend/common/errs"
	"github.com/frhorschig/kant-search-backend/core/upload/internalnew/common/model"
	"github.com/frhorschig/kant-search-backend/core/upload/internalnew/common/util"
)

func ExtractMetadata(works []model.Work, footnotes []model.Footnote, summaries []model.Summary) errs.UploadError {
	latestPage := int32(1)
	for i := range works {
		err := processParagraphs(works[i].Paragraphs, &latestPage)
		if err.HasError {
			return err
		}
		err = processSections(works[i].Sections, &latestPage)
		if err.HasError {
			return err
		}
	}
	err := processFootnotes(footnotes)
	if err.HasError {
		return err
	}
	err = processSummaries(summaries)
	if err.HasError {
		return err
	}
	return errs.Nil()
}

func processSections(sections []model.Section, latestPage *int32) errs.UploadError {
	for i := range sections {
		h := &sections[i].Heading
		err := extractMetadata(h.Text, &h.Pages, &h.FnRefs, latestPage)
		if err.HasError {
			return err
		}
		err = processParagraphs(sections[i].Paragraphs, latestPage)
		if err.HasError {
			return err
		}
		err = processSections(sections[i].Sections, latestPage)
		if err.HasError {
			return err
		}
	}
	return errs.Nil()
}

func processParagraphs(paragraphs []model.Paragraph, latestPage *int32) errs.UploadError {
	for i := range paragraphs {
		p := &paragraphs[i]
		err := extractMetadata(p.Text, &p.Pages, &p.FnRefs, latestPage)
		if err.HasError {
			return err
		}
	}
	return errs.Nil()
}

func extractMetadata(text string, pages *[]int32, fnRefs *[]string, latestPage *int32) errs.UploadError {
	pgs, err := extractPages(text)
	if err.HasError {
		return err
	}
	*pages = pgs
	*fnRefs = extractFnRefs(text)
	supplementPages(pages, text, latestPage)
	return errs.Nil()
}

func extractFnRefs(text string) []string {
	re := regexp.MustCompile(util.FnRefMatch)
	matches := re.FindAllStringSubmatch(text, -1)
	result := []string{}
	for _, match := range matches {
		result = append(result, match[1])
	}
	return result
}

func extractPages(text string) ([]int32, errs.UploadError) {
	re := regexp.MustCompile(util.PageMatch)
	matches := re.FindAllStringSubmatch(text, -1)

	result := []int32{}
	for _, match := range matches {
		nStr := match[1]
		n, err := strconv.ParseInt(nStr, 10, 32)
		if err != nil {
			return nil, errs.New(fmt.Errorf("can't convert page string '%s' to number", nStr), nil)
		}
		result = append(result, int32(n))
	}

	return result, errs.Nil()
}

func supplementPages(pages *[]int32, text string, latestPage *int32) {
	if len(*pages) > 0 {
		firstPage := (*pages)[0]
		pageRef := util.FmtPage(firstPage)
		if !startsWithPageRef(text, pageRef) {
			*pages = append([]int32{firstPage - 1}, *pages...)
		}
		lastPage := (*pages)[len(*pages)-1]
		if lastPage > *latestPage {
			*latestPage = lastPage
		}
	} else {
		// this happens when a heading is fully inside a page and at least on line away from the page start and end
		*pages = []int32{*latestPage}
	}
}

func startsWithPageRef(text, pageRef string) bool {
	index := strings.Index(text, pageRef)
	if index < 0 {
		return false
	}
	cleaned := util.RemoveTags(text[:index])
	return cleaned == "" // in this case the text before page ref is only formatting code, so the actual text content starts with the page ref
}

func processFootnotes(footnotes []model.Footnote) errs.UploadError {
	for i, f := range footnotes {
		pages, err := util.ExtractPages(f.Text)
		if err.HasError {
			return err
		}
		refPage, err := getPageFromRef(f.Ref)
		if err.HasError {
			return err
		}
		if len(pages) == 0 {
			pages = []int32{refPage}
		} else if !startsWithPageRef(f.Text, util.FmtPage(pages[0])) {
			pages = append([]int32{pages[0] - 1}, pages...)
		}
		if pages[0] != refPage {
			return errs.New(fmt.Errorf("footnote page %d does not match the first page of the footnote %d", refPage, pages[0]), nil)
		}
		footnotes[i].Pages = pages
	}
	return errs.Nil()
}

func processSummaries(summaries []model.Summary) errs.UploadError {
	for i, s := range summaries {
		pages, err := util.ExtractPages(s.Text)
		if err.HasError {
			return err
		}
		refPage, err := getPageFromRef(s.Ref)
		if err.HasError {
			return err
		}
		if len(pages) == 0 {
			pages = []int32{refPage}
		} else if !startsWithPageRef(s.Text, util.FmtPage(pages[0])) {
			pages = append([]int32{pages[0] - 1}, pages...)
		}
		if pages[0] != refPage {
			return errs.New(fmt.Errorf("summary page %d does not match the first page of the summary %d", refPage, pages[0]), nil)
		}
		summaries[i].Pages = pages
		summaries[i].FnRefs = extractFnRefs(s.Text)
	}
	return errs.Nil()
}

func getPageFromRef(ref string) (int32, errs.UploadError) {
	parts := strings.Split(ref, ".")
	page, err := strconv.ParseInt(parts[0], 10, 32)
	if err != nil {
		return 0, errs.New(nil, fmt.Errorf("unable to convert page '%s' from ref to number", parts[0]))
	}
	return int32(page), errs.Nil()
}
