package extract

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/frhorschig/kant-search-backend/common/errors"
	"github.com/frhorschig/kant-search-backend/core/upload/internal/model"
)

func ExtractFnRefs(text string) []string {
	re := regexp.MustCompile(model.FnRefExtract)
	matches := re.FindAllStringSubmatch(text, -1)
	result := []string{}
	for _, match := range matches {
		result = append(result, match[1])
	}
	return result
}

func ExtractPages(text string) ([]int32, errors.ErrorNew) {
	re := regexp.MustCompile(model.PageExtract)
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
