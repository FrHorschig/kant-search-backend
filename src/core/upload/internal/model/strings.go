package model

const (
	LineFmt    = `<ks-meta-line>%d</ks-meta-line>`
	LineMatch  = `<ks-meta-line>(\d+)</ks-meta-line>`
	PageFmt    = `<ks-meta-page>%d</ks-meta-page>`
	PageMatch  = `<ks-meta-page>(\d+)</ks-meta-page>`
	FnRefFmt   = `<ks-meta-fnref>%d.%d</ks-meta-fnref>`
	FnRefMatch = `<ks-meta-fnref>(\d+\.\d+)</ks-meta-fnref>`
	LangFmt    = `%s`
)

const (
	SummaryFmt = "<ks-fmt-summary>%s</ks-fmt-summary>"
	EmphFmt    = "<ks-fmt-emph>%s</ks-fmt-emph>"
	Emph2Fmt   = "<ks-fmt-emph2>%s</ks-fmt-emph2>"
	BoldFmt    = "<ks-fmt-bold>%s</ks-fmt-bold>"
	FormulaFmt = "<ks-fmt-formula>%s</ks-fmt-formula>"
	TrackedFmt = "<ks-fmt-tracked>%s</ks-fmt-tracked>"
)

const (
	ImageFmt      = `{extract-image src="%s" desc="%s"}`
	ImageRefMatch = `{extract-image src=".+" desc=".+"}`
	TableMatch    = `{extract-table}`
)
