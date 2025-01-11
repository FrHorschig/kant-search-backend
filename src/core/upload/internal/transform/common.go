package transform

import (
	"fmt"

	"github.com/beevik/etree"
)

func Seite(el *etree.Element) string {
	// TODO improvement: handle 'satz' attribute
	return fmt.Sprintf(
		"<ks-page>%s</ks-page>",
		el.SelectAttrValue("nr", "MISSING_PAGE_NUMBER"),
	)
}

func em1(el *etree.Element) string {
	// TODO implement me
	return ""
}

func fett(el *etree.Element) string {
	// TODO implement me
	return ""
}

func fr(el *etree.Element) string {
	// TODO implement me
	return ""
}

func fremdsprache(el *etree.Element) string {
	// TODO implement me
	return ""
}

func gesperrt(el *etree.Element) string {
	// TODO implement me
	return ""
}

func name(el *etree.Element) string {
	// TODO implement me
	return ""
}

func romzahl(el *etree.Element) string {
	// TODO implement me
	return ""
}

func zeile(el *etree.Element) string {
	return fmt.Sprintf(
		"<ks-line>%s</ks-line>",
		el.SelectAttrValue("nr", "MISSING_LINE_NUMBER"),
	)
}
