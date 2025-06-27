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
		pages = &[]int32{*latestPage}
	}
}

func startsWithPageRef(text, pageRef string) bool {
	index := strings.Index(text, pageRef)
	cleaned := util.RemoveTags(text[:index])
	return cleaned == "" // in this case the text before page ref is only formatting code, so the actual text content starts with the page ref
}
