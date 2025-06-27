package texttransform

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/beevik/etree"
	"github.com/frhorschig/kant-search-backend/common/errs"
	"github.com/frhorschig/kant-search-backend/core/upload/internal/model"
	"github.com/rs/zerolog/log"
)

func level(el *etree.Element) model.TreeLevel {
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

func extractText(elem *etree.Element, switchFn func(el *etree.Element) (string, errs.UploadError)) (string, errs.UploadError) {
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
		} else if childEl, ok := ch.(*etree.Comment); ok {
			log.Debug().Msgf("Comment: '%s'", childEl.Data)
			continue
		} else {
			return "", errs.New(nil, fmt.Errorf("unknown child type '%v' in tag '%v', it is neither CharData nor Element nor Comment", ch.(*etree.Element).Tag, elem.Tag))
		}
		text += " "
	}
	return contractSpaces(text), errs.Nil()
}

func contractSpaces(s string) string {
	re := regexp.MustCompile(`\s+`)
	s = re.ReplaceAllString(s, " ")
	re = regexp.MustCompile(`\s+([.,:;?!])`)
	s = re.ReplaceAllString(s, `$1`)
	return strings.TrimSpace(s)
}
