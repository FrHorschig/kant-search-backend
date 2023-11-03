package syntax

import (
	"github.com/frhorschig/kant-search-backend/core/errors"
	"github.com/frhorschig/kant-search-backend/core/syntax/internal"
)

func CheckSyntax(searchTerms string) (string, *errors.Error) {
	tokens, err := internal.Tokenize(searchTerms)
	if err != nil {
		return "", err
	}
	err = internal.Parse(tokens)
	if err != nil {
		return "", err
	}
	return internal.GetSearchString(tokens), nil
}
