package syntax

import (
	"strings"
)

func ruleHasAtLeastOneTerm(words []word) bool {
	return len(words) > 0
}

func ruleWordContainsOperator(words []word) (bool, string) {
	for _, w := range words {
		if containsOperator(w.text) {
			return false, w.text
		}
	}
	return true, ""
}

func ruleIsBeginOk(parens []parenthesis, words []word) bool {
	if len(parens) > 0 && parens[0].isOpen {
		return true
	}
	if len(words) > 0 && words[0].index == 0 {
		return true
	}
	return false
}

func ruleIsEndOk(parens []parenthesis, words []word, lastIndex int) bool {
	if len(parens) > 0 && parens[len(parens)-1].index == lastIndex && parens[len(parens)-1].isClose {
		return true
	}
	if len(words) > 0 && words[len(words)-1].index == lastIndex {
		return true
	}
	return false
}

func ruleIsOperatorOrderOk(ops []operator) (bool, string) {
	for i := 0; i < len(ops)-1; i++ {
		if (ops[i+1].index - ops[i].index) < 2 {
			return false, "TODO"
		}
	}
	return true, ""
}

func containsOperator(text string) bool {
	opArr := []string{"&", "|", "!", "(", ")"}
	for _, o := range opArr {
		if strings.Contains(text, o) {
			return true
		}
	}
	return false
}
