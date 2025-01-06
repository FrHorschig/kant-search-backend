package transform

import (
	"regexp"
)

func Simplify(xml string) string {
	// We don't replace <seite ...> with {p...}, because the <seite> elements are sometimes on the same level as the headings and paragraphs. Replacing them here would make finding them later more difficult, so we do the <seite> replacement manually at another place.
	reZeile := regexp.MustCompile(`<zeile\s+nr="(\d+)"\s*/>`)
	xml = reZeile.ReplaceAllString(xml, `{l$1}`)
	reTrenn := regexp.MustCompile(`<trenn\s*/>`)
	xml = reTrenn.ReplaceAllString(xml, "")
	return xml
}
