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
			case "seite":
				text += Seite(el)
			case "op":
				continue
			case "zeile":
				text += zeile(el)
			case "fr":
				text += fr(el)
			case "fremdsprache":
				text += fremdsprache(el)
			case "romzahl":
				text += romzahl(el)
			case "gesperrt":
				text += gesperrt(el)
			case "antiqua":
				text += antiqua(el)
			case "name":
				text += name(el)
			case "fett":
				text += fett(el)
			case "formel":
				text += formel(el)
			case "em1":
				text += em1(el)
			case "trenn":
				continue
			default:
				return "", errors.NewError(fmt.Errorf("unknown tag '%s' in hu element", el.Tag), nil)
			}
		}
		text += " "
	}
	return strings.TrimSpace(text), errors.NilError()
}

func Table(table *etree.Element) (string, errors.ErrorNew) {
	// TODO implement me
	return "", errors.NilError()
}

func antiqua(antiqua *etree.Element) string {
	// TODO implement me
	return ""
}

func formel(formel *etree.Element) string {
	// TODO implement me
	return ""
}

func bild(bild *etree.Element) string {
	// TODO implement me
	return ""
}

func em2(em2 *etree.Element) string {
	// TODO implement me
	return ""
}

func bildverweis(bildverweis *etree.Element) string {
	// TODO implement me
	return ""
}
