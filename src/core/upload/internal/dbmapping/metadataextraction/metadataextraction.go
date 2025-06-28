package metadataextraction

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/frhorschig/kant-search-backend/common/errs"
	"github.com/frhorschig/kant-search-backend/core/upload/internal/common/util"
	"github.com/frhorschig/kant-search-backend/dataaccess/model"
)

func ExtractMetadata(contents []model.Content) errs.UploadError {
	for i := range contents {
		c := &contents[i]
		wordIndexMap, err := findWordIndexMap(c.SearchText, c.FmtText)
		if err != nil {
			return errs.New(nil, fmt.Errorf("unable to create word index map: %v", err.Error()))
		}
		c.WordIndexMap = wordIndexMap
		pageByIndex, err := findNumberByIndex(c.FmtText, util.PageMatch)
		if err != nil {
			return errs.New(nil, fmt.Errorf("unable to create index->page map: %v", err.Error()))
		}
		c.PageByIndex = pageByIndex
		lineByIndex, err := findNumberByIndex(c.FmtText, util.LineMatch)
		if err != nil {
			return errs.New(nil, fmt.Errorf("unable to create index->line map: %v", err.Error()))
		}
		c.LineByIndex = lineByIndex
	}
	return errs.Nil()
}

func findWordIndexMap(rawText string, fmtText string) (map[int32]int32, error) {
	rawWords := extractWordData(rawText)
	fmtWords := extractWordData(util.MaskTags(fmtText)) // we mask the tags here so that no word from a tag in fmtText may be accidentally matched to a normal word in rawText; this could happen if there are tags with attributes (like image or table tags)
	if len(rawWords) != len(fmtWords) {
		// and because we mask all known tags, we (should) get the exact same number of words in both cases ...
		return nil, fmt.Errorf("unequal number of words in searchText and fmtText: {%s} vs {%s}", fmtText, rawText)
	}

	result := make(map[int32]int32)
	for i, rawWord := range rawWords {
		fmtWord := fmtWords[i]
		if rawWord.Text != fmtWord.Text {
			return nil, fmt.Errorf("unequal matched words '%s' at index %d in searchText and '%s' at %d in fmtText", rawWord.Text, rawWord.Index, fmtWord.Text, fmtWord.Index)
		}
		// ... which then makes it super simple to map the indices
		result[rawWord.Index] = fmtWord.Index
	}
	return result, nil
}

type WordData struct {
	Text  string
	Index int32
}

func extractWordData(text string) []WordData {
	words := []WordData{}
	var currentWord strings.Builder
	startIndex := 0
	i := 0
	for _, r := range text {
		if unicode.IsLetter(r) {
			if currentWord.Len() == 0 {
				startIndex = i
			}
			currentWord.WriteRune(r)
		} else if currentWord.Len() > 0 {
			words = append(words, WordData{
				Text:  currentWord.String(),
				Index: int32(startIndex),
			})
			currentWord.Reset()
		}
		i++
	}
	if currentWord.Len() > 0 {
		words = append(words, WordData{
			Text:  currentWord.String(),
			Index: int32(startIndex),
		})
	}
	return words
}

func findNumberByIndex(text string, matcher string) ([]model.IndexNumberPair, error) {
	result := []model.IndexNumberPair{}
	re := regexp.MustCompile(matcher)
	for _, match := range re.FindAllStringSubmatchIndex(text, -1) {
		num, err := strconv.ParseInt(text[match[2]:match[3]], 10, 32)
		if err != nil {
			return nil, err
		}
		result = append(result, model.IndexNumberPair{
			I:   int32(utf8.RuneCountInString(text[:match[0]])),
			Num: int32(num),
		})
	}
	return result, nil
}
