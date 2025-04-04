package transform

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/beevik/etree"
	"github.com/frhorschig/kant-search-backend/common/errors"
	"github.com/frhorschig/kant-search-backend/core/upload/internal/model"
)

func ExtractNumericAttribute(el *etree.Element, attr string) (int32, errors.ErrorNew) {
	defaultStr := "DEFAULT_STRING"
	nStr := strings.TrimSpace(el.SelectAttrValue(attr, defaultStr))
	if nStr == defaultStr {
		return 0, errors.NewError(fmt.Errorf("missing '%s' attribute in '%s' element", attr, el.Tag), nil)
	}
	n, err := strconv.ParseInt(nStr, 10, 32)
	if err != nil {
		return 0, errors.NewError(nil, fmt.Errorf("error converting string '%s' to number: %v", nStr, err.Error()))
	}
	return int32(n), errors.NilError()
}

func level(el *etree.Element) model.Level {
	switch el.Tag {
	case "h1":
		return model.HWork
	case "h2":
		return model.H1
	case "h3":
		return model.H2
	case "h4":
		return model.H3
	case "h5":
		return model.H4
	case "h6":
		return model.H5
	case "h7":
		return model.H6
	case "h8":
		return model.H7
	}
	return model.H8
}

func arabicToRoman(number int64) string {
	conversions := []struct {
		value int64
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

func extractText(elem *etree.Element, switchFn func(el *etree.Element) (string, errors.ErrorNew)) (string, errors.ErrorNew) {
	text := ""
	for _, ch := range elem.Child {
		if str, ok := ch.(*etree.CharData); ok {
			text += strings.TrimSpace(str.Data)
		} else if childEl, ok := ch.(*etree.Element); ok {
			extracted, err := switchFn(childEl)
			if err.HasError {
				return "", err
			}
			if extracted == "" {
				continue
			}
			text += extracted
		} else {
			return "", errors.NewError(nil, fmt.Errorf("unknown child type in tag '%v', it is neither CharData nor Element", elem.Tag))
		}
		text += " "
	}
	return strings.TrimSpace(text), errors.NilError()
}
