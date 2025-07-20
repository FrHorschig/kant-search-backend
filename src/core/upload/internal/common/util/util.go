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
	fnRefFmt    = `<ks-meta-fnref>%d.%d</ks-meta-fnref>`
	FnRefMatch  = `<ks-meta-fnref>(\d+\.\d+)</ks-meta-fnref>`
	imgFmt      = `<fmt-meta-img src="%s" desc="%s"/>`
	ImgMatch    = `<fmt-meta-img src=".+" desc=".+"/>`
	imgRefFmt   = `<fmt-meta-imgref src="%s" desc="%s"/>`
	ImgRefMatch = `<fmt-meta-imgref src=".+" desc=".+"/>`
	lineFmt     = `<ks-meta-line>%d</ks-meta-line>`
	LineMatch   = `<ks-meta-line>(\d+)</ks-meta-line>`
	pageFmt     = `<ks-meta-page>%d</ks-meta-page>`
	PageMatch   = `<ks-meta-page>(\d+)</ks-meta-page>`
)

func FmtFnRef(page int32, nr int32) string {
	return fmt.Sprintf(fnRefFmt, page, nr)
}

func FmtImg(src, desc string) string {
	return "" // ignore for now
}

func FmtImgRef(src, desc string) string {
	return "" // ignore for now
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
	headMatch       = headMatchStart + "%s" + headMatchEnd
	langFmt         = "%s"
	nameFmtStart    = "<ks-fmt-name>"
	nameFmtEnd      = "</ks-fmt-name>"
	nameFmt         = nameFmtStart + "%s" + nameFmtEnd
	parHeadFmtStart = "<ks-fmt-hpar>"
	parHeadFmtEnd   = "</ks-fmt-hpar>"
	parHeadFmt      = parHeadFmtStart + "%s" + parHeadFmtEnd
	tableFmtStart   = `<ks-fmt-table>`
	tableFmtEnd     = `</ks-fmt-table>`
	tableFmt        = tableFmtStart + "%s" + tableFmtEnd
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

func FmtName(content string) string {
	return fmt.Sprintf(nameFmt, content)
}

func FmtParHeading(content string) string {
	return fmt.Sprintf(parHeadFmt, content)
}

func FmtTable(content string) string {
	return fmt.Sprintf(tableFmt, content)
}

func FmtTracked(content string) string {
	return fmt.Sprintf(trackedFmt, content)
}

func RemoveTags(input string) string {
	input = regexp.MustCompile(FnRefMatch).ReplaceAllString(input, "")
	input = regexp.MustCompile(ImgMatch).ReplaceAllString(input, "")
	input = regexp.MustCompile(ImgRefMatch).ReplaceAllString(input, "")
	input = regexp.MustCompile(LineMatch).ReplaceAllString(input, "")
	input = regexp.MustCompile(PageMatch).ReplaceAllString(input, "")
	input = regexp.MustCompile(headMatchStart).ReplaceAllString(input, "")
	input = regexp.MustCompile(headMatchEnd).ReplaceAllString(input, "")

	input = strings.ReplaceAll(input, boldFmtStart, "")
	input = strings.ReplaceAll(input, boldFmtEnd, "")
	input = strings.ReplaceAll(input, emphFmtStart, "")
	input = strings.ReplaceAll(input, emphFmtEnd, "")
	input = strings.ReplaceAll(input, emph2FmtStart, "")
	input = strings.ReplaceAll(input, emph2FmtEnd, "")
	input = strings.ReplaceAll(input, formulaFmtStart, "")
	input = strings.ReplaceAll(input, formulaFmtEnd, "")
	input = strings.ReplaceAll(input, nameFmtStart, "")
	input = strings.ReplaceAll(input, nameFmtEnd, "")
	input = strings.ReplaceAll(input, parHeadFmtStart, "")
	input = strings.ReplaceAll(input, parHeadFmtEnd, "")
	input = strings.ReplaceAll(input, tableFmtStart, "")
	input = strings.ReplaceAll(input, tableFmtEnd, "")
	input = strings.ReplaceAll(input, "<tr>", "")
	input = strings.ReplaceAll(input, "</tr>", "")
	input = regexp.MustCompile(`<td(?:\s+(?:colspan|rowspan)="[^"]*")*\s*>`).ReplaceAllString(input, "")
	input = strings.ReplaceAll(input, "</td>", "")
	input = strings.ReplaceAll(input, trackedFmtStart, "")
	input = strings.ReplaceAll(input, trackedFmtEnd, "")

	input = regexp.MustCompile(`<[^>]*>`).ReplaceAllString(input, "")
	input = regexp.MustCompile(`\s+`).ReplaceAllString(input, " ")
	return strings.TrimSpace(input)
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
	input = strings.ReplaceAll(input, nameFmtStart, mask(trackedFmtStart))
	input = strings.ReplaceAll(input, nameFmtEnd, mask(trackedFmtEnd))
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
