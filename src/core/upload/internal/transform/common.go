package transform

import (
	"fmt"
	"strings"

	"github.com/beevik/etree"
)

func Seite(seite *etree.Element) string {
	// TODO improvement: handle 'satz' attribute
	return fmt.Sprintf(
		"<ks-page>%s</ks-page>",
		seite.SelectAttrValue("nr", "MISSING_PAGE_NUMBER"),
	)
}

func em1(em1 *etree.Element) string {
	return fmt.Sprintf(
		"<em>%s</em>",
		strings.TrimSpace(em1.Text()),
	)
}

func fett(fett *etree.Element) string {
	// TODO implement me: zeile, trenn, seite
	return fmt.Sprintf(
		"<b>%s</b>",
		strings.TrimSpace(fett.Text()),
	)
}

func fr(fr *etree.Element) string {
	return fmt.Sprintf(
		"<ks-fn-ref>%s.%s</ks-fn-ref>",
		strings.TrimSpace(fr.SelectAttrValue("seite", "0")),
		strings.TrimSpace(fr.SelectAttrValue("nr", "0")),
	)
}

func fremdsprache(fs *etree.Element) string {
	// TODO implement me: seite, zeile, trenn, fremdsprache, fr, romzahl, gesperrt, name, fett, formel, bild, em1, em2, bildverweis
	// attribute: sprache, zeichen, umschrift
	return ""
}

func gesperrt(gesperrt *etree.Element) string {
	return fmt.Sprintf(
		"<ks-tracked>%s</em>",
		strings.TrimSpace(gesperrt.Text()),
	)
}

func name(name *etree.Element) string {
	// TODO implement me: zeile, trenn, seite
	return strings.TrimSpace(name.Text())
}

func romzahl(romzahl *etree.Element) string {
	// TODO implement me: only text
	return ""
}

func zeile(zeile *etree.Element) string {
	return fmt.Sprintf(
		"<ks-line>%s</ks-line>",
		zeile.SelectAttrValue("nr", "MISSING_LINE_NUMBER"),
	)
}
