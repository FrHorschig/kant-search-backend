package syntax

import "errors"

func CheckSyntax(searchTerms string) (string, error) {
	if wrongBeginChar(searchTerms[0]) {
		return "", errors.New("search input must not start with &, | or )")
	}
	if wrongEndChar(searchTerms[len(searchTerms)-1]) {
		return "", errors.New("search input must not end with &, |, ! or (")
	}
	// tokens := tokenizer.Tokenize(searchTerms)
	return "", nil // TODO frhorsch: implement me
}

func wrongBeginChar(c byte) bool {
	return c == '&' || c == '|' || c == ')'
}

func wrongEndChar(c byte) bool {
	return c == '&' || c == '|' || c == '!' || c == '('
}
