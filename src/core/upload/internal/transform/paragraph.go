package transform

import (
	"fmt"
	"strings"

	"github.com/beevik/etree"
	"github.com/frhorschig/kant-search-backend/common/errors"
)

func P(p *etree.Element) (string, errors.ErrorNew) {
	// TODO improvement: handle 'ausrichtung' attribute
	text := ""
	for _, ch := range p.Child {
		if str, ok := ch.(*etree.CharData); ok {
			text += strings.TrimSpace(str.Data)
		} else if el, ok := ch.(*etree.Element); ok {
			switch el.Tag {
			case "antiqua":
				text += antiqua(el)
			case "em1":
				text += em1(el)
			case "em2":
				text += em2(el)
			case "fett":
				text += fett(el)
			case "formel":
				text += formel(el)
			case "fr":
				text += fr(el)
			case "fremdsprache":
				text += fremdsprache(el)
			case "gesperrt":
				text += gesperrt(el)
			case "name":
				text += name(el)
			case "op":
				continue
			case "romzahl":
				text += romzahl(el)
			case "seite":
				text += Seite(el)
			case "trenn":
				continue
			case "zeile":
				text += zeile(el)
			default:
				return "", errors.NewError(fmt.Errorf("unknown tag '%s' in hu element", el.Tag), nil)
			}
		}
		text += " "
	}
	return strings.TrimSpace(text), errors.NilError()
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
