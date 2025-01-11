package transform

import (
	"strings"

	"github.com/beevik/etree"
	"github.com/frhorschig/kant-search-backend/common/errors"
)

func arabicToRoman(number int32) string {
	conversions := []struct {
		value int32
		digit string
	}{
		{1000, "M"},
		{900, "CM"},
		{500, "D"},
		{400, "CD"},
		{100, "C"},
		{90, "XC"},
		{50, "L"},
		{40, "XL"},
		{10, "X"},
		{9, "IX"},
		{5, "V"},
		{4, "IV"},
		{1, "I"},
	}
	var roman strings.Builder
	for _, conversion := range conversions {
		for number >= conversion.value {
			roman.WriteString(conversion.digit)
			number -= conversion.value
		}
	}
	return roman.String()
}

func extractText(element *etree.Element, switchFn func(el *etree.Element) (string, errors.ErrorNew)) (string, errors.ErrorNew) {
	text := ""
	for _, ch := range element.Child {
		if str, ok := ch.(*etree.CharData); ok {
			text += strings.TrimSpace(str.Data)
		} else if childEl, ok := ch.(*etree.Element); ok {
			extracted, err := switchFn(childEl)
			if extracted == "" {
				continue
			}
			if err.HasError {
				return "", err
			}
			text += extracted
		}
		text += " "
	}
	return strings.TrimSpace(text), errors.NilError()
}
