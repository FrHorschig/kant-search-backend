package extract

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/frhorschig/kant-search-backend/core/upload/errors"
	"github.com/frhorschig/kant-search-backend/core/upload/internal/util"
)

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
			return nil, errors.NewError(nil, fmt.Errorf("error converting string '%s' to number: %v", nStr, err.Error()))
		}
		result = append(result, int32(n))
	}

	return result, errors.NilError()
}

func RemoveTags(text string) string {
	re := regexp.MustCompile(util.FnRefMatch)
	text = re.ReplaceAllString(text, "")
	re = regexp.MustCompile(util.LineMatch)
	text = re.ReplaceAllString(text, "")
	re = regexp.MustCompile(util.PageMatch)
	text = re.ReplaceAllString(text, "")
	re = regexp.MustCompile(`<[^>]*>`)
	text = re.ReplaceAllString(text, "")
	re = regexp.MustCompile(`\s+`)
	text = re.ReplaceAllString(text, " ")
	return strings.TrimSpace(text)
}
