package transform

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/beevik/etree"
	"github.com/frhorschig/kant-search-backend/common/errors"
)

func antiqua(elem *etree.Element) (string, errors.ErrorNew) {
	switchFn := func(el *etree.Element) (string, errors.ErrorNew) {
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
			return "", errors.NilError()
		case "zeile":
			return zeile(el)
		default:
			return "", errors.NewError(fmt.Errorf("unknown tag '%s' in %s element", el.Tag, elem.Tag), nil)
		}
	}
	return extractText(elem, switchFn)
}

func bildBildverweis(elem *etree.Element) string {
	return fmt.Sprintf(
		`{image-extract src="%s" desc="%s"}`,
		strings.TrimSpace(elem.SelectAttrValue("src", "MISSING_IMG_SRC")),
		strings.TrimSpace(elem.SelectAttrValue("beschreibung", "MISSING_IMG_DESC")),
	)
}

func em1(elem *etree.Element) string {
	return fmt.Sprintf(
		"<ks-fmt-emph>%s</ks-fmt-emph>",
		strings.TrimSpace(elem.Text()),
	)
}

func em2(elem *etree.Element) (string, errors.ErrorNew) {
	switchFn := func(el *etree.Element) (string, errors.ErrorNew) {
		switch el.Tag {
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
			return "", errors.NilError()
		case "zeile":
			return zeile(el)
		default:
			return "", errors.NewError(fmt.Errorf("unknown tag '%s' in %s element", el.Tag, elem.Tag), nil)
		}
	}
	extracted, err := extractText(elem, switchFn)
	if err.HasError {
		return "", err
	}
	return fmt.Sprintf(
		"<ks-fmt-emph2>%s</ks-fmt-emph2>",
		extracted,
	), errors.NilError()
}

func fett(elem *etree.Element) (string, errors.ErrorNew) {
	switchFn := func(el *etree.Element) (string, errors.ErrorNew) {
		switch el.Tag {
		case "seite":
			return seite(el)
		case "zeile":
			return zeile(el)
		case "trenn":
			return "", errors.NilError()
		default:
			return "", errors.NewError(fmt.Errorf("unknown tag '%s' in %s element", el.Tag, elem.Tag), nil)
		}
	}
	extracted, err := extractText(elem, switchFn)
	if err.HasError {
		return "", err
	}
	return fmt.Sprintf("<ks-fmt-bold>%s</ks-fmt-bold>", extracted), errors.NilError()
}

func formel(elem *etree.Element) (string, errors.ErrorNew) {
	switchFn := func(el *etree.Element) (string, errors.ErrorNew) {
		switch el.Tag {
		case "em1":
			return em1(el), errors.NilError()
		default:
			return "", errors.NewError(fmt.Errorf("unknown tag '%s' in %s element", el.Tag, elem.Tag), nil)
		}
	}
	extracted, err := extractText(elem, switchFn)
	if err.HasError {
		return "", err
	}
	return fmt.Sprintf(
		"<ks-fmt-formula>%s</ks-fmt-formula>",
		extracted,
	), errors.NilError()
}

func fr(elem *etree.Element) (string, errors.ErrorNew) {
	page, err := extractNumericAttribute(elem, "seite")
	if err.HasError {
		return "", err
	}
	nr, err := extractNumericAttribute(elem, "nr")
	if err.HasError {
		return "", err
	}
	return fmt.Sprintf(
		"<ks-fmt-fnref>%d.%d</ks-fmt-fnref>",
		page,
		nr,
	), errors.NilError()
}

func fremdsprache(elem *etree.Element) (string, errors.ErrorNew) {
	switchFn := func(el *etree.Element) (string, errors.ErrorNew) {
		switch el.Tag {
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
			return "", errors.NilError()
		case "zeile":
			return zeile(el)
		default:
			return "", errors.NewError(fmt.Errorf("unknown tag '%s' in %s element", el.Tag, elem.Tag), nil)
		}
	}
	extracted, err := extractText(elem, switchFn)
	if err.HasError {
		return "", err
	}
	return fmt.Sprintf(
		`<ks-meta-lang%s>%s</ks-meta-lang>`,
		extractForeignLangAttrs(elem),
		extracted,
	), errors.NilError()
}

func gesperrt(elem *etree.Element) (string, errors.ErrorNew) {
	switchFn := func(el *etree.Element) (string, errors.ErrorNew) {
		switch el.Tag {
		case "fett":
			return fett(el)
		case "name":
			return name(el)
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
	extracted, err := extractText(elem, switchFn)
	if err.HasError {
		return "", err
	}
	return fmt.Sprintf("<ks-fmt-tracked>%s</ks-fmt-tracked>", extracted), errors.NilError()
}

func name(elem *etree.Element) (string, errors.ErrorNew) {
	switchFn := func(el *etree.Element) (string, errors.ErrorNew) {
		switch el.Tag {
		case "seite":
			return seite(el)
		case "zeile":
			return zeile(el)
		case "trenn":
			return "", errors.NilError()
		default:
			return "", errors.NewError(fmt.Errorf("unknown tag '%s' in %s element", el.Tag, elem.Tag), nil)
		}
	}
	return extractText(elem, switchFn)
}

func romzahl(elem *etree.Element) (string, errors.ErrorNew) {
	content := strings.TrimSpace(elem.Text())
	re := regexp.MustCompile(`^(\d+)(\.)?$`)
	matches := re.FindStringSubmatch(content)
	if len(matches) == 0 {
		return "", errors.NilError()
	}
	num, err := strconv.ParseInt(matches[1], 10, 64)
	if err != nil {
		return "", errors.NewError(nil, fmt.Errorf("error converting string '%s' to number: %v", matches[1], err.Error()))
	}
	return arabicToRoman(num) + matches[2], errors.NilError()
}

func zeile(elem *etree.Element) (string, errors.ErrorNew) {
	line, err := extractNumericAttribute(elem, "nr")
	if err.HasError {
		return "", err
	}
	return fmt.Sprintf(
		"<ks-meta-line>%d</ks-meta-line>",
		line,
	), errors.NilError()
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
