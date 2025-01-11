package transform

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/beevik/etree"
	"github.com/frhorschig/kant-search-backend/common/errors"
)

func Seite(seite *etree.Element) string {
	// TODO improvement: handle 'satz' attribute
	return fmt.Sprintf(
		"<ks-page>%s</ks-page>",
		seite.SelectAttrValue("nr", "MISSING_PAGE_NUMBER"),
	)
}

func em1(em1 *etree.Element) string {
	return fmt.Sprintf(
		"<em>%s</em>",
		strings.TrimSpace(em1.Text()),
	)
}

func fett(fett *etree.Element) (string, errors.ErrorNew) {
	switchFn := func(el *etree.Element) (string, errors.ErrorNew) {
		switch el.Tag {
		case "seite":
			return Seite(el), errors.NilError()
		case "zeile":
			return zeile(el), errors.NilError()
		case "trenn":
			return "", errors.NilError()
		default:
			return "", errors.NewError(fmt.Errorf("unknown tag '%s' in %s element", el.Tag, fett.Tag), nil)
		}
	}
	extracted, err := extractText(fett, switchFn)
	if err.HasError {
		return "", err
	}
	return fmt.Sprintf("<b>%s</b>", extracted), errors.NilError()
}

func fr(fr *etree.Element) string {
	return fmt.Sprintf(
		"<ks-fn-ref>%s.%s</ks-fn-ref>",
		strings.TrimSpace(fr.SelectAttrValue("seite", "MISSING_FR_PAGE")),
		strings.TrimSpace(fr.SelectAttrValue("nr", "MISSING_FR_NUMBER")),
	)
}

func fremdsprache(fs *etree.Element) (string, errors.ErrorNew) {
	// TODO implement me
	// attribute: sprache, zeichen, umschrift
	switchFn := func(el *etree.Element) (string, errors.ErrorNew) {
		switch el.Tag {
		case "bild":
			return bild(el), errors.NilError()
		case "bildverweis":
			return bildverweis(el), errors.NilError()
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
			return Seite(el), errors.NilError()
		case "trenn":
			return "", errors.NilError()
		case "zeile":
			return zeile(el), errors.NilError()
		default:
			return "", errors.NewError(fmt.Errorf("unknown tag '%s' in %s element", el.Tag, fs.Tag), nil)
		}
	}
	return extractText(fs, switchFn)
}

func gesperrt(gesperrtEl *etree.Element) (string, errors.ErrorNew) {
	switchFn := func(el *etree.Element) (string, errors.ErrorNew) {
		switch el.Tag {
		case "fett":
			return fett(el)
		case "name":
			return name(el)
		case "seite":
			return Seite(el), errors.NilError()
		case "trenn":
			return "", errors.NilError()
		case "zeile":
			return zeile(el), errors.NilError()
		default:
			return "", errors.NewError(fmt.Errorf("unknown tag '%s' in %s element", el.Tag, gesperrtEl.Tag), nil)
		}
	}
	extracted, err := extractText(gesperrtEl, switchFn)
	if err.HasError {
		return "", err
	}
	return fmt.Sprintf("<ks-tracked>%s</em>", extracted), errors.NilError()
}

func name(name *etree.Element) (string, errors.ErrorNew) {
	switchFn := func(el *etree.Element) (string, errors.ErrorNew) {
		switch el.Tag {
		case "seite":
			return Seite(el), errors.NilError()
		case "zeile":
			return zeile(el), errors.NilError()
		case "trenn":
			return "", errors.NilError()
		default:
			return "", errors.NewError(fmt.Errorf("unknown tag '%s' in %s element", el.Tag, name.Tag), nil)
		}
	}
	return extractText(name, switchFn)
}

func romzahl(romzahl *etree.Element) (string, errors.ErrorNew) {
	content := strings.TrimSpace(romzahl.Text())
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

func zeile(zeile *etree.Element) string {
	return fmt.Sprintf(
		"<ks-line>%s</ks-line>",
		zeile.SelectAttrValue("nr", "MISSING_LINE_NUMBER"),
	)
}
