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
	// TODO implement me: zeile, trenn, seite, gesperrt, name, fett
	return "", errors.NilError()
}

func bildBildverweis(elem *etree.Element) string {
	// TODO implement me
	// attributes: src, beschreibung, typ, Ort, z_anfang, z_ende, ausrichtung
	return fmt.Sprintf(
		"<ks-img>%s</ks-img>",
		strings.TrimSpace(elem.SelectAttrValue("src", "MISSING_IMG_SRC")),
	)
}

func em1(elem *etree.Element) string {
	return fmt.Sprintf(
		"<em>%s</em>",
		strings.TrimSpace(elem.Text()),
	)
}

func em2(elem *etree.Element) string {
	// TODO implement me
	text := elem.Text()
	// TODO check the AA scans to maybe find something better
	return fmt.Sprintf(
		"<ks-tracked>%s</ks-tracked>",
		strings.TrimSpace(text),
	)
}

func fett(elem *etree.Element) (string, errors.ErrorNew) {
	switchFn := func(el *etree.Element) (string, errors.ErrorNew) {
		switch el.Tag {
		case "seite":
			return seite(el), errors.NilError()
		case "zeile":
			return zeile(el), errors.NilError()
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
	return fmt.Sprintf("<b>%s</b>", extracted), errors.NilError()
}

func formel(elem *etree.Element) string {
	// TODO implement me: em1
	return fmt.Sprintf(
		"<ks-formula>%s</ks-formula>",
		strings.TrimSpace(elem.Text()),
	)
}

func fr(elem *etree.Element) string {
	return fmt.Sprintf(
		"<ks-fn-ref>%s.%s</ks-fn-ref>",
		strings.TrimSpace(elem.SelectAttrValue("seite", "MISSING_FR_PAGE")),
		strings.TrimSpace(elem.SelectAttrValue("nr", "MISSING_FR_NUMBER")),
	)
}

func fremdsprache(elem *etree.Element) (string, errors.ErrorNew) {
	// TODO implement me
	// attributes: sprache, zeichen, umschrift
	switchFn := func(el *etree.Element) (string, errors.ErrorNew) {
		switch el.Tag {
		case "bild":
			return bildBildverweis(el), errors.NilError()
		case "bildverweis":
			return bildBildverweis(el), errors.NilError()
		case "em1":
			return em1(el), errors.NilError()
		case "em2":
			return em2(el), errors.NilError()
		case "fett":
			return fett(el)
		case "formel":
			return formel(el), errors.NilError()
		case "fr":
			return fr(el), errors.NilError()
		case "fremdsprache":
			return fremdsprache(el)
		case "gesperrt":
			return gesperrt(el)
		case "name":
			return name(el)
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

func gesperrt(elem *etree.Element) (string, errors.ErrorNew) {
	switchFn := func(el *etree.Element) (string, errors.ErrorNew) {
		switch el.Tag {
		case "fett":
			return fett(el)
		case "name":
			return name(el)
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
	extracted, err := extractText(elem, switchFn)
	if err.HasError {
		return "", err
	}
	return fmt.Sprintf("<ks-tracked>%s</ks-tracked>", extracted), errors.NilError()
}

func name(elem *etree.Element) (string, errors.ErrorNew) {
	switchFn := func(el *etree.Element) (string, errors.ErrorNew) {
		switch el.Tag {
		case "seite":
			return seite(el), errors.NilError()
		case "zeile":
			return zeile(el), errors.NilError()
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
		return "", errors.NewError(nil, fmt.Errorf("error converting number: %v", err.Error()))
	}
	return arabicToRoman(num) + matches[2], errors.NilError()
}

func zeile(elem *etree.Element) string {
	return fmt.Sprintf(
		"<ks-line>%s</ks-line>",
		strings.TrimSpace(elem.SelectAttrValue("nr", "MISSING_LINE_NUMBER")),
	)
}
