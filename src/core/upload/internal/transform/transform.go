package transform

//go:generate mockgen -source=$GOFILE -destination=mocks/xml_transformator.go -package=mocks

import (
	"fmt"
	"strings"

	"github.com/beevik/etree"
	"github.com/frhorschig/kant-search-backend/common/errors"
	"github.com/frhorschig/kant-search-backend/core/upload/internal/model"
)

type XmlTransformator interface {
	Hx(el *etree.Element) (model.Heading, errors.ErrorNew)
	Hu(el *etree.Element) (string, errors.ErrorNew)
	P(el *etree.Element) (string, errors.ErrorNew)
	Seite(el *etree.Element) string
	Table() string
}

type XmlTransformatorImpl struct {
}

func NewXmlTransformator() XmlTransformator {
	impl := XmlTransformatorImpl{}
	return &impl
}

func (rec *XmlTransformatorImpl) Hx(el *etree.Element) (model.Heading, errors.ErrorNew) {
	return hx(el)
}

func (rec *XmlTransformatorImpl) Hu(el *etree.Element) (string, errors.ErrorNew) {
	return hu(el)
}

func (rec *XmlTransformatorImpl) P(el *etree.Element) (string, errors.ErrorNew) {
	return p(el)
}

func (rec *XmlTransformatorImpl) Seite(el *etree.Element) string {
	return seite(el)
}

func (rec *XmlTransformatorImpl) Table() string {
	return table()
}

func hx(elem *etree.Element) (model.Heading, errors.ErrorNew) {
	textTitle := ""
	tocTitle := ""
	for _, ch := range elem.Child {
		if str, ok := ch.(*etree.CharData); ok {
			textTitle += strings.TrimSpace(str.Data)
			tocTitle += str.Data
		} else if el, ok := ch.(*etree.Element); ok {
			switch el.Tag {
			case "fett":
				fett, err := fett(el)
				if err.HasError {
					return model.Heading{}, err
				}
				tocTitle += fett
				textTitle += fett
			case "fr":
				textTitle += fr(el)
			case "fremdsprache":
				fremdsprache, err := fremdsprache(el)
				if err.HasError {
					return model.Heading{}, err
				}
				tocTitle += fremdsprache
				textTitle += fremdsprache
			case "gesperrt":
				gesperrt, err := gesperrt(el)
				if err.HasError {
					return model.Heading{}, err
				}
				tocTitle += gesperrt
				textTitle += gesperrt
			case "hi":
				tocTitle += strings.TrimSpace(el.Text())
			case "hu":
				hu, err := hu(el)
				if err.HasError {
					return model.Heading{}, err
				}
				textTitle += hu
			case "name":
				name, err := name(el)
				if err.HasError {
					return model.Heading{}, err
				}
				tocTitle += name
				textTitle += name
			case "op":
				continue
			case "romzahl":
				romzahl, err := romzahl(el)
				if err.HasError {
					return model.Heading{}, err
				}
				tocTitle += romzahl
				textTitle += romzahl
			case "seite":
				textTitle += seite(el)
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
		Level:     level(elem),
	}, errors.NilError()
}

func hu(elem *etree.Element) (string, errors.ErrorNew) {
	switchFn := func(el *etree.Element) (string, errors.ErrorNew) {
		switch el.Tag {
		case "em1":
			return em1(el), errors.NilError()
		case "fett":
			return fett(el)
		case "fr":
			return fr(el), errors.NilError()
		case "fremdsprache":
			return fremdsprache(el)
		case "gesperrt":
			return gesperrt(el)
		case "name":
			return name(el)
		case "op":
			return "", errors.NilError()
		case "romzahl":
			return romzahl(el)
		case "seite":
			return seite(el), errors.NilError()
		case "trenn":
			return "", errors.NilError()
		case "zeile":
			return zeile(el), errors.NilError()
		default:
			return "", errors.NewError(fmt.Errorf("unknown tag '%s' in %s element", el.Tag, elem.Tag), nil)
		}
	}
	return extractText(elem, switchFn)
}

func p(elem *etree.Element) (string, errors.ErrorNew) {
	switchFn := func(el *etree.Element) (string, errors.ErrorNew) {
		switch el.Tag {
		case "antiqua":
			return antiqua(el)
		case "bild":
			return bildBildverweis(el), errors.NilError()
		case "bildverweis":
			return bildBildverweis(el), errors.NilError()
		case "em1":
			return em1(el), errors.NilError()
		case "em2":
			return em2(el)
		case "fett":
			return fett(el)
		case "formel":
			return formel(el)
		case "fr":
			return fr(el), errors.NilError()
		case "fremdsprache":
			return fremdsprache(el)
		case "gesperrt":
			return gesperrt(el)
		case "name":
			return name(el)
		case "op":
			return "", errors.NilError()
		case "romzahl":
			return romzahl(el)
		case "table":
			return table(), errors.NilError()
		case "seite":
			return seite(el), errors.NilError()
		case "trenn":
			return "", errors.NilError()
		case "zeile":
			return zeile(el), errors.NilError()
		default:
			return "", errors.NewError(fmt.Errorf("unknown tag '%s' in %s element", el.Tag, elem.Tag), nil)
		}
	}
	return extractText(elem, switchFn)
}

func seite(elem *etree.Element) string {
	return fmt.Sprintf(
		"<ks-meta-page>%s</ks-meta-page>",
		strings.TrimSpace(elem.SelectAttrValue("nr", "MISSING_PAGE_NUMBER")),
	)
}

func table() string {
	return "{table-extract}"
}
