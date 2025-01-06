package mapping

import (
	"regexp"
)

func Simplify(xml string) string {
	reZeile := regexp.MustCompile(`<zeile\s+nr="(\d+)"\s*/>`)
	xml = reZeile.ReplaceAllString(xml, `{l$1}`)
	reSeite := regexp.MustCompile(`<seite\s*[^>]\s*nr="(\d+)"\s*[^>]*\s*/>`)
	xml = reSeite.ReplaceAllString(xml, `{p$1}`)
	reTrenn := regexp.MustCompile(`<trenn\s*/>`)
	xml = reTrenn.ReplaceAllString(xml, "")
	return xml
}
