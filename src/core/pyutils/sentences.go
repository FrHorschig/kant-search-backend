package pyutils

import (
	"bytes"
	"encoding/json"
	"os/exec"

	"github.com/FrHorschig/kant-search-backend/database/model"
)

func SplitIntoSentences(paragraphs []model.Paragraph) (map[int32][]string, error) {
	inputData, err := json.Marshal(paragraphs)
	if err != nil {
		return nil, err
	}

	pythonPath := "src_py/.venv/bin/python3"
	cmd := exec.Command(pythonPath, "src_py/split_into_sentences.py")
	cmd.Stdin = bytes.NewReader(inputData)
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	result := make(map[int32][]string)
	err = json.Unmarshal(output, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
