package pyutil

//go:generate mockgen -source=$GOFILE -destination=mocks/python_util_mock.go -package=mocks

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"

	"github.com/frhorschig/kant-search-backend/common/model"
)

type PythonUtil interface {
	SplitIntoSentences(paragraphs []model.Paragraph) (map[int32][]string, error)
}

type pythonUtilImpl struct {
}

func NewPythonUtil() PythonUtil {
	impl := pythonUtilImpl{}
	return &impl
}

// Returns a map of paragraph ids to sentences and error
func (util *pythonUtilImpl) SplitIntoSentences(paragraphs []model.Paragraph) (map[int32][]string, error) {
	inputData, err := json.Marshal(paragraphs)
	if err != nil {
		return nil, err
	}

	pyPath := os.Getenv("KSGO_PYTHON_PATH")
	if pyPath == "" {
		pyPath = "src_py"
	}

	cmd := exec.Command(pyPath+"/.venv/bin/python3", pyPath+"/split_into_sentences.py")
	cmd.Stdin = bytes.NewReader(inputData)
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	sentencesById := make(map[int32][]string)
	err = json.Unmarshal(output, &sentencesById)
	if err != nil {
		return nil, err
	}
	return sentencesById, nil
}
