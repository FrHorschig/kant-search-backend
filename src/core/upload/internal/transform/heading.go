package transform

import (
	"fmt"
	"strings"

	"github.com/beevik/etree"
	"github.com/frhorschig/kant-search-backend/common/errors"
	commonmodel "github.com/frhorschig/kant-search-backend/common/model"
	"github.com/frhorschig/kant-search-backend/core/upload/internal/model"
)

func Hx(hx *etree.Element) (model.Heading, errors.ErrorNew) {
	textTitle := ""
	tocTitle := ""
	for _, ch := range hx.Child {
		if str, ok := ch.(*etree.CharData); ok {
			textTitle += str.Data
			tocTitle += str.Data
		} else if el, ok := ch.(*etree.Element); ok {
			switch el.Tag {
			case "hi":
				tocTitle += el.Text()
			case "hu":
				hu, err := Hu(el)
				if err.HasError {
					return model.Heading{}, err
				}
				textTitle += hu
			case "zeile":
				textTitle += zeile(el)
			case "seite":
				textTitle += Seite(el)
			case "op":
				continue
			case "fremdsprache":
				fremdsprache := fremdsprache(el)
				tocTitle += fremdsprache
				textTitle += fremdsprache
			case "romzahl":
				romzahl := romzahl(el)
				tocTitle += romzahl
				textTitle += romzahl
			case "gesperrt":
				gesperrt := gesperrt(el)
				tocTitle += gesperrt
				textTitle += gesperrt
			case "name":
				name := name(el)
				tocTitle += name
				textTitle += name
			case "fett":
				fett := fett(el)
				tocTitle += fett
				textTitle += fett
			case "fr":
				textTitle += fr(el)
			case "trenn":
				continue
			default:
				return model.Heading{}, errors.NewError(fmt.Errorf("unknown tag '%s' in hauptteil element", el.Tag), nil)
			}
		}
	}
	return model.Heading{
		TocTitle:  tocTitle,
		TextTitle: textTitle,
		Level:     level(hx),
	}, errors.NilError()
}

func Hu(hu *etree.Element) (string, errors.ErrorNew) {
	text := ""
	for _, ch := range hu.Child {
		if str, ok := ch.(*etree.CharData); ok {
			text += strings.TrimSpace(str.Data)
		} else if el, ok := ch.(*etree.Element); ok {
			switch el.Tag {
			case "seite":
				text += Seite(el)
			case "zeile":
				text += zeile(el)
			case "fremdsprache":
				text += fremdsprache(el)
			case "romzahl":
				text += romzahl(el)
			case "gesperrt":
				text += gesperrt(el)
			case "name":
				text += name(el)
			case "fett":
				text += fett(el)
			case "em1":
				text += em1(el)
			case "fr":
				text += fr(el)
			case "op":
				continue
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

func level(el *etree.Element) commonmodel.Level {
	switch el.Tag {
	case "h1":
		return commonmodel.H1
	case "h2":
		return commonmodel.H2
	case "h3":
		return commonmodel.H3
	case "h4":
		return commonmodel.H4
	case "h5":
		return commonmodel.H5
	case "h6":
		return commonmodel.H6
	case "h7":
		return commonmodel.H7
	case "h8":
		return commonmodel.H8
	}
	return commonmodel.H9
}
