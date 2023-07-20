package textprocessing

import (
	"encoding/json"
	"os/exec"
	"strings"
)

func SplitIntoParagraphs(text string) map[int]string {
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

func SplitIntoSentences(text string) ([]string, error) {
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
