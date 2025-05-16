package extract

import (
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/beevik/etree"
	"github.com/frhorschig/kant-search-backend/core/upload/errors"
	"github.com/frhorschig/kant-search-backend/core/upload/internal/util"
)

func ExtractNumericAttribute(el *etree.Element, attr string) (int32, errors.UploadError) {
	defaultStr := "DEFAULT_STRING"
	nStr := strings.TrimSpace(el.SelectAttrValue(attr, defaultStr))
	if nStr == defaultStr {
		return 0, errors.New(fmt.Errorf("missing '%s' attribute in '%s' element", attr, el.Tag), nil)
	}

	// TODO make this configurable at runtime
	if slices.Contains([]string{"272a", "272c", "272d"}, nStr) {
		nStr = "272"
	}

	n, err := strconv.ParseInt(nStr, 10, 32)
	if err != nil {
		return 0, errors.New(fmt.Errorf("can't convert attribute string '%s' to number", nStr), nil)
	}
	return int32(n), errors.Nil()
}

func ExtractFnRefs(text string) []string {
	re := regexp.MustCompile(util.FnRefMatch)
	matches := re.FindAllStringSubmatch(text, -1)
	result := []string{}
	for _, match := range matches {
		result = append(result, match[1])
	}
	return result
}

func ExtractPages(text string) ([]int32, errors.UploadError) {
	re := regexp.MustCompile(util.PageMatch)
	matches := re.FindAllStringSubmatch(text, -1)

	result := []int32{}
	for _, match := range matches {
		nStr := match[1]
		n, err := strconv.ParseInt(nStr, 10, 32)
		if err != nil {
			return nil, errors.New(fmt.Errorf("can't convert page string '%s' to number", nStr), nil)
		}
		result = append(result, int32(n))
	}

	return result, errors.Nil()
}

func RemoveTags(text string) string {
	re := regexp.MustCompile(util.FnRefMatch)
	text = re.ReplaceAllString(text, "")
	re = regexp.MustCompile(util.LineMatch)
	text = re.ReplaceAllString(text, "")
	re = regexp.MustCompile(util.PageMatch)
	text = re.ReplaceAllString(text, "")
	re = regexp.MustCompile(util.SummaryRefMatch)
	text = re.ReplaceAllString(text, "")
	re = regexp.MustCompile(`<[^>]*>`)
	text = re.ReplaceAllString(text, "")
	re = regexp.MustCompile(`\s+`)
	text = re.ReplaceAllString(text, " ")
	return strings.TrimSpace(text)
}
