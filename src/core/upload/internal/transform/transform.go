package transform

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/beevik/etree"
	"github.com/frhorschig/kant-search-backend/core/upload/errors"
	"github.com/frhorschig/kant-search-backend/core/upload/internal/extract"
	model "github.com/frhorschig/kant-search-backend/core/upload/internal/treemodel"
	"github.com/frhorschig/kant-search-backend/core/upload/internal/util"
)

func Hx(el *etree.Element) (model.Heading, errors.UploadError) {
	return hx(el)
}

func Hu(el *etree.Element) (string, errors.UploadError) {
	return hu(el)
}

func P(el *etree.Element) (string, errors.UploadError) {
	return p(el)
}

func Seite(el *etree.Element) (string, errors.UploadError) {
	return seite(el)
}

func Table() string {
	return table()
}

func Summary(el *etree.Element) (model.Summary, errors.UploadError) {
	return summary(el)
}

func Footnote(el *etree.Element) (model.Footnote, errors.UploadError) {
	return footnote(el)
}

func hx(elem *etree.Element) (model.Heading, errors.UploadError) {
	textTitle := ""
	tocTitle := ""
	for _, ch := range elem.Child {
		if str, ok := ch.(*etree.CharData); ok {
			textTitle += strings.TrimSpace(str.Data)
			tocTitle += " " + str.Data
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
				fr, err := fr(el)
				if err.HasError {
					return model.Heading{}, err
				}
				textTitle += fr
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
				return model.Heading{}, errors.New(fmt.Errorf("unknown tag '%s' in hauptteil element", el.Tag), nil)
			}
		}
		tocTitle += " "
		textTitle += " "
	}
	tocTitle = extract.RemoveTags(tocTitle)
	tocTitle = strings.TrimSpace(tocTitle)
	tocTitle = removePunctuation(tocTitle)
	tocTitle = fixCapitalization(tocTitle)
	return model.Heading{
		Level:     level(elem),
		TocTitle:  tocTitle,
		TextTitle: strings.TrimSpace(textTitle),
	}, errors.Nil()
}

func hu(elem *etree.Element) (string, errors.UploadError) {
	switchFn := func(el *etree.Element) (string, errors.UploadError) {
		switch el.Tag {
		case "em1":
			return em1(el), errors.Nil()
		case "fett":
			return fett(el)
		case "fr":
			return fr(el)
		case "fremdsprache":
			return fremdsprache(el)
		case "gesperrt":
			return gesperrt(el)
		case "name":
			return name(el)
		case "op":
			return "", errors.Nil()
		case "romzahl":
			return romzahl(el)
		case "seite":
			return seite(el)
		case "trenn":
			return "", errors.Nil()
		case "zeile":
			return zeile(el)
		default:
			return "", errors.New(fmt.Errorf("unknown tag '%s' in %s element", el.Tag, elem.Tag), nil)
		}
	}
	return extractText(elem, switchFn)
}

func p(elem *etree.Element) (string, errors.UploadError) {
	switchFn := func(el *etree.Element) (string, errors.UploadError) {
		switch el.Tag {
		case "antiqua":
			return antiqua(el)
		case "bild":
			return bildBildverweis(el), errors.Nil()
		case "bildverweis":
			return bildBildverweis(el), errors.Nil()
		case "em1":
			return em1(el), errors.Nil()
		case "em2":
			return em2(el)
		case "fett":
			return fett(el)
		case "formel":
			return formel(el)
		case "fr":
			return fr(el)
		case "fremdsprache":
			return fremdsprache(el)
		case "gesperrt":
			return gesperrt(el)
		case "name":
			return name(el)
		case "op":
			return "", errors.Nil()
		case "romzahl":
			return romzahl(el)
		case "table":
			return table(), errors.Nil()
		case "seite":
			return seite(el)
		case "trenn":
			return "", errors.Nil()
		case "zeile":
			return zeile(el)
		default:
			return "", errors.New(fmt.Errorf("unknown tag '%s' in %s element", el.Tag, elem.Tag), nil)
		}
	}
	return extractText(elem, switchFn)
}

func seite(elem *etree.Element) (string, errors.UploadError) {
	page, err := extract.ExtractNumericAttribute(elem, "nr")
	if err.HasError {
		return "", err
	}
	return util.FmtPage(page), errors.Nil()
}

func table() string {
	return util.TableMatch
}

func footnote(elem *etree.Element) (model.Footnote, errors.UploadError) {
	switchFn := func(el *etree.Element) (string, errors.UploadError) {
		switch el.Tag {
		case "p":
			return p(el)
		default:
			return "", errors.New(fmt.Errorf("unknown tag '%s' in %s element", el.Tag, elem.Tag), nil)
		}
	}
	text, err := extractText(elem, switchFn)
	if err.HasError {
		return model.Footnote{}, err
	}
	page, err := extract.ExtractNumericAttribute(elem, "seite")
	if err.HasError {
		return model.Footnote{}, err
	}
	nr, err := extract.ExtractNumericAttribute(elem, "nr")
	if err.HasError {
		return model.Footnote{}, err
	}
	return model.Footnote{
		Page: page,
		Nr:   nr,
		Text: text,
	}, errors.Nil()
}

func summary(elem *etree.Element) (model.Summary, errors.UploadError) {
	switchFn := func(el *etree.Element) (string, errors.UploadError) {
		switch el.Tag {
		case "p":
			return p(el)
		default:
			return "", errors.New(fmt.Errorf("unknown tag '%s' in %s element", el.Tag, elem.Tag), nil)
		}
	}
	text, err := extractText(elem, switchFn)
	if err.HasError {
		return model.Summary{}, err
	}
	page, err := extract.ExtractNumericAttribute(elem, "seite")
	if err.HasError {
		return model.Summary{}, err
	}
	line, err := extract.ExtractNumericAttribute(elem, "anfang")
	if err.HasError {
		return model.Summary{}, err
	}
	return model.Summary{
		Page: page,
		Line: line,
		Text: text,
	}, errors.Nil()
}

func removePunctuation(s string) string {
	var result []rune
	length := len(s)
	for i, r := range s {
		if r == ':' {
			if i > 0 && i < length-1 {
				result = append(result, r)
			}
			continue
		}
		if unicode.IsPunct(r) {
			continue
		}
		result = append(result, r)
	}
	return string(result)
}

func fixCapitalization(s string) string {
	if len(s) == 0 {
		return s
	}
	runes := []rune(s)
	if isAllUpperCase(s) {
		return string(runes[:1]) + strings.ToLower(string(runes[1:]))
	}
	return strings.ToUpper(string(runes[:1])) + string(runes[1:])
}

func isAllUpperCase(s string) bool {
	for _, r := range s {
		if unicode.IsLetter(r) && !unicode.IsUpper(r) {
			return false
		}
	}
	return true
}
