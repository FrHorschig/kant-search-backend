package syntax

import (
	"fmt"
)

type operator struct {
	isAnd bool
	isOr  bool
	isNot bool
	index int
}

func newAnd(index int) operator {
	return operator{isAnd: true, index: index}
}
func newOr(index int) operator {
	return operator{isOr: true, index: index}
}
func newNot(index int) operator {
	return operator{isNot: true, index: index}
}

type parenthesis struct {
	isOpen  bool
	isClose bool
	index   int
}

func newOpen(index int) parenthesis {
	return parenthesis{isOpen: true, index: index}
}
func newClose(index int) parenthesis {
	return parenthesis{isClose: true, index: index}
}

type word struct {
	text  string
	index int
}

func CheckSyntax(searchTerms string) (string, error) {
	ops, parens, words, lastIndex := tokenize(searchTerms)
	if !ruleHasAtLeastOneTerm(words) {
		return "", fmt.Errorf("you must provide at least one search term")
	}
	if ok, hint := ruleWordContainsOperator(words); !ok {
		return "", fmt.Errorf("search terms cannot contain operators: %s", hint)
	}
	if !ruleIsBeginOk(parens, words) {
		return "", fmt.Errorf("input must begin with opening bracket or a search term")
	}
	if !ruleIsEndOk(parens, words, lastIndex) {
		return "", fmt.Errorf("input must end with closing bracket or a search term")
	}
	if ok, hint := ruleIsOperatorOrderOk(ops); !ok {
		return "", fmt.Errorf("search terms cannot contain operators: %s", hint)
	}
	return "", nil // TODO frhorsch: implement me
}
