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
	Seite(el *etree.Element) (string, errors.ErrorNew)
	Table() string
	Randtext(el *etree.Element) (model.Randtext, errors.ErrorNew)
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

func (rec *XmlTransformatorImpl) Seite(el *etree.Element) (string, errors.ErrorNew) {
	return seite(el)
}

func (rec *XmlTransformatorImpl) Table() string {
	return table()
}

func (rec *XmlTransformatorImpl) Randtext(el *etree.Element) (model.Randtext, errors.ErrorNew) {
	return randtext(el)
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
				page, err := seite(el)
				if err.HasError {
					return model.Heading{}, err
				}
				textTitle += page
			case "trenn":
				continue
			case "zeile":
				line, err := zeile(el)
				if err.HasError {
					return model.Heading{}, err
				}
				textTitle += line
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
			return seite(el)
		case "trenn":
			return "", errors.NilError()
		case "zeile":
			return zeile(el)
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
			return seite(el)
		case "trenn":
			return "", errors.NilError()
		case "zeile":
			return zeile(el)
		default:
			return "", errors.NewError(fmt.Errorf("unknown tag '%s' in %s element", el.Tag, elem.Tag), nil)
		}
	}
	return extractText(elem, switchFn)
}

func seite(elem *etree.Element) (string, errors.ErrorNew) {
	page, err := extractNumericAttribute(elem, "nr")
	if err.HasError {
		return "", err
	}
	return fmt.Sprintf(
		"<ks-meta-page>%d</ks-meta-page>",
		page,
	), errors.NilError()
}

func table() string {
	return "{table-extract}"
}

func randtext(elem *etree.Element) (model.Randtext, errors.ErrorNew) {
	switchFn := func(el *etree.Element) (string, errors.ErrorNew) {
		switch el.Tag {
		case "p":
			return p(el)
		default:
			return "", errors.NewError(fmt.Errorf("unknown tag '%s' in %s element", el.Tag, elem.Tag), nil)
		}
	}
	paragraph, err := extractText(elem, switchFn)
	if err.HasError {
		return model.Randtext{}, err
	}
	page, err := extractNumericAttribute(elem, "seite")
	if err.HasError {
		return model.Randtext{}, err
	}
	line, err := extractNumericAttribute(elem, "anfang")
	if err.HasError {
		return model.Randtext{}, err
	}
	return model.Randtext{
		Page: page,
		Line: line,
		Text: paragraph,
	}, errors.NilError()
}
