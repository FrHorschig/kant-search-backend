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
			textTitle += strings.TrimSpace(str.Data)
			tocTitle += str.Data
		} else if el, ok := ch.(*etree.Element); ok {
			switch el.Tag {
			case "fett":
				fett := fett(el)
				tocTitle += fett
				textTitle += fett
			case "fr":
				textTitle += fr(el)
			case "fremdsprache":
				fremdsprache := fremdsprache(el)
				tocTitle += fremdsprache
				textTitle += fremdsprache
			case "gesperrt":
				gesperrt := gesperrt(el)
				tocTitle += gesperrt
				textTitle += gesperrt
			case "hi":
				tocTitle += strings.TrimSpace(el.Text())
			case "hu":
				hu, err := Hu(el)
				if err.HasError {
					return model.Heading{}, err
				}
				textTitle += hu
			case "name":
				name := name(el)
				tocTitle += name
				textTitle += name
			case "op":
				continue
			case "romzahl":
				romzahl := romzahl(el)
				tocTitle += romzahl
				textTitle += romzahl
			case "seite":
				textTitle += Seite(el)
			case "trenn":
				continue
			case "zeile":
				textTitle += zeile(el)
			default:
				return model.Heading{}, errors.NewError(fmt.Errorf("unknown tag '%s' in hauptteil element", el.Tag), nil)
			}
		}
		tocTitle += " "
		textTitle += " "
	}
	return model.Heading{
		TocTitle:  strings.TrimSpace(tocTitle),
		TextTitle: strings.TrimSpace(textTitle),
		Level:     level(hx),
	}, errors.NilError()
}

func Hu(hu *etree.Element) (string, errors.ErrorNew) {
	switchFn := func(el *etree.Element) (string, errors.ErrorNew) {
		switch el.Tag {
		case "seite":
			return Seite(el), errors.NilError()
		case "zeile":
			return zeile(el), errors.NilError()
		case "fremdsprache":
			return fremdsprache(el), errors.NilError()
		case "romzahl":
			return romzahl(el), errors.NilError()
		case "gesperrt":
			return gesperrt(el), errors.NilError()
		case "name":
			return name(el), errors.NilError()
		case "fett":
			return fett(el), errors.NilError()
		case "em1":
			return em1(el), errors.NilError()
		case "fr":
			return fr(el), errors.NilError()
		case "op":
			return "", errors.NilError()
		case "trenn":
			return "", errors.NilError()
		default:
			return "", errors.NewError(fmt.Errorf("unknown tag '%s' in %s element", el.Tag, hu.Tag), nil)
		}
	}
	return extractText(hu, switchFn)
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
