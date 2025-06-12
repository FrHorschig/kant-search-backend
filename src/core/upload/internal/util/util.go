package util

import (
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/beevik/etree"
	"github.com/frhorschig/kant-search-backend/common/errs"
)

const (
	fnRefFmt   = `<ks-meta-fnref>%d.%d</ks-meta-fnref>`
	FnRefMatch = `<ks-meta-fnref>(\d+\.\d+)</ks-meta-fnref>`
	lineFmt    = `<ks-meta-line>%d</ks-meta-line>`
	LineMatch  = `<ks-meta-line>(\d+)</ks-meta-line>`
	pageFmt    = `<ks-meta-page>%d</ks-meta-page>`
	PageMatch  = `<ks-meta-page>(\d+)</ks-meta-page>`
)

func FmtFnRef(page int32, nr int32) string {
	return fmt.Sprintf(fnRefFmt, page, nr)
}

func FmtLine(line int32) string {
	return fmt.Sprintf(lineFmt, line)
}

func FmtPage(page int32) string {
	return fmt.Sprintf(pageFmt, page)
}

const (
	boldFmtStart    = "<ks-fmt-bold>"
	boldFmtEnd      = "</ks-fmt-bold>"
	boldFmt         = boldFmtStart + "%s" + boldFmtEnd
	emphFmtStart    = "<ks-fmt-emph>"
	emphFmtEnd      = "</ks-fmt-emph>"
	emphFmt         = emphFmtStart + "%s" + emphFmtEnd
	emph2FmtStart   = "<ks-fmt-emph2>"
	emph2FmtEnd     = "</ks-fmt-emph2>"
	emph2Fmt        = emph2FmtStart + "%s" + emph2FmtEnd
	formulaFmtStart = "<ks-fmt-formula>"
	formulaFmtEnd   = "</ks-fmt-formula>"
	formulaFmt      = formulaFmtStart + "%s" + formulaFmtEnd
	headingFmt      = "<ks-fmt-h%d>%s</ks-fmt-h%d>"
	headMatchStart  = `<ks-fmt-h\d>`
	headMatchEnd    = `</ks-fmt-h\d>`
	headMatch       = headMatchStart + `%s` + headMatchEnd
	langFmt         = `%s`
	parHeadFmtStart = "<ks-fmt-hpar>"
	parHeadFmtEnd   = "</ks-fmt-hpar>"
	parHeadFmt      = parHeadFmtStart + "%s" + parHeadFmtEnd
	trackedFmtStart = "<ks-fmt-tracked>"
	trackedFmtEnd   = "</ks-fmt-tracked>"
	trackedFmt      = trackedFmtStart + "%s" + trackedFmtEnd
)

func FmtBold(content string) string {
	return fmt.Sprintf(boldFmt, content)
}

func FmtEmph(content string) string {
	return fmt.Sprintf(emphFmt, content)
}

func FmtEmph2(content string) string {
	return fmt.Sprintf(emph2Fmt, content)
}

func FmtFormula(content string) string {
	return fmt.Sprintf(formulaFmt, content)
}

func FmtHeading(level int32, content string) string {
	return fmt.Sprintf(headingFmt, level, content, level)
}

func FmtLang(lang string) string {
	return lang
}

func FmtParHeading(content string) string {
	return fmt.Sprintf(parHeadFmt, content)
}

func FmtTracked(content string) string {
	return fmt.Sprintf(trackedFmt, content)
}

const (
	imageFmt      = `{extract-image src="%s" desc="%s"}`
	ImageRefMatch = `{extract-image src=".+" desc=".+"}`
	TableMatch    = `` // ignore for now
)

func FmtImage(src, desc string) string {
	return "" // ignore for now
}

func RemoveTags(text string) string {
	re := regexp.MustCompile(FnRefMatch)
	text = re.ReplaceAllString(text, "")
	re = regexp.MustCompile(LineMatch)
	text = re.ReplaceAllString(text, "")
	re = regexp.MustCompile(PageMatch)
	text = re.ReplaceAllString(text, "")
	re = regexp.MustCompile(`<[^>]*>`)
	text = re.ReplaceAllString(text, "")
	re = regexp.MustCompile(`\s+`)
	text = re.ReplaceAllString(text, " ")
	return strings.TrimSpace(text)
}

func MaskTags(input string) string {
	re := regexp.MustCompile(FnRefMatch)
	input = re.ReplaceAllStringFunc(input, mask)
	re = regexp.MustCompile(LineMatch)
	input = re.ReplaceAllStringFunc(input, mask)
	re = regexp.MustCompile(PageMatch)
	input = re.ReplaceAllStringFunc(input, mask)
	re = regexp.MustCompile(headMatchStart)
	input = re.ReplaceAllStringFunc(input, mask)
	re = regexp.MustCompile(headMatchEnd)
	input = re.ReplaceAllStringFunc(input, mask)

	input = strings.ReplaceAll(input, boldFmtStart, mask(boldFmtStart))
	input = strings.ReplaceAll(input, boldFmtEnd, mask(boldFmtEnd))
	input = strings.ReplaceAll(input, emphFmtStart, mask(emphFmtStart))
	input = strings.ReplaceAll(input, emphFmtEnd, mask(emphFmtEnd))
	input = strings.ReplaceAll(input, emph2FmtStart, mask(emph2FmtStart))
	input = strings.ReplaceAll(input, emph2FmtEnd, mask(emph2FmtEnd))
	input = strings.ReplaceAll(input, formulaFmtStart, mask(formulaFmtStart))
	input = strings.ReplaceAll(input, formulaFmtEnd, mask(formulaFmtEnd))
	input = strings.ReplaceAll(input, parHeadFmtStart, mask(parHeadFmtStart))
	input = strings.ReplaceAll(input, parHeadFmtEnd, mask(parHeadFmtEnd))
	input = strings.ReplaceAll(input, trackedFmtStart, mask(trackedFmtStart))
	input = strings.ReplaceAll(input, trackedFmtEnd, mask(trackedFmtEnd))
	return input
}

func mask(s string) string {
	return strings.Repeat("*", len(s))
}

func ExtractNumericAttribute(el *etree.Element, attr string) (int32, errs.UploadError) {
	defaultStr := "DEFAULT_STRING"
	nStr := strings.TrimSpace(el.SelectAttrValue(attr, defaultStr))
	if nStr == defaultStr {
		return 0, errs.New(fmt.Errorf("missing '%s' attribute in '%s' element", attr, el.Tag), nil)
	}

	// TODO how to best handle this special case?
	if slices.Contains([]string{"272a", "272c", "272d"}, nStr) {
		nStr = "272"
	}

	n, err := strconv.ParseInt(nStr, 10, 32)
	if err != nil {
		return 0, errs.New(fmt.Errorf("can't convert attribute string '%s' to number", nStr), nil)
	}
	return int32(n), errs.Nil()
}

func ExtractFnRefs(text string) []string {
	re := regexp.MustCompile(FnRefMatch)
	matches := re.FindAllStringSubmatch(text, -1)
	result := []string{}
	for _, match := range matches {
		result = append(result, match[1])
	}
	return result
}

func ExtractPages(text string) ([]int32, errs.UploadError) {
	re := regexp.MustCompile(PageMatch)
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
