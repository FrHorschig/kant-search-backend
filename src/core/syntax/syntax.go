package syntax

import "github.com/FrHorschig/kant-search-backend/core/syntax/internal"

func CheckSyntax(searchTerms string) (string, error) {
	tokens, err := internal.Tokenize(searchTerms)
	if err != nil {
		return "", err
	}
	err = internal.CheckSyntax(&tokens)
	if err != nil {
		return "", err
	}
	return internal.GetSearchString(tokens), nil
}
