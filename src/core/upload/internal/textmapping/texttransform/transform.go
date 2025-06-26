package texttransform

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/beevik/etree"
	"github.com/frhorschig/kant-search-backend/common/errs"
	"github.com/frhorschig/kant-search-backend/core/upload/internal/util"
)

// returns the fmtText, the tocText and an error
func Hx(xml string) (string, string, errs.UploadError) {
	return hx(createElement(xml))
}

func Hu(xml string) (string, errs.UploadError) {
	return hu(createElement(xml))
}

func P(xml string) (string, errs.UploadError) {
	return p(createElement(xml))
}

func Seite(xml string) (string, errs.UploadError) {
	return seite(createElement(xml))
}

func Table(xml string) (string, errs.UploadError) {
	return table(createElement(xml))
}

// return text, ref and error
func Summary(xml string) (string, string, errs.UploadError) {
	return summary(createElement(xml))
}

// return text, ref and error
func Footnote(xml string) (string, string, errs.UploadError) {
	return footnote(createElement(xml))
}

func hx(elem *etree.Element) (string, string, errs.UploadError) {
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
					return "", "", err
				}
				tocTitle += fett
				textTitle += fett
			case "fr":
				fr, err := fr(el)
				if err.HasError {
					return "", "", err
				}
				textTitle += fr
			case "fremdsprache":
				fremdsprache, err := fremdsprache(el)
				if err.HasError {
					return "", "", err
				}
				tocTitle += fremdsprache
				textTitle += fremdsprache
			case "gesperrt":
				gesperrt, err := gesperrt(el)
				if err.HasError {
					return "", "", err
				}
				tocTitle += gesperrt
				textTitle += gesperrt
			case "hi":
				tocTitle += strings.TrimSpace(el.Text())
			case "hu":
				hu, err := hu(el)
				if err.HasError {
					return "", "", err
				}
				textTitle += hu
			case "name":
				name, err := name(el)
				if err.HasError {
					return "", "", err
				}
				tocTitle += name
				textTitle += name
			case "op":
				continue
			case "romzahl":
				romzahl, err := romzahl(el)
				if err.HasError {
					return "", "", err
				}
				tocTitle += romzahl
				textTitle += romzahl
			case "seite":
				page, err := seite(el)
				if err.HasError {
					return "", "", err
				}
				textTitle += page
			case "trenn":
				continue
			case "zeile":
				line, err := zeile(el)
				if err.HasError {
					return "", "", err
				}
				textTitle += line
			default:
				return "", "", errs.New(fmt.Errorf("unknown tag '%s' in '%s' element", el.Tag, elem.Tag), nil)
			}
		}
		tocTitle += " "
		textTitle += " "
	}
	tocTitle = util.RemoveTags(tocTitle)
	tocTitle = strings.TrimSpace(tocTitle)
	tocTitle = removeTrailingPunctuation(tocTitle)
	tocTitle = fixCapitalization(tocTitle)
	return strings.TrimSpace(textTitle), tocTitle, errs.Nil()
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
			return table(el)
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

func table(elem *etree.Element) (string, errs.UploadError) {
	return util.TableMatch, errs.Nil()
}

func footnote(elem *etree.Element) (string, string, errs.UploadError) {
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
		return "", "", err
	}
	page, err := util.ExtractNumericAttribute(elem, "seite")
	if err.HasError {
		return "", "", err
	}
	nr, err := util.ExtractNumericAttribute(elem, "nr")
	if err.HasError {
		return "", "", err
	}
	return text, fmt.Sprintf("%d.%d", page, nr), errs.Nil()
}

func summary(elem *etree.Element) (string, string, errs.UploadError) {
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
		return "", "", err
	}
	page, err := util.ExtractNumericAttribute(elem, "seite")
	if err.HasError {
		return "", "", err
	}
	line, err := util.ExtractNumericAttribute(elem, "anfang")
	if err.HasError {
		return "", "", err
	}
	return text, fmt.Sprintf("%d.%d", page, line), errs.Nil()
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

func createElement(xml string) *etree.Element {
	doc := etree.NewDocument()
	doc.ReadFromString(xml)
	return doc.Root()
}
