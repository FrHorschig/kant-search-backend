package transform

import (
	"fmt"
	"strings"

	"github.com/beevik/etree"
	"github.com/frhorschig/kant-search-backend/common/errors"
)

func P(p *etree.Element) (string, errors.ErrorNew) {
	// TODO improvement: handle 'ausrichtung' attribute
	switchFn := func(el *etree.Element) (string, errors.ErrorNew) {
		switch el.Tag {
		case "antiqua":
			return antiqua(el), errors.NilError()
		case "em1":
			return em1(el), errors.NilError()
		case "em2":
			return em2(el), errors.NilError()
		case "fett":
			return fett(el), errors.NilError()
		case "formel":
			return formel(el), errors.NilError()
		case "fr":
			return fr(el), errors.NilError()
		case "fremdsprache":
			return fremdsprache(el), errors.NilError()
		case "gesperrt":
			return gesperrt(el), errors.NilError()
		case "name":
			return name(el), errors.NilError()
		case "op":
			return "", errors.NilError()
		case "romzahl":
			return romzahl(el), errors.NilError()
		case "seite":
			return Seite(el), errors.NilError()
		case "trenn":
			return "", errors.NilError()
		case "zeile":
			return zeile(el), errors.NilError()
		default:
			return "", errors.NewError(fmt.Errorf("unknown tag '%s' in %s element", el.Tag, p.Tag), nil)
		}
	}
	return extractText(p, switchFn)
}

func Table(table *etree.Element) (string, errors.ErrorNew) {
	// TODO implement me: many things ...
	return "", errors.NilError()
}

func antiqua(antiqua *etree.Element) string {
	// TODO implement me: zeile, trenn, seite, gesperrt, name, fett
	return ""
}

func bild(bild *etree.Element) string {
	// TODO implement me
	// attributes: src, beschreibung, typ, Ort, z_anfang, z_ende, ausrichtung
	return ""
}

func bildverweis(bildverweis *etree.Element) string {
	// TODO implement me
	// attributes: src, beschreibung, typ, Ort, z_anfang, z_ende, ausrichtung
	return ""
}
func em2(em2El *etree.Element) string {
	// TODO implement me
	text := ""
	// TODO check the AA scans to maybe find something better
	return fmt.Sprintf(
		"<ks-tracked>%s</ks-tracked>",
		strings.TrimSpace(text),
	)
}

func formel(formel *etree.Element) string {
	// TODO implement me: em1
	return strings.TrimSpace(formel.Text())
}
