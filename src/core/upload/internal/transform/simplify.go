package transform

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func Simplify(xml string) string {
	// We don't replace <seite ...> with {p...}, because the <seite> elements are sometimes on the same level as the headings and paragraphs. Replacing them here would make finding them later more difficult.
	reZeile := regexp.MustCompile(`<zeile\s+nr="(\d+)"\s*/>`)
	xml = reZeile.ReplaceAllString(xml, `{l$1}`)
	reTrenn := regexp.MustCompile(`<trenn/>`)
	xml = reTrenn.ReplaceAllString(xml, "")

	reRomzahl := regexp.MustCompile(`<romzahl>\s*(\d+)\.\s*</romzahl>`)
	xml = reRomzahl.ReplaceAllStringFunc(xml, func(match string) string {
		numStr := reRomzahl.FindStringSubmatch(match)[1]
		num, err := strconv.ParseInt(numStr, 10, 32)
		if err != nil {
			return match
		}
		return fmt.Sprintf(" %s. ", arabicToRoman(int32(num)))
	})

	reSpace := regexp.MustCompile(` +`)
	xml = reSpace.ReplaceAllString(xml, " ")
	return xml
}

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
