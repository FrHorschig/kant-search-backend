package util

import (
	"fmt"
)

const (
	fnRefFmt        = `<ks-meta-fnref>%d.%d</ks-meta-fnref>`
	FnRefMatch      = `<ks-meta-fnref>(\d+\.\d+)</ks-meta-fnref>`
	lineFmt         = `<ks-meta-line>%d</ks-meta-line>`
	LineMatch       = `<ks-meta-line>(\d+)</ks-meta-line>`
	pageFmt         = `<ks-meta-page>%d</ks-meta-page>`
	PageMatch       = `<ks-meta-page>(\d+)</ks-meta-page>`
	summaryRefFmt   = `<ks-meta-sumref>%s</ks-meta-sumref>`
	SummaryRefMatch = `<ks-meta-sumref>.*?</ks-meta-sumref>`
)

func FmtFnRef(page, nr int32) string {
	return fmt.Sprintf(fnRefFmt, page, nr)
}

func FmtLine(line int32) string {
	return fmt.Sprintf(lineFmt, line)
}

func FmtPage(page int32) string {
	return fmt.Sprintf(pageFmt, page)
}

func FmtSummaryRef(name string) string {
	return fmt.Sprintf(summaryRefFmt, name)
}

const (
	boldFmt       = "<ks-fmt-bold>%s</ks-fmt-bold>"
	emphFmt       = "<ks-fmt-emph>%s</ks-fmt-emph>"
	emph2Fmt      = "<ks-fmt-emph2>%s</ks-fmt-emph2>"
	formulaFmt    = "<ks-fmt-formula>%s</ks-fmt-formula>"
	headingFmt    = "<ks-fmt-h%d>%s</ks-fmt-h%d>"
	langFmt       = `%s`
	parHeadingFmt = "ks-fmt-hpar>%s</ks-fmt-hpar>"
	trackedFmt    = "<ks-fmt-tracked>%s</ks-fmt-tracked>"
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
	return fmt.Sprintf(parHeadingFmt, content)
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
