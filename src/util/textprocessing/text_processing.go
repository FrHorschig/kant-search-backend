package textprocessing

import (
	"encoding/json"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

func GetParagraphsAndSentences(text string) (paragraphByNumber map[int]string, sentencesByParagraphNumber map[int][]string, err error) {
	paragraphByNumber = splitIntoParagraphs(text)
	sentencesByParagraphNumber = make(map[int][]string)
	for paragraphNumber, paragraph := range paragraphByNumber {
		sentences, err := splitIntoSentences(paragraph)
		if err != nil {
			return nil, nil, err
		}
		sentencesByParagraphNumber[paragraphNumber] = sentences
	}
	return paragraphByNumber, sentencesByParagraphNumber, nil
}

func GetPages(paragraph string) ([]int32, error) {
	r := regexp.MustCompile(`\{p(\d+)\}`)
	matches := r.FindAllStringSubmatch(paragraph, -1)
	pages := make([]int32, 0)
	for _, match := range matches {
		num, err := strconv.Atoi(match[1])
		if err != nil {
			return nil, err
		}
		pages = append(pages, int32(num))
	}
	return pages, nil
}

func splitIntoParagraphs(text string) map[int]string {
	split := strings.Split(text, "{r}")
	textByNumber := make(map[int]string)
	index := 0
	for _, substring := range split {
		trimmed := strings.TrimSpace(substring)
		if trimmed != "" {
			textByNumber[index] = trimmed
			index++
		}
	}
	return textByNumber
}

func splitIntoSentences(text string) ([]string, error) {
	pythonPath := "src_py/.venv/bin/python3"
	cmd := exec.Command(pythonPath, "src_py/split_into_sentences.py", text)
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	sentences := make([]string, 0)
	err = json.Unmarshal(output, &sentences)
	if err != nil {
		return nil, err
	}
	return sentences, nil
}
