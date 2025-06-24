package transform

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/beevik/etree"
	"github.com/frhorschig/kant-search-backend/common/errs"
	"github.com/frhorschig/kant-search-backend/core/upload/internal/model"
	"github.com/frhorschig/kant-search-backend/core/upload/internal/util"
)

func Hx(el *etree.Element) (model.TreeHeading, errs.UploadError) {
	return hx(el)
}

func Hu(el *etree.Element) (string, errs.UploadError) {
	return hu(el)
}

func P(el *etree.Element) (string, errs.UploadError) {
	return p(el)
}

func Seite(el *etree.Element) (string, errs.UploadError) {
	return seite(el)
}

func Table() string {
	return table()
}

func Summary(el *etree.Element) (model.TreeSummary, errs.UploadError) {
	return summary(el)
}

func Footnote(el *etree.Element) (model.TreeFootnote, errs.UploadError) {
	return footnote(el)
}

func hx(elem *etree.Element) (model.TreeHeading, errs.UploadError) {
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
					return model.TreeHeading{}, err
				}
				tocTitle += fett
				textTitle += fett
			case "fr":
				fr, err := fr(el)
				if err.HasError {
					return model.TreeHeading{}, err
				}
				textTitle += fr
			case "fremdsprache":
				fremdsprache, err := fremdsprache(el)
				if err.HasError {
					return model.TreeHeading{}, err
				}
				tocTitle += fremdsprache
				textTitle += fremdsprache
			case "gesperrt":
				gesperrt, err := gesperrt(el)
				if err.HasError {
					return model.TreeHeading{}, err
				}
				tocTitle += gesperrt
				textTitle += gesperrt
			case "hi":
				tocTitle += strings.TrimSpace(el.Text())
			case "hu":
				hu, err := hu(el)
				if err.HasError {
					return model.TreeHeading{}, err
				}
				textTitle += hu
			case "name":
				name, err := name(el)
				if err.HasError {
					return model.TreeHeading{}, err
				}
				tocTitle += name
				textTitle += name
			case "op":
				continue
			case "romzahl":
				romzahl, err := romzahl(el)
				if err.HasError {
					return model.TreeHeading{}, err
				}
				tocTitle += romzahl
				textTitle += romzahl
			case "seite":
				page, err := seite(el)
				if err.HasError {
					return model.TreeHeading{}, err
				}
				textTitle += page
			case "trenn":
				continue
			case "zeile":
				line, err := zeile(el)
				if err.HasError {
					return model.TreeHeading{}, err
				}
				textTitle += line
			default:
				return model.TreeHeading{}, errs.New(fmt.Errorf("unknown tag '%s' in '%s' element", el.Tag, elem.Tag), nil)
			}
		}
		tocTitle += " "
		textTitle += " "
	}
	tocTitle = util.RemoveTags(tocTitle)
	tocTitle = strings.TrimSpace(tocTitle)
	tocTitle = removeTrailingPunctuation(tocTitle)
	tocTitle = fixCapitalization(tocTitle)
	return model.TreeHeading{
		Level:     level(elem),
		TocTitle:  tocTitle,
		TextTitle: strings.TrimSpace(textTitle),
	}, errs.Nil()
}

func hu(elem *etree.Element) (string, errs.UploadError) {
	switchFn := func(el *etree.Element) (string, errs.UploadError) {
		switch el.Tag {
		case "em1":
			return em1(el), errs.Nil()
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
			return "", errs.Nil()
		case "romzahl":
			return romzahl(el)
		case "seite":
			return seite(el)
		case "trenn":
			return "", errs.Nil()
		case "zeile":
			return zeile(el)
		default:
			return "", errs.New(fmt.Errorf("unknown tag '%s' in '%s' element", el.Tag, elem.Tag), nil)
		}
	}
	return extractText(elem, switchFn)
}

func p(elem *etree.Element) (string, errs.UploadError) {
	switchFn := func(el *etree.Element) (string, errs.UploadError) {
		switch el.Tag {
		case "antiqua":
			return antiqua(el)
		case "bild":
			return bildBildverweis(el), errs.Nil()
		case "bildverweis":
			return bildBildverweis(el), errs.Nil()
		case "em1":
			return em1(el), errs.Nil()
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
			return "", errs.Nil()
		case "romzahl":
			return romzahl(el)
		case "table":
			return table(), errs.Nil()
		case "seite":
			return seite(el)
		case "trenn":
			return "", errs.Nil()
		case "zeile":
			return zeile(el)
		default:
			return "", errs.New(fmt.Errorf("unknown tag '%s' in '%s' element", el.Tag, elem.Tag), nil)
		}
	}
	return extractText(elem, switchFn)
}

func seite(elem *etree.Element) (string, errs.UploadError) {
	page, err := util.ExtractNumericAttribute(elem, "nr")
	if err.HasError {
		return "", err
	}
	return util.FmtPage(page), errs.Nil()
}

func table() string {
	return util.TableMatch
}

func footnote(elem *etree.Element) (model.TreeFootnote, errs.UploadError) {
	switchFn := func(el *etree.Element) (string, errs.UploadError) {
		switch el.Tag {
		case "p":
			return p(el)
		default:
			return "", errs.New(fmt.Errorf("unknown tag '%s' in '%s' element", el.Tag, elem.Tag), nil)
		}
	}
	text, err := extractText(elem, switchFn)
	if err.HasError {
		return model.TreeFootnote{}, err
	}
	page, err := util.ExtractNumericAttribute(elem, "seite")
	if err.HasError {
		return model.TreeFootnote{}, err
	}
	nr, err := util.ExtractNumericAttribute(elem, "nr")
	if err.HasError {
		return model.TreeFootnote{}, err
	}
	return model.TreeFootnote{
		Page: page,
		Nr:   nr,
		Text: text,
	}, errs.Nil()
}

func summary(elem *etree.Element) (model.TreeSummary, errs.UploadError) {
	switchFn := func(el *etree.Element) (string, errs.UploadError) {
		switch el.Tag {
		case "p":
			return p(el)
		default:
			return "", errs.New(fmt.Errorf("unknown tag '%s' in '%s' element", el.Tag, elem.Tag), nil)
		}
	}
	text, err := extractText(elem, switchFn)
	if err.HasError {
		return model.TreeSummary{}, err
	}
	page, err := util.ExtractNumericAttribute(elem, "seite")
	if err.HasError {
		return model.TreeSummary{}, err
	}
	line, err := util.ExtractNumericAttribute(elem, "anfang")
	if err.HasError {
		return model.TreeSummary{}, err
	}
	return model.TreeSummary{
		Page: page,
		Line: line,
		Text: text,
	}, errs.Nil()
}

func removeTrailingPunctuation(s string) string {
	runes := []rune(s)
	end := len(runes)
	for end > 0 && unicode.IsPunct(runes[end-1]) {
		end--
	}
	return string(runes[:end])
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
