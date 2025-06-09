package transform

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/beevik/etree"
	"github.com/frhorschig/kant-search-backend/core/upload/errs"
	"github.com/frhorschig/kant-search-backend/core/upload/internal/util"
)

func antiqua(elem *etree.Element) (string, errs.UploadError) {
	switchFn := func(el *etree.Element) (string, errs.UploadError) {
		switch el.Tag {
		case "fett":
			return fett(el)
		case "gesperrt":
			return gesperrt(el)
		case "name":
			return name(el)
		case "seite":
			return seite(el)
		case "trenn":
			return "", errs.Nil()
		case "zeile":
			return zeile(el)
		default:
			return "", errs.New(fmt.Errorf("unknown tag '%s' in %s element", el.Tag, elem.Tag), nil)
		}
	}
	return extractText(elem, switchFn)
}

func bildBildverweis(elem *etree.Element) string {
	return util.FmtImage(
		strings.TrimSpace(elem.SelectAttrValue("src", "MISSING_IMG_SRC")),
		strings.TrimSpace(elem.SelectAttrValue("beschreibung", "MISSING_IMG_DESC")),
	)
}

func em1(elem *etree.Element) string {
	return util.FmtEmph(strings.TrimSpace(elem.Text()))
}

func em2(elem *etree.Element) (string, errs.UploadError) {
	switchFn := func(el *etree.Element) (string, errs.UploadError) {
		switch el.Tag {
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
		case "romzahl":
			return romzahl(el)
		case "seite":
			return seite(el)
		case "trenn":
			return "", errs.Nil()
		case "zeile":
			return zeile(el)
		default:
			return "", errs.New(fmt.Errorf("unknown tag '%s' in %s element", el.Tag, elem.Tag), nil)
		}
	}
	extracted, err := extractText(elem, switchFn)
	if err.HasError {
		return "", err
	}
	return util.FmtEmph2(extracted), errs.Nil()
}

func fett(elem *etree.Element) (string, errs.UploadError) {
	switchFn := func(el *etree.Element) (string, errs.UploadError) {
		switch el.Tag {
		case "seite":
			return seite(el)
		case "zeile":
			return zeile(el)
		case "trenn":
			return "", errs.Nil()
		default:
			return "", errs.New(fmt.Errorf("unknown tag '%s' in %s element", el.Tag, elem.Tag), nil)
		}
	}
	extracted, err := extractText(elem, switchFn)
	if err.HasError {
		return "", err
	}
	return util.FmtBold(extracted), errs.Nil()
}

func formel(elem *etree.Element) (string, errs.UploadError) {
	switchFn := func(el *etree.Element) (string, errs.UploadError) {
		switch el.Tag {
		case "em1":
			return em1(el), errs.Nil()
		default:
			return "", errs.New(fmt.Errorf("unknown tag '%s' in %s element", el.Tag, elem.Tag), nil)
		}
	}
	extracted, err := extractText(elem, switchFn)
	if err.HasError {
		return "", err
	}
	return util.FmtFormula(extracted), errs.Nil()
}

func fr(elem *etree.Element) (string, errs.UploadError) {
	page, err := util.ExtractNumericAttribute(elem, "seite")
	if err.HasError {
		return "", err
	}
	nr, err := util.ExtractNumericAttribute(elem, "nr")
	if err.HasError {
		return "", err
	}
	return util.FmtFnRef(page, nr), errs.Nil()
}

func fremdsprache(elem *etree.Element) (string, errs.UploadError) {
	switchFn := func(el *etree.Element) (string, errs.UploadError) {
		switch el.Tag {
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
		case "romzahl":
			return romzahl(el)
		case "seite":
			return seite(el)
		case "trenn":
			return "", errs.Nil()
		case "zeile":
			return zeile(el)
		default:
			return "", errs.New(fmt.Errorf("unknown tag '%s' in %s element", el.Tag, elem.Tag), nil)
		}
	}
	extracted, err := extractText(elem, switchFn)
	if err.HasError {
		return "", err
	}
	return util.FmtLang(extracted), errs.Nil()
}

func gesperrt(elem *etree.Element) (string, errs.UploadError) {
	switchFn := func(el *etree.Element) (string, errs.UploadError) {
		switch el.Tag {
		case "fett":
			return fett(el)
		case "name":
			return name(el)
		case "seite":
			return seite(el)
		case "trenn":
			return "", errs.Nil()
		case "zeile":
			return zeile(el)
		default:
			return "", errs.New(fmt.Errorf("unknown tag '%s' in %s element", el.Tag, elem.Tag), nil)
		}
	}
	extracted, err := extractText(elem, switchFn)
	if err.HasError {
		return "", err
	}
	return util.FmtTracked(extracted), errs.Nil()
}

func name(elem *etree.Element) (string, errs.UploadError) {
	switchFn := func(el *etree.Element) (string, errs.UploadError) {
		switch el.Tag {
		case "seite":
			return seite(el)
		case "zeile":
			return zeile(el)
		case "trenn":
			return "", errs.Nil()
		default:
			return "", errs.New(fmt.Errorf("unknown tag '%s' in %s element", el.Tag, elem.Tag), nil)
		}
	}
	return extractText(elem, switchFn)
}

func romzahl(elem *etree.Element) (string, errs.UploadError) {
	content := strings.TrimSpace(elem.Text())
	re := regexp.MustCompile(`^(\d+)(\.)?$`)
	matches := re.FindStringSubmatch(content)
	if len(matches) == 0 {
		return "", errs.Nil()
	}
	num, err := strconv.ParseInt(matches[1], 10, 64)
	if err != nil {
		return "", errs.New(fmt.Errorf("can't convert string '%s' to number", matches[1]), nil)
	}
	return arabicToRoman(num) + matches[2], errs.Nil()
}

func zeile(elem *etree.Element) (string, errs.UploadError) {
	line, err := util.ExtractNumericAttribute(elem, "nr")
	if err.HasError {
		return "", err
	}
	return util.FmtLine(line), errs.Nil()
}

func extractForeignLangAttrs(el *etree.Element) string {
	result := ""
	lang := strings.TrimSpace(el.SelectAttrValue("sprache", ""))
	if lang != "" {
		result = result + ` lang="` + lang + `"`
	}
	alphabet := strings.TrimSpace(el.SelectAttrValue("zeichen", ""))
	if alphabet != "" {
		result = result + ` alphabet="` + alphabet + `"`
	}
	transcript := strings.TrimSpace(el.SelectAttrValue("umschrift", ""))
	if transcript != "" {
		result = result + ` transcript="` + transcript + `"`
	}
	return result
}
