package extract

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/frhorschig/kant-search-backend/common/errors"
	"github.com/frhorschig/kant-search-backend/core/upload/internal/model"
	dbmodel "github.com/frhorschig/kant-search-backend/dataaccess/model"
)

func ExtractFnRefs(text string) []string {
	re := regexp.MustCompile(model.FnRefMatch)
	matches := re.FindAllStringSubmatch(text, -1)
	result := []string{}
	for _, match := range matches {
		result = append(result, match[1])
	}
	return result
}

func ExtractPages(text string) ([]int32, errors.ErrorNew) {
	re := regexp.MustCompile(model.PageMatch)
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

	parts := strings.Split(text, strings.Split(model.PageFmt, "%")[0])
	if parts[0] != "" && len(result) > 0 {
		result = append([]int32{result[0] - 1}, result...)
	}

	return result, errors.NilError()
}

func FindParagraph(sec *dbmodel.Section, page int32, line int32) *dbmodel.Paragraph {
	for iPar := range sec.Paragraphs {
		par := &sec.Paragraphs[iPar]
		for _, pg := range par.Pages {
			if page == pg {
				parts := strings.Split(par.Text, fmt.Sprintf(model.PageFmt, page))
				part := ""
				if len(parts) == 1 {
					part = par.Text
				} else if len(parts) > 1 {
					part = strings.Split(parts[1], fmt.Sprintf(model.PageFmt, page+1))[0]
				}
				if strings.Contains(part, fmt.Sprintf(model.LineFmt, line)) {
					return par
				}
			}
		}
	}
	for i := range sec.Sections {
		return FindParagraph(&sec.Sections[i], page, line)
	}
	return nil
}
