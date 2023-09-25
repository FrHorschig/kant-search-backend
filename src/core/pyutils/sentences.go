package pythonutils

import (
	"encoding/json"
	"os/exec"
)

func SplitIntoSentences(text string) ([]string, error) {
	pythonPath := "src_py/.venv/bin/python3"
	cmd := exec.Command(pythonPath, "src_py/split_into_sentences.py", text)
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var sentences []string
	err = json.Unmarshal(output, &sentences)
	if err != nil {
		return nil, err
	}
	return sentences, nil
}
